package mqttserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/iotopo/topmq/cache"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/internal/utils"
	"github.com/iotopo/topmq/web"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/system"
	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"
	"log/slog"
	"os"
)

var server *mqtt.Server

func GetInfo() *system.Info {
	if server != nil {
		return server.Info
	}
	return nil
}

func Init() {
	setupMonitoringRouter(web.Router)
}

func Run(ctx context.Context) error {
	conf := config.Conf.MQTTServer

	if conf.TLS {
		if !utils.PathExists("./mqtt/certs") { // 自动创建自签名证书
			generateCerts()
		}
	}

	server = mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	level := new(slog.LevelVar)
	level.Set(slog.LevelWarn)
	server.Log = slog.New(sloglogrus.Option{Level: level, Logger: logrus.StandardLogger()}.NewLogrusHandler())

	err := server.AddHook(newAuthHook(ctx), nil)
	if err != nil {
		return err
	}
	//	err := server.AddHook(new(auth.AllowHook), nil)
	//	if err != nil {
	//		return err
	//	}

	if conf.Persist {
		err := server.AddHook(newRedisStorageHook(cache.GetClient()), &RedisOptions{
			HPrefix: config.Conf.Redis.KeyPrefix + "mqtt:storage-",
		})
		if err != nil {
			return err
		}
	}

	var tlsConf *tls.Config
	if conf.TLS {
		// 只信任本地 CA 颁发的证书
		caCert, err := os.ReadFile(caCertFile)
		if err != nil {
			return fmt.Errorf("加载 CA 证书失败: %v", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		serverCer, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
		if err != nil {
			return fmt.Errorf("failed to read mqtt tls cert: %w", err)
		}
		tlsConf = &tls.Config{
			ClientCAs:    caCertPool,
			Certificates: []tls.Certificate{serverCer},
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
	}

	if conf.TCPPort <= 0 {
		conf.TCPPort = 1883
	}
	tcpAddr := fmt.Sprintf(":%d", conf.TCPPort)

	serverConf := listeners.Config{
		Type:      listeners.TypeTCP,
		ID:        "tcp",
		Address:   tcpAddr,
		TLSConfig: tlsConf,
	}

	tcp := listeners.NewTCP(serverConf)
	if err := server.AddListener(tcp); err != nil {
		return err
	}

	if conf.WSPort > 0 {
		wsAddr := fmt.Sprintf(":%d", conf.WSPort)
		ws := listeners.NewWebsocket(listeners.Config{
			Type:      listeners.TypeWS,
			ID:        "ws",
			Address:   wsAddr,
			TLSConfig: tlsConf,
		})
		if err := server.AddListener(ws); err != nil {
			return err
		}
	}
	//listeners.TypeSysInfo

	if conf.ExporterPort > 0 {
		go runExporter(fmt.Sprintf(":%d", conf.ExporterPort))
	}
	defer stopServer()

	return server.Serve()
}

//func SetupRouter(_ db.Database, router *gin.Engine) {
//	router.GET("api/v1/mqtt_server/info", func(c *gin.Context) {
//		if server == nil {
//			web.SuccessResponse(c, nil)
//			return
//		}
//		web.SuccessResponse(c, server.Info)
//	})
//}

func stopServer() {
	if httpServer != nil {
		err := httpServer.Close()
		if err != nil {
			logrus.WithError(err).Warn("failed to stop mqtt server exporter")
		}
	}
	if server != nil {
		_ = server.Close()
	}
}

//	func Publish(topic string, qos byte, retained bool, payload []byte) error {
//		return server.Publish(topic, payload, retained, qos)
//	}
//
//	func Subscribe(topic string, qos byte, handler func(topic string, msg []byte, time time.Time)) error {
//		subscriptionId++
//		return server.Subscribe(topic, subscriptionId, func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
//			handler(pk.TopicName, pk.Payload, time.Now())
//		})
//	}
//
// var subscriptionId = 0
