package metrics

import (
	"fmt"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/config/secrets"
	"strings"
	"sync"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var errInfluxDBNotInited = fmt.Errorf("influxdb client is nil")
var errInfluxDBNotConnected = fmt.Errorf("influxdb client not connected")

// InfluxDBClient influxdb logger
type InfluxDBClient struct {
	sync.Mutex
	client    influxdb.Client
	connected bool
	inited    bool
	running   bool
	database  string
	retention int
}

func (conn *InfluxDBClient) Write(datas []MetricData) error {
	if !conn.connected {
		return errInfluxDBNotConnected
	}
	if conn.client == nil {
		return errInfluxDBNotInited
	}
	if len(datas) == 0 {
		return nil
	}

	conn.Lock()
	defer conn.Unlock()

	if !conn.inited {
		if err := conn.init(); err != nil {
			fmt.Println(err.Error())
		}
		conn.inited = true
	}

	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Precision: "s",
		Database:  conn.database,
	})
	if err != nil {
		return err
	}
	for i := range datas {
		data := datas[i]
		tags := data.Tags
		if p, err := influxdb.NewPoint(data.Name, tags, data.Fields, data.Time); err == nil {
			bp.AddPoint(p)
		} else {
			logrus.Errorf("influxdb create point error: %v", err)
		}
	}
	if err := conn.client.Write(bp); err != nil {
		logrus.WithError(err).Warn("write system metrics to influxdb")
	}
	return nil
}

func (conn *InfluxDBClient) RawQuery(database, cmd string) (res []influxdb.Result, err error) {
	return conn.query(database, cmd)
}

// Query convenience function to query the database
func (conn *InfluxDBClient) query(database, cmd string) (res []influxdb.Result, err error) {
	client := conn.client

	if client == nil {
		err = errInfluxDBNotInited
		return
	}

	if !conn.connected {
		err = errInfluxDBNotConnected
		return
	}

	q := influxdb.Query{
		Command:  cmd,
		Database: database,
	}
	logrus.Info(cmd)
	if response, err := client.Query(q); err == nil {
		if response.Error() != nil {
			logrus.WithError(response.Error()).Error("influxdb query data error")
			return res, errors.Wrap(response.Error(), "influxdb response data error")
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func (conn *InfluxDBClient) QueryByGroup(metricName string, ag, group, interval string, tags map[string]string, start, end time.Time) (seriesData []Series, err error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, `select %s(*)`, ag)
	fmt.Fprintf(&sb, ` from "%s"`, metricName)

	//config.Conf.InstanceID != "" && config.Conf.InstanceID != "standalone"
	if len(tags) > 0 {
		var cnd []string
		for k, v := range tags {
			cnd = append(cnd, fmt.Sprintf(`"%s"='%s'`, k, v))
		}
		fmt.Fprintf(&sb, ` where %s and time > '%s' and time <= '%s'`, strings.Join(cnd, " and "), start.Format(time.RFC3339), end.Format(time.RFC3339))
	} else {
		fmt.Fprintf(&sb, ` where time > '%s' and time <= '%s'`, start.Format(time.RFC3339), end.Format(time.RFC3339))
	}

	fmt.Fprintf(&sb, ` group by "%s",time(%s) fill(previous)`, group, interval)

	results, err := conn.query(conn.database, sb.String())
	if err != nil {
		return nil, fmt.Errorf("query metrics from tsdb error: %w", err)
	}

	if len(results) > 0 && len(results[0].Series) > 0 {
		for _, series := range results[0].Series {
			tagValue := series.Tags[group]
			var rows []map[string]any
			for _, values := range series.Values {
				row := map[string]any{}
				for i, col := range series.Columns {
					row[strings.TrimPrefix(col, ag+"_")] = values[i]
				}
				rows = append(rows, row)
			}
			seriesData = append(seriesData, Series{
				Name:  tagValue,
				Items: rows,
			})
		}
	}
	return
}

func (conn *InfluxDBClient) Query(metricName string, ag, interval string, tags map[string]string, start, end time.Time) (rows []map[string]any, err error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, `select %s(*) from "%s"`, ag, metricName)

	if len(tags) > 0 {
		var cnd []string
		for k, v := range tags {
			cnd = append(cnd, fmt.Sprintf(`"%s"='%s'`, k, v))
		}
		fmt.Fprintf(&sb, ` where %s and time > '%s' and time <= '%s'`, strings.Join(cnd, " and "), start.Format(time.RFC3339), end.Format(time.RFC3339))
	} else {
		fmt.Fprintf(&sb, ` where time > '%s' and time <= '%s'`, start.Format(time.RFC3339), end.Format(time.RFC3339))
	}

	fmt.Fprintf(&sb, ` group by time(%s) fill(previous)`, interval)

	results, err := conn.query(conn.database, sb.String())
	if err != nil {
		return nil, fmt.Errorf("query metrics from tsdb error: %w", err)
	}

	if len(results) > 0 && len(results[0].Series) > 0 {
		series := results[0].Series[0]
		for _, values := range series.Values {
			row := map[string]any{}
			for i, col := range series.Columns {
				row[strings.TrimPrefix(col, ag+"_")] = values[i]
			}
			rows = append(rows, row)
		}
	}
	return
}

func (conn *InfluxDBClient) Close() {
	logrus.Info("close influxdb logger")

	conn.Lock()
	defer conn.Unlock()

	// 停止ping
	conn.running = false
	if conn.client != nil {
		_ = conn.client.Close()
	}
	logrus.Info("influxdb client close")
}

func (conn *InfluxDBClient) init() error {
	if conn.client != nil && conn.connected {
		retention := conn.retention
		if retention <= 0 {
			retention = 7
		}
		// DURATION 不能低于 SHARD DURATION（默认 7d）
		shardGroupDuration := 1
		if retention >= 7 {
			shardGroupDuration = 7
		}
		sql := fmt.Sprintf(`CREATE DATABASE %q WITH DURATION %dd SHARD DURATION %dd`, conn.database, retention, shardGroupDuration)
		if _, err := conn.client.Query(influxdb.Query{Command: sql}); err != nil {
			logrus.Infof("influxdb create database error: %v", err)
		}
	}
	return nil
}

func (conn *InfluxDBClient) ping() {
	if conn.client != nil {
		_, _, err := conn.client.Ping(5 * time.Second) // if this takes more than 5 seconds then influxdb is probably down
		if err != nil {
			conn.connected = false
			logrus.WithError(err).Error("Error connecting to InfluxDB")
		} else {
			conn.connected = true
		}
	}
}
func (conn *InfluxDBClient) startPing() {
	conn.ping()
	go func() {
		time.Sleep(10 * time.Second)
		for {
			conn.ping()
			if !conn.running {
				return
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

func NewInfluxDBClient(conf *config.Metrics) *InfluxDBClient {
	addr := conf.Address
	if !strings.HasPrefix(addr, "http://") {
		addr = "http://" + addr
	}

	password := conf.Password
	if password != "" {
		var err error
		password, err = secrets.Decode(password)
		if err != nil {
			logrus.WithError(err).Error("metrics decode influxdb password error")
		}
	}
	client, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     addr,
		Username: conf.Username,
		Password: password,
	})
	if err != nil {
		logrus.WithError(err).Fatal("create influxdb client error")
	}

	conn := &InfluxDBClient{
		running:   true,
		database:  conf.Database,
		retention: conf.Retention,
	}
	conn.client = client
	conn.startPing()
	return conn
}
