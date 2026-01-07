package mqttserver

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/iotopo/topmq/web"
)

func setupMonitoringRouter(router gin.IRouter) {
	api := router.Group("api/v1/monitoring")
	api.GET("clients", func(c *gin.Context) {
		var req ClientRequest
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		//total := server.GetClientsTotal()
		items := GetClients(req)
		web.SuccessResponse(c, items)
	})
	api.GET("client_details/:id", func(c *gin.Context) {
		clientID := c.Param("id")
		if clientID == "" {
			web.BadRequestResponse(c, "client id is required")
			return
		}
		web.SuccessResponse(c, GetClientDetails(clientID))
	})
	api.POST("close_client/:id", func(c *gin.Context) {
		clientID := c.Param("id")
		if clientID == "" {
			web.BadRequestResponse(c, "client id is required")
			return
		}
		err := CloseClient(clientID)
		if err != nil {
			web.InternalErrorResponse(c, err)
			return
		}
		web.SuccessResponse(c, nil)
	})
	api.GET("subscriptions", func(c *gin.Context) {
		var req SubscriptionRequest
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		items := GetSubscriptions(req)
		web.SuccessResponse(c, items)
	})
	api.GET("retained", func(c *gin.Context) {
		var req struct {
			Filter string `form:"filter"`
		}
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		items := GetRetained(req.Filter)
		web.SuccessResponse(c, items)
	})
	api.GET("retained_payload", func(c *gin.Context) {
		var req struct {
			Topic string `form:"topic" binding:"required"`
		}
		if err := c.ShouldBind(&req); err != nil {
			web.BadRequestResponse(c, err.Error())
			return
		}
		payload := GetRetainedPayload(req.Topic)
		web.SuccessResponse(c, gin.H{"payload": base64.StdEncoding.EncodeToString(payload)})
	})
	api.GET("metrics/overview", func(c *gin.Context) {
		web.SuccessResponse(c, GetOverviewInfo())
	})
}
