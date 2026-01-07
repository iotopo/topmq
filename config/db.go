package config

type DBConfig struct {
	ShowSQL            bool   `yaml:"showSQL"`
	DatabaseName       string `yaml:"databaseName"`
	DatabasePath       string `yaml:"databasePath"`
	AutoCreateDatabase bool   `yaml:"autoCreateDatabase,omitempty"`
	AutoCreateAdmin    bool   `yaml:"autoCreateAdmin"`

	// 连接池配置
	Pool PoolConfig `yaml:"pool,omitempty"`

	SQLite *SQLiteConfig `yaml:"sqlite,omitempty"`
}

type PoolConfig struct {
	MaxOpenConns    int `yaml:"maxOpenConns,omitempty"`    // 最大打开连接数
	MaxIdleConns    int `yaml:"maxIdleConns,omitempty"`    // 最大空闲连接数
	ConnMaxLifetime int `yaml:"connMaxLifetime,omitempty"` // 连接的最大生命周期（分钟）
	ConnMaxIdleTime int `yaml:"connMaxIdleTime,omitempty"` // 连接的最大空闲时间（分钟）
}
type SQLiteConfig struct {
}

var defaultDBConf = DBConfig{
	DatabaseName: "topmp",
	DatabasePath: "./db",
	Pool: PoolConfig{
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 60, // 1小时
		ConnMaxIdleTime: 30, // 30分钟
	},
}
