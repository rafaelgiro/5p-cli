package champion

import (
	"fmt"
	"regexp"
	"strings"
)

func HandleTooltip(ttp string, spl SpellDataResource) (string, error) {
	c := initialCleanup(ttp)

	dataValues(&c, spl.DataValues)
	effectAmount(&c, spl.EffectAmount)
	cooldown(&c, spl)

	f := finalCleanup(c)
	return f, nil
}

func initialCleanup(input string) string {
	reSpecial := regexp.MustCompile(`@[^@]*?(?:Postfix|Prefix)@`)
	result := reSpecial.ReplaceAllString(input, "")
	return result
}

func finalCleanup(input string) string {
	re := regexp.MustCompile(`@(\d+(\.\d+)?(?:/\d+(\.\d+)?)*)@`)
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		numStr := match[1 : len(match)-1]

		return numStr
	})

	return result
}

func dataValues(ttp *string, spl []SpellDataValue) {
	for _, val := range spl {
		for i, item := range val.Values {
			// Handle Scaling values on strings
			old := fmt.Sprintf("%s%d", val.Name, i+1)
			new := fmt.Sprint(item)
			n := strings.Replace(*ttp, old, new, -1)

			// Handle Single value on strings
			// old = fmt.Sprintf("@%s@", val.Name)
			// strValues := make([]string, len(val.Values))
			// for i, v := range val.Values {
			// 	strValues[i] = fmt.Sprint(v)
			// }
			// new = strings.Join(strValues, "/")
			// n = strings.Replace(n, old, fmt.Sprintf("@%s@", new), -1)

			*ttp = n
		}
	}
}

func effectAmount(ttp *string, spl []SpellEffectAmount) {
	for ei, val := range spl {
		for i, item := range val.Value {
			// Handle Scaling values on strings
			old := fmt.Sprintf("@Effect%dAmount%d@", ei+1, i+1)
			new := fmt.Sprint(item)
			n := strings.Replace(*ttp, old, fmt.Sprintf("@%s@", new), -1)

			// Handle Single value on strings
			old = fmt.Sprintf("@Effect%dAmount@", ei+1)
			strValues := make([]string, len(val.Value))
			for i, v := range val.Value {
				strValues[i] = fmt.Sprint(v)
			}
			new = strings.Join(strValues, "/")
			n = strings.Replace(n, old, fmt.Sprintf("@%s@", new), -1)

			*ttp = n
		}
	}
}

func cooldown(ttp *string, spl SpellDataResource) {
	for i, cd := range spl.CooldownTime {
		old := fmt.Sprintf("Cooldown%d", i)
		new := fmt.Sprint(cd)
		n := strings.Replace(*ttp, old, new, -1)

		*ttp = n
	}
}
