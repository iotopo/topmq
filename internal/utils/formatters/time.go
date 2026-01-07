package formatters

import (
	"time"

	"github.com/itchyny/timefmt-go"
	"gopkg.in/guregu/null.v4"
)

// 封装一个 time.Time，使其默认字符串格式为 time.DateTime
// 并使用 strftime(3) 格式化
type Time time.Time

func (t Time) String() string {
	return time.Time(t).Format(time.DateTime)
}

func (t Time) Format(format string) string {
	return timefmt.Format(time.Time(t), format)
}

func (t Time) Year() int {
	return time.Time(t).Year()
}

func (t Time) Month() int {
	// 返回数值，以免 Month 类型默认格式化为英文
	return int(time.Time(t).Month())
}

func (t Time) Day() int {
	return time.Time(t).Day()
}

func (t Time) Hour() int {
	return time.Time(t).Hour()
}

func (t Time) Minute() int {
	return time.Time(t).Minute()
}

func (t Time) Second() int {
	return time.Time(t).Second()
}

func (t Time) Weekday() int {
	return int(time.Time(t).Weekday())
}

// 为 null.Time 增加字符串格式和格式化
type NullTime null.Time

func (t NullTime) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(time.DateTime)
}

func (t NullTime) Format(format string) string {
	if !t.Valid {
		return ""
	}
	return timefmt.Format(t.Time, format)
}

func (t NullTime) Year() any {
	if !t.Valid {
		return ""
	}
	return t.Time.Year()
}

func (t NullTime) Month() any {
	if !t.Valid {
		return ""
	}
	// 返回数值，以免 Month 类型默认格式化为英文
	return int(t.Time.Month())
}

func (t NullTime) Day() any {
	if !t.Valid {
		return ""
	}
	return t.Time.Day()
}

func (t NullTime) Hour() any {
	if !t.Valid {
		return ""
	}
	return t.Time.Hour()
}

func (t NullTime) Minute() any {
	if !t.Valid {
		return ""
	}
	return t.Time.Minute()
}

func (t NullTime) Second() any {
	if !t.Valid {
		return ""
	}
	return t.Time.Second()
}

func (t NullTime) Weekday() any {
	if !t.Valid {
		return ""
	}
	return int(t.Time.Weekday())
}

func NewTime(t time.Time, valid bool) NullTime {
	return NullTime(null.NewTime(t, valid))
}

func TimeFrom(t time.Time) NullTime {
	return NullTime(null.TimeFrom(t))
}

func TimeFromPtr(t *time.Time) NullTime {
	return NullTime(null.TimeFromPtr(t))
}
