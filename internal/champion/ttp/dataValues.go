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

func (spl SpellValues) toTooltip(ttp *string) {
	for _, val := range spl {
		val.dataValue(ttp)
	}
}

func (dv SpellDataValue) dataValue(ttp *string) {
	dv.Values.combine(ttp, dv.Name)

	for i, item := range dv.Values {
		item.scaling(ttp, dv.Name, i)
	}
}

// Handle Scaling values on strings @Name1@
func (val DataVal) scaling(ttp *string, name string, i int) {
	old := fmt.Sprintf("%s%d", name, i)
	new := fmt.Sprint(val)
	n := strings.Replace(*ttp, old, new, -1)
	*ttp = n
}

// Handle Single value on strings @Name@
func (val DataVals) combine(ttp *string, name string) {
	old := fmt.Sprintf("@%s@", name)
	strValues := make([]string, len(val))

	for j, v := range val {
		strValues[j] = fmt.Sprint(v)
	}

	new := strings.Join(strValues, "/")
	n := strings.Replace(*ttp, old, new, -1)

	// weird multiplications values on strings @Name*100@
	old = fmt.Sprintf("@%s*", name)
	n = strings.Replace(n, old, fmt.Sprintf("@%s*", new), -1)
	*ttp = n
}
