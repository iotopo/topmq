<script lang="ts" setup>
import dayjs from 'dayjs'
import Split from 'split.js'
import {
  getCurrentInstance,
  onBeforeUnmount,
  onMounted,
  ref,
  useSlots,
  watch,
} from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const props = withDefaults(
  defineProps<{
    modalClass?: string
    padding?: number
    margin?: number
    background?: string
    disableMenuDrag?: boolean
  }>(),
  {
    modalClass: '',
    padding: 16,
    margin: 16,
    background: '',
  },
)
// export default defineComponent({
//   name: 'Page',
//   setup() {
//     return {
//       show,
//       menuShow,
//       year,
//       showClick,
//       locale,
//       t,
//     }
//   },
// })

const emits = defineEmits(['menuVisibleChange', 'menuWidthChange'])
const year = ref(dayjs().year())

const pageMenuRef = ref()
const pageContentRef = ref()

const show = ref(false)
const menuShow = ref(false)
const menuAnimate = ref(false)
function showClick() {
  show.value = !show.value
  menuAnimate.value = true

  if (!show.value) {
    setTimeout(() => {
      menuShow.value = show.value
      const myEvent = new Event('resize')
      window.dispatchEvent(myEvent)
      emits('menuVisibleChange', show.value)
    }, 200)
  }
  else {
    menuShow.value = show.value
    emits('menuVisibleChange', show.value)
  }
  setTimeout(() => {
    menuAnimate.value = false
    const myEvent = new Event('resize')
    window.dispatchEvent(myEvent)
  }, 300)
}

let splitInstance: any
const slots = useSlots()
function initMenuDraggable() {
  if (slots.menu && !props.disableMenuDrag) {
    const sizes = updateSplitInstanceConfig() || [10, 90]
    splitInstance = Split([pageMenuRef.value, pageContentRef.value], {
      gutterSize: 6,
      minSize: 220,
      maxSize: [450, 100000],
      sizes,
      onDragEnd: () => {
        const myEvent = new Event('resize')
        window.dispatchEvent(myEvent)
      },
    })

    window.addEventListener('beforeunload', () => {
      updateSplitInstanceConfig()
    })
  }
}
let resizeObserverLeft
let usedMenuWidth = false
const leftWidth = ref(0)

watch(leftWidth, (newWidth) => {
  if (usedMenuWidth) {
    emits('menuWidthChange', newWidth)
  }
})

onMounted(() => {
  initMenuDraggable()
  const instance = getCurrentInstance()
  usedMenuWidth = !!instance?.vnode.props?.onMenuWidthChange
  if (usedMenuWidth) {
    resizeObserverLeft = new ResizeObserver((entries) => {
      for (const entry of entries) {
        leftWidth.value = Math.round(entry.contentRect.width)
      }
    })
    if (pageMenuRef.value)
      resizeObserverLeft.observe(pageMenuRef.value)
  }
})

const route = useRoute()
const currPath = route.path
function updateSplitInstanceConfig() {
  if (slots.menu && !props.disableMenuDrag) {
    const splitData = localStorage.getItem('splitInstanceConfig')
    const config = JSON.parse(splitData || '{}')
    if (!config.sizes) {
      config.sizes = {}
    }
    if (splitInstance) {
      const sizes = splitInstance.getSizes()
      config.sizes[currPath] = sizes
      localStorage.setItem('splitInstanceConfig', JSON.stringify(config))
      // console.log('upd', currPath, config.sizes[currPath])
    }
    else {
      // console.log('get', currPath, config.sizes[currPath])
    }
    return config.sizes[currPath] as [number, number] | undefined
  }
}

onBeforeUnmount(() => {
  updateSplitInstanceConfig()
})
</script>

