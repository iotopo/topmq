package mqttserver

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

func registerPrometheusMetrics(registry prometheus.Registerer) {
	i := server.Info
	if registry == nil {
		registry = prometheus.DefaultRegisterer
	}

	type metrics struct {
		metricType string
		name       string
		help       string
		value      *int64
	}

	metricsList := []metrics{
		{"c", "bytes_received", "A count of total number of bytes received", &i.BytesReceived},
		{"c", "bytes_sent", "A counter total number of bytes sent", &i.BytesSent},
		{"g", "clients_connected", "A gauge of number of currently connected clients", &i.ClientsConnected},
		{"g", "clients_disconnected", "A gauge of total number of persistent clients", &i.ClientsDisconnected},
		{"c", "clients_maximum", "A count of maximum number of active clients that have been connected", &i.ClientsMaximum},
		{"g", "clients_total", "A gauge of total number of connected and disconnected clients with a persistent session currently connected and registered", &i.ClientsTotal},
		{"c", "messages_received", "A counter of total number of publish messages received", &i.MessagesReceived},
		{"c", "messages_sent", "A counter of total number of publish messages sent", &i.MessagesSent},
		{"c", "messages_dropped", "A counter of total number of publish messages dropped to slow subscriber", &i.MessagesDropped},
		{"g", "retained", "A gauge of total number of retained messages active on the broker", &i.Retained},
		{"g", "inflight", "A gauge of the number of messages currently in-flight", &i.Inflight},
		//{"c", "inflight_dropped", "A", &i.InflightDropped},
		{"g", "subscriptions", "A gauge of total number of subscriptions active on the broker", &i.Subscriptions},
		{"c", "packets_received", "A counter of the total number of packets received", &i.PacketsReceived},
		{"c", "packets_sent", "A counter of the total number of packets sent", &i.PacketsSent},
	}

	for _, m := range metricsList {
		m := m
		fn := func() float64 {
			return float64(atomic.LoadInt64(m.value))
		}

		switch m.metricType {
		case "c":
			registry.MustRegister(
				prometheus.NewCounterFunc(
					prometheus.CounterOpts{
						Name: m.name,
						Help: m.help,
					},
					fn,
				),
			)
		case "g":
			registry.MustRegister(
				prometheus.NewGaugeFunc(
					prometheus.GaugeOpts{
						Name: m.name,
						Help: m.help,
					},
					fn,
				),
			)
		}
	}

	buildInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			// Namespace: AppName,
			Name: "build_info",
			Help: "Build Information",
		},
		[]string{"goversion", "version"},
	)
	prometheus.MustRegister(buildInfo)
	buildInfo.With(prometheus.Labels{"goversion": runtime.Version(), "version": i.Version}).Set(1)
}

var httpServer *http.Server

func runExporter(address string) {
	registerPrometheusMetrics(prometheus.DefaultRegisterer)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		out, err := json.MarshalIndent(server.Info, "", "\t")
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
		}

		_, _ = w.Write(out)
	})

	mux.Handle("/metrics", promhttp.Handler())

	httpServer = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:         address,
		Handler:      mux,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		logrus.WithError(err).Fatal("failed to start mqtt server exporter")
	}
}
