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

type NullTypedJSON[T any] struct {
	JSON  T
	Valid bool
}

// impl database/sql.Scanner
func (j *NullTypedJSON[T]) Scan(value interface{}) error {
	if value == nil {
		j.Valid = false
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		if len(v) > 0 {
			bytes = make([]byte, len(v))
			copy(bytes, v)
		}
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("typed json must scan from string or bytes value")
	}

	if err := json.Unmarshal(bytes, &j.JSON); err != nil {
		return err
	}
	j.Valid = true
	return nil
}

// impl encoding/json.Marshaler
func (j NullTypedJSON[T]) MarshalJSON() ([]byte, error) {
	if !j.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(j.JSON)
}

// impl encoding/json.Unmarshaler
func (j *NullTypedJSON[T]) UnmarshalJSON(data []byte) error {
	switch len(data) {
	case 4:
		if string(data) != "null" {
			break
		}
		fallthrough
	case 0:
		j.Valid = false
		return nil
	default:
		break
	}

	if err := json.Unmarshal(data, &j.JSON); err != nil {
		return err
	}
	j.Valid = true
	return nil
}

// GormDataType gorm common data type
func (NullTypedJSON[T]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (NullTypedJSON[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (j NullTypedJSON[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if !j.Valid {
		return gorm.Expr("?", nil)
	}
	data, _ := j.MarshalJSON()

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

// 支持 Go 1.24 的 omitzero 标签
func (j NullTypedJSON[T]) IsZero() bool {
	return !j.Valid
}
