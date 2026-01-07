import qs from 'qs'
import { BaseApi } from './base'

interface QueryRequest {
  metric?: string
  start: Date | string
  end: Date | string
  ag?: string
  group?: string
  host?: boolean
}

class MetricsService extends BaseApi {
  constructor() {
    super('/api/v1/metrics')
  }

  getOverview() {
    return this.http.get(`overview`).json<{
      cpu_percent: number
      mem_percent: number
      disk_percent: number
      net_in_rate: number
      disks_usage: Array<{
        name: string
        usagePercent: number
      }>
    }>()
  }

  query<T>(query: QueryRequest) {
    const queryString = qs.stringify(query, { arrayFormat: 'repeat' })
    return this.http.get(`query`, { searchParams: queryString }).json<Array<T>>()
  }

  queryByTag<T>(query: QueryRequest) {
    const searchParams = qs.stringify(query, { arrayFormat: 'repeat' })
    return this.http.get(`query`, { searchParams }).json<Array<{
      name: string
      items: Array<T>
    }>>()
  }
}

export const metricsService = new MetricsService()
