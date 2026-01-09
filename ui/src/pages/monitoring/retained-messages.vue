<script setup lang="ts">
import type { Column } from 'element-plus'
import type { RetainedMessage } from '@/api/monitoring'
import { Delete, Document } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElButton, ElIcon, ElMessage, ElMessageBox } from 'element-plus'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import monitoringService from '@/api/monitoring'

const router = useRouter()
const { t } = useI18n()

const queryParams = ref({
  filter: '',
})

const isLoading = ref(false)
const tableData = ref<RetainedMessage[]>([])

// 为每行添加唯一 key
const tableDataWithKey = computed(() => {
  return tableData.value.map((row, index) => ({
    ...row,
    _rowKey: `${row.topic}-${index}`,
  }))
})

// Payload 对话框相关
const payloadDialogVisible = ref(false)
const currentPayloadTopic = ref('')
const payloadFormat = ref<'json' | 'plaintext' | 'base64' | 'hex'>('plaintext')
const rawPayload = ref('')
const formattedPayload = ref('')

async function loadData() {
  isLoading.value = true
  try {
    const res = await monitoringService.getRetained({
      filter: queryParams.value.filter || undefined,
    })
    tableData.value = res || []
  } catch (error) {
    console.error(error)
  } finally {
    isLoading.value = false
  }
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

async function viewPayload(row: RetainedMessage) {
  currentPayloadTopic.value = row.topic
  payloadFormat.value = 'plaintext'
  payloadDialogVisible.value = true

  try {
    const resp = await monitoringService.getRetainedPayload(row.topic)
    rawPayload.value = resp.payload || ''
    updateFormattedPayload()
  } catch (error) {
    console.error(error)
    ElMessage.error(t('common.error.load'))
    payloadDialogVisible.value = false
  }
}

function updateFormattedPayload() {
  try {
    switch (payloadFormat.value) {
      case 'json': {
        // 尝试将 base64 解码后解析为 JSON
        const decoded = atob(rawPayload.value)
        try {
          const jsonObj = JSON.parse(decoded)
          formattedPayload.value = JSON.stringify(jsonObj, null, 2)
        } catch {
          // 如果不是有效的 JSON，直接显示解码后的文本
          formattedPayload.value = decoded
        }
        break
      }
      case 'plaintext': {
        // 将 base64 解码为文本
        formattedPayload.value = atob(rawPayload.value)
        break
      }
      case 'base64': {
        // 直接显示原始 base64
        formattedPayload.value = rawPayload.value
        break
      }
      case 'hex': {
        // 将 base64 解码后转换为十六进制
        const decoded = atob(rawPayload.value)
        const hexArray = Array.from(decoded).map(char =>
          char.charCodeAt(0).toString(16).padStart(2, '0'),
        )
        formattedPayload.value = hexArray.join(' ').toUpperCase()
        break
      }
    }
  } catch (error) {
    console.error('Error formatting payload:', error)
    formattedPayload.value = rawPayload.value
  }
}

watch(payloadFormat, () => {
  if (rawPayload.value) {
    updateFormattedPayload()
  }
})

function copyPayload() {
  navigator.clipboard.writeText(formattedPayload.value)
  ElMessage.success(t('common.success.copy'))
}

async function deleteRetained(row: RetainedMessage) {
  try {
    await ElMessageBox.confirm(
      t('delete_retained_confirm', { topic: row.topic }),
      t('common.confirm_delete'),
      {
        confirmButtonText: t('common.ok'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
    )
    await monitoringService.deleteRetained(row.topic)
    ElMessage.success(t('global.messages.delete_success'))
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error(t('common.error.operate'))
    }
  }
}

function columns(tableWidth: number): Column<RetainedMessage & { _rowKey: string }>[] {
  const actionWidth = 150
  // 主题列占2倍宽度，其他4列各占1倍，操作列固定100px
  // 总共6个单位：主题2 + QoS1 + 客户端ID1 + 发布时间1 + 过期时间1 = 6个单位
  const colWidth = Math.floor((tableWidth - actionWidth) / 6)
  return [
    {
      key: 'topic',
      dataKey: 'topic',
      title: t('topic'),
      width: colWidth * 2,
      cellRenderer: ({ rowData }: { rowData: RetainedMessage }) => {
        return h(
          'div',
          {
            style: {
              padding: '8px',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
            },
          },
          rowData.topic || '-',
        )
      },
    },
    {
      key: 'qos',
      dataKey: 'qos',
      title: t('qos'),
      width: colWidth,
      cellRenderer: ({ rowData }: { rowData: RetainedMessage }) => {
        return h(
          'div',
          {
            style: {
              padding: '8px',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
            },
          },
          String(rowData.qos ?? '-'),
        )
      },
    },
    {
      key: 'clientID',
      dataKey: 'clientID',
      title: t('clientID'),
      width: colWidth,
      cellRenderer: ({ rowData }: { rowData: RetainedMessage }) => {
        return h(
          'div',
          {
            style: {
              padding: '8px',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
              color: 'var(--el-color-primary)',
              cursor: 'pointer',
            },
            class: 'client-id-link',
            onClick: () => {
              if (rowData.clientID) {
                router.push({
                  name: 'ClientDetail',
                  params: { clientID: rowData.clientID },
                })
              }
            },
          },
          rowData.clientID,
        )
      },
    },
    {
      key: 'createdAt',
      dataKey: 'createdAt',
      title: t('publish_time'),
      width: colWidth,
      cellRenderer: ({ rowData }: { rowData: RetainedMessage }) => {
        return h(
          'div',
          {
            style: {
              padding: '8px',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
            },
          },
          rowData.createdAt ? formatTime(rowData.createdAt) : '-',
        )
      },
    },
    {
      key: 'expiredAt',
      dataKey: 'expiredAt',
      title: t('expired_at'),
      width: colWidth,
      cellRenderer: ({ rowData }: { rowData: RetainedMessage }) => {
        return h(
          'div',
          {
            style: {
              padding: '8px',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
            },
          },
          rowData.expiredAt ? formatTime(rowData.expiredAt) : '-',
        )
      },
    },
    {
      key: 'actions',
      dataKey: 'topic',
      title: t('common.operation'),
      width: 150,
      cellRenderer: ({ rowData }: { rowData: RetainedMessage }) => {
        return h(
          'div',
          {
            style: {
              padding: '8px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            },
          },
          [
            h(
              ElButton,
              {
                link: true,
                size: 'small',
                title: t('view_payload'),
                onClick: () => viewPayload(rowData),
              },
              {
                default: () =>
                  h(
                    ElIcon,
                    {},
                    {
                      default: () => h(Document),
                    },
                  ),
              },
            ),
            h(
              ElButton,
              {
                link: true,
                size: 'small',
                title: t('common.delete'),
                onClick: () => deleteRetained(rowData),
              },
              {
                default: () =>
                  h(
                    ElIcon,
                    {},
                    {
                      default: () => h(Delete),
                    },
                  ),
              },
            ),
          ],
        )
      },
    },
  ]
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <page>
    <template #header>
      <search>
        <search-input
          v-model="queryParams.filter"
          :placeholder="t('topic_wildcard_hint')"
          style="width: 300px"
          @enter="loadData"
        />
        <search-button @click="loadData" />
      </search>
    </template>
    <template #headerRight>
      <!-- <el-tooltip :content="t('common.refresh')" placement="top">
        <el-button style="padding: 8px" @click="loadData">
          <template #icon>
            <el-icon>
              <i-ep-refresh />
            </el-icon>
          </template>
        </el-button>
      </el-tooltip> -->
    </template>
    <div v-loading="isLoading" relative h-full w-full>
      <el-auto-resizer class="absolute inset-0">
        <template #default="{ height, width }">
          <el-table-v2
            :columns="columns(width)"
            :data="tableDataWithKey"
            :width="width"
            :height="height"
            :header-height="40"
            :row-height="48"
            fixed
            row-key="_rowKey"
          />
        </template>
      </el-auto-resizer>
    </div>

    <!-- 查看 Payload 对话框 -->
    <el-dialog
      v-model="payloadDialogVisible"
      :title="t('view_payload')"
      width="600px"
      :close-on-click-modal="false"
    >
      <div style="margin-bottom: 16px">
        <span style="color: var(--el-text-color-secondary); margin-right: 8px">
          {{ t('topic') }}:
        </span>
        <span>{{ currentPayloadTopic }}</span>
      </div>
      <div
        style="
          margin-bottom: 16px;
          display: flex;
          align-items: center;
          gap: 8px;
        "
      >
        <span style="color: var(--el-text-color-secondary)">{{ t('payload_format') }}:</span>
        <el-select v-model="payloadFormat" style="width: 150px">
          <el-option :label="t('plaintext')" value="plaintext" />
          <el-option :label="t('json')" value="json" />
          <el-option :label="t('base64')" value="base64" />
          <el-option :label="t('hex')" value="hex" />
        </el-select>
        <ElButton @click="copyPayload">
          {{ t('copy') }}
        </ElButton>
      </div>
      <div style="border: var(--el-border)">
        <monaco-editor
          v-model="formattedPayload"
          :language="payloadFormat === 'json' ? 'json' : 'text'"
          :options="{ readOnly: true, lineNumbers: 'on' }"
          height="300px"
        />
      </div>
    </el-dialog>
  </page>
</template>

<style lang="scss" scoped>
:deep(.el-table-v2__row-cell) {
  border-bottom: 1px solid var(--el-border-color-lighter);
}

:deep(.el-table-v2__header-row) {
  background-color: var(--el-table-header-bg-color);
}

:deep(.el-table-v2__row:hover) {
  background-color: var(--el-table-row-hover-bg-color);
}

:deep(.client-id-link) {
  &:hover {
    text-decoration: underline;
  }
}
</style>

<i18n lang="yaml">
zh:
  retained_messages: 保留消息
  topic_wildcard_hint: 主题（支持通配符搜索）
  publish_time: 发布时间
  expired_at: 过期时间
  view_payload: 查看 Payload
  payload_format: Payload 格式
  json: JSON
  plaintext: Plaintext
  base64: Base64
  hex: Hex
  copy: 复制
  topic: 主题
  qos: QoS
  clientID: 客户端 ID
  delete_retained_confirm: 确定要删除主题 "{topic}" 的保留消息吗？
en:
  retained_messages: Retained Messages
  topic_wildcard_hint: Topic (wildcard search supported)
  publish_time: Publish Time
  expired_at: Expired At
  view_payload: View Payload
  payload_format: Payload Format
  json: JSON
  plaintext: Plaintext
  base64: Base64
  hex: Hex
  copy: Copy
  topic: Topic
  qos: QoS
  clientID: Client ID
  delete_retained_confirm: Are you sure you want to delete the retained message for topic "{topic}"?
es:
  retained_messages: Mensajes Retenidos
  topic_wildcard_hint: Tema (búsqueda con comodines soportada)
  publish_time: Tiempo de Publicación
  expired_at: Tiempo de Expiración
  view_payload: Ver Payload
  payload_format: Formato de Payload
  json: JSON
  plaintext: Texto Plano
  base64: Base64
  hex: Hexadecimal
  copy: Copiar
  topic: Tema
  qos: QoS
  clientID: Cliente ID
  delete_retained_confirm: ¿Está seguro de que desea eliminar el mensaje retenido para el tema "{topic}"?
</i18n>
