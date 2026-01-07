<script setup lang="ts">
const emit = defineEmits(['click'])
</script>

<template>
  <el-button text bg style="padding: 8px" @click="$emit('click')">
    <el-icon>
      <i-ep-refresh-right />
    </el-icon>
  </el-button>
</template>

<style lang="scss" scoped></style>
