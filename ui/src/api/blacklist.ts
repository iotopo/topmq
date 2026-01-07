import { BaseApi } from './base'

export interface Blacklist {
  id: string
  username: string
  clientID: string
  remote: string
  expiredAt: string
  description: string
  createdAt: string
  updatedAt: string
}

class BlacklistService extends BaseApi {
  constructor() {
    super('/api/v1')
  }

  getAll(params: {
    pageNum: number
    pageSize: number
    search?: string
  }): Promise<{
    total: number
    items: Blacklist[]
  }> {
    return this.http.get('black_list', { searchParams: params }).json()
  }

  delete(id: string): Promise<void> {
    return this.http.delete(`black_list/${id}`).json()
  }

  create(data: {
    username: string
    clientID: string
    remote: string
    expiredAt: Date | string
    description: string
  }): Promise<void> {
    return this.http.post('black_list', { json: data }).json()
  }

  update(id: string, data: {
    username: string
    clientID: string
    remote: string
    expiredAt: Date | string
    description: string
  }): Promise<void> {
    return this.http.put(`black_list/${id}`, { json: data }).json()
  }
}

export default new BlacklistService()
