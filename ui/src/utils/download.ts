import * as qs from 'qs'
import clone from 'clone'
import selfService from '../api/self'

interface DownloadOptions {
  url: string
  auth?: boolean // 是否增加临时token
  params?: Record<string, any>
  download?: string
}

export async function openLink(options: DownloadOptions): Promise<void> {
  try {
    options = clone(options)
    if (options.auth) {
      const resp = await selfService.tempToken()
      if (options.params) {
        options.params._csrf_token_ = resp.token
      } else {
        options.params = {
          _csrf_token_: resp.token,
        }
      }
    }
    let url = ''
    if (options.params) {
      url = `${options.url}?${qs.stringify(options.params, { indices: false })}`
    } else {
      url = options.url
    }
    window.open(url)
  } catch (error) {
    error && console.warn(error)
  }
}
