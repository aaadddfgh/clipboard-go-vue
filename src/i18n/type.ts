import zhCN from './zh-cn.json'

// Type-define 'en-US' as the master schema for the resource
export type MessageSchema = typeof zhCN
export type NumberSchema = {
    currency: {
      style: 'currency',
      currencyDisplay: 'symbol'
      currency: string
    }
  }