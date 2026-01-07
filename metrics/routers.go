package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/iotopo/topmq/web"
	"time"
)

type QueryRequest struct {
	Metric string    `json:"metric" form:"metric"`
	Start  time.Time `json:"start" form:"start" binding:"required"`
	End    time.Time `json:"end" form:"end" binding:"required"`
	Ag     string    `json:"ag" form:"ag"`
	Group  string    `json:"group,omitempty" form:"group"` // 分组条件
	Host   bool      `json:"host,omitempty" form:"host"`   // 是否只查询当前 host
}

func getInterval(duration time.Duration) string {
	if duration >= 14*time.Hour*24 { // 14 天
		return "120m"
	} else if duration >= 7*24*time.Hour { // 7d
		return "60m"
	} else if duration >= 3*24*time.Hour { // 3d
		return "30m"
	} else if duration >= 1*24*time.Hour { // 1d
		return "10m"
	} else if duration >= 12*time.Hour { // 12h
		return "5m"
	} else if duration >= 6*time.Hour { // 6h
		return "1m"
	} else if duration >= 3*time.Hour { // 3h
		return "1m"
	} else { // 1h
		return "15s"
	}
}

func getHostTags() map[string]string {
	var tags map[string]string
	//if config.Conf.InstanceID != "" && config.Conf.InstanceID != "standalone" {
	//	tags = map[string]string{
	//		"host": config.Conf.InstanceID,
	//	}
	//}
	return tags
}

func Init() {
	router := web.Router
	api := router.Group("api/v1/metrics")

	api.GET("overview", func(c *gin.Context) {
		web.SuccessResponse(c, gin.H{
			"cpu_percent":  CurrentCpuPercent,
			"mem_percent":  CurrentMemPercent, // 内存使用率
			"disk_percent": CurrentDiskPercent,
			"net_in_rate":  CurrentNetRate,
			"disks_usage":  CurrentDisksUsage,
		})
	})
	api.GET("cpu", func(c *gin.Context) {
		var req QueryRequest
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		if req.Ag == "" {
			req.Ag = "mean"
		}
		interval := getInterval(req.End.Sub(req.Start))
		data, err := Query("cpu", req.Ag, interval, getHostTags(), req.Start, req.End)
		if err != nil {
			web.InternalErrorResponse(c, err)
			return
		}
		web.SuccessResponse(c, data)
	})
	api.GET("mem", func(c *gin.Context) {
		var req QueryRequest
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		if req.Ag == "" {
			req.Ag = "mean"
		}
		interval := getInterval(req.End.Sub(req.Start))
		data, err := Query("mem", req.Ag, interval, getHostTags(), req.Start, req.End)
		if err != nil {
			web.InternalErrorResponse(c, err)
			return
		}
		web.SuccessResponse(c, data)
	})
	api.GET("avg_load", func(c *gin.Context) {
		var req QueryRequest
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		if req.Ag == "" {
			req.Ag = "mean"
		}
		interval := getInterval(req.End.Sub(req.Start))
		data, err := Query("avg_load", req.Ag, interval, getHostTags(), req.Start, req.End)
		if err != nil {
			web.InternalErrorResponse(c, err)
			return
		}
		web.SuccessResponse(c, data)
	})

	api.GET("disk", func(c *gin.Context) {
		var req QueryRequest
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		if req.Ag == "" {
			req.Ag = "mean"
		}
		interval := getInterval(req.End.Sub(req.Start))
		data, err := QueryByGroup("disk", req.Ag, "name", interval, getHostTags(), req.Start, req.End)
		if err != nil {
			web.InternalErrorResponse(c, err)
			return
		}
		web.SuccessResponse(c, data)
	})

	api.GET("disk_speed", func(c *gin.Context) {
		var req QueryRequest
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		if req.Ag == "" {
			req.Ag = "mean"
		}
		interval := getInterval(req.End.Sub(req.Start))
		data, err := QueryByGroup("disk_speed", req.Ag, "name", interval, getHostTags(), req.Start, req.End)
		if err != nil {
			web.InternalErrorResponse(c, err)
			return
		}
		web.SuccessResponse(c, data)
	})

	api.GET("disk_count", func(c *gin.Context) {
		var req QueryRequest
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		if req.Ag == "" {
			req.Ag = "mean"
		}
		interval := getInterval(req.End.Sub(req.Start))
		data, err := QueryByGroup("disk_count", req.Ag, "name", interval, getHostTags(), req.Start, req.End)
		if err != nil {
			web.InternalErrorResponse(c, err)
			return
		}
		web.SuccessResponse(c, data)
	})

	api.GET("net", func(c *gin.Context) {
		var req QueryRequest
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		if req.Ag == "" {
			req.Ag = "mean"
		}
		interval := getInterval(req.End.Sub(req.Start))

		data, err := QueryByGroup("net", req.Ag, "name", interval, getHostTags(), req.Start, req.End)
		if err != nil {
			web.InternalErrorResponse(c, err)
			return
		}
		web.SuccessResponse(c, data)
	})

	api.GET("query", func(c *gin.Context) {
		var req QueryRequest
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		if req.Ag == "" {
			req.Ag = "mean"
		}
		interval := getInterval(req.End.Sub(req.Start))
		//var tags map[string]string
		//if req.Host {
		//	tags = getHostTags()
		//}
		tags := getHostTags()
		if req.Group == "" {
			data, err := Query(req.Metric, req.Ag, interval, tags, req.Start, req.End)
			if err != nil {
				web.InternalErrorResponse(c, err)
				return
			}
			web.SuccessResponse(c, data)
		} else {
			data, err := QueryByGroup(req.Metric, req.Ag, req.Group, interval, tags, req.Start, req.End)
			if err != nil {
				web.InternalErrorResponse(c, err)
				return
			}
			web.SuccessResponse(c, data)
		}
	})
}
