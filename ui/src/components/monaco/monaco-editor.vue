<script lang="ts">
import type { StyleValue } from 'vue'
import { useResizeObserver } from '@vueuse/core'
import * as monaco from 'monaco-editor'
import {
  computed,
  defineComponent,
  onMounted,
  onUnmounted,
  ref,
  toRefs,
  watch,
} from 'vue'
import './monaco'

export default defineComponent({
  name: 'MonacoEditor',
  props: {
    width: {
      type: String,
      default: '100%',
    },
    height: {
      type: String,
      default: '100%',
    },
    modelValue: {
      type: String,
      default: '',
    },
    language: {
      type: String,
      default: 'javascript',
    },
    theme: { type: String, default: 'vs' },

    options: {
      type: Object,
    },
  },
  emits: ['update:modelValue', 'editorWillMount', 'editorDidMount', 'change'],
  setup(props, { emit }) {
    const { width, height, modelValue, language, options } = toRefs(props)
    const style = computed(() => {
      return {
        'width': width.value,
        'height': height.value,
        'text-align': 'left',
      } as StyleValue
    })

    const monacoLoading = ref(false)

    let editor: any
    const container = ref<HTMLDivElement | null>()
    let mo: MutationObserver | null = null
    let editorObserver: { stop: () => void } | null = null

    // 清理函数
    const cleanup = () => {
      mo?.disconnect()
      editor?.dispose()
      editorObserver?.stop()
    }

    // 在 setup 开始时注册生命周期钩子
    onUnmounted(cleanup)

    const insertValue = (value: string): string => {
      // 插入值方法
      const editorInstance = editor.getModel()
      if (editorInstance) {
        editorInstance.setValue(value)
        return ''
      }
      else {
        console.error('couldn\'t get editor instance')
        return 'couldn\'t get editor instance'
      }
    }

    const getValue = () => {
      return editor.getValue()
    }

    const getModelMarkers = () => {
      return monaco.editor.getModelMarkers({ resource: editor.getModel()!.uri })
    }

    watch(modelValue, (val) => {
      if (!editor) {
        return
      }
      if (editor.getValue() !== val) {
        editor.setValue(val)
        editor.trigger('anyString', 'editor.action.formatDocument', null)
        editor.setValue(editor.getValue())
      }
    })
    watch(language, (lan) => {
      monaco?.editor.setModelLanguage(editor.getModel()!, lan)
    })
    // watch(theme, (val) => {
    //   monaco.editor.setTheme(val)
    // })
    watch(options, (val) => {
      editor?.updateOptions(val as any)
    })

    onMounted(async () => {
      monacoLoading.value = true
      emit('editorWillMount', monaco)

      const theme = document.documentElement.classList.contains('dark')
        ? 'vs-dark'
        : 'vs'

      mo = new MutationObserver(() => {
        const theme = document.documentElement.classList.contains('dark')
          ? 'vs-dark'
          : 'vs'
        editor.updateOptions({
          theme,
        })
      })
      mo.observe(document.documentElement, {
        attributes: true,
        attributeFilter: ['class'],
      })

      const opt = options.value || {}

      editor = monaco.editor.create(container.value!, {
        value: modelValue.value,
        language: language.value,
        theme,
        minimap: {
          enabled: false,
        },
        lineNumbersMinChars: 3,
        lineNumbers: 'off',
        folding: false,
        ...opt,
      })

      // this.diffEditor && this._setModel(this.value, this.original)
      // @event `change`
      editor.onDidChangeModelContent((event: any) => {
        const value = editor.getValue()
        if (modelValue.value !== value) {
          emit('update:modelValue', value)
          emit('change', value, event)
        }
      })
      editorObserver = useResizeObserver(container, () => {
        editor.layout()
      })

      monacoLoading.value = false
      emit('editorDidMount', editor)
    })
    return {
      container,
      style,
      monacoLoading,
      insertValue,
      getValue,
      getModelMarkers,
    }
  },
})
</script>

<template>
  <div ref="container" v-loading="monacoLoading" class="monaco-editor-vue3" :style="style" />
</template>

<style lang="scss" scoped>
.monaco-editor-vue3 {
  border: var(--el-border);
  // border-radius: var(--el-border-radius-base);
  :hover {
    border-color: var(--el-border-color-hover);
  }
  :focus {
    border-color: var(--el-border-color-focus);
  }
}
</style>
