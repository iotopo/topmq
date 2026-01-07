package websocket

import (
	"errors"
	"fmt"
	"github.com/iotopo/topmq/auth"
	"github.com/iotopo/topmq/internal/utils"
	"github.com/iotopo/topmq/internal/utils/cmap"
	"github.com/iotopo/topmq/web/middlewares"
	"github.com/panjf2000/ants/v2"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	writeWait      = 10 * time.Second            // Milliseconds until write times out.
	pongWait       = 60 * time.Second            // Timeout for waiting on pong.
	pingPeriod     = (60 * time.Second * 9) / 10 // Milliseconds between pings.
	maxMessageSize = int64(2 << 20)              // Maximum size in bytes of a message. 2M.

	typeSubscribe    = "subscribe"
	typeUnsubscribe  = "unsubscribe"
	typeFetch        = "fetch"   // 主动请求数据
	typePublish      = "publish" // 主动发布数据
	typeNotification = "notification"
)

type Subscriber interface {
	Unsubscribe() error
}

// 客户端消息，用于订阅或取消订阅
type clientMessage struct {
	Type    string      `json:"type,omitempty"`
	Topic   string      `json:"topic,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// 服务器端消息，用于返回响应数据
type serverMessage struct {
	Type    string      `json:"type,omitempty"`
	Topic   string      `json:"topic,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

type wsHandler struct {
	// conn net.Conn
	id            string
	conn          *websocket.Conn
	writeStop     chan struct{}
	writeBuf      []serverMessage
	writeChan     chan serverMessage
	batchSize     int
	subscriberMap cmap.Map[string, Subscriber] //sync.Map
	doClose       sync.Once
	running       bool
	closed        bool
	context       *gin.Context
	pool          *ants.Pool
}

func newWSHandler(conn *websocket.Conn, c *gin.Context) (*wsHandler, error) {
	handler := &wsHandler{
		id:            utils.UUIDWithoutDash(),
		conn:          conn,
		subscriberMap: cmap.New[string, Subscriber](),
		writeStop:     make(chan struct{}, 1),
		writeChan:     make(chan serverMessage, 100),
		batchSize:     200,
		running:       true,
		context:       c,
	}
	handler.pool, _ = ants.NewPool(100)
	conn.SetReadLimit(maxMessageSize)
	err := conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return nil, err
	}

	conn.SetPongHandler(func(string) error {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		logger.Debug("ws ping pong")
		return nil
	})
	conn.SetCloseHandler(func(code int, text string) error {
		logger.Info("ws closed")
		handler.closed = true
		return nil
	})
	return handler, nil
}

// processUnsubscribe 处理取消订阅请求
func (handler *wsHandler) processUnsubscribe(msg *clientMessage) {
	if v, ok := handler.subscriberMap.Pop(msg.Topic); ok {
		logrus.WithField("id", handler.id).Debug("websocket unsubscribe: ", msg.Topic)
		if subscriber, ok := v.(Subscriber); ok {
			subscriber.Unsubscribe()
		}
	}
}

func (handler *wsHandler) processFetch(msg *clientMessage) {
	logrus.WithField("id", handler.id).Debug("websocket fetch: ", msg.Topic)
}

func (handler *wsHandler) processPublish(msg *clientMessage) {
	logrus.WithField("id", handler.id).Debug("websocket publish: ", msg.Topic)
	for _, router := range routers {
		if router.Match(msg.Topic) {
			err := router.OnPublish(msg.Topic, msg.Payload, handler)
			if err != nil {
				logger.WithError(err).Error("websocket handler publish msg failed")
			}
			break
		}
	}
}

func (handler *wsHandler) parseMessage(msg *clientMessage) {
	defer func() {
		if r := recover(); r != nil {
			logger.WithField("id", handler.id).Errorf("websocket handler parseMessage panic: %v\n, %v", r, string(debug.Stack()))
		}
	}()
	switch msg.Type {
	case typeSubscribe:
		handler.processSubscribe(msg)
	case typeUnsubscribe:
		handler.processUnsubscribe(msg)
	case typeFetch:
		handler.processFetch(msg)
	case typePublish:
		handler.processPublish(msg)
	}
}

// processSubscribe 处理订阅请求
func (handler *wsHandler) processSubscribe(msg *clientMessage) {
	logger.WithField("id", handler.id).Debug("websocket subscribe: ", msg.Topic)

	for _, router := range routers {
		if router.Match(msg.Topic) {
			err := router.OnSubscribe(msg.Topic, msg.Payload, handler)
			if err != nil {
				logger.WithError(err).Error("websocket subscribe failed")
			}
			break
		}
	}
}

func (handler *wsHandler) AddSubscriber(topic string, createSubscriber func() (Subscriber, error)) {
	handler.addSubscriber(topic, createSubscriber)
}

func (handler *wsHandler) addSubscriber(topic string, createSubscriber func() (Subscriber, error)) {
	if handler.closed || !handler.running {
		return
	}
	if ok := handler.subscriberMap.Has(topic); !ok {
		subscriber, err := createSubscriber()
		if err != nil {
			handler.writeLog("create subscriber error: %v", err)
			return
		}
		if subscriber == nil {
			return
		}
		// 避免重复订阅, 如果已经订阅过，则取消刚才新建的订阅
		if _, exist := handler.subscriberMap.LoadOrStore(topic, subscriber); exist || !handler.running {
			subscriber.Unsubscribe()
		}
	}
}

