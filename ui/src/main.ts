import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { createApp } from 'vue'

import i18n from '@/i18n'
import { router } from '@/router'
import { api } from './api/base'
import App from './App.vue'
import { getActiveStore, pinia } from './store'
import '@/styles/index.scss'
import 'uno.css'

getActiveStore()

api.get('config').json<{
  metrics: boolean
  minPwdLen: number
  version: string
}>().then((data) => {
  const store = getActiveStore()
  store.metrics = data.metrics
  store.minPwdLen = data.minPwdLen
  store.version = data.version

  const app = createApp(App)
  app.use(i18n)
  app.use(pinia)
  // .use(router)
  app.use(ElementPlus, {
  // zIndex: 3000,
    locale: zhCn,
  })
  app.use(router)
  app.mount('#app')
})
