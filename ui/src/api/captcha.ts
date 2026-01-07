import { BaseApi } from './base'

class CaptchaService extends BaseApi {
  constructor() {
    super('/api/v1/captcha')
  }

  isCaptchaShown(): Promise<{
    show: boolean
  }> {
    return this.http.get('show').json()
  }

  getCaptchaID(): Promise<{
    captchaID: string
  }> {
    return this.http.get('refresh').json()
  }

  verifyCaptcha(captchaID: string, captchaValue: string): Promise<void> {
    return this.http.post('verify', { json: { captchaID, captchaValue } }).json()
  }
}

export default new CaptchaService()
