<script lang="ts" setup>
import dayjs from 'dayjs'
import { ref, watch } from 'vue'

interface TimeRange {
  label: string
  value: string
}

interface Props {
  modelValue?: [Date, Date]
}

interface Emits {
  (e: 'update:modelValue', value: [Date, Date]): void
  (e: 'change'): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => [dayjs().subtract(1, 'hour').toDate(), dayjs().toDate()],
})

const emit = defineEmits<Emits>()

const timeRange = ref('1h')
const customTime = ref<[Date, Date]>(props.modelValue)

const timeRanges = ref<TimeRange[]>([
  { label: '1小时', value: '1h' },
  { label: '3小时', value: '3h' },
  { label: '6小时', value: '6h' },
  { label: '12小时', value: '12h' },
  { label: '1天', value: '1d' },
  { label: '3天', value: '3d' },
  { label: '7天', value: '7d' },
  { label: '14天', value: '14d' },
  { label: '自定义', value: 'custom' },
])

// 监听外部传入的 modelValue 变化
watch(
  () => props.modelValue,
  (newValue) => {
    customTime.value = newValue
  },
  { deep: true },
)

// 监听 customTime 变化，向父组件发送更新
watch(
  customTime,
  (newValue) => {
    emit('update:modelValue', newValue)
  },
  { deep: true },
)

function handleTimeRangeChange() {
  if (timeRange.value === 'custom') {
    return
  }

  let start: Date
  let end: Date

  if (timeRange.value === '1h') {
    start = dayjs().subtract(1, 'hour').toDate()
    end = dayjs().toDate()
  } else if (timeRange.value === '3h') {
    start = dayjs().subtract(3, 'hour').toDate()
    end = dayjs().toDate()
  } else if (timeRange.value === '6h') {
    start = dayjs().subtract(6, 'hour').toDate()
    end = dayjs().toDate()
  } else if (timeRange.value === '12h') {
    start = dayjs().subtract(12, 'hour').toDate()
    end = dayjs().toDate()
  } else if (timeRange.value === '1d') {
    start = dayjs().subtract(1, 'day').toDate()
    end = dayjs().toDate()
  } else if (timeRange.value === '3d') {
    start = dayjs().subtract(3, 'day').toDate()
    end = dayjs().toDate()
  } else if (timeRange.value === '7d') {
    start = dayjs().subtract(7, 'day').toDate()
    end = dayjs().toDate()
  } else if (timeRange.value === '14d') {
    start = dayjs().subtract(14, 'day').toDate()
    end = dayjs().toDate()
  } else {
    return
  }

  customTime.value = [start, end]
  emit('change')
}

function handleCustomTimeChange() {
  emit('change')
}

// 暴露当前的时间范围值
defineExpose({
  timeRange,
  customTime,
  reset() {
    handleCustomTimeChange()
  },
})
</script>

<template>
  <div class="time-range-selector">
    <el-radio-group v-model="timeRange" @change="handleTimeRangeChange">
      <el-radio-button
        v-for="item in timeRanges"
        :key="item.value"
        :value="item.value"
        :label="item.label"
      />
    </el-radio-group>
    <el-date-picker
      v-if="timeRange === 'custom'"
      v-model="customTime"
      type="datetimerange"
      style="flex-grow: unset; margin-left: 8px"
      @change="handleCustomTimeChange"
    />
  </div>
</template>

<style lang="scss" scoped>
.time-range-selector {
  display: flex;
  align-items: center;
}
</style>
