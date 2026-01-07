package datatypes

import (
	"database/sql/driver"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// decimal.NullDecimal 不能被 gorm 用于 migrate
type NullDecimal struct {
	Decimal decimal.Decimal
	Valid   bool

	// 如果为 0，则 MarshalJSON 为字符串
	// 如果为 -1，则 MarshalJSON 为数值
	// 如果为 1，则 MarshalJSON 为字符串
	//     且 MarshalJSON 和 String() 会固定小数位数为 MarshalFixedPoint
	//
	// 注意：设为 -1 时，由于 Decimal 值的精度可能超过 IEEE754 支持的范围
	// 导致反序列化失败或丢失精度
	//
	// 所有返回 NullDecimal 的方法，都保留 MarshalMethod 和 MarshalFixedPoint 不变
	MarshalMethod     int8
	MarshalFixedPoint int32
}

func (d NullDecimal) String() string {
	if d.Valid {
		if d.MarshalMethod > 0 {
			return d.Decimal.StringFixed(d.MarshalFixedPoint)
		} else {
			return d.Decimal.String()
		}
	} else {
		return "null"
	}
}

func (d *NullDecimal) Scan(value any) error {
	if value == nil {
		d.Valid = false
		return nil
	}
	d.Valid = true
	return d.Decimal.Scan(value)
}

func (d NullDecimal) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Decimal.Value()
}

func (d NullDecimal) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	var str string
	if d.MarshalMethod < 0 {
		str = d.Decimal.String()
	} else if d.MarshalMethod > 0 {
		str = "\"" + d.Decimal.StringFixed(d.MarshalFixedPoint) + "\""
	} else {
		str = "\"" + d.Decimal.String() + "\""
	}
	return []byte(str), nil
}

func (d *NullDecimal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		d.Valid = false
		return nil
	}
	err := d.Decimal.UnmarshalJSON(data)
	d.Valid = err == nil
	return err
}

// GormDataType gorm common data type
func (NullDecimal) GormDataType() string {
	return string(schema.Float)
}

// GormDBDataType gorm db data type
func (NullDecimal) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

// 返回一个 null 值
func (d NullDecimal) WithoutDecimal() NullDecimal {
	return NullDecimal{
		Valid:             false,
		Decimal:           decimal.Decimal{},
		MarshalMethod:     d.MarshalMethod,
		MarshalFixedPoint: d.MarshalFixedPoint,
	}
}

// 返回 val
func (d NullDecimal) WithDecimal(val decimal.Decimal) NullDecimal {
	return NullDecimal{
		Valid:             true,
		Decimal:           val,
		MarshalMethod:     d.MarshalMethod,
		MarshalFixedPoint: d.MarshalFixedPoint,
	}
}

// 返回 val
func (d NullDecimal) With(val NullDecimal) NullDecimal {
	if val.Valid {
		return d.WithDecimal(val.Decimal)
	} else {
		return d.WithoutDecimal()
	}
}

// 求和，除非都是 null，否则 null 被视为 0
func (d NullDecimal) Add(rhs NullDecimal) NullDecimal {
	if d.Valid {
		if rhs.Valid {
			return d.WithDecimal(d.Decimal.Add(rhs.Decimal))
		} else {
			return d
		}
	} else {
		return d.With(rhs)
	}
}

// 和 Add 相同，但可以对多个值求和
func (d NullDecimal) Sum(xs ...NullDecimal) NullDecimal {
	val := d
	for _, x := range xs {
		val = val.Add(x)
	}
	return val
}

// 求和，如有 null 则值不变
func (d NullDecimal) MaybeAdd(rhs NullDecimal) NullDecimal {
	if d.Valid && rhs.Valid {
		return d.WithDecimal(d.Decimal.Add(rhs.Decimal))
	}
	return d
}

// 求差，如有 null 则值不变
func (d NullDecimal) MaybeSub(rhs NullDecimal) NullDecimal {
	if d.Valid && rhs.Valid {
		return d.WithDecimal(d.Decimal.Sub(rhs.Decimal))
	}
	return d
}

// 求积，除非都是 null，否则 null 被视为 1
func (d NullDecimal) Mul(rhs NullDecimal) NullDecimal {
	if d.Valid {
		if rhs.Valid {
			return d.WithDecimal(rhs.Decimal.Mul(rhs.Decimal))
		} else {
			return d
		}
	} else {
		return d.With(rhs)
	}
}

// 求积，如有 null 则值不变
func (d NullDecimal) MaybeMul(rhs NullDecimal) NullDecimal {
	if d.Valid && rhs.Valid {
		return d.WithDecimal(d.Decimal.Mul(rhs.Decimal))
	}
	return d
}

// 求商，如有 null 则值不变
// 除 0 会得到 null
func (d NullDecimal) MaybeDiv(rhs NullDecimal) NullDecimal {
	if d.Valid && rhs.Valid {
		if rhs.Decimal.IsZero() {
			return d.WithoutDecimal()
		}
		return d.WithDecimal(d.Decimal.Div(rhs.Decimal))
	}
	return d
}

// places > 0 四舍五入到小数点后 places 位
// places = 0 四舍五入取整
// places = -1 取整到 10 的整数倍
// places = -2 取整到 100 的整数倍
// ......
// places = -N 取整到 10^N 的整数倍
func (d NullDecimal) Round(places int32) NullDecimal {
	if d.Valid {
		return NullDecimal{Valid: true, Decimal: d.Decimal.Round(places)}
	}
	return d
}

func NewNullDecimal(val decimal.Decimal, valid ...bool) NullDecimal {
	if len(valid) > 0 {
		return NullDecimal{Valid: valid[0], Decimal: val}
	}
	return NullDecimal{Valid: true, Decimal: val}
}

// 支持 Go 1.24 的 omitzero 标签
func (d NullDecimal) IsZero() bool {
	return !d.Valid
}
