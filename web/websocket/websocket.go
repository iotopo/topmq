package websocket

import (
	"github.com/iotopo/topmq/auth"
	"github.com/iotopo/topmq/web"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var logger = logrus.StandardLogger()

var upgrader = &websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
	HandshakeTimeout: 20 * time.Second,
}

func SetupRouter() {
	//router gin.IRouter
	router := web.Router
	api := router.Group("ws", auth.Authorize)
	// websocket message
	api.GET("/msg/v1", func(c *gin.Context) {
		if !c.IsAborted() {
			logger.Debug("websocket init")

			conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
			if err != nil {
				logger.Infof("websocket upgrade error: %v", err)
				return
			}
			handler, err := newWSHandler(conn, c)
			if err != nil {
				logger.WithError(err).Error("websocket init error")
				return
			}
			//handler.user = user
			handler.run()
		} else {
			logger.Warn("websocket auth failed")
		}
	})
	//route.GET("/ws/info", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"sessionLen": Router.Len(),
	//	})
	//})
}
