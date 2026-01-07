<script setup lang="ts">
import type { Client } from '@/api/monitoring'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import monitoringService from '@/api/monitoring'
import { useStore } from '@/store'
import SubscriptionList from '../components/subscription-list.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const store = useStore()

const clientID = computed(() => route.params.clientID as string)
const isLoading = ref(false)
const activeTab = ref('info')
const clientDetail = ref<Client | null>(null)

async function loadClientDetail() {
  if (!clientID.value)
    return
  isLoading.value = true
  try {
    const res = await monitoringService.getClientDetail(clientID.value)
    clientDetail.value = res
  } catch (error) {
    console.error(error)
    ElMessage.error(t('common.error.load'))
  } finally {
    isLoading.value = false
  }
}

function goBack() {
  router.back()
}

function copyClientID() {
  if (clientDetail.value?.clientID) {
    navigator.clipboard.writeText(clientDetail.value.clientID)
    ElMessage.success(t('common.success.copy'))
  }
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  loadClientDetail()
})
</script>

<template>
  <page>
    <template #header>
      <div style="display: flex; align-items: center; gap: 8px">
        <el-button link @click="goBack">
          <el-icon><i-ep-arrow-left /></el-icon>
        </el-button>
        <span>{{ clientID }} {{ t('client_detail') }}</span>
      </div>
    </template>
    <div v-loading="isLoading" style="height: 100%; overflow-y: auto">
      <el-tabs v-model="activeTab" class="detail-tabs" type="border-card">
        <el-tab-pane :label="t('information')" name="info">
          <div style="display: flex; gap: 16px; padding: 16px">
            <!-- 连接信息 -->
            <el-card shadow="always" class="info-card">
              <template #header>
                <div class="card-header">
                  <span class="title">{{ t('connection_info') }}</span>
                </div>
              </template>
              <div class="card-body">
                <div class="item-container">
                  <span>{{ t('status') }}</span>
                  <span style="display: flex; align-items: center; gap: 8px">
                    <span
                      class="status-dot" :class="[
                        clientDetail?.closed ? 'disconnected' : 'connected',
                      ]"
                    />
                    {{
                      clientDetail?.closed ? t('disconnected') : t('connected')
                    }}
                  </span>
                </div>
                <div class="item-container">
                  <span>{{ t('clientID') }}</span>
                  <span style="display: flex; align-items: center; gap: 8px">
                    {{ clientDetail?.clientID || '-' }}
                    <el-button
                      v-if="clientDetail?.clientID"
                      link
                      size="small"
                      @click="copyClientID"
                    >
                      <el-icon><i-ep-copy-document /></el-icon>
                    </el-button>
                  </span>
                </div>
                <div class="item-container">
                  <span>{{ t('username') }}</span>
                  <span>{{ clientDetail?.username || '-' }}</span>
                </div>
                <div class="item-container">
                  <span>{{ t('protocol') }}</span>
                  <span>{{ clientDetail?.protocolVersion || '-' }}</span>
                </div>
                <div class="item-container">
                  <span>{{ t('clientIP') }}</span>
                  <span>{{ clientDetail?.remote || '-' }}</span>
                </div>
                <div class="item-container">
                  <span>Keepalive</span>
                  <span>{{ clientDetail?.keepalive ?? '-' }}</span>
                </div>
                <div class="item-container">
                  <span>Clean Start</span>
                  <span>{{ String(clientDetail?.clean ?? false) }}</span>
                </div>
                <!-- <div class="item-container">
                  <span>{{ t('connected_at') }}</span>
                  <span>{{ clientDetail?.connectedAt || '-' }}</span>
                </div> -->
                <div v-if="clientDetail?.disconnectedAt" class="item-container">
                  <span>{{ t('disconnected_at') }}</span>
                  <span>{{
                    clientDetail.disconnectedAt
                      ? formatTime(clientDetail.disconnectedAt)
                      : ''
                  }}</span>
                </div>
              </div>
            </el-card>

            <!-- 会话信息 -->
            <el-card shadow="always" class="info-card">
              <template #header>
                <div class="card-header">
                  <span class="title">{{ t('session_info') }}</span>
                </div>
              </template>
              <div class="card-body">
                <div class="item-container">
                  <span>{{ t('session_expiry_interval') }}</span>
                  <span>{{ clientDetail?.sessionExpiryInterval ?? '-' }}</span>
                </div>
                <!-- <div class="item-container">
                  <span>{{ t('session_created_at') }}</span>
                  <span>{{ clientDetail?.sessionCreatedAt || '-' }}</span>
                </div> -->
                <div class="item-container">
                  <span>{{ t('subscriptions') }}</span>
                  <span>{{ clientDetail?.subscriptions ?? '-' }}</span>
                </div>
                <!-- <div class="item-container">
                  <span>{{ t('message_queue') }}</span>
                  <span>{{ clientDetail?.messageQueue || '-' }}</span>
                </div>
                <div class="item-container">
                  <span>{{ t('in_flight_window') }}</span>
                  <span>{{ clientDetail?.inFlightWindow || '-' }}</span>
                </div>
                <div class="item-container">
                  <span>{{ t('qos2_receive_queue') }}</span>
                  <span>{{ clientDetail?.qos2ReceiveQueue || '-' }}</span>
                </div> -->
              </div>
            </el-card>
          </div>
        </el-tab-pane>
        <!-- <el-tab-pane
          v-if="store.extra.metric"
          :label="t('metrics')"
          name="metrics"
        >
          <div
            style="
              padding: 20px;
              text-align: center;
              color: var(--el-text-color-placeholder);
            "
          >
            {{ t('coming_soon') }}
          </div>
        </el-tab-pane> -->
        <el-tab-pane
          :label="t('subscriptions')"
          name="subscriptions"
          :lazy="true"
        >
          <SubscriptionList :client-i-d="clientID" />
        </el-tab-pane>
      </el-tabs>
    </div>
  </page>
