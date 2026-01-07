import type { ResponseData } from './base'
import { BaseApi } from './base'

export interface Account {
  id: string
  username: string
  password: string
  clientID: string
  remote: string
  disabled: boolean
  description: string
  createdAt: string
  updatedAt: string
}

class AccountService extends BaseApi {
  constructor() {
    super('/api/v1')
  }

  getAccounts(params: {
    pageNum: number
    pageSize: number
    search?: string
  }): Promise<{
    total: number
    items: Account[]
  }> {
    return this.http.get('account', { searchParams: params }).json()
  }

  deleteAccount(id: string): Promise<void> {
    return this.http.delete(`${id}`).json()
  }

  createAccount(account: {
    username: string
    password: string
    clientID: string
    remote: string
    description: string
  }): Promise<void> {
    return this.http.post('account', { json: account }).json()
  }

  updateAccount(id: string, account: {
    username: string
    password: string
    clientID: string
    remote: string
    description: string
  }): Promise<void> {
    return this.http.put(id, { json: account }).json()
  }

  enableAccount(id: string): Promise<void> {
    return this.http.post(`account/${id}/enable`).json()
  }

  disableAccount(id: string): Promise<void> {
    return this.http.post(`account/${id}/disable`).json()
  }
}

export default new AccountService()
