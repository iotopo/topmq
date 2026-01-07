<script lang="ts" setup>
import type { Ref } from 'vue'
import dayjs from 'dayjs'
import { computed, inject, nextTick, ref, watch } from 'vue'
import TimeRangeSelector from './time-range-selector.vue'

const props = defineProps<{
  title: string
}>()

const emit = defineEmits(['reload'])

const visible = ref(false)

const timeRangeSelectorRef = ref()

const customTime = inject<Ref<[Date, Date]>>('dialogCustomTime')

const interval = computed(() => {
  if (!customTime) {
    return '15秒'
  }
  const duration = dayjs(customTime.value[1]).diff(dayjs(customTime.value[0]), 'hour', false)
  if (duration >= 14 * 24) {
    return '120分钟'
  } else if (duration >= 7 * 24) {
    return '60分钟'
  } else if (duration >= 3 * 24) {
    return '30分钟'
  } else if (duration >= 24) {
    return '10分钟'
  } else if (duration >= 12) {
    return '5分钟'
  } else if (duration >= 6) {
    return '1分钟'
  } else if (duration >= 3) {
    return '1分钟'
  } else {
    return '15秒'
  }
})

function reloadCharts() {
  if (!customTime) {
    return
  }
  nextTick(() => {
    emit('reload')
  })
}

const chartRef = ref()
const chartOption = ref({})

defineExpose({
  show() {
    visible.value = true
    nextTick(() => {
      timeRangeSelectorRef.value?.reset()
    })
  },
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
  <el-dialog v-model="visible" :title="title" width="60%">
    <div>
      <div mb-4 flex items-center>
        <!-- <el-radio-group v-model="timeRange">
          <el-radio-button
            v-for="item in timeRanges"
            :key="item.value"
            :value="item.value"
            :label="item.label"
          >
          </el-radio-button>
        </el-radio-group>
        <el-date-picker
          v-if="timeRange === 'custom'"
          v-model="customTime"
          type="datetimerange"
          class="ml-2"
          style="flex-grow: unset"
          @change="reloadCharts"
        /> -->
        <TimeRangeSelector ref="timeRangeSelectorRef" v-model="customTime" @change="reloadCharts" />
        <div flex-1 />
        <div flex items-center style="padding: 0 10px">
          <span
            style="
              color: var(--el-text-color-secondary);
              margin-right: 10px;
              padding-bottom: 2px;
            "
          >周期：{{ interval }}</span>
          <slot name="tools" />
        </div>
        <el-button @click="reloadCharts">
          <el-icon><i-mdi-reload /></el-icon>
        </el-button>
      </div>
      <div style="border: var(--el-border-color) 1px solid; height: 500px">
        <echarts-viewer
          ref="chartRef"
          :option="chartOption"
          style="width: 100%; height: 100%; padding: 10px"
        />
      </div>
    </div>
  </el-dialog>
</template>