</template>

<style lang="scss" scoped>
.detail-tabs {
  height: 100%;
  :deep(.el-tabs__content) {
    padding: 0;
    height: calc(100% - 39px);
    overflow: auto;
    .el-tab-pane {
      height: 100%;
    }
  }
}

.info-card {
  width: 100%;
  --el-card-padding: 0;
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 14px;
    font-weight: 700;
    padding: 18px 20px;
  }
  .card-body {
    padding: 20px;
    .item-container {
      padding: 8px 0;
      display: flex;
      & > :first-child {
        color: var(--el-text-color-secondary);
        margin-right: 20px;
        min-width: 120px;
      }
    }
  }
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
  &.connected {
    background-color: #67c23a;
  }
  &.disconnected {
    background-color: #909399;
  }
}
</style>

<i18n lang="yaml">
zh:
  client_detail: 客户端详情
  connection_info: 连接信息
  session_info: 会话信息
  clientID: 客户端 ID
  username: 用户名
  clientIP: 客户端 IP
  status: 状态
  protocol: 协议
  connected_at: 连接时间
  disconnected_at: 断开连接时间
  session_created_at: 会话创建时间
  session_expiry_interval: 会话过期时间（秒）
  message_queue: 消息队列
  in_flight_window: 飞行窗口
  qos2_receive_queue: QoS 2 报文接收队列
  connected: 已连接
  disconnected: 已断开
  information: 信息
  metrics: 指标
  subscriptions: 订阅
  coming_soon: 即将推出
en:
  client_detail: Client Details
  connection_info: Connection Information
  session_info: Session Information
  clientID: Client ID
  username: Username
  clientIP: Client IP
  status: Status
  protocol: Protocol
  connected_at: Connection Time
  disconnected_at: Disconnection Time
  session_created_at: Session Creation Time
  session_expiry_interval: Session Expiry Interval (seconds)
  message_queue: Message Queue
  in_flight_window: In-flight Window
  qos2_receive_queue: QoS 2 Message Receive Queue
  connected: Connected
  disconnected: Disconnected
  information: Information
  metrics: Metrics
  subscriptions: Subscriptions
  coming_soon: Coming Soon
es:
  client_detail: Detalles del Cliente
  connection_info: Información de Conexión
  session_info: Información de Sesión
  clientID: Cliente ID
  username: Usuario
  clientIP: IP del Cliente
  status: Estado
  protocol: Protocolo
  connected_at: Tiempo de Conexión
  disconnected_at: Tiempo de Desconexión
  session_created_at: Tiempo de Creación de Sesión
  session_expiry_interval: Intervalo de Expiración de Sesión (segundos)
  message_queue: Cola de Mensajes
  in_flight_window: Ventana en Vuelo
  qos2_receive_queue: Cola de Recepción de Mensajes QoS 2
  connected: Conectado
  disconnected: Desconectado
  information: Información
  metrics: Métricas
  subscriptions: Suscripciones
  coming_soon: Próximamente
</i18n>
