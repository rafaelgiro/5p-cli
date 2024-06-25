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
	strValues := make([]string, len(ef.Value))
	for i, v := range ef.Value {
		strValues[i] = fmt.Sprint(v)
	}

	old := fmt.Sprintf("@Effect%dAmount@", ei+1)
	new := strings.Join(strValues, "/")
	n := strings.Replace(ttp.ToString(), old, new, -1)

	// Additional replace to handle multiplication values @Effect1Ammount*100@
	old = fmt.Sprintf("Effect%dAmount", ei+1)
	new = strings.Join(strValues, "/")
	n = strings.Replace(n, old, new, -1)

	*ttp = Tooltip(n)
}

func (ef Effect) toString() string {
	val := ef.Value

	firstValue := val[0]
	allSame := true
	for _, v := range val {
		if v != firstValue {
			allSame = false
			break
		}
	}

	if allSame {
		return fmt.Sprint(firstValue)
	}

	strValues := make([]string, len(val))
	for i, v := range val {
		strValues[i] = fmt.Sprint(v)
	}

	return strings.Join(strValues, "/")
}
