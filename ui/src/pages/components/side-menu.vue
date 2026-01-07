<script lang="ts" setup>
import type { RouteRecordRaw } from 'vue-router'
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import LogoImage from '@/assets/logo.png'
import dictIcon from '@/icons'
import { routes } from '@/router'
import { useStore } from '@/store'
import MenuItem from './menu-item.vue'

const props = defineProps<{
  collapsed: boolean
}>()

const { t } = useI18n()
const route = useRoute()
// const router = useRouter()
const store = useStore()

const activedMenu = ref((route.meta?.activeMenu || route.path) as string)

// 监听路由变化
watch(
  () => route.path,
  (newPath) => {
    activedMenu.value = newPath
  },
)

/**
 * 是否为菜单
 */
function isMenu(item: RouteRecordRaw) {
  if (item.meta?.alwaysShow) {
    // 永远显示根
    return false
  }
  return Array.isArray(item.children) && item.children?.length > 0
}

function icon(name: string) {
  return (dictIcon as any)[name]
}
</script>

<template>
  <div h-full class="menu-container">
    <div class="menu-header">
      <img alt="" :src="LogoImage" height="28">
      <h2 v-if="!props.collapsed" class="menu-h2">
        TopMQ
      </h2>
    </div>
    <div w-full style="height: calc(100% - var(--tp-header-height) - 32px)">
      <el-scrollbar w-full>
        <el-menu :default-active="activedMenu" router w-full class="aside-menu" :collapse="props.collapsed" style="border-right: 0">
          <template v-for="item in routes" :key="item.path">
            <template
              v-if="!item.meta?.hidden && (typeof item.meta?.visibleFunc === 'function' ? item.meta.visibleFunc() : true)"
            >
              <el-sub-menu v-if="isMenu(item)" :index="item.path">
                <template #title>
                  <el-icon v-if="item.meta?.icon && icon(item.meta.icon as string)">
                    <component
                      :is="icon(item.meta.icon as string)"
                    />
                  </el-icon>
                  <span>{{ t(item.meta?.title as string) }}</span>
                </template>

                <template v-for="subItem in item.children" :key="subItem.path">
                  <MenuItem
                    v-if="!subItem.meta?.hidden && (typeof subItem.meta?.visibleFunc === 'function' ? subItem.meta.visibleFunc() : true)"
                    :path="`${item.path}/${subItem.path}`" :title="t(subItem.meta?.title as string)"
                  />
                </template>
              </el-sub-menu>
              <MenuItem v-else :path="item.path" :title="t(item.meta?.title as string)" :icon="item.meta?.icon as string" />
            </template>
          </template>
        </el-menu>
      </el-scrollbar>
    </div>
    <div v-if="!props.collapsed" class="menu-footer">
      {{ t('version', { version: store.version }) }}
    </div>
  </div>
</template>

<i18n lang="yaml" src="../nav.yaml"></i18n>

<style lang="scss" scoped>
.menu-container {
  background-color: var(--tp-menu-bg-color);
  border-right: solid 1px var(--el-menu-border-color);

  .menu-header {
    color: var(--el-text-color-regular);
    padding: 0 16px;
    position: relative;
    display: flex;
    flex-direction: row;
    align-items: center;
    height: var(--tp-header-height);

    // width: 200px;
    .menu-title-icon {
      color: var(--tp-menu-text-color);
    }

    .menu-h2 {
      margin-left: 10px;
      font-size: 16px;
      color: #fff;
    }
  }

  .menu-footer {
    height: 32px;
    color: var(--tp-menu-text-color);
    padding: 0 16px;
    position: relative;
    display: flex;
    flex-direction: row;
    justify-content: flex-start;
  }
  .aside-menu {
    --el-menu-base-level-padding: 16px;
    --el-menu-level-padding: 16px;
    --el-menu-bg-color: var(--tp-menu-bg-color);
    --el-menu-hover-bg-color: var(--tp-menu-hover-bg-color);
    --el-menu-text-color: var(--tp-menu-text-color);
    --el-menu-hover-text-color: var(--tp-menu-text-color);
    --el-menu-active-color: var(--tp-menu-active-color);
    --el-menu-item-height: var(--tp-menu-item-height);
    --el-menu-sub-item-height: var(--tp-menu-sub-item-height);

    .el-menu-item.is-active {
      background-color: var(--el-color-primary);
    }

    .is-opened {
      --el-menu-bg-color: var(--tp-submenu-bg-color);
      background-color: var(--tp-submenu-bg-color);
    }
  }

  .aside-menu:not(.el-menu--collapse) {
    width: 200px;
  }
}
</style>
