package metrics

import (
	"context"
	"errors"
	"fmt"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/config/secrets"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/openGemini/opengemini-client-go/opengemini"
	"github.com/sirupsen/logrus"
)

var errOpenGeminiNotCreated = fmt.Errorf("open_gemini client created")
var errOpenGeminiNotConnected = fmt.Errorf("open_gemini client not connected")

type OpenGeminiClient struct {
	sync.Mutex
	client    opengemini.Client
	connected bool
	inited    bool
	running   bool
	database  string
	retention int
}

func (conn *OpenGeminiClient) Write(datas []MetricData) error {
	if !conn.connected {
		return errors.New("open_gemini not connected")
	}
	if conn.client == nil {
		return errors.New("could not connect to open_gemini")
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

	var pointList []*opengemini.Point
	for i := range datas {
		data := datas[i]
		tags := data.Tags
		p := &opengemini.Point{
			Measurement: data.Name,
			Tags:        tags,
			Fields:      data.Fields,
			Precision:   opengemini.PrecisionSecond,
			Timestamp:   data.Time.Round(time.Second).UnixNano(),
			//Time:        data.Time,
		}
		pointList = append(pointList, p)
	}
	err := conn.client.WriteBatchPoints(context.Background(), conn.database, pointList)
	if err != nil {
		logrus.WithError(err).Warn("system metrics write to open_gemini")
	}

	return nil
}

func (conn *OpenGeminiClient) RawQuery(database, cmd string) (res []*opengemini.SeriesResult, err error) {
	return conn.query(database, cmd)
}

// Query convenience function to query the database
func (conn *OpenGeminiClient) query(database, cmd string) (res []*opengemini.SeriesResult, err error) {
	client := conn.client

	if client == nil {
		return nil, errOpenGeminiNotCreated
	}

	if !conn.connected {
		return nil, errOpenGeminiNotConnected
	}

	q := opengemini.Query{
		Command:  cmd,
		Database: database,
	}
	logrus.Info(cmd)
	if response, err := client.Query(q); err == nil {
		if response.Error != "" {
			logrus.Errorf("open_gemini query data error: %s", response.Error)
			return res, fmt.Errorf("open_gemini response data error: %s", response.Error)
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func (conn *OpenGeminiClient) QueryByGroup(metricName string, ag, group, interval string, tags map[string]string, start, end time.Time) (seriesData []Series, err error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, `select %s(*)`, ag)
	fmt.Fprintf(&sb, ` from "%s"`, metricName)
	if tags != nil {
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

func (conn *OpenGeminiClient) Query(metricName string, ag, interval string, tags map[string]string, start, end time.Time) (rows []map[string]any, err error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, `select %s(*) from "%s"`, ag, metricName)
	if tags != nil {
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

// Close 关闭
func (conn *OpenGeminiClient) Close() {
	conn.Lock()
	defer conn.Unlock()

	// 停止ping
	conn.running = false
	if conn.client != nil {
		_ = conn.client.Close()
	}
}

func (conn *OpenGeminiClient) init() error {
	//database := conn.database
	//logrus.Debug("opengemini online, init logger:", database)
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
		err := conn.client.CreateDatabaseWithRp(conn.database, opengemini.RpConfig{
			Name:               "autogen",
			Duration:           fmt.Sprintf("%dd", retention),
			ShardGroupDuration: fmt.Sprintf("%dd", shardGroupDuration),
		})
		if err != nil {
			logrus.Infof("opengemini create database error: %v", err)
		}
	}
	return nil
}

func NewOpenGeminiClient(conf *config.Metrics) *OpenGeminiClient {
	clientConf := &opengemini.Config{
		ContentType: opengemini.ContentTypeMsgPack,
	}
	addrs := strings.Split(conf.Address, ",")
	for _, addr := range addrs {
		// NOTE: 使用 net.SplitHostPort() 而不是用 ":" 拆解字符串以支持 ipv6
		// 例如 "[::1]:80" => "::1", "80"
		host, portStr, err := net.SplitHostPort(addr)
		var port int
		if err == nil {
			port, err = strconv.Atoi(portStr)
		}
		if err == nil {
			clientConf.Addresses = append(clientConf.Addresses, opengemini.Address{
				Host: host,
				Port: port,
			})
		} else {
			logrus.WithError(err).Errorf("invalid open_gemini address")
		}
	}

	if conf.Token != "" {
		clientConf.AuthConfig = &opengemini.AuthConfig{
			AuthType: opengemini.AuthTypePassword,
			Token:    conf.Token,
		}
	} else if conf.Username != "" && conf.Password != "" {
		password := conf.Password
		if password != "" {
			var err error
			password, err = secrets.Decode(password)
			if err != nil {
				logrus.WithError(err).Error("metrics decode influxdb password error")
			}
		}
		clientConf.AuthConfig = &opengemini.AuthConfig{
			AuthType: opengemini.AuthTypeToken,
			Username: conf.Username,
			Password: password,
		}
	}

	client, err := opengemini.NewClient(clientConf)
	if err != nil {
		logrus.WithError(err).Fatal("create open_gemini client error")
	}

	conn := &OpenGeminiClient{
		running:   true,
		database:  conf.Database,
		retention: conf.Retention,
	}
	//if conf.Mode == "Point" {
	//	logger.newMode = true
	//}
	conn.client = client
	return conn
}
