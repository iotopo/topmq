package metrics

import "time"

type MetricData struct {
	Name   string
	Fields map[string]any
	Tags   map[string]string
	Time   time.Time
}
type MetricDatabase interface {
	Query(metricName string, ag, interval string, tags map[string]string, start, end time.Time) (rows []map[string]any, err error)
	QueryByGroup(metricName string, ag, group, interval string, tags map[string]string, start, end time.Time) (seriesData []Series, err error)
	Write(data []MetricData) error
	Close()
}

type Series struct {
	Name  string           `json:"name,omitempty"`
	Items []map[string]any `json:"items,omitempty"`
}

type DiskUsage struct {
	Name         string  `json:"name"`
	UsagePercent float64 `json:"usagePercent"`
}
