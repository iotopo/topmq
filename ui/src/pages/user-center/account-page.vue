<script lang="ts" setup>
const props = defineProps<{
  title: string
}>()
</script>

<template>
  <div class="account-page">
    <slot>
      <el-skeleton />
    </slot>
  </div>
</template>

<style scoped>
.account-page {
  height: 100%;
}
</style>
