package ttp

import (
	"fmt"
	"strings"
)

type Cooldowns []Cooldown
type Cooldown float64

// Handle @Cooldown@ and @Cooldown1@ strings
func (c Cooldowns) toTooltip(ttp *Tooltip) {
	cds := []string{}

	for i, cd := range c {
		old := fmt.Sprintf("@Cooldown%d@", i)
		new := fmt.Sprint(cd)
		n := strings.Replace(ttp.ToString(), old, new, -1)
		cds = append(cds, new)
		*ttp = Tooltip(n)
	}

	n := strings.Replace(ttp.ToString(), "@Cooldown@", strings.Join(cds, "/"), -1)
	*ttp = Tooltip(n)
}
