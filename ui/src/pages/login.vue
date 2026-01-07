<script lang="ts" setup>
import dayjs from 'dayjs'
import { computed, getCurrentInstance, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import authService from '@/api/auth'
import captchaService from '@/api/captcha'
import { useStore } from '@/store'
import LoginBackground from '../assets/login-bg.png'

const { t } = useI18n()
const store = useStore()
const route = useRoute()
const router = useRouter()
const currentInstance = getCurrentInstance()

function setLocale(locale: string) {
  store.setLocale(locale as 'zh' | 'en')
  currentInstance?.proxy?.$forceUpdate()
}

const show1 = ref(false)
const show2 = ref(false)
const show3 = ref(false)
const form = ref({
  username: '',
  password: '',
  captcha: '',
  verifyCode: '',
})
const showCaptcha = ref(false)
const captchaId = ref('')
const captchaImg = ref('')
const redirect = ref('') // 重定向到新地址
const redirect_blacklist = ['/ui/account/security']

// `Copyright © 2018 - 2023 IOTOPO. All Rights Reserved. 图扑物联 版权所有`

// const copyright = ref('Copyright © 2018 - 2026 IOTOPO. All Rights Reserved. 图扑物联 版权所有')
const copyright = computed(() => {
  return `技术支持：图扑物联科技有限公司（www.iotopo.com）`
})
const copyrightLink = ref('https://www.iotopo.com')
const loginBackgroundColor = ref('rgb(52, 86, 234)')
const loginBgUrl = LoginBackground

// const authApi = api.extend(options => ({
//   prefixUrl: `${options.prefixUrl}/auth`,
// }))
const isFormValidate = computed(() => {
  if (!form.value.username || !form.value.password) {
    return false
  }
  return true
})
async function updateCaptcha() {
  try {
    const showCaptchaRes = await captchaService.isCaptchaShown()
    showCaptcha.value = showCaptchaRes.show

    const captchIdRes = await captchaService.getCaptchaID()
    if (captchIdRes.captchaID) {
      captchaId.value = captchIdRes.captchaID
      captchaImg.value = `/captcha/server/${captchaId.value}.png`
    }
  } catch (error) {
    console.error(error)
  }
}

function loginSuccess() {
  if (redirect.value) {
    if (redirect_blacklist.includes(redirect.value.split('?')[0])) {
      router.push({ path: '/' })
    } else {
      location.href = redirect.value
    }
  } else {
    router.push({ path: '/' })
  }
}

async function login() {
  if (!isFormValidate.value) {
    return
  }

  const data = {
    username: form.value.username,
    // password: form.value.password,
    password: window.btoa(form.value.password),
    captchaID: showCaptcha.value ? captchaId.value : '',
    captchaValue: form.value.captcha,
    // solution_id: solutionID.value, // 解决方案
  }

  try {
    const resp = await authService.login(data)
    const token = resp.token
    localStorage.setItem('token', `Bearer ${token}`)
    store.setUser({
      id: resp.id,
      name: resp.name,
      username: resp.username,
      isSuperuser: resp.isSuperuser,
    })
    loginSuccess()
  } catch {
    updateCaptcha()
  }
}

function wait(time: number) {
  return new Promise((next) => {
    setTimeout(() => {
      next(null)
    }, time)
  })
}

onMounted(async () => {
  // let search = window.location.search
  await wait(400)
  show1.value = true
  await wait(300)
  show2.value = true
  await wait(420)
  show3.value = true

  redirect.value = route.query?.redirect as string
  if (redirect.value) {
    redirect.value = route.query?.redirect as string
  }

  updateCaptcha()
})
</script>

<template>
  <div class="login">
    <div class="left" :style="{ background: loginBackgroundColor }">
      <transition name="el-fade-in">
        <div v-show="show1" class="title">
          图扑物联 TopMQ
        </div>
      </transition>
      <transition name="el-fade-in">
        <div v-show="show2" class="sub-title">
          MQTT 5.0 消息中间件
        </div>
      </transition>
      <div class="bg">
        <img :src="loginBgUrl">
      </div>
    </div>
    <div class="right">
      <div style="position: absolute; top: 20px; right: 20px">
        <el-button size="large" link @click="setLocale('zh')">
          中文
        </el-button>
        <el-divider direction="vertical" />
        <el-button size="large" link style="margin-left: 0" @click="setLocale('en')">
          EN
        </el-button>
      </div>
      <transition name="el-fade-in">
        <div v-show="show3" class="login-view">
          <div class="title">
            {{ t("username_login") }}
          </div>
          <!-- <div class="sub-title">
            {{ t('no_account_tips') }}
          </div> -->
          <el-form :model="form" class="login-form" label-position="top" label-width="80px" size="large">
            <el-form-item :label="t('username')">
              <el-input v-model="form.username" :placeholder="t('placeholder.username')" @keyup.enter="login" />
            </el-form-item>
            <el-form-item :label="t('password')">
              <el-input
                v-model="form.password" type="password" :placeholder="t('placeholder.password')" show-password
                @keyup.enter="login"
              />
            </el-form-item>

            <el-form-item v-if="showCaptcha" class="password-item" :label="t('verify_code')">
              <el-input v-model="form.captcha" maxlength="4" @keyup.enter="login" />
              <img :src="captchaImg" @click="updateCaptcha">
            </el-form-item>
            <el-form-item>
              <el-button size="large" :disabled="!isFormValidate" class="login-btn" type="primary" @click="login">
                {{ t("enter") }}
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </transition>
      <div v-show="!!copyright" class="co">
        <template v-if="!!copyrightLink">
          <el-link type="info" :href="copyrightLink" target="_blank">
            {{
              copyright
            }}
          </el-link>
        </template>
        <template v-else>
          {{ copyright }}
        </template>
      </div>
    </div>
  </div>
</template>

<i18n lang="yaml">
zh:
  welcome_to: 欢迎来到
  welcome_login: 欢迎登录！
  no_account_tips: 没有账号？请联系管理员
  primary_account: 主账号
  username: 用户名
  password: 密码
  enter: 登录
  other_login: 其他登录方式
  sso_login: SSO 账号登录
  placeholder:
    project_id: 请输入项目主账号
    username: 请输入用户名
    platform_username: 请输入平台管理账号
    password: 请输入密码
    verify_code: 请输入短信验证码
    mobile: 请输入手机号
  forget_the_password: 忘记密码
  project_user_login: 项目用户登录
  platform_user_login: 平台用户登录
  subuser: 项目子用户
  username_login: 账号登录
  verify_code_login: 短信登录
  init_pass_info: 您的密码为初始密码，请修改密码
  skip: 跳过
  tips: 提示
  mobile_number: 手机号
  verify_code: 验证码
  reacquire: 重新获取
  get_code: 获取验证码
  err_mobile_not_found: 账号不存在或已被删除。
  send_sms_tips: 短信已发送至您的手机, 5分钟内有效

en:
  welcome_to: Welcome to
  welcome_login: Account Login
  no_account_tips: Don't have an account? Please contact    the administrator.
  primary_account: Primary Account
  username: Username
  password: Password
  enter: Login
  other_login: Other Login
  sso_login: SSO Login
  placeholder:
    project_id: Please enter the project primary account
    username: Please enter the username
    platform_username: Please enter the platform management account
    password: Please enter the password
    verify_code: Please enter the SMS verification code
    mobile: Please enter the mobile number
  forget_the_password: Forget Password
  project_user_login: Project User Login
  platform_user_login: Platform User Login
  subuser: Project Subuser
  username_login: Account Login
  verify_code_login: SMS Login
  init_pass_info: Your password is the initial password,    please change it
  skip: Skip
  tips: Tips
  mobile_number: Mobile Number
  verify_code: Verification Code
  reacquire: Reacquire
  get_code: Get Verification Code
  err_mobile_not_found: The account does not exist or  may be deleted.
  send_sms_tips: SMS has been sent to your mobile phone, valid for 5 minutes.
</i18n>

<style lang="scss">
.login {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: row;

  .left {
    flex: 1;
    background: rgb(52, 86, 234);
    box-sizing: content-box;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    position: relative;

    .title {
      padding-top: 80px;
      padding-left: 60px;
      padding-right: 60px;
      font-size: 28px;
      color: #fff;
      font-weight: bold;
    }

    .sub-title {
      padding-top: 10px;
      padding-left: 60px;
      padding-right: 60px;
      font-size: 28px;
      color: #fff;
      font-weight: bold;
    }

    .bg {
      position: absolute;
      bottom: 5%;
      top: 15%;
      flex: 1;
      padding: 80px;
      text-align: center;
      display: flex;
      justify-content: flex-start;
      align-items: center;

      img {
        display: block;
        width: 100%;
      }
    }
  }

  .right {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
    position: relative;

    .login-view {
      padding: 80px;
      width: 520px;

      .title {
        font-size: 20px;
        font-weight: bold;
      }

      .sub-title {
        padding-top: 10px;
        font-size: 16px;
        font-weight: bold;
      }

      .login-form {
        margin-top: 40px;

        .el-link {
          padding: 4px 0;
        }

        :deep(.el-form-item__label) {
          font-weight: bold;
          color: #000;
        }

        .password-item {
          :deep(.el-form-item__content) {
            display: flex;

            img {
              height: 36px;
              margin-left: 10px;
              cursor: pointer;
            }
          }
        }

        .login-btn {
          // margin-top: 30px;
          width: 100%;
          // font-size: 18px;
          // font-weight: bold;
          // padding-top: 12px;
          // padding-bottom: 12px;
        }

        .tools {
          display: flex;
          align-items: center;
          justify-content: space-between;
          width: 100%;
          line-height: 32px;
        }
      }
    }

    .co {
      position: absolute;
      bottom: 16px;
      font-size: 14px;
      color: #9b9b9b;
    }
  }
}
</style>
