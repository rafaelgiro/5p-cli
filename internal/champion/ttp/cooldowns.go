package ttp

import (
	"fmt"
	"strings"
)

type Cooldowns []Cooldown
type Cooldown float64

// Handle @Cooldown@ and @Cooldown1@ strings
func (c Cooldowns) toTooltip(ttp *Tooltip) {
	for i, cd := range c {
		old := fmt.Sprintf("@Cooldown%d@", i)
		new := fmt.Sprint(cd)
		n := strings.Replace(ttp.ToString(), old, new, -1)
		*ttp = Tooltip(n)
	}

	n := strings.Replace(ttp.ToString(), "@Cooldown@", c.toString(), -1)
	*ttp = Tooltip(n)
}

func (c Cooldowns) toString() string {
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
