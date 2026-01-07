package config

import (
	"strings"
)

type WebConfig struct {
	Debug          bool     `json:"debug"`
	Port           uint16   `yaml:"port"`
	AllowOrigins   []string `yaml:"allowOrigins,omitempty"`
	TLS            bool     `yaml:"tls"`
	AutoCert       bool     `yaml:"autoCert,omitempty"`
	SessionTimeout int      `yaml:"sessionTimeout,omitempty"` // 分钟
	HostPolicy     []string `yaml:"hostPolicy,omitempty"`     // 域名白名单
	DisableCache   bool     `yaml:"staticCache,omitempty"`    // 是否禁用静态文件缓存
	CheckIPChange  bool     `yaml:"checkIPChange,omitempty"`  // 是否检查用户IP变更
	PProf          bool     `yaml:"pprof,omitempty"`

	RateDuration int `yaml:"duration,omitempty"`  // 流量控制时长(分钟)
	RateLimit    int `yaml:"rateLimit,omitempty"` // 流量控制访问次数限制

	// 一些场景（发送通知、回调）需要将服务器的接口地址发送给第三方
	BaseUrl string `yaml:"baseUrl,omitempty"`

	// gin 需要将代理网关的 IP 加到白名单中，才会尝试读取 X-Forwarded-For, X-Read-IP，以防客户端伪造这两个请求头
	// 如果不设置，gin 会直接将代理网关的 IP 视为客户端 IP
	// 这将导致操作记录中的 IP 错误，并导致三次登录失败即显示验证码的次数被所有客户端共享
	TrustedProxies []string `yaml:"trustedProxies,omitempty"`

	// ClientMaxBodySize 最大文件上传大小
	ClientMaxBodySize int64 `yaml:"clientMaxBodySize,omitempty"` // 50M

	UsernameBlacklist []string `yaml:"usernameBlacklist,omitempty"`
}

func (conf *WebConfig) validate() {
	if len(conf.UsernameBlacklist) > 0 {
		UsernameBlacklist = make(map[string]bool, len(conf.UsernameBlacklist))
		for _, username := range conf.UsernameBlacklist {
			UsernameBlacklist[strings.TrimSpace(username)] = true
		}
	}
	if conf.ClientMaxBodySize <= 0 {
		conf.ClientMaxBodySize = 50
	}
}

var defaultWebConf = WebConfig{
	Port:              8080,
	ClientMaxBodySize: 50,
	SessionTimeout:    60,
}
