<script lang="ts" setup>
import type { ElForm } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import selfService from '@/api/self'

type FormInstance = InstanceType<typeof ElForm>

const formRef = ref<FormInstance>()

const { t } = useI18n()

const form = reactive({
  username: '',
  name: '',
  tel: '',
  mobile: '',
  email: '',
})

const emailPattern = /^\w+([+.-]\w+)*@\w+([.-]\w+)*\.\w+([.-]\w+)*$/
const rules = reactive({
  name: [
    {
      required: true,
      message: t('common.validation.required'),
      trigger: 'change',
    },
  ],
  email: {
    // type: 'email',
    message: t('validation.email'),
    trigger: 'change',
    pattern: emailPattern,
  },
})
onMounted(async () => {
  const info = await selfService.getSelfInfo()
  Object.assign(form, info)
})

function submitForm(formEl: FormInstance | undefined) {
  if (!formEl)
    return
  formEl.validate(async (valid) => {
    if (valid) {
      try {
        await selfService.setSelfInfo({
          name: form.name,
          tel: form.tel,
        })
        ElMessage.success(t('common.success.save'))
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
    }
  })
}
</script>

<template>
  <div>
    <div class="container">
      <div class="title">
        <span>{{ t('title.basic') }}</span>
      </div>
      <div>
        <el-form
          ref="formRef"
          size="large"
          label-width="auto"
          style="width: 500px; margin: 0 auto"
          :model="form"
          :rules="rules"
        >
          <el-form-item prop="username" :label="t('username')">
            <el-input
              v-model="form.username"
              disabled
              autocomplete="off"
            />
          </el-form-item>
          <el-form-item prop="name" :label="t('name')">
            <el-input
              v-model="form.name"
              autocomplete="off"
              :placeholder="
                t('global.placeholders.please_input', {
                  text: t('name').toLocaleLowerCase(),
                })
              "
            />
          </el-form-item>
          <el-form-item prop="mobile" :label="t('mobile')">
            <el-input
              v-model="form.mobile"
              autocomplete="off"
            />
          </el-form-item>
          <el-form-item prop="email" :label="t('email')">
            <el-input
              v-model="form.email"
              autocomplete="off"
            />
          </el-form-item>
          <el-form-item prop="tel" :label="t('tel')">
            <el-input
              v-model="form.tel"
              :placeholder="
                t('global.placeholders.please_input', {
                  text: t('tel').toLocaleLowerCase(),
                })
              "
              maxlength="20"
              autocomplete="off"
            />
          </el-form-item>
        </el-form>
        <el-row
          style="flex-direction: column; align-items: center; margin-top: 20px"
        >
          <el-button
            size="large"
            round
            style="width: 120px"
            @click="submitForm(formRef)"
          >
            {{ t('common.save') }}
          </el-button>
        </el-row>
      </div>
    </div>
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
    margin-bottom: 20px;
  }
}
</style>
