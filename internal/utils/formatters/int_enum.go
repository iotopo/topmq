package formatters

import (
	"strconv"
)

// 为一个数值枚举提供字符串格式和翻译
type IntEnum struct {
	Value  int
	Labels []string
}

func (n IntEnum) String() string {
	if len(n.Labels) > n.Value {
		return n.Labels[n.Value]
	}
	return strconv.Itoa(n.Value)
}

func (n IntEnum) Translate(translatedLabels ...string) string {
	if len(translatedLabels) > n.Value {
		return translatedLabels[n.Value]
	}
	return n.String()
}

// 为一个布尔提供字符串格式和翻译
type Bool struct {
	Value  bool
	Labels []string
}

func (b Bool) String() string {
	if b.Value {
		if len(b.Labels) > 1 {
			return b.Labels[1]
		}
	} else {
		if len(b.Labels) > 0 {
			return b.Labels[0]
		}
	}
	return strconv.FormatBool(b.Value)
}

func (b Bool) Translate(translatedLabels ...string) string {
	if b.Value {
		if len(translatedLabels) > 1 {
			return translatedLabels[1]
		}
	} else {
		if len(translatedLabels) > 0 {
			return translatedLabels[0]
		}
	}
	return b.String()
}
