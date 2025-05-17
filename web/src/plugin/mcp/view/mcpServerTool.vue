<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        :rules="searchRule"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="工具名称" prop="funcName">
          <el-input v-model="searchInfo.funcName" placeholder="搜索条件" />
        </el-form-item>

        <el-form-item label="工具描述" prop="funcDescription">
          <el-input
            v-model="searchInfo.funcDescription"
            placeholder="搜索条件"
          />
        </el-form-item>

        <el-form-item label="工具类型" prop="funcType">
          <el-select
            v-model="searchInfo.funcType"
            filterable
            clearable
            placeholder="请选择"
            @clear="
              () => {
                searchInfo.funcType = undefined
              }
            "
          >
            <el-option
              v-for="(item, key) in mcp_serv_typeOptions"
              :key="key"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit"
            >查询</el-button
          >
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button
            link
            type="primary"
            icon="arrow-down"
            @click="showAllQuery = true"
            v-if="!showAllQuery"
            >展开</el-button
          >
          <el-button
            link
            type="primary"
            icon="arrow-up"
            @click="showAllQuery = false"
            v-else
            >收起</el-button
          >
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()"
          >新增</el-button
        >
        <el-button
          icon="delete"
          style="margin-left: 10px"
          :disabled="!multipleSelection.length"
          @click="onDelete"
          >删除</el-button
        >
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="light"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        @expand-change="handleExpandChange"
      >
        <el-table-column type="expand">
          <template #default="props">
            <div style="padding: 10px 20px">
              <template v-if="props.row.funcType === 'http_api'">
                <el-tabs
                  v-model="activeParameterTab[props.row.ID]"
                  @tab-click="
                    (tab) => handleTabClick(props.row.ID, tab.props.name)
                  "
                  class="parameter-tabs"
                >
                  <el-tab-pane label="Param" name="Param">
                    <div class="tab-content-area">
                      <div
                        style="
                          margin-bottom: 10px;
                          display: flex;
                          justify-content: flex-start;
                          align-items: center;
                        "
                      >
                        <el-button
                          type="primary"
                          icon="plus"
                          @click="openParameterDialog(props.row, null, 'Param')"
                          >添加 Query 参数</el-button
                        >
                      </div>
                      <el-table
                        :data="filterParameters(props.row.parameters, 'Param')"
                        border
                        style="width: 100%"
                      >
                        <el-table-column label="参数名称" prop="paramName" />
                        <el-table-column
                          label="参数描述"
                          prop="paramDescription"
                          show-overflow-tooltip
                        />
                        <el-table-column
                          label="参数类型"
                          prop="paramDataType"
                        />
                        <el-table-column label="是否必选" prop="paramRequired">
                          <template #default="scope">
                            {{ scope.row.paramRequired ? '是' : '否' }}
                          </template>
                        </el-table-column>
                        <el-table-column
                          label="是否敏感"
                          prop="isSecure"
                          width="100"
                        >
                          <template #default="scope">
                            {{ scope.row.isSecure ? '是' : '否' }}
                          </template>
                        </el-table-column>
                        <el-table-column label="默认值" prop="defaultValue">
                          <template #default="scope">
                            {{
                              scope.row.isSecure
                                ? '******'
                                : scope.row.defaultValue
                            }}
                          </template>
                        </el-table-column>
                        <el-table-column label="操作" width="150" fixed="left">
                          <template #default="paramScope">
                            <el-button
                              type="primary"
                              link
                              icon="edit"
                              @click="
                                openParameterDialog(props.row, paramScope.row)
                              "
                              >编辑</el-button
                            >
                            <el-button
                              type="danger"
                              link
                              icon="delete"
                              @click="
                                deleteParameterFunc(
                                  props.row,
                                  paramScope.row.ID
                                )
                              "
                              >删除</el-button
                            >
                          </template>
                        </el-table-column>
                      </el-table>
                      <p
                        v-if="
                          filterParameters(props.row.parameters, 'Param')
                            .length === 0
                        "
                      >
                        Params下暂无参数。
                      </p>
                    </div>
                  </el-tab-pane>
                  <el-tab-pane label="Body" name="Body">
                    <div class="tab-content-area">
                      <div
                        style="
                          margin-bottom: 10px;
                          display: flex;
                          justify-content: flex-start;
                          align-items: center;
                        "
                      >
                        <el-button
                          type="primary"
                          icon="plus"
                          @click="openParameterDialog(props.row, null, 'Body')"
                          >添加 Body 参数</el-button
                        >
                      </div>
                      <el-table
                        :data="filterParameters(props.row.parameters, 'Body')"
                        border
                        style="width: 100%"
                      >
                        <el-table-column label="参数名称" prop="paramName" />
                        <el-table-column
                          label="参数描述"
                          prop="paramDescription"
                          show-overflow-tooltip
                        />
                        <el-table-column
                          label="参数类型"
                          prop="paramDataType"
                        />
                        <el-table-column label="是否必选" prop="paramRequired">
                          <template #default="scope">
                            {{ scope.row.paramRequired ? '是' : '否' }}
                          </template>
                        </el-table-column>
                        <el-table-column
                          label="是否敏感"
                          prop="isSecure"
                          width="100"
                        >
                          <template #default="scope">
                            {{ scope.row.isSecure ? '是' : '否' }}
                          </template>
                        </el-table-column>
                        <el-table-column label="默认值" prop="defaultValue">
                          <template #default="scope">
                            {{
                              scope.row.isSecure
                                ? '******'
                                : scope.row.defaultValue
                            }}
                          </template>
                        </el-table-column>
                        <el-table-column label="操作" width="150" fixed="left">
                          <template #default="paramScope">
                            <el-button
                              type="primary"
                              link
                              icon="edit"
                              @click="
                                openParameterDialog(props.row, paramScope.row)
                              "
                              >编辑</el-button
                            >
                            <el-button
                              type="danger"
                              link
                              icon="delete"
                              @click="
                                deleteParameterFunc(
                                  props.row,
                                  paramScope.row.ID
                                )
                              "
                              >删除</el-button
                            >
                          </template>
                        </el-table-column>
                      </el-table>
                      <p
                        v-if="
                          filterParameters(props.row.parameters, 'Body')
                            .length === 0
                        "
                      >
                        Body下暂无参数。
                      </p>
                    </div>
                  </el-tab-pane>
                  <el-tab-pane label="Header" name="Header">
                    <div class="tab-content-area">
                      <div
                        style="
                          margin-bottom: 10px;
                          display: flex;
                          justify-content: flex-start;
                          align-items: center;
                        "
                      >
                        <el-button
                          type="primary"
                          icon="plus"
                          @click="
                            openParameterDialog(props.row, null, 'Header')
                          "
                          >添加 Header 参数</el-button
                        >
                      </div>
                      <el-table
                        :data="filterParameters(props.row.parameters, 'Header')"
                        border
                        style="width: 100%"
                      >
                        <el-table-column label="参数名称" prop="paramName" />
                        <el-table-column
                          label="参数描述"
                          prop="paramDescription"
                          show-overflow-tooltip
                        />
                        <el-table-column
                          label="参数类型"
                          prop="paramDataType"
                        />
                        <el-table-column label="是否必选" prop="paramRequired">
                          <template #default="scope">
                            {{ scope.row.paramRequired ? '是' : '否' }}
                          </template>
                        </el-table-column>
                        <el-table-column
                          label="是否敏感"
                          prop="isSecure"
                          width="100"
                        >
                          <template #default="scope">
                            {{ scope.row.isSecure ? '是' : '否' }}
                          </template>
                        </el-table-column>
                        <el-table-column label="默认值" prop="defaultValue">
                          <template #default="scope">
                            {{
                              scope.row.isSecure
                                ? '******'
                                : scope.row.defaultValue
                            }}
                          </template>
                        </el-table-column>
                        <el-table-column label="操作" width="150" fixed="left">
                          <template #default="paramScope">
                            <el-button
                              type="primary"
                              link
                              icon="edit"
                              @click="
                                openParameterDialog(props.row, paramScope.row)
                              "
                              >编辑</el-button
                            >
                            <el-button
                              type="danger"
                              link
                              icon="delete"
                              @click="
                                deleteParameterFunc(
                                  props.row,
                                  paramScope.row.ID
                                )
                              "
                              >删除</el-button
                            >
                          </template>
                        </el-table-column>
                      </el-table>
                      <p
                        v-if="
                          filterParameters(props.row.parameters, 'Header')
                            .length === 0
                        "
                      >
                        Header下暂无参数。
                      </p>
                    </div>
                  </el-tab-pane>
                </el-tabs>
                <p
                  v-if="
                    (!props.row.parameters ||
                      props.row.parameters.length === 0) &&
                    (!activeParameterTab[props.row.ID] ||
                      filterParameters(
                        props.row.parameters,
                        activeParameterTab[props.row.ID]
                      ).length === 0)
                  "
                >
                  当前工具暂无任何参数配置。
                </p>
              </template>
              <template v-else>
                <div
                  style="
                    margin-bottom: 10px;
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                  "
                >
                  <el-button
                    type="primary"
                    icon="plus"
                    @click="openParameterDialog(props.row, null)"
                    >添加参数</el-button
                  >
                </div>
                <el-table
                  :data="props.row.parameters"
                  border
                  style="width: 100%"
                >
                  <el-table-column label="参数名称" prop="paramName" />
                  <el-table-column
                    label="参数描述"
                    prop="paramDescription"
                    show-overflow-tooltip
                  />
                  <el-table-column label="参数类型" prop="paramDataType" />
                  <!-- 请求类型列在非 funcType 'http_api' 时不应该显示，
                       或者说，系统函数的参数没有请求类型的概念 -->
                  <el-table-column label="是否必选" prop="paramRequired">
                    <template #default="scope">
                      {{ scope.row.paramRequired ? '是' : '否' }}
                    </template>
                  </el-table-column>
                  <el-table-column label="是否敏感" prop="isSecure" width="100">
                    <template #default="scope">
                      {{ scope.row.isSecure ? '是' : '否' }}
                    </template>
                  </el-table-column>
                  <el-table-column label="默认值" prop="defaultValue">
                    <template #default="scope">
                      {{
                        scope.row.isSecure ? '******' : scope.row.defaultValue
                      }}
                    </template>
                  </el-table-column>
                  <el-table-column label="操作" width="150" fixed="left">
                    <template #default="paramScope">
                      <el-button
                        type="primary"
                        link
                        icon="edit"
                        @click="openParameterDialog(props.row, paramScope.row)"
                        >编辑</el-button
                      >
                      <el-button
                        type="danger"
                        link
                        icon="delete"
                        @click="
                          deleteParameterFunc(props.row, paramScope.row.ID)
                        "
                        >删除</el-button
                      >
                    </template>
                  </el-table-column>
                </el-table>
                <p
                  v-if="
                    !props.row.parameters || props.row.parameters.length === 0
                  "
                >
                  暂无参数配置。
                </p>
              </template>
            </div>
          </template>
        </el-table-column>
        <el-table-column type="selection" width="55" />
        <el-table-column
          align="left"
          label="工具名称"
          prop="funcName"
          width="180"
        />
        <el-table-column
          align="left"
          label="工具描述"
          prop="funcDescription"
          width="250"
          show-overflow-tooltip
        />
        <el-table-column
          align="left"
          label="API 地址/来源"
          prop="apiUrl"
          width="250"
          show-overflow-tooltip
        >
          <template #default="scope">
            <span v-if="scope.row.funcType === 'system_func'">-</span>
            <span v-else>{{ scope.row.apiUrl }}</span>
          </template>
        </el-table-column>
        <el-table-column
          sortable
          align="left"
          label="工具类型"
          prop="funcType"
          width="120"
        >
          <template #default="scope">
            <el-tag
              type="success"
              v-if="scope.row.funcType === 'system_func'"
              >{{
                filterDict(scope.row.funcType, mcp_serv_typeOptions)
              }}</el-tag
            >
            <el-tag
              type="primary"
              v-else-if="scope.row.funcType === 'http_api'"
              >{{
                filterDict(scope.row.funcType, mcp_serv_typeOptions)
              }}</el-tag
            >
            <!-- TODO: Add more types or a generic display if not system_func or http_api -->
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="请求方法"
          prop="requestMethod"
          width="120"
        >
          <template #default="scope">
            <span v-if="scope.row.funcType === 'system_func'">-</span>
            <span v-else>{{ scope.row.requestMethod }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="创建时间"
          prop="CreatedAt"
          width="180"
        >
          <template #default="scope">
            {{ formatDate(new Date(scope.row.CreatedAt)) }}
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="更新时间"
          prop="UpdatedAt"
          width="180"
        >
          <template #default="scope">
            {{ formatDate(new Date(scope.row.UpdatedAt)) }}
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="操作"
          fixed="right"
          min-width="240"
        >
          <template #default="scope">
            <!-- <el-button
              type="primary"
              link
              class="table-button"
              @click="getDetails(scope.row)"
              ><el-icon style="margin-right: 5px"><InfoFilled /></el-icon
              >查看</el-button
            > -->
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateMcpServerToolFunc(scope.row)"
              >编辑</el-button
            >
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <el-drawer
      destroy-on-close
      size="800"
      v-model="dialogFormVisible"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{
            type === 'create' ? '新增工具' : '编辑工具'
          }}</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog"
              >确 定</el-button
            >
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form
        :model="formData"
        label-position="top"
        ref="elFormRef"
        :rules="rule"
        label-width="80px"
      >
        <el-form-item label="工具名称:" prop="funcName">
          <el-input
            v-model="formData.funcName"
            :clearable="true"
            placeholder="请输入工具名称"
          />
        </el-form-item>
        <el-form-item label="工具描述:" prop="funcDescription">
          <el-input
            v-model="formData.funcDescription"
            type="textarea"
            :autosize="{ minRows: 3, maxRows: 10 }"
            :clearable="true"
            placeholder="请输入工具描述"
          />
        </el-form-item>
        <el-form-item label="工具类型:" prop="funcType">
          <el-select
            v-model="formData.funcType"
            placeholder="请选择工具类型"
            style="width: 100%"
            filterable
            :clearable="true"
          >
            <el-option
              v-for="(item, key) in mcp_serv_typeOptions"
              :key="key"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          label="请求方法:"
          prop="requestMethod"
          v-if="formData.funcType !== 'system_func'"
        >
          <el-select
            v-model="formData.requestMethod"
            placeholder="请选择请求方法"
            style="width: 100%"
            :clearable="true"
          >
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="DELETE" value="DELETE" />
            <el-option label="PUT" value="PUT" />
          </el-select>
        </el-form-item>
        <el-form-item
          label="API 地址/来源:"
          prop="apiUrl"
          v-if="formData.funcType !== 'system_func'"
        >
          <el-input
            v-model="formData.apiUrl"
            :clearable="true"
            placeholder="请输入 API 地址或来源标识"
          />
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- <el-drawer
      destroy-on-close
      size="800"
      v-model="detailShow"
      :show-close="true"
      :before-close="closeDetailShow"
      title="查看工具详情"
    >
      <el-descriptions :column="1" border label-width="120">
        <el-descriptions-item label="ID">
          {{ detailFrom.ID }}
        </el-descriptions-item>
        <el-descriptions-item label="工具名称">
          {{ detailFrom.funcName }}
        </el-descriptions-item>
        <el-descriptions-item label="工具描述">
          {{ detailFrom.funcDescription }}
        </el-descriptions-item>
        <el-descriptions-item label="API 地址/来源">
          {{ detailFrom.apiUrl }}
        </el-descriptions-item>
        <el-descriptions-item label="请求方法">
          {{ detailFrom.requestMethod }}
        </el-descriptions-item>
        <el-descriptions-item label="工具类型">
          {{ filterDict(detailFrom.funcType, mcp_serv_typeOptions) }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ formatDate(new Date(detailFrom.CreatedAt)) }}
        </el-descriptions-item>
        <el-descriptions-item label="更新时间">
          {{ formatDate(new Date(detailFrom.UpdatedAt)) }}
        </el-descriptions-item>
        <el-descriptions-item label="删除时间" v-if="detailFrom.DeletedAt">
          {{ formatDate(new Date(detailFrom.DeletedAt)) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer> -->

    <!-- Drawer for Adding/Editing Parameters -->
    <el-drawer
      destroy-on-close
      size="600"
      v-model="parameterDialogVisible"
      :show-close="false"
      :before-close="closeParameterDialog"
      direction="rtl"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{
            parameterFormType === 'create' ? '添加参数' : '编辑参数'
          }}</span>
          <div>
            <el-button
              :loading="parameterBtnLoading"
              type="primary"
              @click="enterParameterDialog"
              >确 定</el-button
            >
            <el-button @click="closeParameterDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form
        :model="parameterFormData"
        label-position="top"
        ref="elParameterFormRef"
        :rules="parameterFormRules"
        label-width="100px"
      >
        <el-form-item label="参数名称:" prop="paramName">
          <el-input
            v-model="parameterFormData.paramName"
            :clearable="true"
            placeholder="请输入参数名称"
          />
        </el-form-item>
        <el-form-item label="参数描述:" prop="paramDescription">
          <el-input
            v-model="parameterFormData.paramDescription"
            type="textarea"
            :autosize="{ minRows: 2, maxRows: 5 }"
            :clearable="true"
            placeholder="请输入参数描述"
          />
        </el-form-item>
        <el-form-item
          v-if="
            currentSelectedTool && currentSelectedTool.funcType === 'http_api'
          "
          label="工具请求方法"
        >
          <el-input :model-value="currentSelectedTool.requestMethod" disabled />
        </el-form-item>
        <el-form-item label="数据类型:" prop="paramDataType">
          <el-select
            v-model="parameterFormData.paramDataType"
            placeholder="请选择数据类型"
            style="width: 100%"
            :clearable="true"
          >
            <el-option label="string" value="string" />
            <el-option label="int" value="int" />
            <el-option label="bool" value="bool" />
          </el-select>
        </el-form-item>
        <el-form-item
          label="请求类型:"
          prop="requestType"
          v-if="
            currentSelectedTool && currentSelectedTool.funcType === 'http_api'
          "
        >
          <el-select
            v-model="parameterFormData.requestType"
            placeholder="请选择请求类型"
            style="width: 100%"
            :clearable="true"
            :disabled="
              parameterFormType === 'create' &&
              currentSelectedTool &&
              currentSelectedTool.funcType === 'http_api'
            "
          >
            <el-option label="Body" value="Body" />
            <el-option label="Param" value="Param" />
            <el-option label="Header" value="Header" />
          </el-select>
        </el-form-item>
        <el-form-item label="是否必选:" prop="paramRequired">
          <el-switch v-model="parameterFormData.paramRequired" />
        </el-form-item>
        <el-form-item label="是否为敏感字段:" prop="isSecure">
          <el-switch v-model="parameterFormData.isSecure" />
        </el-form-item>
        <el-form-item label="默认值:" prop="defaultValue">
          <el-input
            v-model="parameterFormData.defaultValue"
            :type="parameterFormData.isSecure ? 'password' : 'text'"
            :show-password="parameterFormData.isSecure"
            :clearable="true"
            placeholder="请输入默认值 (可选)"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createMcpServerTool,
    deleteMcpServerTool,
    deleteMcpServerToolByIds,
    updateMcpServerTool,
    findMcpServerTool,
    getMcpServerToolList
  } from '@/plugin/mcp/api/mcpServerTool'

  import {
    createMcpServerParam,
    updateMcpServerParam,
    deleteMcpServerParam
  } from '@/plugin/mcp/api/mcpServerParam.js' // Ensure .js extension if that's your convention

  // 全量引入格式化工具 请按需保留
  import { formatDate, filterDict, getDictFunc } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref, reactive, watch } from 'vue'

  defineOptions({
    name: 'McpServerTool'
  })

  // 提交按钮loading
  const btnLoading = ref(false)

  // 控制更多查询条件显示/隐藏状态
  const showAllQuery = ref(false)
  const mcp_serv_typeOptions = ref([])
  // 自动化生成的字典（可能为空）以及字段
  const formData = ref({
    ID: undefined,
    funcName: '',
    funcDescription: '',
    funcType: '',
    apiUrl: '',
    requestMethod: '',
    parameters: []
  })

  // 验证规则
  const rule = reactive({
    funcName: [
      {
        required: true,
        message: '请输入工具名称',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: '工具名称不能为空格',
        trigger: ['input', 'blur']
      }
    ],
    funcDescription: [
      {
        required: true,
        message: '请输入工具描述',
        trigger: ['input', 'blur']
      },
      {
        whitespace: true,
        message: '工具描述不能为空格',
        trigger: ['input', 'blur']
      }
    ],
    funcType: [
      {
        required: true,
        message: '请选择工具类型',
        trigger: ['input', 'blur']
      }
    ],
    apiUrl: [
      {
        validator: (rule, value, callback) => {
          // Only require apiUrl if funcType is not 'system_func' (System Function)
          if (formData.value.funcType !== 'system_func' && !value) {
            callback(new Error('请输入API地址/来源'))
          } else if (
            formData.value.funcType !== 'system_func' &&
            value &&
            !value.trim()
          ) {
            callback(new Error('API地址/来源不能为空格'))
          } else {
            callback()
          }
        },
        trigger: ['input', 'blur']
      }
    ],
    requestMethod: [
      {
        validator: (rule, value, callback) => {
          // Only require requestMethod if funcType is not 'system_func' (System Function)
          if (formData.value.funcType !== 'system_func' && !value) {
            callback(new Error('请选择请求方法'))
          } else {
            callback()
          }
        },
        trigger: ['change', 'blur']
      }
    ]
  })

  const searchRule = reactive({
    CreatedAt: [
      {
        validator: (rule, value, callback) => {
          if (
            searchInfo.value.startCreatedAt &&
            !searchInfo.value.endCreatedAt
          ) {
            callback(new Error('请填写结束日期'))
          } else if (
            !searchInfo.value.startCreatedAt &&
            searchInfo.value.endCreatedAt
          ) {
            callback(new Error('请填写开始日期'))
          } else if (
            searchInfo.value.startCreatedAt &&
            searchInfo.value.endCreatedAt &&
            (searchInfo.value.startCreatedAt.getTime() ===
              searchInfo.value.endCreatedAt.getTime() ||
              searchInfo.value.startCreatedAt.getTime() >
                searchInfo.value.endCreatedAt.getTime())
          ) {
            callback(new Error('开始日期应当早于结束日期'))
          } else {
            callback()
          }
        },
        trigger: 'change'
      }
    ]
  })

  const elFormRef = ref()
  const elSearchFormRef = ref()

  // =========== 表格控制部分 ===========
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const searchInfo = ref({})
  // 重置
  const onReset = () => {
    searchInfo.value = {}
    getTableData()
  }

  // 搜索
  const onSubmit = () => {
    elSearchFormRef.value?.validate(async (valid) => {
      if (!valid) return
      page.value = 1
      getTableData()
    })
  }

  // 分页
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  // 修改页面容量
  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // 查询
  const getTableData = async () => {
    const table = await getMcpServerToolList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      table.data.list.forEach((item) => {
        if (item.parameters === undefined) {
          item.parameters = null
        }
        item._fullyLoadedParameters = false // Initialize loading state tracker
      })
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  getTableData()

  // ============== 表格控制部分结束 ===============

  // 获取需要的字典 可能为空 按需保留
  const setOptions = async () => {
    mcp_serv_typeOptions.value = await getDictFunc('mcp_serv_type')
  }

  // 获取需要的字典 可能为空 按需保留
  setOptions()

  // Watch for changes in funcType to clear requestMethod if it becomes 'System Function'
  watch(
    () => formData.value.funcType,
    (newType) => {
      if (newType === 'system_func') {
        // Assuming 'system_func' is the value for 'System Function'
        formData.value.requestMethod = ''
        formData.value.apiUrl = '' // Also clear apiUrl
      }
      // If form is available, re-validate requestMethod and apiUrl as their requirement might have changed
      if (elFormRef.value) {
        elFormRef.value.validateField('requestMethod')
        elFormRef.value.validateField('apiUrl')
      }
    }
  )

  // 多选数据
  const multipleSelection = ref([])
  // 多选
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }

  // 删除行
  const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      deleteMcpServerToolFunc(row)
    })
  }

  // 多选删除
  const onDelete = async () => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const ids = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map((item) => {
          ids.push(item.ID)
        })
      const res = await deleteMcpServerToolByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
  }

  // 行为控制标记（弹窗内部需要增还是改）
  const type = ref('')

  // 更新行
  const updateMcpServerToolFunc = async (row) => {
    const res = await findMcpServerTool({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
      formData.value = res.data
      dialogFormVisible.value = true
    }
  }

  // 删除行
  const deleteMcpServerToolFunc = async (row) => {
    const res = await deleteMcpServerTool({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  }

  // 弹窗控制标记
  const dialogFormVisible = ref(false)

  // 打开弹窗
  const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
  }

  // 关闭弹窗
  const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
      ID: undefined,
      funcName: '',
      funcDescription: '',
      funcType: '',
      apiUrl: '',
      requestMethod: '',
      parameters: []
    }
  }
  // 弹窗确定
  const enterDialog = async () => {
    btnLoading.value = true
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return (btnLoading.value = false)
      let res
      switch (type.value) {
        case 'create':
          res = await createMcpServerTool(formData.value)
          break
        case 'update':
          res = await updateMcpServerTool(formData.value)
          break
        default:
          res = await createMcpServerTool(formData.value)
          break
      }
      btnLoading.value = false
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '创建/更改成功'
        })

        // 获取工具ID
        const toolID = type.value === 'update' ? formData.value.ID : res.data.ID

        // 保存当前工具的扩展状态
        const expandedRows = Array.from(
          document.querySelectorAll('.el-table__expand-icon--expanded')
        )
          .map((el) => el.closest('tr'))
          .filter((tr) => tr)
          .map((tr) => {
            const rowID = tr.getAttribute('data-row-key')
            return rowID
          })

        closeDialog()
        await getTableData()

        // 如果是更新操作，确保更新后的工具保持参数的可见性
        if (type.value === 'update' && toolID) {
          // 查找更新后的工具所在行
          const updatedRowIndex = tableData.value.findIndex(
            (item) => item.ID === toolID
          )
          if (updatedRowIndex !== -1) {
            // 获取完整的工具数据包括参数
            const refreshRes = await findMcpServerTool({ ID: toolID })
            if (refreshRes.code === 0 && refreshRes.data) {
              const updatedToolData = { ...tableData.value[updatedRowIndex] }
              updatedToolData.parameters = [
                ...(refreshRes.data.parameters || [])
              ]
              updatedToolData._fullyLoadedParameters = true
              tableData.value[updatedRowIndex] = updatedToolData

              // 如果当前工具之前是展开的，保持展开状态
              if (expandedRows.includes(String(toolID))) {
                setTimeout(() => {
                  const expandBtns = document.querySelectorAll(
                    '.el-table__expand-icon'
                  )
                  if (expandBtns[updatedRowIndex]) {
                    expandBtns[updatedRowIndex].click()
                  }
                }, 100)
              }
            }
          }
        }
      }
    })
  }

  const detailFrom = ref({})

  // 查看详情控制标记
  const detailShow = ref(false)

  // 打开详情弹窗
  const openDetailShow = () => {
    detailShow.value = true
  }

  // 打开详情
  const getDetails = async (row) => {
    // 打开弹窗
    const res = await findMcpServerTool({ ID: row.ID })
    if (res.code === 0) {
      detailFrom.value = res.data
      openDetailShow()
    }
  }

  // 关闭详情弹窗
  const closeDetailShow = () => {
    detailShow.value = false
    detailFrom.value = {}
  }

  // Refs for parameter dialog
  const parameterDialogVisible = ref(false)
  const parameterFormType = ref('create') // 'create' or 'update'
  const parameterBtnLoading = ref(false)
  const elParameterFormRef = ref(null)
  const currentSelectedTool = ref(null) // To store the tool for which parameter is being added/edited

  const parameterFormData = ref({
    ID: undefined,
    toolID: undefined,
    paramName: '',
    paramDescription: '',
    paramDataType: 'string',
    paramRequired: false,
    defaultValue: '',
    requestType: '',
    isSecure: false
  })

  // Parameter item validation rules
  const parameterFormRules = reactive({
    paramName: [
      { required: true, message: '参数名称不能为空', trigger: 'blur' }
    ],
    paramDataType: [
      { required: true, message: '请选择参数数据类型', trigger: 'change' }
    ],
    requestType: [
      {
        validator: (rule, value, callback) => {
          // Skip validation if creating a parameter for a type 'http_api' tool, as the value is forced.
          if (
            parameterFormType.value === 'create' &&
            currentSelectedTool.value &&
            currentSelectedTool.value.funcType === 'http_api'
          ) {
            callback()
            return
          }
          // Original validation logic for other cases (e.g., editing)
          if (
            currentSelectedTool.value &&
            currentSelectedTool.value.funcType === 'http_api' &&
            !value
          ) {
            callback(new Error('请选择请求类型'))
          } else {
            callback()
          }
        },
        trigger: ['change', 'blur'] // Added blur trigger
      }
    ]
    // paramRequired is a switch
    // paramDescription and defaultValue are optional.
  })

  // Parameter Dialog Functions
  const openParameterDialog = (
    toolRow,
    parameter,
    forcedRequestType = null
  ) => {
    currentSelectedTool.value = toolRow
    if (parameter && parameter.ID) {
      // Editing existing parameter
      parameterFormType.value = 'update'
      const baseData = JSON.parse(JSON.stringify(parameter))
      parameterFormData.value = {
        ID: baseData.ID,
        toolID: baseData.toolID || toolRow.ID,
        paramName: baseData.paramName || '',
        paramDescription: baseData.paramDescription || '',
        paramDataType: baseData.paramDataType || 'string',
        paramRequired: baseData.paramRequired || false,
        defaultValue: baseData.defaultValue || '',
        requestType: baseData.requestType || '',
        isSecure: baseData.isSecure || false
      }
      if (toolRow.funcType === 'http_api' && forcedRequestType) {
        parameterFormData.value.requestType = forcedRequestType
      }
      if (!parameterFormData.value.toolID && toolRow.ID) {
        parameterFormData.value.toolID = toolRow.ID
      }
    } else {
      // Adding new parameter
      parameterFormType.value = 'create'
      parameterFormData.value = {
        ID: undefined,
        toolID: toolRow.ID,
        paramName: '',
        paramDescription: '',
        paramDataType: 'string',
        paramRequired: false,
        defaultValue: '',
        requestType: toolRow.funcType === 'http_api' ? forcedRequestType : '',
        isSecure: false
      }
    }
    console.log(
      'Opening Parameter Dialog:',
      'Type:',
      parameterFormType.value,
      'Tool Type:',
      currentSelectedTool.value?.funcType,
      'Forced Request Type:',
      forcedRequestType,
      'Parameter Data:',
      JSON.parse(JSON.stringify(parameterFormData.value))
    )
    parameterDialogVisible.value = true
    if (elParameterFormRef.value) {
      elParameterFormRef.value.clearValidate()
    }
  }

  const closeParameterDialog = () => {
    parameterDialogVisible.value = false
    if (elParameterFormRef.value) {
      elParameterFormRef.value.resetFields()
    }
    parameterFormData.value = {
      ID: undefined,
      toolID: undefined,
      paramName: '',
      paramDescription: '',
      paramDataType: 'string',
      paramRequired: false,
      defaultValue: '',
      requestType: '',
      isSecure: false
    }
    currentSelectedTool.value = null
  }

  const enterParameterDialog = async () => {
    if (!elParameterFormRef.value) return

    elParameterFormRef.value.validate(async (valid) => {
      if (valid) {
        parameterBtnLoading.value = true
        let res
        let operationSuccessful = false

        // Ensure toolID is correctly assigned
        if (
          !parameterFormData.value.toolID &&
          currentSelectedTool.value &&
          currentSelectedTool.value.ID
        ) {
          parameterFormData.value.toolID = currentSelectedTool.value.ID
        } else if (!parameterFormData.value.toolID) {
          console.error('ToolID is missing for parameter operation.')
          ElMessage({ type: 'error', message: '工具ID缺失，操作失败' })
          parameterBtnLoading.value = false
          return
        }

        try {
          if (parameterFormType.value === 'create') {
            const payload = { ...parameterFormData.value }
            delete payload.ID // ID should not be sent for creation
            res = await createMcpServerParam(payload)
          } else {
            res = await updateMcpServerParam(parameterFormData.value)
          }

          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message:
                parameterFormType.value === 'create'
                  ? '参数添加成功'
                  : '参数更新成功'
            })
            operationSuccessful = true

            const toolIdToRefresh = currentSelectedTool.value
              ? currentSelectedTool.value.ID
              : null // Capture ID before closing dialog

            closeParameterDialog() // Close dialog first

            // Refresh parameters for the affected tool
            if (toolIdToRefresh) {
              // Use captured ID
              const refreshedToolRes = await findMcpServerTool({
                ID: toolIdToRefresh
              })
              if (refreshedToolRes.code === 0 && refreshedToolRes.data) {
                const rowIndex = tableData.value.findIndex(
                  (item) => item.ID === toolIdToRefresh
                )
                if (rowIndex !== -1) {
                  // Ensure Vue detects the change by creating a new object reference for the row
                  const updatedToolData = { ...tableData.value[rowIndex] }
                  updatedToolData.parameters = [
                    ...(refreshedToolRes.data.parameters || [])
                  ]
                  updatedToolData._fullyLoadedParameters = true // Mark as reloaded
                  tableData.value[rowIndex] = updatedToolData
                }
              }
            } else {
              getTableData() // Fallback to full table refresh if context is lost
            }
          } else {
            ElMessage({
              type: 'error',
              message:
                res.msg ||
                (parameterFormType.value === 'create'
                  ? '参数添加失败'
                  : '参数更新失败')
            })
          }
        } catch (error) {
          console.error('Error in enterParameterDialog:', error)
          ElMessage({ type: 'error', message: '操作失败，请查看控制台' })
        } finally {
          parameterBtnLoading.value = false
        }
      }
    })
  }

  // Re-add handleExpandChange for lazy loading parameters
  const handleExpandChange = async (row, expandedRows) => {
    const isExpanded = expandedRows.some((r) => r.ID === row.ID)

    // Initialize active tab for funcType 'http_api' when expanded
    if (isExpanded && row.funcType === 'http_api') {
      initializeActiveTab(row.ID)
    }

    // Load parameters if:
    // 1. The row is being expanded.
    // 2. Parameters haven't been loaded yet (_fullyLoadedParameters is false).
    // 3. Or parameters are explicitly null (initial state before any load attempt).
    if (
      isExpanded &&
      (!row._fullyLoadedParameters || row.parameters === null)
    ) {
      try {
        // Optional: Indicate loading state for the specific row if UI needs it
        // row._isLoadingParameters = true;
        const res = await findMcpServerTool({ ID: row.ID })
        // row._isLoadingParameters = false;

        const rowIndex = tableData.value.findIndex((item) => item.ID === row.ID)
        if (rowIndex !== -1) {
          // Ensure row is still in tableData
          if (res.code === 0 && res.data) {
            tableData.value[rowIndex].parameters = res.data.parameters || []
          } else {
            // Handle API error or no parameters found case by setting to empty array
            tableData.value[rowIndex].parameters = []
            console.error(
              'Failed to load parameters or no parameters found for tool ID',
              row.ID,
              ':',
              res.msg
            )
          }
          tableData.value[rowIndex]._fullyLoadedParameters = true // Mark as loaded (or attempted to load)
        }
      } catch (error) {
        // row._isLoadingParameters = false;
        console.error(
          'Error fetching parameters for tool ID',
          row.ID,
          ':',
          error
        )
        const rowIndex = tableData.value.findIndex((item) => item.ID === row.ID)
        if (rowIndex !== -1) {
          tableData.value[rowIndex].parameters = [] // Set to empty on error
          tableData.value[rowIndex]._fullyLoadedParameters = true // Mark as attempted to load
        }
      }
    }
  }

  // New function to delete parameter
  const deleteParameterFunc = async (toolRow, parameterID) => {
    ElMessageBox.confirm('确定要删除此参数吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        const res = await deleteMcpServerParam({ ID: parameterID })
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '参数删除成功' })
          // Refresh parameters for the current tool
          if (toolRow && toolRow.ID) {
            const refreshedToolRes = await findMcpServerTool({ ID: toolRow.ID })
            if (refreshedToolRes.code === 0 && refreshedToolRes.data) {
              const rowIndex = tableData.value.findIndex(
                (item) => item.ID === toolRow.ID
              )
              if (rowIndex !== -1) {
                // Ensure Vue detects the change by creating a new object reference for the row
                const updatedToolData = { ...tableData.value[rowIndex] }
                updatedToolData.parameters = [
                  ...(refreshedToolRes.data.parameters || [])
                ]
                updatedToolData._fullyLoadedParameters = true // Parameters are reloaded
                tableData.value[rowIndex] = updatedToolData
              }
            }
          }
        } else {
          ElMessage({ type: 'error', message: res.msg || '参数删除失败' })
        }
      } catch (error) {
        console.error('Error deleting parameter:', error)
        ElMessage({ type: 'error', message: '参数删除失败' })
      }
    })
  }

  // Additional reactive state and functions for tabbed parameters
  const activeParameterTab = ref({}) // Key: tool.ID, Value: 'Param', 'Body', 'Header'

  const initializeActiveTab = (toolId) => {
    if (!activeParameterTab.value[toolId]) {
      activeParameterTab.value[toolId] = 'Param' // Default to 'Param'
    }
  }

  const handleTabClick = (toolId, tabName) => {
    activeParameterTab.value[toolId] = tabName
    // Potentially trigger a re-render or data refresh if needed, though v-model should handle it
  }

  const filterParameters = (parameters, requestType) => {
    if (!Array.isArray(parameters)) return []
    return parameters.filter((p) => p.requestType === requestType)
  }
</script>

<style scoped>
  .tab-content-area {
    margin-top: 20px;
  }
</style>
