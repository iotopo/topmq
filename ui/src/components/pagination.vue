<script lang="ts" setup>
import { toRefs } from 'vue'

interface Props {
  total: number
  pageNum: number
  pageSize: number
  layout?: string
  teleported?: boolean
  size?: '' | 'default' | 'small' | 'large' | undefined
}

const props = withDefaults(defineProps<Props>(), {
  size: undefined,
  total: 0,
  page: 1,
  pageSize: 10,
  layout: 'total, sizes, prev, pager, next, jumper',
})
const emit = defineEmits(['update:pageNum', 'update:pageSize', 'pageChange'])
const { total, pageNum, pageSize } = toRefs(props)

function pageSizeChange(val: number) {
  emit('update:pageSize', val)
  emit('pageChange', {
    page: pageNum.value,
    pageSize: val,
  })
}
function pageChange(val: number) {
  emit('update:pageNum', val)
  emit('pageChange', {
    pageNum: val,
    pageSize: pageSize.value,
  })
}
</script>

<template>
  <el-pagination
    v-bind="$attrs"
    class="pagination"
    :current-page="pageNum"
    :page-size="pageSize"
    :page-sizes="[10, 20, 50, 100]"
    :total="total"
    :background="true"
    :layout="layout"
    :teleported="teleported"
    :size="size"
    @size-change="pageSizeChange"
    @current-change="pageChange"
  />
</template>

<style lang="scss" scoped>
.pagination {
  :deep(button) {
    background-color: var(--el-fill-color-blank) !important;
  }
  :deep(.el-pager) {
    li {
      background: var(--el-fill-color-blank);
    }
  }
}
</style>
