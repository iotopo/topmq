<script setup lang="ts">
import type { OverviewMetrics } from '@/api/monitoring'
import dayjs from 'dayjs'
import { useI18n } from 'vue-i18n'
import monitoringService from '@/api/monitoring'
import { useStore } from '@/store'
import BytesChart from './metrics/bytes-chart.vue'
import ClientsChart from './metrics/clients-chart.vue'
import TimeRangeSelector from './metrics/components/time-range-selector.vue'
import MessagesChart from './metrics/messages-chart.vue'
import PacketsChart from './metrics/packets-chart.vue'
import SubscriptionsChart from './metrics/subscriptions-chart.vue'

const store = useStore()
const { t } = useI18n()

const overview = ref<OverviewMetrics>({
  clientsConnected: 0,
  clientsTotal: 0,
  retained: 0,
  subscriptions: 0,
  messagesSent: 0,
  messagesReceived: 0,
  messagesDropped: 0,
})

let intervalId: number | undefined
function loadOverviewData() {
  monitoringService.getOverviewMetrics().then((res) => {
    overview.value = res
  })
}

const customTime = ref<[Date, Date]>([
  dayjs().subtract(1, 'hour').toDate(),
  dayjs().toDate(),
])
const itemsPerRow = ref(3)
provide('customTime', customTime)

const gridStyle = computed(() => ({
  display: 'grid',
  gap: '16px',
  gridTemplateColumns: `repeat(${itemsPerRow.value}, 1fr)`,
}))

// const groupTitleStyle = computed(() => ({
//   gridColumn: `span ${itemsPerRow.value}`,
// }))

const interval = computed(() => {
  const duration = dayjs(customTime.value[1]).diff(
    dayjs(customTime.value[0]),
    'hour',
    false,
  )
  if (duration >= 14 * 24) {
    return t('intervals.120min')
  } else if (duration >= 7 * 24) {
    return t('intervals.60min')
  } else if (duration >= 3 * 24) {
    return t('intervals.30min')
  } else if (duration >= 24) {
    return t('intervals.10min')
  } else if (duration >= 12) {
    return t('intervals.5min')
  } else if (duration >= 6) {
    return t('intervals.1min')
  } else if (duration >= 3) {
    return t('intervals.1min')
  } else {
    return t('intervals.15sec')
  }
})
provide('interval', interval)

const bytesChartRef = ref()
const packetsChartRef = ref()
const messagesChartRef = ref()
const clientsChartRef = ref()
const subscriptionsChartRef = ref()
function resizeCharts() {
  nextTick(() => {
    bytesChartRef.value?.resize()
    packetsChartRef.value?.resize()
    messagesChartRef.value?.resize()
    clientsChartRef.value?.resize()
    subscriptionsChartRef.value?.resize()
  })
}

function reloadCharts() {
  nextTick(() => {
    bytesChartRef.value?.reload()
    packetsChartRef.value?.reload()
    messagesChartRef.value?.reload()
    clientsChartRef.value?.reload()
    subscriptionsChartRef.value?.reload()
  })
}

onMounted(() => {
  loadOverviewData()
  intervalId = window.setInterval(loadOverviewData, 15000)
})

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId)
  }
})
</script>

