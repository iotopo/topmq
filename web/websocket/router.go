package websocket

import "github.com/iotopo/topmq/web/middlewares"

type IHandler interface {
	Publish(topic string, payload any)
	Log(format string, a ...any)
	AddSubscriber(topic string, createSubscriber func() (Subscriber, error))
	GetUser() *middlewares.LoginUser
	GetSessionID() string
}

type Router interface {
	Match(topic string) bool
	OnSubscribe(topic string, payload any, handler IHandler) error
	OnPublish(topic string, payload any, handler IHandler) error
}

//func (r *Router) Run() {
//
//}

var routers []Router

func RegisterRouter(router Router) {
	routers = append(routers, router)
}
