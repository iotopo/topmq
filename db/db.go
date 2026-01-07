package db

import (
	"database/sql"
	"github.com/glebarez/sqlite"
	"github.com/iotopo/topmq/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

var DB *gorm.DB

const (
	DBTypePG        = "postgres"
	DBTypeMySQL     = "mysql"
	DBTypeDM        = "dm"
	DBTypeSQLite    = "sqlite"
	DBTypeGauss     = "gauss"
	DBTypeOceanbase = "oceanbase"
)

func initConnectionPool(pool *sql.DB) {
	conf := config.Conf.DB
	// 设置连接池参数
	if conf.Pool.MaxOpenConns == 0 {
		conf.Pool.MaxOpenConns = 10
	}
	if conf.Pool.MaxIdleConns == 0 {
		conf.Pool.MaxIdleConns = 5
	}
	if conf.Pool.ConnMaxLifetime == 0 {
		conf.Pool.ConnMaxLifetime = 60
	}
	if conf.Pool.ConnMaxIdleTime == 0 {
		conf.Pool.ConnMaxIdleTime = 30
	}
	pool.SetMaxOpenConns(conf.Pool.MaxOpenConns)                                    // 最大打开连接数
	pool.SetMaxIdleConns(conf.Pool.MaxIdleConns)                                    // 最大空闲连接数
	pool.SetConnMaxLifetime(time.Duration(conf.Pool.ConnMaxLifetime) * time.Minute) // 连接的最大生命周期
	pool.SetConnMaxIdleTime(time.Duration(conf.Pool.ConnMaxIdleTime) * time.Minute) // 连接的最大空闲时间
}

func initLogger() {
	loggerConf := gormLogger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  gormLogger.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	}
	if config.Conf.DB.ShowSQL {
		loggerConf.LogLevel = gormLogger.Info
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(logrus.StandardLogger().Formatter)
	// 所有 sql 语句输出到 info 级别 logrus 中
	gormLogger.Default = gormLogger.New(log.New(logger.Writer(), "\r\n", log.LstdFlags), loggerConf)
}

func setupSQLite() *gorm.DB {
	conf := config.Conf.DB
	databaseName := conf.DatabaseName
	databasePath := conf.DatabasePath
	if databasePath == "" {
		databasePath = "./db"
	}

	fullDatabasePath, err := filepath.Abs(databasePath)
	if err != nil {
		logrus.WithError(err).Fatalf("failed resolving sqlite database path (%s)", databasePath)
	}

	// "unable to open database file: out of memory (14)" when parent directory does not exist
	// https://gitlab.com/cznic/sqlite/-/issues/102
	if err = os.MkdirAll(fullDatabasePath, 0755); err != nil {
		logrus.WithError(err).Fatalf("failed creating sqlite database path (%s)", fullDatabasePath)
	}

	dsn := path.Join(databasePath, databaseName+".db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(50000)&_pragma=read_uncommitted(true)")
	logrus.Infof("connecting to sqlite: %s", dsn)
	gormConf := &gorm.Config{
		CreateBatchSize: 20,
	}

	db, err := gorm.Open(sqlite.Open(dsn), gormConf)
	if err != nil {
		logrus.WithError(err).Fatal("failed creating sqlite database dialector")
	}
	//if sqlDB, err := db.DB(); err == nil {
	//	sqlDB.SetMaxOpenConns(1)
	//}
	return db
}
func Init() {
	initLogger()

	conf := config.Conf.DB

	conn := setupSQLite()
	if conf.ShowSQL {
		conn = conn.Debug()
	}

	//if config.UptraceEnabled {
	//	//https://github.com/uptrace/opentelemetry-go-extra/blob/main/otelgorm/example/main.go
	//	if err = conn.Use(otelgorm.NewPlugin(otelgorm.WithDBName("topstack"))); err != nil {
	//		err = fmt.Errorf("failed to setup otelgorm: %w", err)
	//		return
	//	}
	//}
	// NOTE: 一些启动步骤依然使用全局变量 DB 而不是依赖注入的值
	// 因此这里设置全局变量位于 lifecycle 之外，而不是设为 StartHook
	DB = conn
	initMigrator(conn)
	return
}

func Close() {
	logrus.Info("close database")
	// 测试发现在程序重启时 sqlite 似乎并不积极地合并 shm 和 wal 文件到 db 文件
	// 在程序退出主动关闭一下数据库连接可能更好
	if DB != nil {
		if conn, err := DB.DB(); err == nil {
			err = conn.Close()
			if err != nil {
				logrus.WithError(err).Error("failed to close database")
			}
		} else {
			logrus.WithError(err).Error("failed to close database")
		}
	}
}
