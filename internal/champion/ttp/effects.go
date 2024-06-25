package ttp

import (
	"fmt"
	"strings"
)

type SpellEffectAmount []Effect
type EffectValue []float64

type Effect struct {
	Type  string      `mapstructure:"__type"`
	Value EffectValue `mapstructure:"value"`
}

func (spl SpellEffectAmount) toTooltip(ttp *Tooltip) {
	for i, ef := range spl {
		ef.scaling(ttp, i)
		ef.combine(ttp, i)
	}
}

// Handle Scaling values on strings @Effect1Amount1@
func (ef Effect) scaling(ttp *Tooltip, ei int) {
	for i, item := range ef.Value {
		old := fmt.Sprintf("@Effect%dAmount%d@", ei+1, i)
		new := fmt.Sprint(item)
		n := strings.Replace(ttp.ToString(), old, new, -1)

		// Additional replace to handle multiplication values @Effect1Ammount1*100@
		old = fmt.Sprintf("Effect%dAmount%d", ei+1, i)
		new = fmt.Sprint(item)
		n = strings.Replace(n, old, new, -1)

		*ttp = Tooltip(n)
	}
}

// Handle Single value on strings @Effect1Amount@
func (ef Effect) combine(ttp *Tooltip, ei int) {
	old := fmt.Sprintf("@Effect%dAmount@", ei+1)
	new := ef.toString(1)
	n := strings.Replace(ttp.ToString(), old, new, -1)

	// Additional replace to handle multiplication values @Effect1Ammount*100@
	old = fmt.Sprintf("Effect%dAmount", ei+1)
	n = strings.Replace(n, old, new, -1)

	*ttp = Tooltip(n)
}

func (ef Effect) toString(mult float64) string {
	val := ef.Value

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
		return fmt.Sprint(firstValue * mult)
	}

	strValues := make([]string, len(val))
	for i, v := range val {
		strValues[i] = fmt.Sprint(v * mult)
	}

	return strings.Join(strValues, "/")
}
