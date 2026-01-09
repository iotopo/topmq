import type { ResponseData } from './base'
import { BaseApi } from './base'

export interface Client {
  clientID: string
  username: string
  remote: string // IP 地址
  keepalive: number
  clean: boolean // Clean Start
  subscriptions: number // 订阅数
  disconnectedAt?: string
  sessionExpiryInterval: number // 会话过期时间（秒）
  protocolVersion: number // 协议版本
  closed: boolean // 是否已关闭
}

export interface Subscription {
  clientID: string
  remote: string // IP 地址
  topic: string
  qos: number
  retain: boolean // 发布时状态保留
  noLocal: boolean // 禁止本地转发
  retainHandling: number // 保留消息处理, 显示数字
}

export interface RetainedMessage {
  clientID: string
  topic: string
  qos: number
  createdAt: string
  expiredAt: string
}

export interface OverviewMetrics {
  clientsConnected: number
  clientsTotal: number
  retained: number
  subscriptions: number
  messagesSent: number // 每秒发送消息数
  messagesReceived: number // 每秒接收消息数
  messagesDropped: number // 每秒丢弃消息数
}

class MonitoringService extends BaseApi {
  constructor() {
    super('/api/v1/monitoring')
  }

  getClients(params: {
    clientID?: string
    username?: string
    remote?: string
  }): Promise<Client[]> {
    return this.http.get('clients', { searchParams: params }).json()
  }

  closeClient(clientID: string): Promise<void> {
    return this.http.post(`close_client/${clientID}`).json()
  }

  getSubscriptions(params: {
    clientID?: string
    username?: string
    topic?: string
  }): Promise<Subscription[]> {
    return this.http.get('subscriptions', { searchParams: params }).json()
  }

  getClientDetail(clientID: string): Promise<Client> {
    return this.http.get(`client_details/${clientID}`).json()
  }

  getRetained(params: {
    filter?: string // 请求过滤
  }): Promise<RetainedMessage[]> {
    return this.http.get('retained', { searchParams: params }).json()
  }

  /**
   * 获取保留消息的 payload
   * @param params { topic: string } 主题
   * @returns base64 编码的 payload
   */
  getRetainedPayload(topic: string): Promise<{
    payload: string
  }> {
    return this.http.get(`retained_payload`, {
      searchParams: {
        topic,
      },
    }).json()
  }

  deleteRetained(topic: string): Promise<void> {
    return this.http.delete(`retained`, {
      searchParams: {
        topic,
      },
    }).json()
  }

  getOverviewMetrics(): Promise<OverviewMetrics> {
    return this.http.get('metrics/overview').json()
  }
}

const monitoringService = new MonitoringService()
export default monitoringService