<template>
  <page2>
    <div
      w-full
      style="
        margin: 16px 0;
        padding: 16px;
        background-color: var(--el-bg-color);
      "
    >
      <el-alert v-if="!store.metrics" :title="t('notEnabled')" type="error" style="max-width: 1200px; margin-bottom: 20px" />
      <div class="overview-grid">
        <div class="metric-item">
          <div class="metric-label">
            {{ t('clientsConnected') }}
          </div>
          <div class="metric-value">
            {{ overview.clientsConnected }}/{{ overview.clientsTotal }}
          </div>
        </div>
        <div class="metric-item">
          <div class="metric-label">
            {{ t('subscriptions') }}
          </div>
          <div class="metric-value">
            {{ overview.subscriptions }}
          </div>
        </div>
        <div class="metric-item">
          <div class="metric-label">
            {{ t('retained') }}
          </div>
          <div class="metric-value">
            {{ overview.retained }}
          </div>
        </div>
        <div class="metric-item">
          <div class="metric-label">
            <span>{{ t('totalTPS') }}</span>
            <el-tooltip placement="right" effect="light">
              <template #content>
                <div style="width: 180px" flex items-center justify-between>
                  <span>{{ t('send') }} </span>
                  <span>{{ overview.messagesSent }}&nbsp;/s</span>
                </div>
                <div style="width: 180px" flex items-center justify-between>
                  <span>{{ t('receive') }} </span>
                  <span>{{ overview.messagesReceived }}&nbsp;/s</span>
                </div>
                <div style="width: 180px" flex items-center justify-between>
                  <span>{{ t('dropped') }} </span>
                  <span>{{ overview.messagesDropped }}&nbsp;/s</span>
                </div>
              </template>
              <el-icon class="metric-icon">
                <i-mdi-information-outline />
              </el-icon>
            </el-tooltip>
          </div>
          <div class="metric-value">
            {{ overview.messagesSent + overview.messagesReceived }}&nbsp;/s
          </div>
        </div>
        <!-- <div class="metric-item">
      <div class="metric-label">
        <span>网络流量</span>
        <el-tooltip placement="right" effect="light" content="网络流入速率">
          <el-icon class="metric-icon"><i-mdi-information-outline /></el-icon>
        </el-tooltip>
      </div>
      <div class="metric-value">{{ netTraffic }}&nbsp;KB/s</div>
    </div> -->
      </div>
      <div>
        <div w-full>
          <div class="header-toolbar">
            <TimeRangeSelector v-model="customTime" @change="reloadCharts" />
            <div style="flex-grow: 1" />
            <el-select
              v-model="itemsPerRow"
              style="width: 120px; margin-right: 8px"
              @change="resizeCharts"
            >
              <el-option
                v-for="item in [1, 2, 3]"
                :key="item"
                :label="t('itemsPerRow', { count: item })"
                :value="item"
              />
            </el-select>
            <el-button @click="reloadCharts">
              <el-icon><i-mdi-reload /></el-icon>
            </el-button>
          </div>
        </div>
        <div :style="gridStyle" mt-16px>
          <BytesChart ref="bytesChartRef" />
          <PacketsChart ref="packetsChartRef" />
          <MessagesChart ref="messagesChartRef" />
          <ClientsChart ref="clientsChartRef" />
          <SubscriptionsChart ref="subscriptionsChartRef" />
        </div>
      </div>
    </div>
  </page2>
</template>

<style lang="scss" scoped>
.overview-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 16px;
  background-color: var(--el-fill-color-lighter);
  border-radius: 4px;
  margin-bottom: 16px;
  border: 1px solid var(--el-border-color);
}

.metric-item {
  display: flex;
  flex-direction: column;
}

.metric-label {
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
}

.metric-icon {
  margin-left: 4px;
  cursor: pointer;
}

.metric-value {
  font-size: 24px;
  font-weight: 500;
}

.header-toolbar {
  display: flex;
  align-items: center;
  width: 100%;
}

.group-title {
  font-size: 16px;
  font-weight: 500;
  margin: 16px 0 0;
}
</style>

<i18n lang="yaml">
zh:
  notEnabled: 指标监控功能未开启。
  clientsConnected: 连接数
  subscriptions: 订阅数
  retained: 保留消息数
  totalTPS: 总 TPS
  send: 发送
  receive: 接收
  dropped: 丢弃
  itemsPerRow: 每行{count}个
  intervals:
    '120min': 120分钟
    '60min': 60分钟
    '30min': 30分钟
    '10min': 10分钟
    '5min': 5分钟
    '1min': 1分钟
    '15sec': 15秒
en:
  notEnabled: Metrics monitoring is not enabled.
  clientsConnected: Connected Clients
  subscriptions: Subscriptions
  retained: Retained Messages
  totalTPS: Total TPS
  send: Send
  receive: Receive
  dropped: Dropped
  itemsPerRow: '{count} per row'
  intervals:
    '120min': 120 minutes
    '60min': 60 minutes
    '30min': 30 minutes
    '10min': 10 minutes
    '5min': 5 minutes
    '1min': 1 minute
    '15sec': 15 seconds
es:
  notEnabled: La supervisión de métricas no está habilitada.
  clientsConnected: Clientes Conectados
  subscriptions: Suscripciones
  retained: Mensajes Retenidos
  totalTPS: TPS Total
  send: Enviar
  receive: Recibir
  dropped: Descartados
  itemsPerRow: '{count} por fila'
  intervals:
    '120min': 120 minutos
    '60min': 60 minutos
    '30min': 30 minutos
    '10min': 10 minutos
    '5min': 5 minutos
    '1min': 1 minuto
    '15sec': 15 segundos
</i18n>
