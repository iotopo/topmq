package db

import (
	"database/sql"
	"github.com/coreos/go-semver/semver"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
)

//type DatabaseMigrator interface {
//	gorm.Migrator
//
//	GetVersion(moduleName string) semver.Version
//	SetVersion(moduleName string, version semver.Version)
//
//	MigrateModel(model interface{})
//	CreateTableIfNotExist(dst interface{})
//
//	DB() *gorm.DB
//	SqlDB() *sql.DB
//}

type versionMigrator struct {
	gorm.Migrator

	db       *gorm.DB
	versions *semver.Version
}

var Migrator *versionMigrator

func initMigrator(conn *gorm.DB) {
	Migrator = &versionMigrator{Migrator: conn.Migrator(), db: conn, versions: semver.New("0.0.0")}

	var VersionRecords = new(VersionRecord)
	if !conn.Migrator().HasTable(VersionRecords) {
		if err := conn.Migrator().CreateTable(VersionRecords); err != nil {
			logrus.Fatal(err)
		}
		return
	}

	var versionRecord VersionRecord
	err := conn.Model(VersionRecords).Limit(1).Find(&versionRecord).Error
	if err != nil {
		logrus.Fatal(err)
	}
	//versions := make(map[string]semver.Version, len(versionRecords))
	//for _, record := range versionRecords {
	//	versions[record.ModuleName] = *semver.New(record.Version)
	//}
	if versionRecord.Version != "" {
		Migrator.versions = semver.New(versionRecord.Version)
	}
}

func (s *versionMigrator) GetVersion() *semver.Version {
	return s.versions
}

func (s *versionMigrator) SetVersion(version semver.Version) {
	err := s.db.Model(VersionRecord{}).Where("1=1").Update("version", version.String()).Error
	if err != nil {
		logrus.Fatal(err)
	}
	s.versions = &version
}

// 同 AutoMigrate，但失败时会 logrus.Fatal
func (s *versionMigrator) MigrateModel(model interface{}) {
	logrus.Infof("MigrateModel %v", reflect.TypeOf(model).Name())
	if err := s.db.AutoMigrate(model); err != nil {
		logrus.Fatal(err)
	}
}

func (s *versionMigrator) CreateTableIfNotExist(dst interface{}) {
	logrus.Infof("CreateTableIfNotExist %v", reflect.TypeOf(dst).Name())
	if !s.HasTable(dst) {
		if err := s.CreateTable(dst); err != nil {
			logrus.Fatal(err.Error())
		}
	}
}

func (s *versionMigrator) DB() *gorm.DB {
	return s.db
}

func (s *versionMigrator) SqlDB() *sql.DB {
	db, err := s.db.DB()
	if err != nil {
		logrus.Fatal(err)
	}
	return db
}
