package mqttserver

import (
	"context"
	"github.com/iotopo/topmq/metrics"
	"github.com/mochi-mqtt/server/v2/system"
	"github.com/sirupsen/logrus"
	"time"
)

var CurrentMessagesReceived float64
var CurrentMessagesSent float64
var CurrentMessagesDropped float64

func saveMetric(ctx context.Context, lastInfo *system.Info) {
	if server == nil {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				logrus.WithError(err).Error("mqtt server metric panic")
			} else {
				logrus.Errorf("mqtt server metric panic: %v", e)
			}
		}
	}()

	period := float64(15)
	info := server.Info

	CurrentMessagesReceived = float64(info.MessagesReceived-lastInfo.MessagesReceived) / period
	CurrentMessagesSent = float64(info.MessagesSent-lastInfo.MessagesSent) / period
	CurrentMessagesDropped = float64(info.MessagesDropped-lastInfo.MessagesDropped) / period

	err := metrics.SaveMetric(metrics.MetricData{
		Name: "mqtt_server",
		Fields: map[string]any{
			"subscriptions": info.Subscriptions,
			"memory_alloc":  info.MemoryAlloc,
			"threads":       info.Threads,
		},
	}, metrics.MetricData{
		Name: "mqtt_server_bytes",
		Fields: map[string]any{
			"bytes_received": float64(info.BytesReceived-lastInfo.BytesReceived) / period,
			"bytes_sent":     float64(info.BytesSent-lastInfo.BytesSent) / period,
		},
	}, metrics.MetricData{
		Name: "mqtt_server_clients",
		Fields: map[string]any{
			"clients_connected":    info.ClientsConnected,
			"clients_disconnected": info.ClientsDisconnected,
		},
	}, metrics.MetricData{
		Name: "mqtt_server_messages",
		Fields: map[string]any{
			"messages_received": CurrentMessagesReceived,
			"messages_sent":     CurrentMessagesSent,
			"messages_dropped":  CurrentMessagesDropped,
			"inflight":          info.Inflight,
			"inflight_dropped":  float64(info.InflightDropped-lastInfo.InflightDropped) / period,
			"retained":          info.Retained,
		},
	}, metrics.MetricData{
		Name: "mqtt_server_packets",
		Fields: map[string]any{
			"packets_received": float64(info.PacketsReceived-lastInfo.PacketsReceived) / period,
			"packets_sent":     float64(info.PacketsSent-lastInfo.PacketsSent) / period,
		},
	})
	if err != nil {
		logrus.WithError(err).Error("failed to save mqtt server metrics")
	}
}

func RunMetrics(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	lastInfo := server.Info.Clone()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			saveMetric(ctx, lastInfo)
			lastInfo = server.Info.Clone()
		}
	}
}
