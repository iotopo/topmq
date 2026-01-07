<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import notfoundUrl from '@/assets/not-found.svg'

const { t } = useI18n()
function goHome() {
  location.href = '/'
}
</script>

<template>
  <div
    class="h-full w-full flex items-center justify-center text-center"
    style="background: var(--el-bg-color)"
  >
    <div class="max-w-740px">
      <img :src="notfoundUrl" class="m-auto block h-260px -mt-15vh">
      <!-- <h1>404</h1> -->
      <p style="color: #777">
        {{ t("core_pages.404_tips") }}
      </p>
      <el-button plain size="default" @click="goHome">
        {{ t("core_pages.back_home") }}
      </el-button>
    </div>
  </div>
</template>
