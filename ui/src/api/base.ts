import type { BeforeErrorHook, KyInstance } from 'ky'
import ky from 'ky'

export interface ResponseData<T> {
  data: T
  code?: string
  msg?: string
  error?: string
  errArgs: any
}

const beforeErrorHook: BeforeErrorHook = async (error) => {
  const request = error.request
  const response = error.response
  //   if (response.status === 500) {
  //     const data = await response.clone().json()
  //     if (data.error) {
  //       ElMessage.error(data.error)
  //     } else if (data.msg) {
  //       ElMessage.error(data.msg)
  //     } else {
  //       console.warn(data)
  //     }
  //   } else
  if (response.status === 401) {
    if (!request.url.endsWith('/api/v1/auth/authenticate')) {
      // 重定向到登录页
      const redirect = `${location.pathname}${location.search}`
      location.href = `/login?redirect=${encodeURIComponent(redirect)}`
    }
  } else {
    console.warn(response.statusText)
  }
  return error
}

export const api = ky.create({
  prefixUrl: '/api/v1',
  parseJson: (text) => {
    if (!text) {
      return text
    }
    const data = JSON.parse(text)
    if (data.code) {
      return Promise.reject(data)
    } else {
      return data.data
    }
  },
  hooks: {
    beforeRequest: [
      (request) => {
        // Only set default auth header on initial request, not on retries
        // (retries may have refreshed token set by beforeRetry)
        const token = localStorage.getItem('token')
        if (token) {
          request.headers.set('Authorization', token)
        }
      },
    ],
    beforeError: [
      beforeErrorHook!,
    ],
    afterResponse: [
    //   async (request, _options, response) => {
    //   },

      // Or force retry based on response body content
    //   async (request, options, response) => {
    //     if (response.status === 200) {
    //       const data = await response.json<{
    //         code?: string
    //         msg?: string
    //         data?: any
    //       }>()
    //       return new Response(JSON.stringify(data), { status: response.status })
    //     }
    //   },
    ],
  },
})

export class BaseApi {
  protected http: KyInstance
  constructor(baseUrl: string) {
    this.http = api.extend({
      prefixUrl: baseUrl,
    })
  }
}
