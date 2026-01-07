package metrics

import (
	"cmp"
	"context"
	"github.com/iotopo/topmq/config"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"
	"time"
)

var CurrentCpuPercent float64
var CurrentMemPercent float64
var CurrentDiskPercent float64
var CurrentNetRate float64
var CurrentDisksUsage []DiskUsage

var lastNetIOCounters map[string]net.IOCountersStat
var lastIOCounters map[string]disk.IOCountersStat

func metricRun() {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				logrus.WithError(err).Error("metric panic")
			} else {
				logrus.Errorf("metric panic: %v", e)
			}
			debug.PrintStack()
		}
	}()
	var items []MetricData
	now := time.Now()
	period := float64(15)

	// 当前 Goroutines 数量
	goroutines := int64(runtime.NumGoroutine())

	var cpuPercent float64

	// 系统级指标（gopsutil）
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		logrus.WithError(err).Error("get process error")
	} else {
		items = append(items, MetricData{
			Name: "process",
			Fields: map[string]any{
				"goroutines": goroutines,
			},
			Time: now,
		})

		// 进程 CPU 使用率
		cpuPercent, err = p.CPUPercent()
		if err != nil {
			logrus.WithError(err).Error("get process cpu percent error")
		}
	}

	c, err := cpu.Percent(0, false)
	if err != nil {
		logrus.WithError(err).Error("get cpu percent error")
	} else {
		CurrentCpuPercent = c[0]
		items = append(items, MetricData{
			Name: "cpu",
			Fields: map[string]any{
				"cpu_percent":         c[0],
				"process_cpu_percent": cpuPercent,
			},
			Time: now,
		})
	}

	//pMem, _ := p.MemoryInfo()
	//pMem.RSS = m.HeapAlloc

	// 当前 Heap 内存使用量
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 内存使用率
	memoryPercent, err := p.MemoryPercent()
	if err != nil {
		logrus.WithError(err).Error("get memory percent error")
	} else {
		v, err := mem.VirtualMemory()
		if err != nil {
			logrus.WithError(err).Error("get virtual memory error")
		} else {
			pm, err := p.MemoryInfo()
			if err != nil {
				logrus.WithError(err).Error("get memory info error")
			} else {
				CurrentMemPercent = v.UsedPercent
				items = append(items, MetricData{
					Name: "mem",
					Fields: map[string]any{
						"heap_alloc": int64(m.HeapAlloc), // 当前 Heap 内存使用量
						//"gc":                  m.NumGC,
						"process_rss": int64(pm.RSS),  // 物理内存大小
						"process_vms": int64(pm.VMS),  // 虚拟内存大小
						"mem_total":   int64(v.Total), // 系统总内存
						"mem_used":    int64(v.Used),  // 占用的内存

						"mem_percent":         v.UsedPercent, // 内存使用率
						"process_mem_percent": memoryPercent,
					},
					Time: now,
				})
			}
		}
	}

	// 获取所有磁盘分区的使用情况
	//partitions, _ := disk.Partitions(false) // false: 不包含虚拟文件系统（如 /proc）
	//for _, par := range partitions {
	//	usage, _ := disk.Usage(par.Mountpoint)    // 传入挂载点（如 "/"、"C:"）
	//	realUsed := usage.Used - usage.InodesUsed // 根据场景调整
	//	diskUsagePercent := float64(realUsed) / float64(usage.Total) * 100
	//	//fmt.Printf("Disk %s (%s) Usage: %.2f%%\n",
	//	//	par.Mountpoint,
	//	//	par.Device,
	//	//	usage.UsedPercent,
	//	//)
	//}

	// 获取所有磁盘的 IO 统计
	ioCounters, _ := disk.IOCounters() // 返回 map[磁盘名称]IO 数据
	if ioCounters != nil && len(ioCounters) > 0 {
		var diskTotal uint64

		var currentDisksUsage []DiskUsage
		for name, io := range ioCounters {
			usage, err := disk.Usage(name) // 传入挂载点（如 "/"、"C:"）
			if err != nil {
				if strings.HasPrefix(err.Error(), "no such file or directory") {
					continue
				}
				logrus.WithError(err).Warnf("get disk(%s) usage error", name)
			} else {
				realUsed := usage.Used - usage.InodesUsed // 根据场景调整
				diskUsagePercent := float64(realUsed) / float64(usage.Total) * 100

				currentDisksUsage = append(currentDisksUsage, DiskUsage{
					Name:         name,
					UsagePercent: diskUsagePercent,
				})
				if usage.Total > diskTotal {
					diskTotal = usage.Total
					CurrentDiskPercent = diskUsagePercent
				}

				items = append(items, MetricData{
					Name: "disk",
					Tags: map[string]string{
						"name": name,
					},
					Fields: map[string]interface{}{
						"disk_total":   int64(usage.Total),
						"disk_used":    int64(realUsed),
						"disk_percent": diskUsagePercent,
						//"read_time":   io.ReadTime,
						//"write_time":  io.WriteTime,
					},
					Time: now,
				})
			}

			lastIO, ok := lastIOCounters[name]
			if !ok {
				continue
			}
			readSpeed := float64(io.ReadBytes-lastIO.ReadBytes) / period
			writeSpeed := float64(io.WriteBytes-lastIO.WriteBytes) / period
			items = append(items, MetricData{
				Name: "disk_speed",
				Tags: map[string]string{
					"name": name,
				},
				Fields: map[string]interface{}{
					"read_speed":  readSpeed,
					"write_speed": writeSpeed,
				},
				Time: now,
			})

			readCount := float64(io.ReadCount-lastIO.ReadCount) / period
			writeCount := float64(io.WriteCount-lastIO.WriteCount) / period
			items = append(items, MetricData{
				Name: "disk_count",
				Tags: map[string]string{
					"name": name,
				},
				Fields: map[string]interface{}{
					"read_count":  readCount,
					"write_count": writeCount,
				},
				Time: now,
			})
		}
		lastIOCounters = ioCounters
		slices.SortFunc(currentDisksUsage, func(a, b DiskUsage) int {
			return cmp.Compare(a.Name, b.Name)
		})
		CurrentDisksUsage = currentDisksUsage
	}

	// 系统平均负载
	avg, err := load.Avg()
	if err != nil {
		logrus.WithError(err).Error("get avg load error")
	} else {
		items = append(items, MetricData{
			Name: "avg_load",
			Fields: map[string]any{
				"load1":  avg.Load1,
				"load5":  avg.Load5,
				"load15": avg.Load15,
			},
			Time: now,
		})
	}

	netCounters, err := net.IOCounters(false)
	if err != nil {
		logrus.WithError(err).Error("get net io counters error")
	} else {
		CurrentNetRate = 0
		for _, n := range netCounters {
			if lastIO, ok := lastNetIOCounters[n.Name]; ok {
				bytesSent := float64(n.BytesSent-lastIO.BytesSent) / period
				bytesRecv := float64(n.BytesRecv-lastIO.BytesRecv) / period
				packetsSent := float64(n.PacketsSent-lastIO.PacketsSent) / period
				packetsRecv := float64(n.PacketsRecv-lastIO.PacketsRecv) / period

				CurrentNetRate += bytesRecv

				items = append(items, MetricData{
					Name: "net",
					Tags: map[string]string{
						"name": n.Name,
					},
					Fields: map[string]interface{}{
						"bytes_sent":   bytesSent,
						"bytes_recv":   bytesRecv,
						"packets_sent": packetsSent,
						"packets_recv": packetsRecv,
						"errin":        float64(n.Errin-lastIO.Errin) / period,
						"errout":       float64(n.Errout-lastIO.Errout) / period,
						"dropin":       float64(n.Dropin-lastIO.Dropin) / period,
						"dropout":      float64(n.Dropout-lastIO.Dropout) / period,
					},
					Time: now,
				})
			}
			lastNetIOCounters[n.Name] = n
		}
	}

	if len(items) > 0 {
		err = metricDatabase.Write(items)
		if err != nil {
			logrus.WithError(err).Error("save metrics error")
		}
	}
}

