<script lang="ts">
import dayjs from 'dayjs'
import { defineComponent, ref } from 'vue'
import { useI18n } from 'vue-i18n'

// const year = dayjs().year()
export default defineComponent({
  name: 'Page',
  setup() {
    const show = ref(false)
    const menuShow = ref(false)
    const year = ref(dayjs().year())
    const { locale, t } = useI18n({
      inheritLocale: true,
    })
    function showClick() {
      show.value = !show.value
      if (!show.value) {
        setTimeout(() => {
          menuShow.value = show.value
        }, 200)
      } else {
        menuShow.value = show.value
      }
    }
    return {
      show,
      menuShow,
      year,
      showClick,
      locale,
      t,
    }
  },
})
</script>

<template>
  <div class="page-container">
    <div
      v-if="$slots.menu"
      class="page-menu"
      :class="{ 'menu-hide': show, 'menu-show': !show }"
    >
      <div v-show="!menuShow" class="menu-scroll">
        <slot name="menu" />
      </div>
    </div>
    <div v-if="$slots.menu" class="menu-btn" @click="showClick">
      <el-button v-if="!show" link type="primary">
        <el-icon style="vertical-align: middle">
          <i-ep-caret-left />
        </el-icon>
      </el-button>
      <el-button v-if="show" link type="primary">
        <el-icon style="vertical-align: middle">
          <i-ep-caret-right />
        </el-icon>
      </el-button>
    </div>
    <div class="page">
      <div
        v-if="$slots.header || $slots.headerLeft || $slots.headerRight"
        class="header"
      >
        <div class="header-left">
          <slot name="headerLeft" />
        </div>
        <div class="header-content">
          <slot name="header" />
        </div>
        <div class="header-right">
          <slot name="headerRight" />
        </div>
      </div>
      <!-- <el-scrollbar class="main-scroll">
        <div class="main">
          <slot>
            <el-skeleton />
          </slot>
        </div>
      </el-scrollbar> -->
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
        <div :class="!$slots.footerLeft ? 'copyright' : ''" class="footer-left">
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
</template>

<style lang="scss" scoped>
.page-container {
  position: relative;
  display: flex;
  width: 100%;
  height: 100%;

  .page-menu {
    // background: #fff;
    width: 220px;
    height: 100%;
    border-right: 1px solid var(--el-border-color-lighter);
    position: relative;
    padding: 16px 0 16px 16px;

    .menu-scroll {
      background: var(--el-bg-color);
      // position: absolute;
      width: 100%;
      height: 100%;
    }
  }

  .menu-btn {
    // position: absolute;
    // right: 0px;
    display: flex;
    height: 100%;
    justify-content: center;
    align-items: center;
    z-index: 1;
    cursor: pointer;
  }

  .menu-hide {
    width: 20px;
    display: none;
    transition: width 0.2s ease-in;
  }
  .menu-show {
    width: 220px;
    transition: width 0.2s ease-in;
  }
}
.page {
  flex: 1;
  // width: 100%;
  // height: 100%;
  position: relative;
  display: flex;
  flex-direction: column;
  font-size: 14px; //页面内普通自己大小
  // padding: 16px;

  .header {
    margin: 16px;
    display: flex;

    .header-right {
      text-align: left;
      display: flex;
      align-items: center;
    }

    .header-content {
      flex: 1;
      text-align: left;
      display: flex;
      align-items: center;
    }

    .header-right {
      text-align: right;
      display: flex;
      align-items: center;
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
      padding: 0 16px;
      width: 100%;
      height: 100%;
      // padding: 16px;
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
