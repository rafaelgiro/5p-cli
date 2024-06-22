package ttp

import (
	"fmt"
	"strings"
)

type Costs []Cost
type Cost float64

// Handle @Cost@ and @Cost1@ strings
func (c Costs) toTooltip(ttp *Tooltip) {
	costs := []string{}

	for i, cost := range c {
		old := fmt.Sprintf("Cost%d", i)
		new := fmt.Sprint(cost)
		n := strings.Replace(ttp.ToString(), old, new, -1)
		costs = append(costs, new)
		*ttp = Tooltip(n)
	}

	n := strings.Replace(ttp.ToString(), "@Cost@", strings.Join(costs, "/"), -1)
	*ttp = Tooltip(n)
}
