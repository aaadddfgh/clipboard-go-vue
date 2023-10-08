import { createI18n } from 'vue-i18n'
import zhCN from './zh-cn.json'
import en  from "./en.json";

type MessageSchema = typeof zhCN

export const i18n = createI18n<[MessageSchema], 'zh-CN'|'en'>({

  locale: 'en',
  legacy: false,
  fallbackLocale: 'en',
  messages: {
    'zh-CN': zhCN,
    'en': en,
  }
})