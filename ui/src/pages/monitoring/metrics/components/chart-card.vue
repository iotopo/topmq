<script lang="ts" setup>
import { inject, ref } from 'vue'

defineProps<{
  title: string
}>()

const emit = defineEmits(['maximize'])

const interval = inject('interval')

const chartRef = ref()
const chartOption = ref({})

defineExpose({
  setOption(option: any) {
    chartOption.value = option
    chartRef.value?.setOption(option)
    chartRef.value?.resize()
  },
  resize() {
    chartRef.value?.resize()
  },
})
</script>

<template>
  <el-card class="chart-card">
    <template #header>
      <div w-full>
        <div w-full flex items-center>
          <span>{{ title }}</span>
          <div class="flex-grow" />
          <span
            style="
              color: var(--el-text-color-secondary);
              margin-right: 10px;
              padding-bottom: 2px;
            "
          >周期：{{ interval }}</span>
          <slot name="tools" />
          <el-button-group>
            <el-button link @click="emit('maximize')">
              <el-icon size="18">
                <i-mdi-launch />
              </el-icon>
            </el-button>
          </el-button-group>
        </div>
        <div flex items-center>
          <!-- <span
            style="color: var(--el-text-color-secondary); margin-right: 10px"
            >周期：{{ interval }}</span
          > -->
          <!-- <slot name="tools" /> -->
        </div>
      </div>
    </template>
    <div class="chart-body">
      <echarts-viewer
        ref="chartRef"
        :option="chartOption"
        style="width: 100%; height: 100%"
      />
    </div>
  </el-card>
</template>

<style lang="scss" scoped>
.chart-card {
  .chart-header {
    display: flex;
    align-items: center;
    width: 100%;
  }
}

.chart-body {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
}
</style>