func Run(ctx context.Context) {
	if !config.Conf.Metrics.Enabled {
		return
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	if metricConf.Type == "OpenGemini" {
		metricDatabase = NewOpenGeminiClient(&metricConf)
	} else {
		metricDatabase = NewInfluxDBClient(&metricConf)
	}

	lastIOCounters, _ = disk.IOCounters() // 返回 map[磁盘名称]IO 数据
	lastNetIOCounters = map[string]net.IOCountersStat{}
	netCounters, _ := net.IOCounters(false)
	for _, n := range netCounters {
		lastNetIOCounters[n.Name] = n
	}

	defer metricDatabase.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metricRun()
		}
	}
}

var metricDatabase MetricDatabase

func Query(metricName string, ag, interval string, tags map[string]string, start, end time.Time) (rows []map[string]any, err error) {
	if metricDatabase == nil {
		return nil, nil
	}
	return metricDatabase.Query(metricName, ag, interval, tags, start, end)
}

func QueryByGroup(metricName string, ag, group, interval string, tags map[string]string, start, end time.Time) (seriesData []Series, err error) {
	if metricDatabase == nil {
		return
	}
	return metricDatabase.QueryByGroup(metricName, ag, group, interval, tags, start, end)
}

func SaveMetric(items ...MetricData) error {
	if metricDatabase == nil {
		return nil
	}
	return metricDatabase.Write(items)
}

var metricConf config.Metrics

func init() {
	conf := config.Conf.Metrics

	if conf.Retention <= 0 {
		conf.Retention = 3
	}
	if conf.Database == "" {
		conf.Database = "topmq"
	}
	metricConf = conf
}
