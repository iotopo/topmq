<script lang="ts" setup>
import { ArrowDown } from '@element-plus/icons-vue'
import { ref } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    options?: { label: string, value: string }[]
  }>(),
  {
    options: () => [
      {
        label: '平均值',
        value: 'mean',
      },
      {
        label: '最大值',
        value: 'max',
      },
      {
        label: '最小值',
        value: 'min',
      },
    ],
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
}>()

const ag = ref(props.modelValue || 'mean')
function getCurrentLabel() {
  const currentOption = props.options.find(item => item.value === ag.value)
  return currentOption ? currentOption.label : ''
}

function handleCommand(command: string) {
  ag.value = command
  emit('update:modelValue', command)
  emit('change', command)
}
</script>

<template>
  <el-dropdown @command="handleCommand">
    <el-button link type="primary">
      {{ getCurrentLabel() }}
      <el-icon class="el-icon--right">
        <ArrowDown />
      </el-icon>
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="item in options"
          :key="item.value"
          :command="item.value"
        >
          {{ item.label }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>
