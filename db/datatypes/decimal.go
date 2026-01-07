package datatypes

import (
	"database/sql/driver"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// decimal.Decimal 不能被 gorm 用于 migrate
type Decimal decimal.Decimal

func (d *Decimal) Scan(value any) error {
	return (*decimal.Decimal)(d).Scan(value)
}

func (d Decimal) Value() (driver.Value, error) {
	return (decimal.Decimal)(d).Value()
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return (decimal.Decimal)(d).MarshalJSON()
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	return (*decimal.Decimal)(d).UnmarshalJSON(data)
}

// GormDataType gorm common data type
func (Decimal) GormDataType() string {
	return string(schema.Float)
}

// GormDBDataType gorm db data type
func (Decimal) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		// https://dev.mysql.com/doc/refman/5.7/en/numeric-type-syntax.html
		return "NUMERIC(65,30)"
	case "postgres":
		// https://www.postgresql.org/docs/14/datatype-numeric.html#DATATYPE-NUMERIC-DECIMAL
		// PostgreSQL 如果同时省略 precision 和 scale
		// 可以在允许的范围内（小数点前 131072 位，小数点后 16383 位）保存任意精度
		return "NUMERIC"
	case "sqlite":
		// https://sqlite.org/datatype3.html
		return "NUMERIC"
	case "dm":
		// https://eco.dameng.com/document/dm/zh-cn/pm/dm8_sql-data-types-operators.html
		return "NUMERIC"
	}
	return "NUMERIC(38,10)"
}
