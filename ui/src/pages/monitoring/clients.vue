<script setup lang="ts">
import type { Column } from 'element-plus'
import type { Client } from '@/api/monitoring'
import { ElCheckbox, ElMessage, ElMessageBox } from 'element-plus'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import monitoringService from '@/api/monitoring'

const router = useRouter()

const { t } = useI18n()

const queryParams = ref<{
  clientID?: string
  username?: string
  remote?: string
}>({
  clientID: '',
  username: '',
  remote: '',
})

const isLoading = ref(false)
const tableData = ref<Client[]>([])
const selectedKeys = ref<Set<string>>(new Set())

const isAllSelected = computed(() => {
  return (
    tableData.value.length > 0
    && selectedKeys.value.size === tableData.value.length
  )
})

const isIndeterminate = computed(() => {
  return (
    selectedKeys.value.size > 0
    && selectedKeys.value.size < tableData.value.length
  )
})

function toggleSelection(clientID: string) {
  if (selectedKeys.value.has(clientID)) {
    selectedKeys.value.delete(clientID)
  } else {
    selectedKeys.value.add(clientID)
  }
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedKeys.value.clear()
  } else {
    selectedKeys.value = new Set(tableData.value.map(row => row.clientID))
  }
}

async function loadData() {
  isLoading.value = true
  try {
    const res = await monitoringService.getClients(queryParams.value)
    tableData.value = res || []
    // 清理无效的选中项
    const validClientIDs = new Set(tableData.value.map(row => row.clientID))
    selectedKeys.value = new Set(
      Array.from(selectedKeys.value).filter(id => validClientIDs.has(id)),
    )
  } catch (error) {
    console.error(error)
  } finally {
    isLoading.value = false
  }
}

function formatProtocolVersion(version: number | undefined): string {
  if (version === undefined || version === null) {
    return '-'
  }
  // MQTT 协议版本号映射
  const versionMap: Record<number, string> = {
    3: 'MQTT 3.1',
    4: 'MQTT 3.1.1',
    5: 'MQTT 5.0',
  }
  return versionMap[version] || `MQTT ${version}`
}

function columns(tableWidth: number): Column<Client>[] {
  const checkboxWidth = 50
  const colWidth = Math.floor((tableWidth - checkboxWidth) / 7)
  return [
    {
      key: 'checkbox',
      dataKey: 'clientID',
      title: '',
      width: checkboxWidth,
      headerCellRenderer: () => {
        return h(
          'div',
          {
            style: {
              padding: '8px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: '100%',
              height: '100%',
            },
          },
          [
            h(ElCheckbox, {
              modelValue: isAllSelected.value,
              indeterminate: isIndeterminate.value,
              onChange: toggleSelectAll,
            }),
          ],
        )
      },
      cellRenderer: ({ rowData }: { rowData: Client }) => {
        return h(
          'div',
          {
            style: {
              padding: '8px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: '100%',
              height: '100%',
            },
          },
          [
            h(ElCheckbox, {
              modelValue: selectedKeys.value.has(rowData.clientID),
              onChange: () => toggleSelection(rowData.clientID),
              onClick: (e: Event) => {
                e.stopPropagation()
              },
            }),
          ],
        )
      },
    },
    {
      key: 'clientID',
      dataKey: 'clientID',
      title: t('clientID'),
      width: colWidth,
      cellRenderer: ({ rowData }: { rowData: Client }) => {
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
              router.push({
                name: 'ClientDetail',
                params: { clientID: rowData.clientID },
              })
            },
          },
          rowData.clientID || '-',
        )
      },
    },
    {
      key: 'username',
      dataKey: 'username',
      title: t('username'),
      width: colWidth,
      cellRenderer: ({ rowData }: { rowData: Client }) => {
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
          rowData.username || '-',
        )
      },
    },
    {
      key: 'remote',
      dataKey: 'remote',
      title: t('clientIP'),
      width: colWidth,
      cellRenderer: ({ rowData }: { rowData: Client }) => {
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
          rowData.remote || '-',
        )
      },
    },
    {
      key: 'keepalive',
      dataKey: 'keepalive',
      title: t('keepalive'),
      width: colWidth,
      cellRenderer: ({ rowData }: { rowData: Client }) => {
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
          String(rowData.keepalive ?? '-'),
        )
      },
    },
    {
      key: 'clean',
      dataKey: 'clean',
      title: t('cleanStart'),
      width: colWidth,
      cellRenderer: ({ rowData }: { rowData: Client }) => {
        const cleanText = String(rowData.clean)
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
          cleanText,
        )
      },
    },
    {
      key: 'sessionExpiryInterval',
      dataKey: 'sessionExpiryInterval',
      title: t('sessionExpiryInterval'),
      width: colWidth,
      cellRenderer: ({ rowData }: { rowData: Client }) => {
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
          String(rowData.sessionExpiryInterval),
        )
      },
    },
    {
      key: 'protocolVersion',
      dataKey: 'protocolVersion',
      title: t('protocolVersion'),
      width: colWidth,
      cellRenderer: ({ rowData }: { rowData: Client }) => {
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
          formatProtocolVersion(rowData.protocolVersion),
        )
      },
    },
  ]
}

