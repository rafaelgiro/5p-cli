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
	spellCalculations(&c, spl)
	cooldown(&c, spl)

	f := finalCleanup(c)
	return f, nil
}

func initialCleanup(input string) string {
	reSpecial := regexp.MustCompile(`@[^@]*?(?:Postfix|Prefix)@`)
	result := reSpecial.ReplaceAllString(input, "")
	result = strings.ReplaceAll(result, "*100.000000", "")
	result = strings.ReplaceAll(result, "*100", "")
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
			old := fmt.Sprintf("%s%d", val.Name, i)
			new := fmt.Sprint(item)
			n := strings.Replace(*ttp, old, new, -1)

			// Handle Single value on strings
			old = fmt.Sprintf("@%s@", val.Name)
			strValues := make([]string, len(val.Values))
			for j, v := range val.Values {
				strValues[j] = fmt.Sprint(v)
			}
			new = strings.Join(strValues, "/")
			n = strings.Replace(n, old, new, -1)

			*ttp = n
		}
	}
}

func effectAmount(ttp *string, spl []SpellEffectAmount) {
	for ei, val := range spl {
		for i, item := range val.Value {
			// Handle Scaling values on strings
			old := fmt.Sprintf("@Effect%dAmount%d@", ei+1, i)
			new := fmt.Sprint(item)
			n := strings.Replace(*ttp, old, new, -1)

			// Handle Single value on strings
			old = fmt.Sprintf("@Effect%dAmount@", ei+1)
			strValues := make([]string, len(val.Value))
			for i, v := range val.Value {
				strValues[i] = fmt.Sprint(v)
			}
			new = strings.Join(strValues, "/")
			n = strings.Replace(n, old, new, -1)

			*ttp = n
		}
	}
}

func cooldown(ttp *string, spl SpellDataResource) {
	cds := []string{}

	for i, cd := range spl.CooldownTime {
		old := fmt.Sprintf("Cooldown%d", i)
		new := fmt.Sprint(cd)
		n := strings.Replace(*ttp, old, new, -1)
		cds = append(cds, new)
		*ttp = n
	}

	n := strings.Replace(*ttp, "@Cooldown@", strings.Join(cds, "/"), -1)
	*ttp = n
}

func spellCalculations(ttp *string, spl SpellDataResource) {
	for key, val := range spl.SpellCalculations {
		for _, fp := range val.FormulaParts {
			if len(fp.DataValue) != 0 {
				dataValCalc(ttp, key, fp.DataValue, spl.DataValues)
			} else if fp.Breakpoints != nil {
				breakpoints(ttp, key, fp)
			}
		}
	}
}

func dataValCalc(ttp *string, ttpKey, dvKey string, spl []SpellDataValue) {
	for _, val := range spl {
		if dvKey == val.Name {
			strValues := make([]string, len(val.Values))
			for i, v := range val.Values {
				strValues[i] = fmt.Sprint(v)
			}

			str := strings.Join(strValues, "/")
			n := strings.Replace(*ttp, fmt.Sprintf("@%s@", ttpKey), str, -1)
			*ttp = n
		}
	}
}

func breakpoints(ttp *string, ttpKey string, fp FormulaPart) {
	base := fp.Level1Value
	vals := []float64{}
	lvs := []float64{}

	for i, bp := range fp.Breakpoints {
		lvs = append(lvs, float64(bp.Level))

		if i == 0 {
			vals = append(vals, base+bp.AdditionalBonusAtThisLevel)
		} else {
			vals = append(vals, vals[i-1]+bp.AdditionalBonusAtThisLevel)
		}
	}

	str := fmt.Sprintf("%s (at level %s)", arrayToString(vals, "/"), arrayToString(lvs, "/"))
	n := strings.Replace(*ttp, fmt.Sprintf("@%s@", ttpKey), str, -1)
	*ttp = n
}

func arrayToString(a []float64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
