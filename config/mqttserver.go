package config

type MQTTServerConfig struct {
	TCPPort      int      `yaml:"tcpPort,omitempty"` //:1883
	WSPort       int      `yaml:"wsPort,omitempty"`
	TLS          bool     `yaml:"tls"`
	DNSNames     []string `yaml:"dnsNames,omitempty"`
	IPAddresses  []string `yaml:"ipAddresses,omitempty"`
	Persist      bool     `yaml:"persist"`
	ExporterPort int      `yaml:"exporterPort,omitempty"` // 7777
	//Users        auth.Users     `yaml:"users,omitempty"`
	//Auth         auth.AuthRules `yaml:"auth,omitempty"`
	//ACL          auth.ACLRules  `yaml:"acl,omitempty"`
}

var defaultMQTTServerConf = MQTTServerConfig{
	TCPPort: 1883,
	WSPort:  1882,
	//Users: auth.Users{
	//	"iotopo": {
	//		Username: "iotopo",
	//		Password: "iotopo",
	//	},
	//},
	//// 没有匹配到用户时再检查 auth rules
	//Auth: auth.AuthRules{
	//	{
	//		Allow:  true,
	//		Remote: "127.0.0.1:*",
	//	},
	//	{
	//		Allow:  true,
	//		Remote: "localhost:*",
	//	},
	//},
}
