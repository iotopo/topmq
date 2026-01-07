package datatypes

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"
)

// 和 gorm.io/datatypes.Date 类似，但转换为 JSON 是 time.DateOnly 格式
type Date time.Time

func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date(nullTime.Time)
	return
}

func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}

// GormDataType gorm common data type
func (date Date) GormDataType() string {
	return "date"
}

func (date Date) GobEncode() ([]byte, error) {
	return time.Time(date).GobEncode()
}

func (date *Date) GobDecode(b []byte) error {
	return (*time.Time)(date).GobDecode(b)
}

func (date Date) MarshalJSON() (b []byte, err error) {
	b = make([]byte, 0, len(time.DateOnly)+2)
	b = append(b, '"')
	b = time.Time(date).AppendFormat(b, time.DateOnly)
	b = append(b, '"')
	return
}

func (date *Date) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return errors.New("Date.UnmarshalJSON: input is not a JSON string")
	}
	b = b[1 : len(b)-1]
	t, err := time.Parse(time.DateOnly, string(b))
	if err != nil {
		return err
	}
	*date = Date(t)
	return nil
}
