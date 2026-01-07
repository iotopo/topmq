import type { Composer } from 'vue-i18n'
import { createI18n } from 'vue-i18n'
import en from './locales/en'
import zh from './locales/zh'

const i18n = createI18n<false>({
  locale: localStorage.getItem('locale') || 'zh',
  fallbackLocale: 'zh',
  legacy: false, // 如果要支持compositionAPI，此项必须设置为false;
  // allowComposition: true,
  // inheritLocale: true,
  missingWarn: false, // 禁用缺失翻译键的警告
  fallbackWarn: false, // 禁用回退语言警告
  warnHtmlMessage: false, // 禁用 HTML 消息警告
  globalInjection: true, // 全局注册$t方法
  messages: {
    zh,
    en,
  },
})

export const t = (i18n.global as Composer).t

// vueUseI18n.prototype.
export default i18n
