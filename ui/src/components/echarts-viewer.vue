<script lang="ts" setup>
import { useDark } from '@vueuse/core'
import * as echarts from 'echarts'
import { onMounted, onUnmounted, ref, watch } from 'vue'

const props = defineProps<{
  option?: any
  theme?: string
}>()

const emit = defineEmits(['init'])

const chartRef = ref<HTMLElement>()
let chart: any = null

const isDark = useDark()
async function initChart() {
  if (!chartRef.value)
    return

  try {
    const theme = props.theme || (isDark.value ? 'dark' : 'walden')
    // TODO: 需要支持多语言, 支持主题切换
    chart = echarts.init(chartRef.value, theme, {
      locale: 'ZH',
    })
    chart.setOption(props.option || {})
    emit('init', chart)
  } catch (error) {
    console.error('Failed to initialize chart:', error)
  }
}

function resizeChart(opts: any) {
  chart?.resize(opts)
}

watch(
  () => props.option,
  (newOption) => {
    if (!chart) {
      return
    }
    chart.setOption(newOption, true, true)
    // chart?.setOption(newOption, true, true)
  },
  // { deep: true }
)

onMounted(async () => {
  await initChart()
  window.addEventListener('resize', resizeChart)
})

onUnmounted(() => {
  window.removeEventListener('resize', resizeChart)
  chart?.dispose()
})

// 暴露方法给父组件
defineExpose({
  getChart: () => chart,
  mergeOption: (newOption: any) => {
    chart?.setOption(newOption)
  },
  setOption: (...args: any[]) => {
    // eslint-disable-next-line prefer-spread
    chart?.setOption.apply(chart, args)
  },
  resize: resizeChart,
  clear: () => {
    chart?.clear()
  },
})
</script>

<template>
  <div ref="chartRef" />
</template>
