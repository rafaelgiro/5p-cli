package ttp

import (
	"fmt"
	"strings"
)

type DataVal float64
type DataVals []DataVal
type SpellValues []SpellDataValue

type SpellDataValue struct {
	Type   string   `mapstructure:"__type"`
	Name   string   `mapstructure:"mName"`
	Values DataVals `mapstructure:"mValues"`
}

func (spl SpellValues) toTooltip(ttp *Tooltip) {
	for _, val := range spl {
		val.dataValue(ttp)
	}
}

func (dv SpellDataValue) dataValue(ttp *Tooltip) {
	dv.Values.combine(ttp, dv.Name)

	for i, item := range dv.Values {
		item.scaling(ttp, dv.Name, i)
	}
}

// Handle Scaling values on strings @Name1@
func (val DataVal) scaling(ttp *Tooltip, name string, i int) {
	old := fmt.Sprintf("@%s%d", name, i)
	new := fmt.Sprint(val)
	n := strings.Replace(string(*ttp), old, fmt.Sprintf("@%s", new), -1)
	*ttp = Tooltip(n)
}

// Handle Single value on strings @Name@
func (val DataVals) combine(ttp *Tooltip, name string) {
	old := fmt.Sprintf("@%s@", name)
	new := val.toString(false)
	n := strings.Replace(string(*ttp), old, new, -1)

	// weird multiplications values on strings @Name*100@
	old = fmt.Sprintf("@%s*", name)
	n = strings.Replace(n, old, fmt.Sprintf("@%s*", new), -1)
	*ttp = Tooltip(n)
}

func (val DataVals) toString(percentage bool) string {
	if len(val) == 0 {
		return ""
	}

	firstValue := val[0]
	allSame := true
	for _, v := range val {
		if v != firstValue {
			allSame = false
			break
		}
	}

	if allSame {
		if percentage {
			return fmt.Sprintf("%.1f%%", firstValue*100)
		} else {
			return fmt.Sprint(firstValue)
		}
	}

	strValues := make([]string, len(val))
	for i, v := range val {
		if percentage {
			strValues[i] = fmt.Sprintf("%.1f%%", v*100)
		} else {
			strValues[i] = fmt.Sprint(v)
		}
	}

	return strings.Join(strValues, "/")
}
