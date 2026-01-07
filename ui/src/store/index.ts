import dayjs from 'dayjs'
import { createPinia, defineStore } from 'pinia'
import authService from '@/api/auth'
import i18n from '@/i18n'

export const pinia = (window as any).$pinia || createPinia()
;(window as any).$pinia = pinia

export interface User {
  id: string
  isSuperuser: boolean
  name: string
  username: string
}
export interface State {
  // isDark?: boolean
  locale: string // 'zh' | 'en'
  user: User | null
  metrics: boolean
  minPwdLen: number
  version: string
}
export const useStore = defineStore('common', {
  // other options...
  // a function that returns a fresh state
  state: (): State => ({
    locale: 'zh',
    user: null,
    metrics: false,
    minPwdLen: 8,
    version: '1.0.0',
  }),

  getters: {
    isLogin(): boolean {
      return !!this.user
    },
  },
  actions: {
    initLocale() {
      let locale = localStorage.getItem('locale')
      if (!locale) {
        locale = navigator.language.slice(0, 2)
        if (
          locale !== 'en'
          && locale !== 'zh'
        ) {
          locale = 'zh'
        }
      }
      this.setLocale(locale)
    },

    setLocale(locale: string) {
      ;(i18n.global.locale.value as string) = locale
      localStorage.setItem('locale', locale)
      if (locale === 'zh') {
        dayjs.locale('zh')
      } else if (locale === 'en') {
        dayjs.locale(locale)
      }
      this.locale = locale
    },
    setUser(user: User | null) {
      this.user = user
    },

    /**
     * 认证登录与权限 成功后将更新 user 和 permissions
     * @return 如果获取用户失败，则返回 null，反之是 user
     */
    async authenticate() {
      const user = await authService.authenticate()
      if (user) {
        this.setUser(user as User)
        return user
      } else {
        return null
      }
    },

    /**
     * 注销登录
     * 清理用户信息、清理权限、清理路由
     */
    async logout() {
      this.setUser(null)
      await authService.logout()
      localStorage.removeItem('token')
    },
  },
})

export function getActiveStore() {
  return useStore(pinia)
}
