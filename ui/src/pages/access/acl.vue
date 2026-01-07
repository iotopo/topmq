<script lang="ts" setup>
import type { FormInstance, FormRules } from 'element-plus'
import type { TableColumnCtx } from 'element-plus/es/components/table/src/table-column/defaults'
import type { AccessControl } from '@/api/acl'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import accessControlService from '@/api/acl'

const { t } = useI18n()

function timeFormatter(row: any, col: TableColumnCtx<any>, value: any) {
  if (!value) {
    return ''
  }
  return dayjs(value).format('YYYY-MM-DD HH:mm:ss') || ''
}

const accessType = ref<'username' | 'clientID' | 'clientIP' | ''>('')
const searchText = ref('')
const tableData = ref<AccessControl[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const isLoading = ref(false)

// 对话框相关状态
const dialogVisible = ref(false)
const isEdit = ref(false)
const isSubmitting = ref(false)
const formRef = ref<FormInstance>()
const form = reactive({
  id: '',
  username: '',
  clientID: '',
  remote: '',
  topic: '',
  access: 'rw' as 'd' | 'r' | 'w' | 'rw',
  description: '',
})

function reloadData() {
  currentPage.value = 1
  loadData()
}
async function loadData() {
  isLoading.value = true
  try {
    const res = await accessControlService.getAccessControlList({
      pageNum: currentPage.value,
      pageSize: pageSize.value,
      search: searchText.value,
      type: accessType.value,
    })
    tableData.value = res?.items ?? []
    total.value = res?.total ?? 0
  }
  catch (error) {
    console.error(error)
  }
  finally {
    isLoading.value = false
  }
}

function showCreateDialog() {
  isEdit.value = false
  Object.assign(form, {
    id: '',
    username: '',
    clientID: '',
    remote: '',
    topic: '',
    access: 'rw' as 'd' | 'r' | 'w' | 'rw',
    description: '',
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
  dialogVisible.value = true
}

function showEditDialog(row: AccessControl) {
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    username: row.username || '',
    clientID: row.clientID || '',
    remote: row.remote || '',
    topic: row.topic,
    access: row.access,
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
        await accessControlService.update(form.id, {
          username: form.username,
          clientID: form.clientID,
          remote: form.remote,
          topic: form.topic,
          access: form.access,
          description: form.description,
        })
        ElMessage.success(t('global.messages.edit_success'))
      }
      else {
        await accessControlService.create({
          username: form.username,
          clientID: form.clientID,
          remote: form.remote,
          topic: form.topic,
          access: form.access,
          description: form.description,
        })
        ElMessage.success(t('global.messages.create_success'))
      }
      dialogVisible.value = false
      loadData()
    }
    catch (error: any) {
      if (error?.code && error.code.startsWith('err_')) {
        ElMessage.error(
          t(`error.${error.code}`, {
            msg: error.msg,
            ...error.errArgs,
          }),
        )
      }
      else {
        error && console.warn(error)
      }
    }
    finally {
      isSubmitting.value = false
    }
  })
}

function deleteAccessControl(row: AccessControl) {
  accessControlService
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
      }
      else {
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
        <el-radio-group v-model="accessType" @change="reloadData">
          <el-radio-button :label="t('username')" value="username" />
          <el-radio-button :label="t('clientID')" value="clientID" />
          <el-radio-button :label="t('clientIP')" value="clientIP" />
          <el-radio-button :label="t('all_user')" value="" />
        </el-radio-group>
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
      <el-table-column
        v-if="accessType === 'username'"
        :label="t('username')"
        prop="username"
        width="150"
      />
      <el-table-column
        v-if="accessType === 'clientID'"
        :label="t('clientID')"
        prop="clientID"
        width="150"
      />
      <el-table-column
        v-if="accessType === 'clientIP'"
        :label="t('clientIP')"
        prop="remote"
        width="150"
      />
      <el-table-column :label="t('topic')" prop="topic" width="150" />
      <el-table-column :label="t('access')" prop="access" width="150" />
      <el-table-column
        :label="t('common.description')"
        prop="description"
        show-overflow-tooltip
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
            @confirm="deleteAccessControl(scoped.row)"
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
      class="dialog-form"
      label-width="auto"
      label-position="top"
    >
      <el-form-item
        v-if="accessType === 'username'"
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
        v-if="accessType === 'clientID'"
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
        v-if="accessType === 'clientIP'"
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
        :label="t('topic')"
        prop="topic"
        class="custom-form-item"
        :rules="[{ required: true, message: t('common.validation.required') }]"
      >
        <template #label>
          <div flex items-center justify-between>
            <div m-r-10px>
              {{ t('topic') }}
            </div>
            <el-tooltip effect="dark" placement="top">
              <template #content>
                <div style="width: 400px">
                  <p>
                    您可以在设置主题时使用占位符，在匹配规则时将当前客户端信息等动态替换到主题中，支持的占位符如下：
                  </p>
                  <div>• {{ t('username') }}: ${username}</div>
                  <div>• {{ t('clientID') }}: ${clientid}</div>
                  <p>
                    如果您想要限制所有用户只允许订阅或者发布特定主题，可以类似这样填写：
                  </p>
                  <div>• 主题: xx/${username}/xxx</div>
                  <div>• 主题: xx/${clientid}/xxx</div>
                </div>
              </template>
              <el-icon class="cursor-pointer">
                <i-ep-question-filled />
              </el-icon>
            </el-tooltip>
          </div>
        </template>
        <el-input
          v-model="form.topic"
          autocomplete="off"
          :placeholder="
            t('global.placeholders.please_input', {
              text: t('topic'),
            })
          "
        />
      </el-form-item>
      <el-form-item :label="t('access')" prop="access" required>
        <el-radio-group v-model="form.access">
          <el-radio value="d">
            {{ t('access_options.d') }}
          </el-radio>
          <el-radio value="r">
            {{ t('access_options.r') }}
          </el-radio>
          <el-radio value="w">
            {{ t('access_options.w') }}
          </el-radio>
          <el-radio value="rw">
            {{ t('access_options.rw') }}
          </el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item :label="t('common.description')" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          :placeholder="
            t('global.placeholders.please_input', {
              text: t('common.description'),
            })
          "
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
  all_user: 所有用户
  topic: 主题
  access: 访问权限
  access_options:
    d: 禁用
    r: 发布
    w: 订阅
    rw: 发布 & 订阅
en:
  username: Username
  clientID: Client ID
  clientIP: Client IP
  all_user: All User
  topic: Topic
  access: Access
  access_options:
    d: Disabled
    r: Publish
    w: Subscribe
    rw: Publish & Subscribe
</i18n>

<style lang="scss" scoped>
.custom-form-item {
  :deep(.el-form-item__label) {
    display: inline-flex;
    align-items: center;
    width: 100% !important;
  }
}
</style>
