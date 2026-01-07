<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import AccountPage from './account-page.vue'

const { t } = useI18n()

const router = useRouter()
const route = useRoute()

const activeMenu = ref<string>()

onMounted(async () => {
  if (route.meta.activeMenu) {
    activeMenu.value = route.meta.activeMenu as string
  }
})

function changeRouter(index: string) {
  router.push(index)
}
</script>

<template>
  <AccountPage :title="t('title.page')">
    <el-container
      style="width: 1200px; height: 100%; margin: 0 auto"
      class="base-container"
    >
      <el-aside width="200px">
        <el-menu
          :default-active="activeMenu"
          style="height: 100%"
          @select="changeRouter"
        >
          <el-menu-item index="basic">
            <el-icon>
              <i-mdi-account />
            </el-icon>
            <span>{{ t('title.basic') }}</span>
          </el-menu-item>
          <el-menu-item index="security">
            <el-icon>
              <i-mdi-security />
            </el-icon>
            <span>{{ t('title.security') }}</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view v-slot="{ Component }">
          <transition mode="out-in" name="router">
            <keep-alive>
              <component :is="Component" />
            </keep-alive>
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </AccountPage>
</template>

<i18n lang="yaml" src="./lang.yaml" />

<style lang="scss" scoped>
.base-container {
  background: var(--el-fill-color-blank);
}
</style>
