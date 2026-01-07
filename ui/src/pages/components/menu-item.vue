<script lang="ts" setup>
const { path, title } = defineProps<{
  path: string
  title: string
}>()
</script>

<template>
  <el-menu-item :index="path">
    <template #title>
      <div class="menu-item-content">
        <span class="menu-item-title">{{ title }}</span>
      </div>
    </template>
  </el-menu-item>
</template>

<style lang="scss" scoped>
.menu-item-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 8px;

  .menu-item-title {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .favorite-icon {
    flex-shrink: 0;
    font-size: 16px;
    color: var(--el-text-color-secondary);
    cursor: pointer;
    transition:
      color 0.2s,
      transform 0.1s,
      opacity 0.2s;
    padding: 2px;
    margin-left: auto;

    // 未收藏的默认隐藏，hover 时显示
    opacity: 0;
    pointer-events: none;

    // 已收藏的始终显示
    &.is-favorite {
      opacity: 1;
      pointer-events: auto;
      color: var(--el-color-warning);
    }

    &:hover {
      color: var(--el-color-warning);
    }

    &:active {
      transform: scale(0.95);
    }
  }
}
</style>

<style lang="scss">
// 菜单项 hover 时显示收藏按钮（使用全局样式，因为需要匹配 el-menu-item）
.el-menu-item:hover .menu-item-content .favorite-icon {
  opacity: 1;
  pointer-events: auto;
}
</style>
