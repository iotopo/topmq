<script setup lang="ts">
import { useVModels } from '@vueuse/core'

const props = withDefaults(
  defineProps<{
    modelValue: string
    placeholder?: string | undefined
  }>(),
  {
    placeholder: undefined,
  },
)
const emit = defineEmits(['update:modelValue', 'enter'])
const { modelValue } = useVModels(props, emit)

let timer: any = null
function onChange() {
  if (timer) {
    clearTimeout(timer)
  }
  timer = setTimeout(() => {
    emit('enter', modelValue)
    timer = null
  }, 300)
}
function onEnter() {
  if (timer) {
    clearTimeout(timer)
    timer = null
  }
  emit('enter', modelValue)
}
</script>

<template>
  <el-input
    v-model="modelValue"
    clearable
    :placeholder="placeholder"
    @keypress.enter="onEnter"
    @input="onChange"
  >
    <template v-if="$slots.prepend" #prepend>
      <slot name="prepend" />
    </template>
    <template #prefix>
      <el-icon>
        <i-ep-search />
      </el-icon>
    </template>
    <template v-if="$slots.suffix" #suffix>
      <slot name="suffix" />
    </template>
    <template v-if="$slots.append" #append>
      <slot name="append" />
    </template>
  </el-input>
</template>