func (handler *wsHandler) run() {
	defer func() {
		if e := recover(); e != nil {
			logger.WithField("id", handler.id).Errorf("ws unknown error: %v", e)
		}
		handler.close()
	}()

	// go handler.writePump()
	var group errgroup.Group

	group.Go(handler.writePump)

	group.Go(func() error {
		defer func() {
			if e := recover(); e != nil {
				logger.WithField("id", handler.id).Errorf("ws unknown error: %v", e)
			}
			handler.close()
		}()
		conn := handler.conn
		for {
			var msg clientMessage
			err := conn.ReadJSON(&msg)

			if err != nil {
				var netError net.Error
				if errors.As(err, &netError) { // 网络通讯链路问题
					if netError.Timeout() {
						continue
					}
				}
				if _, ok := err.(*websocket.CloseError); ok {
					//logger.Warnf("ws closed: %+v", err)
					return nil
				}
				return fmt.Errorf("ws read error: %+v", err)
			}

			handler.pool.Submit(func() {
				handler.parseMessage(&msg)
			})
		}
	})
	if err := group.Wait(); err != nil {
		logger.Error(err.Error())
	}
}

func (handler *wsHandler) close() {
	handler.doClose.Do(func() {
		logger.WithField("id", handler.id).Info("ws disconnected")
		handler.running = false
		handler.pool.Release()
		for _, subscriber := range handler.subscriberMap.Iter() {
			subscriber.Unsubscribe()
		}
		close(handler.writeChan)
		handler.writeStop <- struct{}{}
		handler.conn.Close()
	})
}

func (handler *wsHandler) GetSessionID() string {
	return handler.id
}

func (handler *wsHandler) GetUser() *middlewares.LoginUser {
	return middlewares.GetLoginUser(handler.context)
}

func (handler *wsHandler) Log(format string, a ...any) {
	handler.writeLog(format, a...)
}

func (handler *wsHandler) writeLog(format string, a ...any) {
	payload := fmt.Sprintf(format, a...)
	logger.Warn(payload)
	handler.writeData(serverMessage{
		Type:    "log",
		Payload: payload,
	})
}

func (handler *wsHandler) Publish(topic string, payload any) {
	handler.writeData(serverMessage{
		Topic:   topic,
		Payload: payload,
	})
}

func (handler *wsHandler) writeData(msg serverMessage) {
	defer func() {
		if e := recover(); e != nil {
			logger.WithField("id", handler.id).Error("ws write to chan error: %v", e)
		}
	}()
	if !handler.running || handler.closed {
		return
	}
	// handler.writeChan <- msg

	select {
	case handler.writeChan <- msg:
	case <-time.After(time.Second * 10):
		// overtime
		logger.WithField("id", handler.id).Warn("maybe blocked in ws write channel")
	}
}

// func (handler *wsHandler) writeJSON(msg interface{}) error {
//	handler.conn.SetWriteDeadline(time.Now().Add(writeWait))
//	return handler.conn.WriteJSON(msg)
// }

// flushBuffer 与 ping 在同一个 goroutine 中，不会存在并发问题
func (handler *wsHandler) flushBuffer() error {
	if len(handler.writeBuf) > 0 {
		if err := handler.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
			return fmt.Errorf("ws SetWriteDeadline error: %w", err)
		}

		err := handler.conn.WriteJSON(serverMessage{
			Type:    "batch",
			Payload: handler.writeBuf,
		})
		handler.writeBuf = []serverMessage{} // handler.writeBuf[:0]
		return err
	}
	return nil
}

// ping 与 flushBuffer 在同一个 goroutine 中，不会存在并发问题
func (handler *wsHandler) ping() error {
	if err := handler.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return fmt.Errorf("ws SetWriteDeadline error: %w", err)
	}
	return handler.conn.WriteMessage(websocket.PingMessage, []byte{})
}

func (handler *wsHandler) writePump() error {
	defer func() {
		if e := recover(); e != nil {
			logger.WithField("id", handler.id).Errorf("ws unknown error: %v", e)
		}
		handler.close()
	}()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	keepAliveTicker := time.NewTicker(1 * time.Minute)
	defer keepAliveTicker.Stop()

	flushTicker := time.NewTicker(100 * time.Millisecond)
	defer flushTicker.Stop()

	for {
		select {
		case <-keepAliveTicker.C:
			if authToken := middlewares.GetAuthToken(handler.context); authToken != "" {
				auth.KeepSession(authToken)
			}
		case msg := <-handler.writeChan:
			handler.writeBuf = append(handler.writeBuf, msg)
			if len(handler.writeBuf) >= handler.batchSize {
				if err := handler.flushBuffer(); err != nil {
					// logger.WithError(err).Error("ws write msg error")
					// return
					return fmt.Errorf("ws write msg error: %v", err)
				}
			}
		case <-flushTicker.C:
			if err := handler.flushBuffer(); err != nil {
				// logger.WithError(err).Error("ws write msg error")
				// return
				return fmt.Errorf("ws write msg error: %v", err)
			}
		case <-handler.writeStop:
			return nil
		case <-ticker.C:
			err := handler.ping()
			if err != nil {
				// logger.WithError(err).Error("ws ping error")
				// return
				return fmt.Errorf("ws ping error: %v", err)
			}
		}
	}
}
