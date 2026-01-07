<script lang="ts" setup>
import type { ElForm } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import authApi from '@/api/auth'
import { useStore } from '@/store'
import { passowrdPattern } from '@/utils/validates'

type FormInstance = InstanceType<typeof ElForm>

const formRef = ref<FormInstance>()
const formDisabled = ref(false)

const { t } = useI18n()

// interface Props {
//   isNew: boolean
// }

// const props = withDefaults(defineProps<Props>(), {
//   isNew: false
// })
// const { isNew } = toRefs(props)

const store = useStore()

const dialogVisible = ref(false)

const form = reactive({
  oldPassword: '',
  password: '',
  passwordConfirm: '',
})

const minPwdLen = computed(() => store.minPwdLen)

const rules = reactive({
  oldPassword: [
    {
      required: true,
      message: t('common.validation.required'),
      trigger: 'change',
    },
  ],
  password: [
    {
      required: true,
      message: t('common.validation.required'),
      trigger: 'change',
    },
    {
      message: t('validation.format'),
      trigger: 'change',
      pattern: passowrdPattern,
    },
    {
      min: minPwdLen.value,
      max: 30,
      message: t('validation.size', { min: minPwdLen.value }),
      trigger: 'blur',
    },
    {
      validator: (rule: any, value: any, callback: any) => {
        const reg = /[\u4E00-\u9FA5]/g
        if (reg.test(value)) {
          callback(new Error(t('validation.chinese')))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
  passwordConfirm: [
    {
      required: true,
      message: t('common.validation.required'),
      trigger: 'change',
    },
    {
      validator: (rule: any, value: any, callback: any) => {
        if (form.password !== value) {
          callback(new Error(t('validation.inconformity')))
          return
        }
        callback()
      },
    },
  ],
})

watchEffect(() => {
  if (!dialogVisible.value) {
    form.oldPassword = ''
    form.password = ''
    form.passwordConfirm = ''
    formRef.value?.clearValidate()
  }
})

// onMounted(() => {
// })
defineExpose({
  showDialog: () => {
    formRef.value?.clearValidate()
    dialogVisible.value = true
  },
})

function submitForm(formEl: FormInstance | undefined) {
  if (!formEl)
    return
  formEl.validate(async (valid) => {
    if (valid) {
      try {
        await authApi.resetPassword({
          oldPassword: form.oldPassword,
          password: form.password,
        })
        // await userService.bindMobile({
        //   id: codeID,
        //   mobile: form.mobile,
        //   code: form.verifyCode,
        // })
        dialogVisible.value = false
        ElMessage.success(t('common.success.operate'))

        // 注销
        await store.logout()
        // const url = router.resolve('/login').href
        // window.location.href = url
      } catch (error: any) {
        console.error(error)
        if (error.code) {
          ElMessage.error(t(`error.${error.code}`))
        }
      }
    }
  })
}
</script>

<template>
  <el-dialog
    v-model="dialogVisible"
    :title="t('title.passwords')"
    width="500px"
  >
    <el-form ref="formRef" size="large" :model="form" :rules="rules">
      <el-form-item prop="oldPassword">
        <el-input
          v-model="form.oldPassword"
          :placeholder="t('ph.old_pw')"
          show-password
        />
      </el-form-item>
      <el-form-item prop="password">
        <el-input
          v-model="form.password"
          :placeholder="t('ph.new_pw')"
          show-password
        />
      </el-form-item>
      <el-form-item prop="passwordConfirm">
        <el-input
          v-model="form.passwordConfirm"
          :placeholder="t('ph.pw_confirm')"
          show-password
        />
      </el-form-item>
    </el-form>

    <div style="margin-bottom: 10px">
      {{ t('msg.pw_requirements') }}：
    </div>
    <div
      class="secondary-text-color"
      style="font-size: var(--el-font-size-small)"
    >
      <!-- <el-icon>
            <i-mdi-check-circle-outline />
      </el-icon> -->
      <span>{{ t('validation.size', { min: minPwdLen }) }}</span>
    </div>
    <div
      class="secondary-text-color"
      style="font-size: var(--el-font-size-small)"
    >
      <!-- <el-icon>
            <i-mdi-check-circle-outline />
      </el-icon> -->
      <span>{{ t('validation.format') }}</span>
    </div>
    <div
      class="secondary-text-color"
      style="margin-top: 20px; font-size: var(--el-font-size-small)"
    >
      {{ t('msg.reuse') }}
    </div>
    <!-- <el-link type="primary" href="/ui/account/retrieve" style="margin: 10px 0">
      {{ t('forget') }}
    </el-link> -->
    <template #footer>
      <span class="dialog-footer">
        <el-button size="large" @click="dialogVisible = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button
          size="large"
          type="primary"
          :disabled="formDisabled"
          @click="submitForm(formRef)"
        >
          {{ t('common.ok') }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<i18n lang="yaml" src="./lang.yaml" />