function closeSelection() {
  ElMessageBox.confirm(t('close_selection_tip'), {
    confirmButtonText: t('common.ok'),
    cancelButtonText: t('common.cancel'),
    type: 'warning',
  })
    .then(() => {
      Promise.all(
        Array.from(selectedKeys.value).map(clientID =>
          monitoringService.closeClient(clientID),
        ),
      )
        .then(() => {
          ElMessage.success(t('common.success.operate'))
          selectedKeys.value.clear()
          loadData()
        })
        .catch((error) => {
          console.error(error)
          // ElMessage.error(t('common.error.operate'))
        })
    })
    .catch(() => {})
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
          v-model="queryParams.clientID!"
          :placeholder="t('clientID')"
          style="width: 300px"
          @enter="loadData"
        />
        <search-input
          v-model="queryParams.username!"
          :placeholder="t('username')"
          @enter="loadData"
        />
        <search-input
          v-model="queryParams.remote!"
          :placeholder="t('clientIP')"
          @enter="loadData"
        />
        <search-button @click="loadData" />
      </search>
    </template>
    <template #headerRight>
      <el-tooltip :content="t('close_selection')" placement="top">
        <el-button
          :disabled="selectedKeys.size === 0"
          style="padding: 8px"
          @click="closeSelection"
        >
          <template #icon>
            <el-icon>
              <i-mdi-link-variant-remove />
            </el-icon>
          </template>
        </el-button>
      </el-tooltip>
    </template>
    <div v-loading="isLoading" relative h-full w-full>
      <el-auto-resizer class="absolute inset-0">
        <template #default="{ height, width }">
          <el-table-v2
            :columns="columns(width)"
            :data="tableData"
            :width="width"
            :height="height"
            :header-height="40"
            :row-height="48"
            fixed
            row-key="clientID"
          />
        </template>
      </el-auto-resizer>
    </div>
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
  clientID: 客户端 ID
  username: 用户名
  clientIP: 客户端 IP
  keepalive: Keepalive
  cleanStart: Clean Start
  sessionExpiryInterval: 会话过期时间（秒）
  protocolVersion: 协议版本
  close_selection: 关闭选中
  close_selection_tip: 您确定要关闭选中的客户端吗？
en:
  clientID: Client ID
  username: Username
  clientIP: Client IP
  keepalive: Keepalive
  cleanStart: Clean Start
  sessionExpiryInterval: Session Expiry Interval (seconds)
  protocolVersion: Protocol Version
  close_selection: Close Selected
  close_selection_tip: Are you sure you want to close the selected clients?
es:
  clientID: Cliente ID
  username: Usuario
  clientIP: Cliente IP
  keepalive: Keepalive
  cleanStart: Clean Start
  sessionExpiryInterval: Intervalo de Expiración de Sesión (segundos)
  protocolVersion: Versión del Protocolo
  close_selection: Cerrar Seleccionados
  close_selection_tip: ¿Estás seguro de querer cerrar los clientes seleccionados?
</i18n>
