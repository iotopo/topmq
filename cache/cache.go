package cache

//
// TODO: rate limiting based on Redis. https://github.com/go-redis/redis_rate
// TODO: Simplified distributed locking implementation using Redis. https://github.com/bsm/redislock
import (
	"context"
	"github.com/go-redis/cache/v9"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/config/secrets"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

var conn redis.Cmdable
var client redis.UniversalClient
var ctx = context.Background()

//func InitClient() redis.UniversalClient {
//	return client
//}

func Init() {
	var err error
	conf := config.Conf.Redis
	username := conf.Username
	password := conf.Password
	if password != "" {
		password, err = secrets.Decode(password)
		if err != nil {
			logrus.WithError(err).Fatal("failed to decode redis password")
		}
	}

	db := conf.DB

	addr := conf.Addr
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	//idleTimeout := time.Duration(idleTimeout) * time.Second
	//connectTimeout := time.Duration(connectTimeout) * time.Second
	if conf.Sentinel != nil {
		//addr := os.Getenv("REDIS_SENTINEL_ADDRS")
		//if addr == "" {
		//	addr = "127.0.0.1:6379"
		//}
		sentinel := conf.Sentinel

		var sentinelPassword string
		sentinelPassword, err = secrets.Decode(sentinel.Password)
		if err != nil {
			logrus.WithError(err).Fatal("failed to decode sentinel password")
		}
		// 哨兵模式
		client = redis.NewFailoverClusterClient(&redis.FailoverOptions{
			// The master name.
			MasterName: sentinel.MasterName,
			// A seed list of host:port addresses of sentinel nodes.
			SentinelAddrs: sentinel.Addrs,
			// Sentinel password from "requirepass <password>" (if enabled) in Sentinel configuration
			SentinelUsername: sentinel.Username,
			SentinelPassword: sentinelPassword,
			ClientName:       sentinel.ClientName,

			// RouteByLatency: true,
			RouteRandomly:   true,
			Username:        username,
			Password:        password,
			PoolSize:        conf.PoolSize,
			MinIdleConns:    conf.MinIdleConns,
			PoolTimeout:     time.Duration(conf.PoolTimeout) * time.Second,
			ConnMaxIdleTime: time.Duration(conf.ConnMaxIdleTime) * time.Second,
		})
	} else {
		if len(conf.Addrs) > 1 { //	集群模式
			client = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:           conf.Addrs,
				ReadOnly:        true,
				RouteRandomly:   true,
				PoolSize:        conf.PoolSize,
				MinIdleConns:    conf.MinIdleConns,
				Username:        username,
				Password:        password,
				PoolTimeout:     time.Duration(conf.PoolTimeout) * time.Second,
				ConnMaxIdleTime: time.Duration(conf.ConnMaxIdleTime) * time.Second,
			})
		} else {
			client = redis.NewClient(&redis.Options{
				Addr:            addr,
				Username:        username,
				Password:        password,
				DB:              db,
				PoolSize:        conf.PoolSize,
				MinIdleConns:    conf.MinIdleConns,
				PoolTimeout:     time.Duration(conf.PoolTimeout) * time.Second,
				ConnMaxIdleTime: time.Duration(conf.ConnMaxIdleTime) * time.Second,
			})
		}
	}

	conn = client
	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(client); err != nil {
		logrus.WithError(err).Fatal("failed to instrument redis tracing")
	}

	if config.UptraceEnabled {
		//https://redis.uptrace.dev/guide/go-redis-monitoring.html#opentelemetry-instrumentation
		// Enable metrics instrumentation.
		if err := redisotel.InstrumentMetrics(client); err != nil {
			logrus.WithError(err).Fatal("failed to instrument redis metrics")
		}
	}
	//_, err := conn.Ping(ctx).Result()
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func GetConn() redis.Cmdable {
	return conn
}

func GetClient() redis.UniversalClient {
	return client
}

func Close() {
	if client != nil {
		client.Close()
	}
}

func Del(key string) error {
	return conn.Del(ctx, key).Err()
}

func Set(key string, value interface{}) error {
	return conn.Set(ctx, key, value, 0).Err()
}

func SetEx(key string, value interface{}, expiration time.Duration) error {
	return conn.SetEx(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return conn.Get(ctx, key).Result()
}

func HGet(key, childKey string) (string, error) {
	return conn.HGet(ctx, key, childKey).Result()
}

func HGetAll(key string) (map[string]string, error) {
	return conn.HGetAll(ctx, key).Result()
}

func HMGetAll(key string, fields ...string) ([]string, error) {
	result, err := conn.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}
	var out []string
	for i := range result {
		out = append(out, result[i].(string))
	}
	return out, nil
}

func HSet(key string, value map[string]interface{}) error {
	return HSetEx(key, value, 0)
}

func HSetEx(key string, value map[string]interface{}, seconds uint) error {
	var data []interface{}
	for k, v := range value {
		data = append(data, k)
		data = append(data, v)
	}
	pipe := conn.Pipeline()
	pipe.HSet(ctx, key, data...)
	if seconds > 0 {
		pipe.Expire(ctx, key, time.Duration(seconds)*time.Second)
	}
	_, err := pipe.Exec(ctx)
	return err
}

type CacheOption struct {
	StatsEnabled   bool
	Marshal        cache.MarshalFunc
	Unmarshal      cache.UnmarshalFunc
	LocalCache     bool
	LocalCacheSize int
	LocalCacheTTL  time.Duration
}

func NewCache(opt *CacheOption) *cache.Cache {
	option := &cache.Options{
		Redis:        conn,
		StatsEnabled: opt.StatsEnabled,
		Marshal:      opt.Marshal,
		Unmarshal:    opt.Unmarshal,
	}
	if opt.LocalCache {
		ttl := opt.LocalCacheTTL
		if ttl == 0 {
			ttl = time.Minute
		}
		size := opt.LocalCacheSize
		if size <= 0 {
			size = 1000
		}
		option.LocalCache = cache.NewTinyLFU(size, ttl)
	}
	return cache.New(option)
}