<template>
  <div class="page-container">
    <div v-if="$slots.navbar" class="page-navbar">
      <slot name="navbar" />
    </div>
    <div class="page-main" :class="{ 'hidden-gutter': menuShow }">
      <div
        v-if="$slots.menu"
        ref="pageMenuRef"
        class="page-menu split-0"
        overflow-hidden
        :class="{
          'menu-hide': show,
          'menu-animate': menuAnimate,
          'menu-minw': !menuAnimate,
        }"
        :style="{
          padding: `${margin}px 0 ${margin}px ${margin}px`,
        }"
      >
        <div v-show="!menuShow" class="menu-scroll">
          <slot name="menu" />
        </div>
      </div>
      <div ref="pageContentRef" class="split-1 flex flex-1 overflow-hidden">
        <div v-if="$slots.menu" class="menu-btn" @click="showClick">
          <el-button v-if="!show" link style="padding: 0">
            <el-icon style="vertical-align: middle" size="16">
              <i-ep-caret-left />
            </el-icon>
          </el-button>
          <el-button v-if="show" link style="padding: 0">
            <el-icon style="vertical-align: middle" size="16">
              <i-ep-caret-right />
            </el-icon>
          </el-button>
        </div>
        <div
          class="page"
          :style="{
            margin: $slots.menu
              ? `${margin}px ${margin}px ${margin}px 0`
              : `${margin}px`,
            padding: `${padding}px`,
            background: background ? background : 'var(--el-bg-color)',
          }"
        >
          <div
            v-if="$slots.header || $slots.headerLeft || $slots.headerRight"
            class="header"
          >
            <div class="header-left" flex flex-items-center>
              <slot name="headerLeft" />
            </div>
            <div class="header-content" flex flex-items-center>
              <slot name="header" />
            </div>
            <div class="header-right" flex flex-items-center>
              <slot name="headerRight" />
            </div>
          </div>
          <div class="main-scroll">
            <div class="main">
              <slot>
                <el-skeleton />
              </slot>
            </div>
          </div>
          <div
            v-if="$slots.footer || $slots.footerLeft || $slots.footerRight"
            class="footer"
          >
            <div
              :class="!$slots.footerLeft ? 'copyright' : ''"
              class="footer-left"
            >
              <slot name="footerLeft">
                <div v-if="false">
                  Copyright © 2018 - {{ year }} Iotopo. All Rights Reserved.
                  图扑物联 版权所有
                </div>
              </slot>
            </div>
            <div class="footer-content">
              <slot name="footer" />
            </div>
            <div class="footer-right">
              <slot name="footerRight" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.page-container {
  position: relative;
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;

  .page-navbar {
    margin: 16px 16px 0 16px;
    // background: var(--el-bg-color);
    // padding: 16px;
  }

  .hidden-gutter {
    > .gutter {
      display: none !important;
    }
  }

  .page-main {
    position: relative;
    display: flex;
    flex: 1;
    width: 100%;
    height: 100%;
    overflow: hidden;

    > .gutter {
      margin: 16px 0;
      background-color: transparent;
      &:hover {
        background-color: var(--el-color-primary-light-7) !important;
      }
      // border: 1px solid red;
    }
  }

  .page-menu {
    height: 100%;
    position: relative;
    .menu-scroll {
      background: var(--el-bg-color);
      width: 100%;
      height: 100%;
    }
  }

  .menu-btn {
    display: flex;
    height: 100%;
    justify-content: center;
    align-items: center;
    z-index: 1;
    cursor: pointer;
  }

  .menu-hide {
    min-width: 0 !important;
    width: 0 !important;
    transition: width 0.2s ease-in;
    padding-left: 0 !important;
  }
  .menu-animate {
    min-width: 0 !important;
    transition: width 0.2s ease-in;
  }
  .menu-minw {
    min-width: 220px;
  }
}
.page {
  flex: 1;
  position: relative;
  display: flex;
  flex-direction: column;
  font-size: 14px; //页面内普通自己大小
  padding: 16px;
  // margin: 16px 16px 16px 0;
  background: var(--el-bg-color);
  overflow: hidden;

  .header {
    margin-bottom: 16px;
    display: flex;

    .header-left {
      text-align: left;
    }

    .header-content {
      flex: 1;
      text-align: left;
    }

    .header-right {
      text-align: right;
    }
  }

  .main-scroll {
    // background: #fff;
    position: relative;
    flex: 1;
    // box-shadow: 0px 0px 3px rgba(0, 0, 0, 0.13);
    // border-radius: 3px;
    overflow: hidden;

    .main {
      position: absolute;
      width: 100%;
      height: 100%;
      overflow-y: auto;
      // padding: 16px;
    }
  }

  .footer {
    display: flex;
    margin-top: 16px;
    align-items: flex-end;

    .footer-left {
      text-align: left;
    }

    .footer-content {
      flex: 1;
      text-align: left;
    }

    .footer-right {
      text-align: right;
    }
  }

  .copyright {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipse;
    font-size: 14px;
    color: #9b9b9b;
    width: auto;
  }
}
</style>
