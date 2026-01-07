import { BaseApi } from './base'

class AuthService extends BaseApi {
  constructor() {
    super('/api/v1')
  }

  authenticate() {
    return this.http.get('authenticate').json<{
      id: string
      name: string
      username: string
      isSuperuser: boolean
    }>()
  }

  login(data: {
    username: string
    password: string
    captchaID?: string
    captchaValue?: string
  }): Promise<{
    id: string
    name: string
    username: string
    isSuperuser: boolean
    token: string
  }> {
    return this.http.post('login', { json: data }).json()
  }

  logout(): Promise<void> {
    return this.http.post('logout').json()
  }

  /**
   * 设置新密码接口
   * @param {object} data - 新密码设置所需的数据
   * @param {string} data.oldPassword - 旧密码
   * @param {string} data.password - 新密码
   */
  resetPassword(data: { oldPassword: string, password: string }) {
    return this.http.put('reset_password', { json: data }).json()
  }
}

export default new AuthService()
