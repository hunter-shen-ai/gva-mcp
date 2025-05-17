<template>
  <div class="mcp-tools-container">
    <h2 class="page-title">可用工具</h2>
    <div v-if="loading" class="loading-text">加载中...</div>
    <div v-if="fetchError" class="error-text">{{ fetchError }}</div>
    <div
      v-if="!loading && !fetchError && (!tools || tools.length === 0)"
      class="no-tools-text"
    >
      暂无可用工具
    </div>
    <div class="tools-grid">
      <div v-for="tool in tools" :key="tool.id" class="tool-item-wrapper">
        <div class="tool-item" @click="openToolModal(tool)">
          <div class="tool-header">
            <div class="tool-info">
              <span class="tool-name">{{ tool.name }}</span>
              <span class="tool-id">{{ tool.id }}</span>
            </div>
          </div>
          <div class="tool-description">{{ tool.description }}</div>
        </div>
      </div>
    </div>

    <!-- 工具详情模态框 -->
    <div
      v-if="activeToolModal"
      class="tool-modal-overlay"
      @click.self="closeToolModal"
    >
      <div class="tool-modal">
        <div class="tool-modal-header">
          <div class="tool-modal-title">
            <span class="tool-name">{{ activeToolModal.name }}</span>
            <span class="tool-id">{{ activeToolModal.id }}</span>
          </div>
          <button class="close-modal-btn" @click="closeToolModal">×</button>
        </div>

        <div class="tool-modal-body">
          <div class="tool-modal-description">
            {{ activeToolModal.description }}
          </div>

          <div class="parameters-title">输入模式:</div>
          <div
            v-if="
              activeToolModal.parameters &&
              activeToolModal.parameters.length > 0
            "
            class="tool-parameters"
          >
            <div
              v-for="(param, index) in activeToolModal.parameters"
              :key="index"
              class="param-item"
            >
              <div class="param-header">
                <span class="param-name">{{ param.name }}</span>
                <span v-if="param.required" class="param-required"
                  >Required</span
                >
              </div>
              <el-input
                :type="param.type === 'number' ? 'number' : 'text'"
                :placeholder="param.description"
                v-model="param.value"
                class="param-input"
                clearable
              />
              <div class="param-details">
                <span class="param-type">{{ param.type }}</span>
                <span class="param-description-text">{{
                  param.description
                }}</span>
              </div>
            </div>
          </div>
          <div v-else class="no-params-text">无输入参数</div>

          <div class="tool-execution-section">
            <el-button
              type="primary"
              @click="handleTestTool(activeToolModal)"
              :disabled="
                !isToolTestable(activeToolModal) || activeToolModal.testing
              "
              :loading="activeToolModal.testing"
            >
              执行测试
            </el-button>
            <div v-if="activeToolModal.testing" class="execution-status-text">
              执行中...
            </div>
            <div
              v-if="activeToolModal.executionError"
              class="execution-result error"
            >
              <pre>{{ activeToolModal.executionError }}</pre>
            </div>
            <div
              v-if="activeToolModal.executionResult"
              class="execution-result success"
            >
              <strong>结果:</strong>
              <pre>{{ activeToolModal.executionResult }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { getMcpClientToolsList, executeTool } from '../api/mcpClientTools'

  // 响应式数据
  const loading = ref(false)
  const fetchError = ref(null)
  const tools = ref([])
  const activeToolModal = ref(null) // 当前激活的工具模态框

  // 获取工具列表数据
  const fetchToolsData = async () => {
    loading.value = true
    fetchError.value = null
    try {
      const response = await getMcpClientToolsList()
      if (response.code === 0 && Array.isArray(response.data)) {
        tools.value = response.data.map((tool) => ({
          id: tool.id,
          name: tool.func_name,
          originalFuncName: tool.func_name,
          description: tool.func_description,
          enabled: tool.enabled === undefined ? true : !!tool.enabled,
          parameters: mapApiParameters(tool.parameters || []),
          testing: false,
          executionResult: null,
          executionError: null
        }))
      } else {
        fetchError.value = `获取工具列表失败: ${
          response.msg || '数据格式不正确'
        }`
        // 如果项目中配置了Element UI的Message组件，可以使用下面的代码
        // ElMessage.error(fetchError.value)
        console.error('API response error or unexpected format:', response)
        tools.value = []
      }
    } catch (err) {
      fetchError.value = `获取工具列表时发生网络错误: ${err.message}`
      // ElMessage.error(fetchError.value)
      console.error('Network error fetching tools:', err)
      tools.value = []
    } finally {
      loading.value = false
    }
  }

  // 映射API参数
  const mapApiParameters = (apiParams) => {
    if (!Array.isArray(apiParams)) return []
    return apiParams.map((p) => ({
      name: p.param_name,
      required: !!p.param_required,
      type: p.param_type || 'string',
      description: p.param_description || '',
      value: p.default_value || ''
    }))
  }

  // 打开工具模态框
  const openToolModal = (tool) => {
    // 清除之前的执行结果
    tool.executionResult = null
    tool.executionError = null
    activeToolModal.value = tool
    // 添加body类阻止背景滚动
    document.body.classList.add('modal-open')
  }

  // 关闭工具模态框
  const closeToolModal = () => {
    activeToolModal.value = null
    // 移除body类允许背景滚动
    document.body.classList.remove('modal-open')
  }

  // 检查工具是否可测试（必填参数是否已填写）
  const isToolTestable = (tool) => {
    if (!tool || !tool.parameters || tool.parameters.length === 0) {
      return true // 没有参数的工具可以直接测试
    }
    return tool.parameters.every((param) => {
      return (
        !param.required ||
        (param.value !== null &&
          param.value !== undefined &&
          param.value.toString().trim() !== '')
      )
    })
  }

  // 处理工具测试
  const handleTestTool = async (tool) => {
    if (!tool || !isToolTestable(tool)) {
      // ElMessage.warning('请填写所有必填参数!')
      console.warn('请填写所有必填参数!')
      return
    }

    tool.testing = true
    tool.executionResult = null
    tool.executionError = null

    const paramsPayload = {}
    tool.parameters.forEach((param) => {
      // 根据参数类型转换
      let value = param.value
      if (param.type === 'number' && value !== '' && !isNaN(Number(value))) {
        value = Number(value)
      } else if (param.type === 'boolean') {
        if (typeof value === 'string') {
          value = value.toLowerCase() === 'true'
        } else {
          value = !!value
        }
      }
      paramsPayload[param.name] = value
    })

    const payload = {
      tool_name: tool.originalFuncName,
      params: paramsPayload
    }

    try {
      const response = await executeTool(payload)
      if (response.code === 0) {
        tool.executionResult =
          typeof response.data === 'object'
            ? JSON.stringify(response.data, null, 2)
            : response.data
        // ElMessage.success(`工具 "${tool.name}" 执行成功!`)
        console.log(`工具 "${tool.name}" 执行成功!`)
      } else {
        tool.executionError = `执行失败 (Code: ${response.code}): ${
          response.msg || '未知错误'
        }`
        if (response.data) {
          tool.executionError += `\nDetails: ${
            typeof response.data === 'object'
              ? JSON.stringify(response.data, null, 2)
              : response.data
          }`
        }
        // ElMessage.error(tool.executionError)
        console.error(tool.executionError)
      }
    } catch (error) {
      console.error('Error executing tool:', error)
      tool.executionError = `执行 "${tool.name}" 时发生网络错误: ${
        error.message || '请检查网络连接或联系管理员'
      }`
      if (error.response && error.response.data) {
        tool.executionError += `\nServer Response: ${
          typeof error.response.data === 'object'
            ? JSON.stringify(error.response.data, null, 2)
            : error.response.data
        }`
      }
      // ElMessage.error(tool.executionError)
      console.error(tool.executionError)
    } finally {
      tool.testing = false
    }
  }

  // 生命周期钩子
  onMounted(() => {
    fetchToolsData()
  })

  // 注释掉了Element UI的Message通知的使用
  // 如果项目中配置了Element UI，请取消以下注释并导入相应组件
  /*
import { ElMessage } from 'element-plus'  // 或者 'element-ui'，取决于项目使用的UI库
*/
</script>

<style scoped>
  .mcp-tools-container {
    padding: 20px;
    background-color: #f7f8fa;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
      'Helvetica Neue', Arial, sans-serif;
  }

  .page-title {
    font-size: 18px;
    color: #303133;
    margin-bottom: 20px;
    font-weight: 500;
  }

  .loading-text,
  .error-text,
  .no-tools-text {
    text-align: center;
    padding: 20px;
    color: #909399;
    font-size: 14px;
  }

  .error-text {
    color: #f56c6c;
  }

  /* 工具网格布局 - 一行两列 */
  .tools-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr); /* 两列等宽 */
    gap: 15px; /* 卡片之间的间距 */
  }

  /* 手机屏幕上改为单列 */
  @media (max-width: 768px) {
    .tools-grid {
      grid-template-columns: 1fr; /* 单列 */
    }
  }

  .tool-item-wrapper {
    margin-bottom: 15px;
  }

  .tool-item {
    background-color: #fff;
    border: 1px solid #ebeef5;
    border-radius: 4px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    height: 100%; /* 保持高度一致 */
    cursor: pointer;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
  }

  .tool-item:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.15);
  }

  .tool-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px;
  }

  .tool-info {
    display: flex;
    align-items: center;
    flex: 1;
  }

  .tool-name {
    font-weight: 500;
    font-size: 16px;
    color: #303133;
  }

  .tool-id {
    font-size: 12px;
    color: #909399;
    margin-left: 8px;
    padding: 2px 0px;
    border-radius: 4px;
  }

  .tool-description {
    font-size: 14px;
    color: #606266;
    padding: 0 15px 15px 15px;
  }

  /* 模态框相关样式 */
  .tool-modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 9999;
  }

  .tool-modal {
    background-color: #fff;
    border-radius: 4px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    width: 90%;
    max-width: 800px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .tool-modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px 20px;
    border-bottom: 1px solid #ebeef5;
  }

  .tool-modal-title {
    display: flex;
    align-items: center;
  }

  .tool-modal-title .tool-name {
    font-size: 18px;
    font-weight: 500;
  }

  .tool-modal-title .tool-id {
    font-size: 12px;
    color: #909399;
    margin-left: 10px;
  }

  .close-modal-btn {
    background: none;
    border: none;
    font-size: 24px;
    color: #909399;
    cursor: pointer;
    padding: 0;
    margin: 0;
    line-height: 1;
  }

  .close-modal-btn:hover {
    color: #409eff;
  }

  .tool-modal-body {
    padding: 20px;
    overflow-y: auto;
    max-height: calc(90vh - 60px); /* 减去header的高度 */
  }

  .tool-modal-description {
    font-size: 14px;
    color: #606266;
    margin-bottom: 20px;
    padding-bottom: 15px;
    border-bottom: 1px solid #ebeef5;
  }

  /* 参数样式 */
  .parameters-title {
    font-size: 14px;
    color: #606266;
    margin-bottom: 15px;
    font-weight: 500;
  }

  .no-params-text {
    font-size: 13px;
    color: #c0c4cc;
    padding: 10px 0;
  }

  .param-item {
    margin-bottom: 15px;
    padding: 10px 0;
    border-radius: 4px;
  }

  .param-header {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
  }

  .param-name {
    font-weight: 500;
    color: #303133;
    font-size: 14px;
  }

  .param-required {
    background-color: #fef0f0;
    color: #f56c6c;
    font-size: 10px;
    font-weight: bold;
    padding: 2px 6px;
    border-radius: 3px;
    margin-left: 8px;
    text-transform: uppercase;
  }

  .param-input {
    width: 100%;
    margin-bottom: 5px;
  }

  .param-details {
    display: flex;
    align-items: center;
    font-size: 12px;
    color: #909399;
    margin-top: 5px;
  }

  .param-type {
    background-color: #f0f2f5;
    color: #909399;
    padding: 2px 6px;
    border-radius: 3px;
    margin-right: 8px;
    font-family: SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono',
      'Courier New', monospace;
  }

  .param-description-text {
    flex-grow: 1;
    white-space: pre-wrap;
    word-break: break-word;
  }

  .tool-execution-section {
    margin-top: 20px;
    padding-top: 15px;
    border-top: 1px dashed #ebeef5;
  }

  .execution-status-text {
    font-size: 13px;
    color: #409eff;
    margin-top: 5px;
  }

  .execution-result {
    margin-top: 10px;
    padding: 10px;
    border-radius: 4px;
    font-size: 13px;
    background-color: #f4f4f5;
    border: 1px solid #e9e9eb;
  }

  .execution-result.error {
    background-color: #fef0f0;
    border-color: #fde2e2;
    color: #f56c6c;
  }

  .execution-result.success {
    background-color: #f0f9eb;
    border-color: #e1f3d8;
    color: #67c23a;
  }

  .execution-result pre {
    white-space: pre-wrap;
    word-wrap: break-word;
    margin: 0;
    font-family: Menlo, Monaco, Consolas, 'Courier New', monospace;
  }

  /* 防止背景滚动 */
  :global(.modal-open) {
    overflow: hidden;
  }
</style>
