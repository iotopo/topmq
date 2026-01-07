package config

type Metrics struct {
	Enabled   bool   `json:"enabled" yaml:"enabled,omitempty"`
	Type      string `json:"type" yaml:"type,omitempty"`
	Address   string `json:"address" yaml:"address,omitempty"`
	Database  string `json:"database" yaml:"database,omitempty"`
	Username  string `json:"username" yaml:"username,omitempty"`
	Password  string `json:"password" yaml:"password,omitempty"`
	Token     string `json:"token" yaml:"token,omitempty"`
	TimeZone  string `yaml:"timezone,omitempty"`                   // Local
	Retention int    `json:"retention" yaml:"retention,omitempty"` // 日志保存天数
}
