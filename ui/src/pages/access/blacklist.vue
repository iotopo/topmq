<script lang="ts" setup>
import type { FormInstance, FormRules } from 'element-plus'
import type { TableColumnCtx } from 'element-plus/es/components/table/src/table-column/defaults'
import type { Blacklist } from '@/api/blacklist'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import blacklistService from '@/api/blacklist'

const { t } = useI18n()

function timeFormatter(row: any, col: TableColumnCtx<any>, value: any) {
  if (!value) {
    return ''
  }
  return dayjs(value).format('YYYY-MM-DD HH:mm:ss') || ''
}

const searchText = ref('')
const tableData = ref<Blacklist[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const isLoading = ref(false)

// 对话框相关状态
const dialogVisible = ref(false)
const isEdit = ref(false)
const isSubmitting = ref(false)
const blacklistType = ref<'username' | 'clientID' | 'remote'>('username')
const formRef = ref<FormInstance>()
const form = reactive({
  id: '',
  username: '',
  clientID: '',
  remote: '',
  expiredAt: '',
  description: '',
})

function onBlacklistTypeChange(value: 'username' | 'clientID' | 'remote') {
  if (value === 'username') {
    form.clientID = ''
    form.remote = ''
  } else if (value === 'clientID') {
    form.username = ''
    form.remote = ''
  } else if (value === 'remote') {
    form.username = ''
    form.clientID = ''
  }
}

function reloadData() {
  currentPage.value = 1
  loadData()
}
async function loadData() {
  isLoading.value = true
  try {
    const res = await blacklistService.getAll({
      pageNum: currentPage.value,
      pageSize: pageSize.value,
      search: searchText.value,
    })
    tableData.value = res?.items ?? []
    total.value = res?.total ?? 0
  } catch (error) {
    console.error(error)
  } finally {
    isLoading.value = false
  }
}

function showCreateDialog() {
  isEdit.value = false
  blacklistType.value = 'username'
  Object.assign(form, {
    id: '',
    username: '',
    clientID: '',
    remote: '',
    expiredAt: '',
    description: '',
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
  dialogVisible.value = true
}

function showEditDialog(row: Blacklist) {
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    username: row.username || '',
    clientID: row.clientID || '',
    remote: row.remote || '',
    expiredAt: row.expiredAt || '',
    description: row.description || '',
  })

  if (row.username) {
    blacklistType.value = 'username'
    row.clientID = ''
    row.remote = ''
  } else if (row.clientID) {
    blacklistType.value = 'clientID'
    row.remote = ''
    row.username = ''
  } else if (row.remote) {
    blacklistType.value = 'remote'
    row.username = ''
    row.clientID = ''
  }

  if (formRef.value) {
    formRef.value.clearValidate()
  }
  dialogVisible.value = true
}

async function submit() {
  if (!formRef.value) {
    return
  }
  await formRef.value.validate(async (valid) => {
    if (!valid) {
      return
    }
    isSubmitting.value = true
    try {
      if (isEdit.value) {
        await blacklistService.update(form.id, {
          username: form.username,
          clientID: form.clientID,
          remote: form.remote,
          expiredAt: form.expiredAt,
          description: form.description,
        })
        ElMessage.success(t('global.messages.edit_success'))
      } else {
        await blacklistService.create({
          username: form.username,
          clientID: form.clientID,
          remote: form.remote,
          expiredAt: form.expiredAt,
          description: form.description,
        })
        ElMessage.success(t('global.messages.create_success'))
      }
      dialogVisible.value = false
      loadData()
    } catch (error: any) {
      if (error?.code && error.code.startsWith('err_')) {
        ElMessage.error(
          t(`error.${error.code}`, {
            msg: error.msg,
            ...error.errArgs,
          }),
        )
      } else {
        error && console.warn(error)
      }
    } finally {
      isSubmitting.value = false
    }
  })
}

function deleteBlacklist(row: Blacklist) {
  blacklistService
    .delete(row.id)
    .then(() => {
      ElMessage.success(t('global.messages.delete_success'))
      loadData()
    })
    .catch((error: any) => {
      if (error?.code && error.code.startsWith('err_')) {
        ElMessage.error(
          t(`error.${error.code}`, {
            msg: error.msg,
            ...error.errArgs,
          }),
        )
      } else {
        error && console.warn(error)
      }
    })
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <page>
    <template #header>
      <search>
        <search-input
          v-model="searchText"
          :placeholder="t('common.search')"
          @enter="reloadData"
        />
        <search-button @click="reloadData" />
      </search>
    </template>
    <template #headerRight>
      <el-button type="primary" @click="showCreateDialog()">
        {{ t('common.add') }}
      </el-button>
    </template>
    <el-table v-loading="isLoading" :data="tableData" border height="100%">
      <el-table-column :label="t('blacklist_type')" width="150">
        <template #default="scoped">
          <span v-if="scoped.row.username">{{ t('username') }}</span>
          <span v-else-if="scoped.row.clientID">{{ t('clientID') }}</span>
          <span v-else-if="scoped.row.remote">{{ t('clientIP') }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('username')" prop="username" width="150" />
      <el-table-column :label="t('clientID')" prop="clientID" width="150" />
      <el-table-column :label="t('clientIP')" prop="remote" width="150" />
      <el-table-column
        :label="t('common.description')"
        prop="description"
        show-overflow-tooltip
      />
      <el-table-column
        :label="t('expired_at')"
        prop="expiredAt"
        :formatter="timeFormatter"
        width="160"
      />
      <el-table-column
        :label="t('common.createdAt')"
        prop="createdAt"
        :formatter="timeFormatter"
        width="160"
      />
      <el-table-column
        :label="t('common.updatedAt')"
        prop="updatedAt"
        :formatter="timeFormatter"
        width="160"
      />
      <!-- 操作 -->
      <el-table-column
        :label="t('global.operation')"
        width="190"
        fixed="right"
        align="right"
      >
        <template #default="scoped">
          <el-button link type="primary" @click="showEditDialog(scoped.row)">
            {{ t('common.edit') }}
          </el-button>
          <el-popconfirm
            :title="t('common.confirm_delete')"
            @confirm="deleteBlacklist(scoped.row)"
          >
            <template #reference>
              <el-button type="danger" link class="important-pr-0">
                {{ t('common.delete') }}
              </el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    <template #footerRight>
      <pagination
        v-model:page-num="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        @page-change="loadData()"
      />
    </template>
  </page>
  <!-- 创建/编辑对话框 -->
  <el-dialog
    v-model="dialogVisible"
    :close-on-click-modal="false"
    :title="isEdit ? t('common.edit') : t('common.add')"
    width="480px"
  >
    <el-form
      ref="formRef"
      :model="form"
      class="dialog-form"
      label-width="auto"
      label-position="top"
    >
      <el-form-item :label="t('blacklist_type')" required>
        <el-select v-model="blacklistType" @change="onBlacklistTypeChange">
          <el-option :label="t('username')" value="username" />
          <el-option :label="t('clientID')" value="clientID" />
          <el-option :label="t('clientIP')" value="remote" />
        </el-select>
      </el-form-item>
      <el-form-item
        v-if="blacklistType === 'username'"
        :label="t('username')"
        prop="username"
        :rules="[{ required: true, message: t('common.validation.required') }]"
      >
        <el-input
          v-model="form.username"
          autocomplete="off"
          :placeholder="
            t('global.placeholders.please_input', {
              text: t('username'),
            })
          "
        />
      </el-form-item>
      <el-form-item
        v-if="blacklistType === 'clientID'"
        :label="t('clientID')"
        prop="clientID"
        :rules="[{ required: true, message: t('common.validation.required') }]"
      >
        <el-input
          v-model="form.clientID"
          autocomplete="off"
          :placeholder="
            t('global.placeholders.please_input', {
              text: t('clientID'),
            })
          "
        />
      </el-form-item>
      <el-form-item
        v-if="blacklistType === 'remote'"
        :label="t('clientIP')"
        prop="remote"
        :rules="[{ required: true, message: t('common.validation.required') }]"
      >
        <el-input
          v-model="form.remote"
          autocomplete="off"
          :placeholder="
            t('global.placeholders.please_input', {
              text: t('clientIP'),
            })
          "
        />
      </el-form-item>
      <el-form-item
        :label="t('expired_at')"
        prop="expiredAt"
        :rules="[{ required: true, message: t('common.validation.required') }]"
      >
        <el-date-picker
          v-model="form.expiredAt"
          type="datetime"
          style="width: 100%"
        />
      </el-form-item>
      <el-form-item :label="t('common.description')" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button :loading="isSubmitting" type="primary" @click="submit">
          {{ t('common.ok') }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<i18n lang="yaml">
zh:
  username: 用户名
  clientID: 客户端 ID
  clientIP: 客户端 IP
  blacklist_type: 黑名单类型
  expired_at: 过期时间
en:
  username: Username
  clientID: Client ID
  clientIP: Client IP
  blacklist_type: Blacklist Type
  expired_at: Expired At
</i18n>
