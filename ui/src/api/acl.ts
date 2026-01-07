import type { ResponseData } from './base'
import { BaseApi } from './base'

export interface AccessControl {
  id: string
  username: string // 用户名，空表示不限制
  clientID: string // 客户端ID，空表示不限制
  remote: string // 客户端IP，空表示不限制
  topic: string // 主题，不可为空
  access: 'd' | 'r' | 'w' | 'rw' // d(禁用) r(发布) w(订阅) rw(发布 & 订阅)
  description: string
  createdAt: string
  updatedAt: string
}

class AccessControlService extends BaseApi {
  constructor() {
    super('/api/v1')
  }

  getAccessControlList(params: {
    pageNum: number
    pageSize: number
    search?: string
    type: 'username' | 'clientID' | 'clientIP' | ''
  }): Promise<{
    total: number
    items: AccessControl[]
  }> {
    return this.http.get('acl', { searchParams: params }).json()
  }

  delete(id: string): Promise<void> {
    return this.http.delete(id).json()
  }

  create(accessControl: {
    username: string
    clientID: string
    remote: string
    topic: string
    access: 'd' | 'r' | 'w' | 'rw'
    description: string
  }): Promise<void> {
    return this.http.post('acl', { json: accessControl }).json()
  }

  update(id: string, accessControl: {
    username: string
    clientID: string
    remote: string
    topic: string
    access: 'd' | 'r' | 'w' | 'rw'
    description: string
  }): Promise<void> {
    return this.http.put(`acl/${id}`, { json: accessControl }).json()
  }
}

export default new AccessControlService()
