<script lang="ts" setup>
import type { FormInstance, FormRules } from 'element-plus'
import type { TableColumnCtx } from 'element-plus/es/components/table/src/table-column/defaults'
import type { Account } from '@/api/account'
import { DocumentCopy } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import accountService from '@/api/account'
import copyText from '@/utils/copyText'

const { t } = useI18n()

function timeFormatter(row: any, col: TableColumnCtx<any>, value: any) {
  if (!value) {
    return ''
  }
  return dayjs(value).format('YYYY-MM-DD HH:mm:ss') || ''
}

const searchText = ref('')
const tableData = ref<Account[]>([])

const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const isLoading = ref(false)

// 对话框相关状态
const dialogVisible = ref(false)
const isEdit = ref(false)
const isSubmitting = ref(false)
const formRef = ref<FormInstance>()
const passwordDialogVisible = ref(false)
const currentPassword = ref('')
const form = reactive({
  id: '',
  username: '',
  password: '',
  clientID: '',
  remote: '',
  description: '',
})

// 表单验证规则
const rules = reactive<FormRules>({
  username: [
    {
      required: true,
      message: t('common.validation.required'),
      trigger: 'blur',
    },
  ],
  password: [
    {
      required: true,
      message: t('common.validation.required'),
      trigger: 'blur',
    },
  ],
  // clientID: [
  //   {
  //     required: true,
  //     message: '请输入客户端ID',
  //     trigger: 'blur',
  //   },
  // ],
})

function reloadData() {
  currentPage.value = 1
  loadData()
}

async function loadData() {
  try {
    isLoading.value = true
    const res = await accountService.getAccounts({
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
  Object.assign(form, {
    id: '',
    username: '',
    password: '',
    clientID: '',
    remote: '',
    description: '',
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
  dialogVisible.value = true
}

function showEditDialog(row: Account) {
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    username: row.username,
    password: row.password,
    clientID: row.clientID,
    remote: row.remote,
    description: row.description || '',
  })
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
        await accountService.updateAccount(form.id, {
          username: form.username,
          password: form.password,
          clientID: form.clientID,
          remote: form.remote,
          description: form.description,
        })
      } else {
        await accountService.createAccount({
          username: form.username,
          password: form.password,
          clientID: form.clientID,
          remote: form.remote,
          description: form.description,
        })
      }
      dialogVisible.value = false
      ElMessage.success(t(isEdit.value ? 'global.messages.edit_success' : 'global.messages.create_success'))
      loadData()
    } catch (error: any) {
      if (error.code && error.code.startsWith('err_')) {
        ElMessage.error(t(`error.${error.code}`))
      } else if (error.msg) {
        ElMessage.error(error.msg)
      } else {
        error && console.warn(error)
      }
    } finally {
      isSubmitting.value = false
    }
  })
}

function deleteAccount(row: Account) {
  accountService.deleteAccount(row.id).then(() => {
    ElMessage.success(t('global.messages.delete_success'))
    loadData()
  })
}

function toggleDisabled(row: Account) {
  ElMessageBox.confirm(
    row.disabled ? t('msg.enable') : t('msg.disable'),
    t('common.tip'),
    {
      confirmButtonText: t('common.ok'),
      cancelButtonText: t('common.cancel'),
      type: 'warning',
    },
  )
    .then(async () => {
      try {
        if (row.disabled) {
          accountService.enableAccount(row.id)
        } else {
          accountService.disableAccount(row.id)
        }
        ElMessage.success(t('common.success.modify'))
        loadData()
      } catch (error) {
        error && console.error(error)
      }
    })
    .catch(() => {
      // Error handled silently
    })
}

function showPassword(password: string) {
  currentPassword.value = password
  passwordDialogVisible.value = true
}

