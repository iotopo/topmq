<script lang="ts" setup>
import dayjs from 'dayjs'
import numeral from 'numeral'
import { useI18n } from 'vue-i18n'
import { metricsService } from '@/api/metrics'
import AgSelect from './components/ag-select.vue'
import ChartCard from './components/chart-card.vue'
import ChartDialog from './components/chart-dialog.vue'

const { t } = useI18n()

interface PacketsSeries {
  time: string
  packets_received: number
  packets_sent: number
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
  const res = await metricsService.query<PacketsSeries>({
    metric: 'mqtt_server_packets',
    start: dialogCustomTime.value[0],
    end: dialogCustomTime.value[1],
    ag: dialogAg.value,
  })
  const option = getChartOption(res || [])
  chartDialogRef.value.setOption(option)
}

function getChartOption(res: PacketsSeries[]) {
  return {
    tooltip: {
      trigger: 'axis',
      valueFormatter: (value: number) => numeral(value).format('0.000'),
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
        formatter: (value: number) => numeral(value).format('0'),
      },
      axisLine: {
        show: false,
      },
    },
    series: [
      {
        type: 'line',
        name: `packets_received`,
        showSymbol: false,
        lineStyle: {
          width: 2,
        },
        data: res.map(item => [item.time, item.packets_received]),
      },
      {
        type: 'line',
        name: `packets_sent`,
        showSymbol: false,
        lineStyle: {
          width: 2,
        },
        data: res.map(item => [item.time, item.packets_sent]),
      },
    ],
  }
}

async function loadData() {
  if (!customTime) {
    return
  }

  const res = await metricsService.query<PacketsSeries>({
    metric: 'mqtt_server_packets',
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
  title: I/O 包数(pps)
en:
  title: I/O Packets (pps)
es:
  title: I/O Packets (pps)
</i18n>

<style lang="scss" scoped></style>
