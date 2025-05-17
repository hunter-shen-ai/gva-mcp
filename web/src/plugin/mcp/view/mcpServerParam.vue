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
        <el-form-item label="函数名称" prop="funcName">
          <el-input v-model="searchInfo.funcName" placeholder="搜索条件" />
        </el-form-item>

        <el-form-item label="参数名称" prop="paramName">
          <el-input v-model="searchInfo.paramName" placeholder="搜索条件" />
        </el-form-item>

        <el-form-item label="参数描述" prop="paramDescription">
          <el-input
            v-model="searchInfo.paramDescription"
            placeholder="搜索条件"
          />
        </el-form-item>

        <el-form-item label="是否必需" prop="paramRequired">
          <el-select
            v-model="searchInfo.paramRequired"
            clearable
            placeholder="请选择"
          >
            <el-option key="true" label="是" value="true"></el-option>
            <el-option key="false" label="否" value="false"></el-option>
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
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column
          align="left"
          label="函数名称"
          prop="funcName"
          width="120"
        />

        <el-table-column
          align="left"
          label="参数名称"
          prop="paramName"
          width="120"
        />

        <el-table-column
          align="left"
          label="参数描述"
          prop="paramDescription"
          width="480"
          show-overflow-tooltip
        />

        <el-table-column
          align="left"
          label="是否必需"
          prop="paramRequired"
          width="120"
        >
          <template #default="scope">{{
            formatBoolean(scope.row.paramRequired)
          }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="创建时间"
          prop="createdAt"
          width="180"
        >
          <template #default="scope">{{
            formatDate(scope.row.createdAt)
          }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="更新时间"
          prop="updatedAt"
          width="180"
        >
          <template #default="scope">{{
            formatDate(scope.row.updatedAt)
          }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="操作"
          fixed="right"
          min-width="240"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              class="table-button"
              @click="getDetails(scope.row)"
              ><el-icon style="margin-right: 5px"><InfoFilled /></el-icon
              >查看</el-button
            >
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateMcpServerParamFunc(scope.row)"
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
          <span class="text-lg">{{ type === 'create' ? '新增' : '编辑' }}</span>
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
        <el-form-item label="函数名称:" prop="funcName">
          <el-input
            v-model="formData.funcName"
            :clearable="true"
            placeholder="请输入函数名称"
          />
        </el-form-item>
        <el-form-item label="参数名称:" prop="paramName">
          <el-input
            v-model="formData.paramName"
            :clearable="true"
            placeholder="请输入参数名称"
          />
        </el-form-item>
        <el-form-item label="参数描述:" prop="paramDescription">
          <el-input
            v-model="formData.paramDescription"
            type="textarea"
            :autosize="{ minRows: 3, maxRows: 10 }"
            :clearable="true"
            placeholder="请输入参数描述"
          />
        </el-form-item>
        <el-form-item label="是否必需:" prop="paramRequired">
          <el-switch
            v-model="formData.paramRequired"
            active-color="#13ce66"
            inactive-color="#ff4949"
            active-text="是"
            inactive-text="否"
            clearable
          ></el-switch>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer
      destroy-on-close
      size="800"
      v-model="detailShow"
      :show-close="true"
      :before-close="closeDetailShow"
      title="查看"
    >
      <el-descriptions :column="1" border label-width="120">
        <el-descriptions-item label="主键 ID">
          {{ detailFrom.id }}
        </el-descriptions-item>
        <el-descriptions-item label="函数名称">
          {{ detailFrom.funcName }}
        </el-descriptions-item>
        <el-descriptions-item label="参数名称">
          {{ detailFrom.paramName }}
        </el-descriptions-item>
        <el-descriptions-item label="参数描述">
          {{ detailFrom.paramDescription }}
        </el-descriptions-item>
        <el-descriptions-item label="是否必需">
          {{ formatBoolean(detailFrom.paramRequired) }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ detailFrom.createdAt }}
        </el-descriptions-item>
        <el-descriptions-item label="更新时间">
          {{ detailFrom.updatedAt }}
        </el-descriptions-item>
        <el-descriptions-item label="删除时间">
          {{ detailFrom.deletedAt }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createMcpServerParam,
    deleteMcpServerParam,
    deleteMcpServerParamByIds,
    updateMcpServerParam,
    findMcpServerParam,
    getMcpServerParamList
  } from '@/plugin/mcp/api/mcpServerParam'

  // 全量引入格式化工具 请按需保留
  import { formatDate, formatBoolean } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref, reactive } from 'vue'

  defineOptions({
    name: 'McpServerParam'
  })

  // 提交按钮loading
  const btnLoading = ref(false)

  // 控制更多查询条件显示/隐藏状态
  const showAllQuery = ref(false)

  // 自动化生成的字典（可能为空）以及字段
  const formData = ref({
    id: undefined,
    createdAt: new Date(),
    updatedAt: new Date(),
    deletedAt: new Date(),
    funcName: '',
    paramName: '',
    paramDescription: '',
    paramRequired: false
  })

  // 验证规则
  const rule = reactive({})

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
      if (searchInfo.value.paramRequired === '') {
        searchInfo.value.paramRequired = null
      }
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
    const table = await getMcpServerParamList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  getTableData()

  // ============== 表格控制部分结束 ===============

  // 获取需要的字典 可能为空 按需保留
  const setOptions = async () => {}

  // 获取需要的字典 可能为空 按需保留
  setOptions()

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
      deleteMcpServerParamFunc(row)
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
          ids.push(item.id)
        })
      const res = await deleteMcpServerParamByIds({ ids })
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
  const updateMcpServerParamFunc = async (row) => {
    const res = await findMcpServerParam({ id: row.id })
    type.value = 'update'
    if (res.code === 0) {
      formData.value = res.data
      dialogFormVisible.value = true
    }
  }

  // 删除行
  const deleteMcpServerParamFunc = async (row) => {
    const res = await deleteMcpServerParam({ id: row.id })
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
      id: undefined,
      createdAt: new Date(),
      updatedAt: new Date(),
      deletedAt: new Date(),
      funcName: '',
      paramName: '',
      paramDescription: '',
      paramRequired: false
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
          res = await createMcpServerParam(formData.value)
          break
        case 'update':
          res = await updateMcpServerParam(formData.value)
          break
        default:
          res = await createMcpServerParam(formData.value)
          break
      }
      btnLoading.value = false
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '创建/更改成功'
        })
        closeDialog()
        getTableData()
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
    const res = await findMcpServerParam({ id: row.id })
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
</script>

<style></style>
