package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

type RedisSentinel struct {
	// The master name.
	MasterName string `yaml:"masterName,omitempty"`
	// A seed list of host:port addresses of sentinel nodes.
	Addrs []string `yaml:"addrs,omitempty"`

	// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
	ClientName string `yaml:"clientName,omitempty"`

	// If specified with SentinelPassword, enables ACL-based authentication (via
	// AUTH <user> <pass>).
	Username string `yaml:"username,omitempty"`
	// Sentinel password from "requirepass <password>" (if enabled) in Sentinel
	// configuration, or, if SentinelUsername is also supplied, used for ACL-based
	// authentication.
	Password string `yaml:"password,omitempty"`
}

type RedisConfig struct {
	KeyPrefix    string `yaml:"keyPrefix,omitempty"`    // 通用前缀，默认为 "topmq:"
	PasswordHash string `yaml:"passwordHash,omitempty"` // MQTT 扩展认证 hash 算法
	EMQ          bool   `yaml:"emq,omitempty"`          // MQTT 扩展认证，是否兼容 EMQ 格式
	Addr         string `yaml:"addr"`
	Username     string `yaml:"username,omitempty"`
	Password     string `yaml:"password,omitempty"`
	DB           int    `yaml:"db,omitempty"`

	DialTimeout  int `yaml:"dialTimeout,omitempty"`
	ReadTimeout  int `yaml:"readTimeout,omitempty"`
	WriteTimeout int `yaml:"writeTimeout,omitempty"`

	PoolSize        int `yaml:"poolSize,omitempty"`
	PoolTimeout     int `yaml:"poolTimeout,omitempty"` //s
	MinIdleConns    int `yaml:"minIdleConns,omitempty"`
	MaxIdleConns    int `yaml:"maxIdleConns,omitempty"`
	ConnMaxIdleTime int `yaml:"connMaxIdleTime,omitempty"`
	ConnMaxLifetime int `yaml:"connMaxLifetime,omitempty"`

	// A seed list of host:port addresses of cluster nodes.
	Addrs []string `yaml:"addrs,omitempty"`

	Sentinel *RedisSentinel `yaml:"sentinel,omitempty"`
	//// The master name.
	//MasterName string `yaml:"masterName,omitempty"`
	//// A seed list of host:port addresses of sentinel nodes.
	//SentinelAddrs []string `yaml:"sentinelAddrs,omitempty"`
	//
	//// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
	//ClientName string `yaml:"clientName,omitempty"`
	//
	//// If specified with SentinelPassword, enables ACL-based authentication (via
	//// AUTH <user> <pass>).
	//SentinelUsername string `yaml:"sentinelUsername,omitempty"`
	//// Sentinel password from "requirepass <password>" (if enabled) in Sentinel
	//// configuration, or, if SentinelUsername is also supplied, used for ACL-based
	//// authentication.
	//SentinelPassword string `yaml:"sentinelPassword,omitempty"`
}

func (conf *RedisConfig) validate() {
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		if conf.Password != "" {
			logrus.Warnf("Redis: use password from env variable")
		}
		conf.Password = password
	}
}

var defaultRedisConf = RedisConfig{
	Addr: "127.0.0.1:6379",
}
