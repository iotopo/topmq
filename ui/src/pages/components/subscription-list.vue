<script setup lang="ts">
import type { Column } from 'element-plus'
import type { Subscription } from '@/api/monitoring'
import { computed, h, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import monitoringService from '@/api/monitoring'

interface Props {
  clientID?: string
}
const props = defineProps<Props>()

const router = useRouter()
const { t } = useI18n()

const topicSearch = ref('')

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
      clientID: props.clientID,
      topic: topicSearch.value,
    })
    tableData.value = res || []
  } catch (error) {
    console.error(error)
  } finally {
    isLoading.value = false
  }
}

function columns(tableWidth: number): Column<Subscription>[] {
  const cols: Column<Subscription>[] = []
  // topic 列占 2 份，其他 4 列各占 1 份，总共 6 份
  const totalParts = 6
  const partWidth = Math.floor(tableWidth / totalParts)
  const colWidth = partWidth

  // 主题列
  cols.push({
    key: 'topic',
    dataKey: 'topic',
    title: t('topic'),
    width: partWidth * 2,
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
  })

  // QoS 列
  cols.push({
    key: 'qos',
    dataKey: 'qos',
    title: t('qos'),
    width: colWidth,
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
  })

  // No Local 列
  cols.push({
    key: 'noLocal',
    dataKey: 'noLocal',
    title: t('no_local'),
    width: colWidth,
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
  })

  // Retain 列
  cols.push({
    key: 'retain',
    dataKey: 'retain',
    title: t('retain'),
    width: colWidth,
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
  })

  // Retain Handling 列
  cols.push({
    key: 'retainHandling',
    dataKey: 'retainHandling',
    title: t('retain_handling'),
    width: colWidth,
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
  })

  return cols
}

// 监听 clientID prop 变化，自动加载数据
watch(
  () => props.clientID,
  () => {
    if (props.clientID) {
      loadData()
    }
  },
  { immediate: true },
)

// 如果没有传入 clientID，在 mounted 时加载数据
onMounted(() => {
  if (!props.clientID) {
    loadData()
  }
})

// 暴露 loadData 方法供父组件调用
defineExpose({
  loadData,
})
</script>

<template>
  <div
    v-loading="isLoading"

    relative h-full w-full flex flex-col overflow-hidden
    style="padding: 16px"
  >
    <div
      style="
        padding-bottom: 16px;
        border-bottom: 1px solid var(--el-border-color-lighter);
      "
    >
      <search>
        <search-input
          v-model="topicSearch"
          :placeholder="t('topic')"
          style="width: 300px"
          @enter="loadData"
        />
        <search-button @click="loadData" />
      </search>
    </div>
    <el-auto-resizer flex-1>
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
</template>

<style lang="scss" scoped>
.el-auto-resizer {
  &.search-mode {
    position: absolute;
    inset: 0;
    top: 73px;
  }
  &.full-mode {
    position: absolute;
    inset: 0;
  }
}

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