function handleCopyPassword() {
  copyText(currentPassword.value)
  ElMessage.success(t('msg.copy_success'))
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
        {{ t('common.create') }}
      </el-button>
    </template>
    <el-table v-loading="isLoading" :data="tableData" border height="100%">
      <el-table-column :label="t('username')" prop="username" width="150" />
      <el-table-column :label="t('password')" width="150">
        <template #default="scope">
          <div flex items-center gap-2>
            <div>••••••••</div>
            <el-button
              link
              type="info"
              @click="showPassword(scope.row.password)"
            >
              <el-icon>
                <i-ep-view />
              </el-icon>
            </el-button>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="t('clientID')" prop="clientID" width="150" />
      <el-table-column :label="t('clientIP')" prop="remote" width="150" />
      <el-table-column
        :label="t('common.disabled')"
        prop="disabled"
        width="100"
      >
        <template #default="scope">
          <el-tooltip
            :content="
              scope.row.disabled
                ? t('tip.click_to_enable')
                : t('tip.click_to_disable')
            "
            placement="left"
          >
            <el-button
              :type="scope.row.locked ? 'warning' : 'success'"
              size="small"
              plain
              round
              @click="toggleDisabled(scope.row)"
            >
              {{
                scope.row.disabled ? t('common.disabled') : t('common.enabled')
              }}
            </el-button>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column
        :label="t('common.description')"
        prop="description"
        show-overflow-tooltip
      />
      <el-table-column
        :label="t('common.createdAt')"
        prop="createdAt"
        :formatter="timeFormatter"
        width="220"
      />
      <el-table-column
        :label="t('common.updatedAt')"
        prop="updatedAt"
        :formatter="timeFormatter"
        width="220"
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
            @confirm="deleteAccount(scoped.row)"
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
    :title="isEdit ? t('common.edit') : t('common.create')"
    width="480px"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      class="dialog-form"
      label-width="auto"
      label-position="top"
    >
      <el-form-item :label="t('username')" prop="username">
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
      <el-form-item :label="t('password')" prop="password">
        <el-input
          v-model="form.password"
          type="password"
          autocomplete="off"
          show-password
          :placeholder="
            t('global.placeholders.please_input', {
              text: t('password'),
            })
          "
        />
      </el-form-item>
      <el-form-item :label="t('clientID')" prop="clientID">
        <el-input v-model="form.clientID" autocomplete="off" />
      </el-form-item>
      <el-form-item :label="t('clientIP')" prop="remote">
        <el-input v-model="form.remote" autocomplete="off" />
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
  <!-- 密码查看对话框 -->
  <el-dialog
    v-model="passwordDialogVisible"
    :title="t('view_password')"
    width="480px"
    :close-on-click-modal="false"
  >
    <div class="password-display-container">
      <div class="password-content">
        <div class="password-value-wrapper">
          <div class="password-value">
            {{ currentPassword }}
          </div>
          <el-button
            :icon="DocumentCopy"
            circle
            type="primary"
            plain
            class="copy-button"
            @click="handleCopyPassword"
          />
        </div>
      </div>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button type="primary" @click="passwordDialogVisible = false">
          {{ t('common.close') }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<i18n lang="yaml">
zh:
  username: 用户名
  password: 密码
  clientID: 客户端 ID
  clientIP: 客户端 IP
  description: 描述
  view_password: 查看密码
  tip:
    click_to_enable: '点击启用'
    click_to_disable: '点击禁用'
  msg:
    enable: '您确定要启用该账号吗？'
    disable: '您确定要禁用该账号吗？'
    copy_success: '复制成功'
en:
  username: Username
  password: Password
  clientID: Client ID
  clientIP: Client IP
  description: Description
  view_password: View password
  tip:
    click_to_enable: 'Click to enable'
    click_to_disable: 'Click to disable'
  msg:
    enable: 'Are you sure you want to enable this account?'
    disable: 'Are you sure you want to disable this account?'
    copy_success: 'Copy success'
</i18n>

<style lang="scss" scoped>
.password-display-container {
  .password-content {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .password-label {
    font-size: 14px;
    color: var(--el-text-color-regular);
    font-weight: 500;
  }

  .password-value-wrapper {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: var(--el-bg-color-page);
    border-radius: 8px;
    border: 1px solid var(--el-border-color-lighter);
    transition: all 0.3s ease;

    &:hover {
      border-color: var(--el-color-primary-light-7);
      background: var(--el-color-primary-light-9);
    }
  }

  .password-value {
    flex: 1;
    font-size: 16px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
    color: var(--el-color-primary);
    font-weight: 600;
    letter-spacing: 1px;
    word-break: break-all;
    user-select: all;
    padding: 4px 0;
  }
}
.copy-button {
  flex-shrink: 0;
  transition: all 0.3s ease;

  &:hover {
    transform: scale(1.1);
  }
}
</style>

<i18n lang="yaml">
zh:
  error:
    err_username_duplicated: '用户名已存在'
en:
  error:
    err_username_duplicated: 'Username duplicated'
</i18n>
