package ttp

import (
	"fmt"
	"strings"
)

type Costs []Cost
type Cost float64

// Handle @Cost@ and @Cost1@ strings
func (c Costs) toTooltip(ttp *Tooltip) {
	for i, cost := range c {
		old := fmt.Sprintf("@Cost%d@", i)
		new := fmt.Sprint(cost)
		n := strings.Replace(ttp.ToString(), old, new, -1)
		old = fmt.Sprintf("@BaseCost%d@", i)
		n = strings.Replace(n, old, new, -1)
		*ttp = Tooltip(n)
	}

	n := strings.Replace(ttp.ToString(), "@Cost@", c.toString(), -1)
	*ttp = Tooltip(n)
}

func (c Costs) toString() string {
	if len(c) == 0 {
		return ""
	}

	firstValue := c[0]
	allSame := true
	for _, v := range c {
		if v != firstValue {
			allSame = false
			break
		}
	}

	if allSame {
		return fmt.Sprint(firstValue)
	}

	strValues := make([]string, len(c))
	for i, v := range c {
		strValues[i] = fmt.Sprint(v)
	}

	return strings.Join(strValues, "/")
}
