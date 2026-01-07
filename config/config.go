package config

import (
	"fmt"
	"github.com/iotopo/topmq/internal/utils"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var Version = "1.0.0"
var ReleaseTime = "2023-12-11"
var AppName = "topmq"

var UptraceDSN string
var UptraceEnabled bool

var UsernameBlacklist map[string]bool

type CommonConfig struct {
	MaxProcs int `yaml:"maxProcs,omitempty"` // 最大协程数，默认为 0 表示不限制
	// 是否显示登录人员水印（默认关闭）
	Watermark bool `yaml:"watermark,omitempty"`

	MinPwdLen int `yaml:"minPwdLen,omitempty"` // 密码最小长度，默认 8

	Logger LoggerConfig `yaml:"logger"`

	MQTTServer MQTTServerConfig `yaml:"mqttServer"`
	Web        WebConfig        `yaml:"web"`
	DB         DBConfig         `yaml:"db"`
	Redis      RedisConfig      `yaml:"redis"`

	Metrics Metrics `yaml:"metrics,omitempty"`
}

var Conf = CommonConfig{
	Logger:     defaultLoggerConf,
	MQTTServer: defaultMQTTServerConf,
	Web:        defaultWebConf,
	DB:         defaultDBConf,
	Redis:      defaultRedisConf,
}

func Save() error {
	data, err := yaml.Marshal(&Conf)
	if err != nil {
		return err
	}
	return os.WriteFile("./config.yml", data, 0666)
}

func workDir() {
	if runtime.GOOS == "windows" {
		file, _ := os.Executable()       // 获取当前可执行文件所在的文件目录
		_ = os.Chdir(filepath.Dir(file)) // [2] 更新当前工作目录
		//} else {
		//	file, _ := exec.LookPath(os.Args[0]) // [1]
		//	_ = os.Chdir(filepath.Dir(file))     // [2] 更新当前工作目录(darwin/linux)
	}
}

func showInfo() {
	fmt.Println("")
	fmt.Println("================================================")
	fmt.Printf("topmq MQTT V5.0 Broker v%s build at %s\n", Version, ReleaseTime)
	fmt.Println("www.iotopo.com")
	fmt.Println("================================================")

	if dir, err := os.Getwd(); err == nil {
		fmt.Println("cwd", dir)
	}
}

func init() {
	if len(os.Args) <= 1 {
		showInfo()
	}

	workingDir := os.Getenv("TOPMQ_WD")
	if workingDir != "" {
		_ = os.Chdir(workingDir)
	} else {
		workDir()
	}

	//if utils.PathExists(".env") {
	//	if err := godotenv.Load(); err != nil {
	//		logrus.WithError(err).Fatal("failed loading .env file")
	//	} else {
	//		logrus.Info(".env file loaded")
	//	}
	//}

	configPath := "./config.yml"
	if utils.PathExists(configPath) {
		data, err := os.ReadFile(configPath)
		if err != nil {
			logrus.WithError(err).Fatal("failed loading config.yml")
		}
		if err := yaml.Unmarshal(data, &Conf); err != nil {
			logrus.WithError(err).Fatal("failed parse config.yml")
		}
	}

	if Conf.MinPwdLen <= 0 {
		Conf.MinPwdLen = 8
	}
	if Conf.Redis.KeyPrefix == "" {
		Conf.Redis.KeyPrefix = "topmq:"
	}

	Conf.Web.validate()

	UptraceDSN = os.Getenv("UPTRACE_DSN")
	if UptraceDSN != "" {
		UptraceEnabled = true
	}
}
