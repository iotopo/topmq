package datatypes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type TypedJSON[T any] struct {
	Value T
}

// Scan scan value into TypedJSON[T], implements sql.Scanner interface
func (j *TypedJSON[T]) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case nil:
		break
	case []byte:
		if len(v) > 0 {
			// Reference types such as []byte are only valid until the next call to Scan
			// and should not be retained. Their underlying memory is owned by the driver.
			// If retention is necessary, copy their values before the next call to Scan.
			bytes = make([]byte, len(v))
			copy(bytes, v)
		}
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("typed json must scan from string or bytes value")
	}

	return json.Unmarshal(bytes, &j.Value)
}

// MarshalJSON to output non base64 encoded []byte
func (j TypedJSON[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Value)
}

// UnmarshalJSON to deserialize []byte
func (j *TypedJSON[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &j.Value)
}

// GormDataType gorm common data type
func (TypedJSON[T]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (TypedJSON[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	case "dm":
		return DMJsonType
	}
	return "BLOB"
}

func (j TypedJSON[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := j.MarshalJSON()

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

func (j *TypedJSON[T]) UnwrapPointer() *T {
	if j == nil {
		return nil
	}
	return &j.Value
}
