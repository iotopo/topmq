<script lang="ts" setup>
import elEn from 'element-plus/es/locale/lang/en'
import elZh from 'element-plus/es/locale/lang/zh-cn'
// import elZhTW from 'element-plus/es/locale/lang/zh-tw'
import { onMounted, ref, watch } from 'vue'
import i18n from './i18n'
import { useStore } from './store'

const store = useStore()
const locale = ref(elZh)

async function loadLocaleFile(val: string) {
  if (val === 'en') {
    locale.value = elEn
  } else if (val === 'zh') {
    locale.value = elZh
  }
  i18n.global.locale.value = val
}
watch(() => store.locale, loadLocaleFile)
loadLocaleFile(store.locale)

onMounted(() => {
})
</script>

<template>
  <div id="main">
    <el-config-provider :locale="locale">
      <router-view v-slot="{ Component }">
        <keep-alive>
          <component :is="Component" />
        </keep-alive>
      </router-view>
    </el-config-provider>
  </div>
</template>

<style lang="scss">
* {
  box-sizing: border-box;
}

html,
body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

body {
  font-size: 14px;
  font-family:
    'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  --el-border-radius-base: 2px;
}

#app {
  // font-size: 14px;
  width: 100%;
  height: 100%;
}

#main {
  width: 100%;
  height: 100%;
}
</style>
