import { BaseApi } from './base'

export interface UserInfo {
  username: string
  name: string
  email: string
  mobile: string
  tel: string
  wechat_mini_bound: boolean
}

class SelfService extends BaseApi {
  constructor() {
    super('/api/v1/self')
  }

  getMyUserGroups() {
    return this.http.get('userGroups').json<{
      id: string
      name: string
      description?: string
    }[]>()
  }

  /**
   * 获取个人信息
   */
  getSelfInfo() {
    return this.http.get('').json<UserInfo>()
  }

  /** 设置个人信息 */
  setSelfInfo(info: {
    name: string
    email?: string
    mobile?: string
    tel?: string
  }) {
    return this.http.put('', { json: info }).json()
  }

  /**
   * 绑定手机号
   * @param data
   * @returns
   */
  bindMobile(data: { mobile: string, id: string, code: string }) {
    return this.http.post('bindMobile', { json: data }).json()
  }

  bindEmail(data: { email: string, id: string, code: string }) {
    return this.http.post('bindEmail', { json: data }).json()
  }

  tempToken() {
    return this.http.post('temp_token').json<{ token: string }>()
  }
}
const selfService = new SelfService()
export default selfService
