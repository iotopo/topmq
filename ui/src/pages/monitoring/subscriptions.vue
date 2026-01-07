<script setup lang="ts">
import type { Column } from 'element-plus'
import type { Subscription } from '@/api/monitoring'
import { h } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import monitoringService from '@/api/monitoring'

const router = useRouter()
const { t } = useI18n()

const queryParams = ref<{
  clientID?: string
  topic?: string
}>({
  clientID: '',
  topic: '',
})

const isLoading = ref(false)
const tableData = ref<Subscription[]>([])

// 为每行添加唯一 key
const tableDataWithKey = computed(() => {
  return tableData.value.map((row, index) => ({
    ...row,
    _rowKey: `${row.clientID}-${row.topic}-${index}`,
  }))
})

async function loadData() {
  isLoading.value = true
  try {
    const res = await monitoringService.getSubscriptions({
      clientID: queryParams.value.clientID || undefined,
      topic: queryParams.value.topic || undefined,
    })
    tableData.value = res || []
  } catch (error) {
    console.error(error)
  } finally {
    isLoading.value = false
  }
}

function columns(tableWidth: number): Column<Subscription>[] {
  // 列宽比例：clientID(1) + topic(2) + qos(1) + noLocal(1) + retain(1) + retainHandling(1) = 7份
  const totalRatio = 7
  const baseWidth = Math.floor(tableWidth / totalRatio)
  const topicWidth = baseWidth * 2
  // 计算剩余宽度，分配给最后一列以避免舍入误差
  const usedWidth = baseWidth * 5 + topicWidth
  const lastColWidth = tableWidth - usedWidth + baseWidth

  return [
    {
      key: 'clientID',
      dataKey: 'clientID',
      title: t('clientID'),
      width: baseWidth,
      cellRenderer: ({ rowData }: { rowData: Subscription }) => {
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
      key: 'topic',
      dataKey: 'topic',
      title: t('topic'),
      width: topicWidth,
      cellRenderer: ({ rowData }: { rowData: Subscription }) => {
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
      width: baseWidth,
      cellRenderer: ({ rowData }: { rowData: Subscription }) => {
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
      key: 'noLocal',
      dataKey: 'noLocal',
      title: t('no_local'),
      width: baseWidth,
      cellRenderer: ({ rowData }: { rowData: Subscription }) => {
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
          String(rowData.noLocal ?? false),
        )
      },
    },
    {
      key: 'retain',
      dataKey: 'retain',
      title: t('retain'),
      width: baseWidth,
      cellRenderer: ({ rowData }: { rowData: Subscription }) => {
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
          String(rowData.retain ?? false),
        )
      },
    },
    {
      key: 'retainHandling',
      dataKey: 'retainHandling',
      title: t('retain_handling'),
      width: lastColWidth,
      cellRenderer: ({ rowData }: { rowData: Subscription }) => {
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
          String(rowData.retainHandling ?? '-'),
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
          v-model="queryParams.clientID!"
          :placeholder="t('clientID')"
          style="width: 300px"
          @enter="loadData"
        />
        <search-input
          v-model="queryParams.topic!"
          :placeholder="t('topic')"
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
  topic: 主题
  qos: QoS
  no_local: 禁止本地转发
  retain: 发布时状态保留
  retain_handling: 保留消息处理
en:
  clientID: Client ID
  topic: Topic
  qos: QoS
  no_local: Disable Local Forwarding
  retain: Retain Status on Publish
  retain_handling: Retained Message Handling
es:
  clientID: Cliente ID
  topic: Tema
  qos: QoS
  no_local: Deshabilitar Reenvío Local
  retain: Estado de Retención al Publicar
  retain_handling: Manejo de Mensajes Retenidos
</i18n>
