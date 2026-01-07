package migrator

import (
	"github.com/coreos/go-semver/semver"
	"github.com/iotopo/topmq/accounts"
	"github.com/iotopo/topmq/acls"
	"github.com/iotopo/topmq/auth"
	"github.com/iotopo/topmq/black_list"
	"github.com/iotopo/topmq/db"
)

var version = *semver.New("1.0.0")

func DatabaseMigrate() {
	migrator := db.Migrator
	currentVersion := migrator.GetVersion()

	switch {
	case !currentVersion.LessThan(version):
		return
	case currentVersion.Equal(*semver.New("0.0.0")):
		updateLatest()
		break
	default:
		break
	}

	migrator.SetVersion(version)
}

func updateLatest() {
	m := db.Migrator
	m.CreateTableIfNotExist(auth.Users)
	m.CreateTableIfNotExist(accounts.Accounts)
	m.CreateTableIfNotExist(black_list.BlackLists)
	m.CreateTableIfNotExist(acls.AccessControls)
}
