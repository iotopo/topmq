package config

type LoggerConfig struct {
	FileLog bool   `yaml:"fileLog"`
	Level   string `yaml:"level"`
}

var defaultLoggerConf = LoggerConfig{
	FileLog: false,
	Level:   "info",
}
