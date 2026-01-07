<script setup lang="ts">
import { getCurrentInstance, provide, ref } from 'vue'

const showAdv = ref(false)

// 提供组件实例给子组件
const instance = getCurrentInstance()
provide('search', instance)
</script>

<template>
  <div class="search">
    <div class="search-main">
      <slot />
      <el-button
        v-if="$slots.adv"
        text
        bg
        :type="showAdv ? 'primary' : 'default'"
        style="padding: 8px"
        @click="showAdv = !showAdv"
      >
        <el-icon>
          <i-ep-filter />
        </el-icon>
      </el-button>
    </div>
    <el-collapse-transition>
      <div v-show="showAdv" class="search-adv">
        <slot name="adv" />
      </div>
    </el-collapse-transition>
  </div>
</template>

<style lang="scss" scoped>
.search {
  :deep(.search-main > .el-input) {
    width: 180px;
    flex: none;
  }

  .search-main {
    display: flex;
    flex-direction: row;
    align-items: center;
    // flex-wrap: wrap;

    & > :deep(*) {
      margin-left: 0 !important;
      margin-right: 10px;
    }
  }

  .search-adv {
    display: flex;
    margin-top: 10px;

    & > :deep(*) {
      margin-left: 0;
      margin-right: 10px;
    }
  }
}
</style>
