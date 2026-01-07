package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestData struct {
	Name   string
	Exp    string
	Points []string
	Attrs  []string
}

func TestSplitProp(t *testing.T) {
	datas := []TestData{
		{
			Exp:    "(data.flow_total+3)/attr.pressure==data.flow_rate",
			Points: []string{"flow_total", "flow_rate"},
			Attrs:  []string{"pressure"},
		},
		{
			Exp:    "data.flow_total>99&&data.flow_rate<20",
			Points: []string{"flow_total", "flow_rate"},
		},
		{
			Exp:    "ABS(data.flow_total+3)<attr.pressure",
			Points: []string{"flow_total"},
			Attrs:  []string{"pressure"},
		},
		{
			Exp:    `ABS(data.flow_total+3)>0&&attr.pressure=="test"`,
			Points: []string{"flow_total"},
			Attrs:  []string{"pressure"},
		},
		{
			Exp:    `(data.flow_total+3>0 || data.flow_total<-1) && attr.pressure=="test"`,
			Points: []string{"flow_total"},
			Attrs:  []string{"pressure"},
		},
		{
			Exp:    `(data["flow_total"]+3>0 || data["flow_total"]<-1) && attr["pressure"]=="test"`,
			Points: []string{"flow_total"},
			Attrs:  []string{"pressure"},
		},
		{
			Exp:    `(data["flow_total"]+3>0 || data["0-flow_total"]<-1) && attr["pressure"]=="test"`,
			Points: []string{"flow_total", "0-flow_total"},
			Attrs:  []string{"pressure"},
		},
	}
	for _, data := range datas {
		t.Log("exp:", data.Exp)
		points := SplitProps("data", data.Exp)
		attrs := SplitProps("attr", data.Exp)

		if len(points) == 0 && len(attrs) == 0 {
			t.Error("split error")
		}
		t.Log(points, attrs)

		assert.ElementsMatch(t, points, data.Points)
		assert.ElementsMatch(t, attrs, data.Attrs)
	}
}
