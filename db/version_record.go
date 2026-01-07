package db

import (
	"time"
)

type VersionRecord struct {
	Version   string
	CreatedAt time.Time
}
