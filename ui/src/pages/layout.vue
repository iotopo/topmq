<script lang="ts" setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import NavTools from './components/nav-tools.vue'
import SideMenu from './components/side-menu.vue'

const { t } = useI18n()

const route = useRoute()
// const router = useRouter()
// const store = useStore()

const collapsed = ref(false)
</script>

<template>
  <div class="main-view" flex>
    <SideMenu :collapsed="collapsed" />
    <div h-full flex-1>
      <div w-full class="page-title">
        <!-- {{ currentPageTitle }} -->
        <el-button link m-r-10px @click="collapsed = !collapsed">
          <el-icon size="18">
            <i-ep-expand v-if="collapsed" />
            <i-ep-fold v-else />
          </el-icon>
        </el-button>
        <el-breadcrumb separator="/">
          <template v-for="item in route.matched" :key="item.path">
            <el-breadcrumb-item
              v-if="!item.meta?.activeMenu"
              :to="{ path: item.path }"
            >
              {{ t(item.meta!.title as string) }}
            </el-breadcrumb-item>
          </template>
        </el-breadcrumb>
        <div style="flex: 1" />
        <NavTools />
      </div>
      <div w-full style="height: calc(100% - 48px)">
        <router-view v-slot="{ Component }">
          <keep-alive>
            <component
              :is="Component"
              v-if="route.meta?.cache"
              :key="route.name"
            />
          </keep-alive>
          <component
            :is="Component"
            v-if="!route.meta?.cache"
            :key="route.name"
          />
        </router-view>
      </div>
    </div>
  </div>
</template>

<i18n lang="yaml" src="./nav.yaml"></i18n>

<style lang="scss" scoped>
$menu-footer-height: 56px;

.main-view {
  // flex: 1;
  position: relative;
  // height: calc(100% - 32px);
  width: 100%;
  height: 100%;
  background-color: var(--el-bg-color-page);

  .page-title {
    height: 48px;
    background: var(--el-bg-color);
    padding-left: 16px;
    display: flex;
    align-items: center;
  }

  .router-content {
    width: 100%;
    height: 100%;
  }
}

.hamburger-container {
  line-height: 46px;
  height: 100%;
  // float: left;
  cursor: pointer;
  transition: background 0.3s;
  -webkit-tap-highlight-color: transparent;
  :deep(.hamburger) {
    // fill: var(--el-text-color-primary) !important;
    fill: #fff !important;
  }
  &:hover {
    background: rgba(0, 0, 0, 0.025);
  }
}
</style>
