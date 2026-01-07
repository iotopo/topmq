<script lang="ts" setup>
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

import { useRoute } from 'vue-router'
import selfService from '@/api/self'
import ResetPassword from './reset-password.vue'

const { t } = useI18n()

const route = useRoute()

const resetPasswordDialog = ref()
const form = reactive({
  username: '',
  name: '',
  mobile: '',
  email: '',
  wechat_web_bound: false,
  wechat_mini_bound: false,
})

async function loadData() {
  try {
    const res = await selfService.getSelfInfo()
    form.username = res.username
    form.name = res.name
    form.mobile = res.mobile
    form.email = res.email
  }
  catch (error: any) {
    if (error?.code && error.code.startsWith('err_')) {
      ElMessage.error(t(error.code, error.msg))
    }
    else {
      error && console.error(error)
    }
    // } finally {
  }
}

onMounted(() => {
  const err = route.query.err
  if (err) {
    ElMessage.error(t(err as string))
  }
  loadData()
  // console.log(route.query)
  if (route.query.changPassword)
    resetPasswordDialog.value.showDialog()
})
</script>

<template>
  <div>
    <div class="container">
      <div class="title">
        <span>{{ t('title.security') }}</span>
      </div>
      <el-row class="basic-item-container">
        <el-button type="primary" size="large" circle>
          <el-icon :size="24">
            <i-mdi-account />
          </el-icon>
        </el-button>
        <div class="middle">
          <div>
            <span>{{ t('username') }}</span>
          </div>
          <div class="secondary-text-color">
            <span>{{ form.username }}</span>
          </div>
        </div>
      </el-row>
      <div class="line-separator" />
      <el-row class="basic-item-container">
        <div class="left">
          <el-button type="warning" size="large" circle>
            <el-icon :size="24">
              <i-mdi-key />
            </el-icon>
          </el-button>
        </div>
        <div class="middle">
          <div>
            <span>{{ t('title.passwords') }}</span>
          </div>
          <div class="secondary-text-color">
            <!-- <span>{{ form.email || '未绑定' }}</span> -->
          </div>
        </div>
        <div class="right">
          <el-button @click="resetPasswordDialog.showDialog()">
            {{
              t('reset')
            }}
          </el-button>
        </div>
      </el-row>
    </div>
    <ResetPassword ref="resetPasswordDialog" />
  </div>
</template>

<i18n lang="yaml" src="./lang.yaml" />

<style lang="scss" scoped>
.container {
  width: 960px;
  background: var(--el-fill-color-blank);
  margin: auto;
  padding: 50px;
  min-height: 500px;
  .title {
    text-align: center;
    margin: 0;
    font-size: 32px;
  }
  .basic-item-container {
    padding: 20px 0;
    .middle {
      flex: 1;
      margin: 0 16px;
      div:first-child {
        margin-bottom: 4px;
      }
    }
  }
  .line-separator {
    margin-left: 54px;
    border-top: 1px solid #e5e5e5;
  }
}
</style>
