<script lang="ts" setup>
import dayjs from 'dayjs'
import numeral from 'numeral'
import { useI18n } from 'vue-i18n'
import { metricsService } from '@/api/metrics'
import AgSelect from './components/ag-select.vue'
import ChartCard from './components/chart-card.vue'
import ChartDialog from './components/chart-dialog.vue'

const { t } = useI18n()

interface BytesSeries {
  time: string
  bytes_received: number
  bytes_sent: number
}
const customTime = inject<Ref<[Date, Date]>>('customTime')
const ag = ref('mean')
const chartRef = ref()
const chartDialogRef = ref()

const dialogCustomTime = ref<[Date, Date]>([
  dayjs().subtract(1, 'hour').toDate(),
  dayjs().toDate(),
])
const dialogAg = ref('mean')
provide('dialogCustomTime', dialogCustomTime)
function showDialog() {
  chartDialogRef.value.show()
  // loadDialogChart()
}

async function loadDialogChart() {
  const res = await metricsService.query<BytesSeries>({
    metric: 'mqtt_server_bytes',
    start: dialogCustomTime.value[0],
    end: dialogCustomTime.value[1],
    ag: dialogAg.value,
  })
  const option = getChartOption(res || [])
  chartDialogRef.value.setOption(option)
}

function getChartOption(res: BytesSeries[]) {
  return {
    tooltip: {
      trigger: 'axis',
      valueFormatter: (value: number) => {
        if (value > 1024) {
          return `${numeral(value / 1024).format('0.000')} M`
        }
        return `${numeral(value).format('0.000')} K`
      },
    },
    legend: {
      type: 'scroll',
      left: 20,
      top: 0,
    },
    grid: {
      top: 40,
      left: 50,
      right: 10,
      bottom: 20,
    },
    xAxis: {
      type: 'time',
      boundaryGap: false,
      splitLine: {
        show: false,
      },
    },
    yAxis: {
      type: 'value',
      min: 0,
      axisLabel: {
        formatter: (value: number) => {
          if (value > 1024) {
            return `${numeral(value / 1024).format('0.0')} M`
          }
          return `${numeral(value).format('0.0')} K`
        },
      },
      axisLine: {
        show: false,
      },
    },
    series: [
      {
        type: 'line',
        name: `bytes_received`,
        showSymbol: false,
        lineStyle: {
          width: 2,
        },
        data: res.map(item => [item.time, item.bytes_received / 1024]),
      },
      {
        type: 'line',
        name: `bytes_sent`,
        showSymbol: false,
        lineStyle: {
          width: 2,
        },
        data: res.map(item => [item.time, item.bytes_sent / 1024]),
      },
    ],
  }
}

async function loadData() {
  if (!customTime) {
    return
  }

  const res = await metricsService.query<BytesSeries>({
    metric: 'mqtt_server_bytes',
    start: customTime.value[0],
    end: customTime.value[1],
    ag: ag.value,
  })

  const option = getChartOption(res || [])
  chartRef.value.setOption(option)
}

onMounted(() => {
  nextTick(async () => {
    await loadData()
  })
})

defineExpose({
  resize() {
    chartRef.value?.resize()
  },
  reload() {
    loadData()
  },
})
</script>

<template>
  <ChartCard ref="chartRef" :title="t('title')" @maximize="showDialog">
    <template #tools>
      <AgSelect v-model="ag" @change="loadData" />
    </template>
  </ChartCard>
  <ChartDialog
    ref="chartDialogRef"
    :title="t('title')"
    @reload="loadDialogChart"
  >
    <template #tools>
      <AgSelect v-model="dialogAg" @change="loadDialogChart" />
    </template>
  </ChartDialog>
</template>

<i18n lang="yaml">
zh:
  title: I/O 字节数 (bps)
en:
  title: I/O Bytes (bps)
es:
  title: I/O Bytes (bps)
</i18n>

<style lang="scss" scoped></style>
