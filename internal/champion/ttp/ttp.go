package ttp

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Tooltip string

func (ttp *Tooltip) Calc(spl SpellDataResource) error {
	initialCleanup(ttp)

	spl.DataValues.toTooltip(ttp)

	// effectAmount(ttp, spl.EffectAmount)
	// spellCalculations(ttp, spl)
	// cooldown(ttp, spl)
	// cost(ttp, spl)

	// f := finalCleanup(c)
	return nil
}

func (ttp Tooltip) ToString() string {
	tp := string(ttp)
	return tp
}

func initialCleanup(ttp *Tooltip) {
	reSpecial := regexp.MustCompile(`@[^@]*?(?:Postfix|Prefix)@`)
	result := reSpecial.ReplaceAllString(ttp.ToString(), "")
	*ttp = Tooltip(result)
}

func effectAmount(ttp *string, spl []SpellEffectAmount) {
	for ei, val := range spl {
		for i, item := range val.Value {
			// Handle Scaling values on strings
			old := fmt.Sprintf("@Effect%dAmount%d@", ei+1, i)
			new := fmt.Sprint(item)
			n := strings.Replace(*ttp, old, new, -1)

			// Additional replace to handle multiplication values
			old = fmt.Sprintf("Effect%dAmount%d", ei+1, i)
			new = fmt.Sprint(item)
			n = strings.Replace(n, old, new, -1)

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

func cost(ttp *string, spl SpellDataResource) {
	costs := []string{}

	for i, cost := range spl.Mana {
		old := fmt.Sprintf("Cost%d", i)
		new := fmt.Sprint(cost)
		n := strings.Replace(*ttp, old, new, -1)
		costs = append(costs, new)
		*ttp = n
	}

	n := strings.Replace(*ttp, "@Cost@", strings.Join(costs, "/"), -1)
	*ttp = n
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
		if val.FormulaParts != nil {
			formulaParts(ttp, key, val.FormulaParts, spl.DataValues)
		}
		if val.ModifiedGameCalculation != "" {
			if val.Multiplier.Number != 0 {
				// TODO
				// multiplierNumber(ttp, key, val.ModifiedGameCalculation, val)
			}
			if val.Multiplier.DataValue != "" {
				// TODO
				// multiplierData(ttp, key, val.ModifiedGameCalculation, spl, val.Multiplier)
			}
		}

	}
}

func formulaParts(ttp *string, key string, fps []FormulaPart, spl []SpellDataValue) {
	for _, fp := range fps {
		if len(fp.DataValue) != 0 {
			dataValCalc(ttp, key, fp.DataValue, spl)
		} else if fp.Breakpoints != nil {
			breakpoints(ttp, key, fp)
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

func finalCleanup(input string) string {
	re := regexp.MustCompile(`@(\d+(\.\d+)?(?:/\d+(\.\d+)?)*)@`)
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		numStr := match[1 : len(match)-1]

		return numStr
	})
	r, err := calculate(result)
	if err != nil {

	}

	result = r

	return result
}

func floatValues(dataVals []interface{}) ([]float64, error) {
	var floatVals []float64
	for _, val := range dataVals {
		switch v := val.(type) {
		case float64:
			floatVals = append(floatVals, v)
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return []float64{}, fmt.Errorf("Error converting string to float:", err)

			}
			floatVals = append(floatVals, f)
		default:
			return []float64{}, fmt.Errorf("Unsupported data type in dataVals")
		}
	}

	return floatVals, nil
}

// formatDecimal formats a floating point number string to two decimal places.
func formatDecimal(number string) (string, error) {
	// Convert the string to float64
	value, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return "", err
	}
	// Format the value to two decimal places
	return fmt.Sprintf("%.2f", value), nil
}

// evaluateExpression takes an expression string, evaluates it, and returns the result as a string.
func evaluateExpression(expression string) (string, error) {
	// Remove the `@` characters
	expression = expression[1 : len(expression)-1]
	// Split the expression into its components (assuming the format is `a*b`)
	parts := regexp.MustCompile(`\*`).Split(expression, -1)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid expression format")
	}
	// Convert the parts to float64
	a, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return "", err
	}
	b, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return "", err
	}
	// Perform the multiplication
	result := a * b
	// Return the result as a string formatted to 2 decimal places
	return fmt.Sprintf("%.2f", result), nil
}

// processString takes the entire input string, formats all decimal numbers, evaluates expressions between `@`, and replaces them with their results.
func calculate(input string) (string, error) {
	// Regular expression to find all decimal numbers
	reDecimal := regexp.MustCompile(`\d+\.\d+`)
	// Function to replace each decimal match with its formatted version
	result := reDecimal.ReplaceAllStringFunc(input, func(match string) string {
		formatted, err := formatDecimal(match)
		if err != nil {
			// If there's an error, return the original match
			return match
		}
		return formatted
	})

	// Regular expression to find expressions between `@`
	reExpression := regexp.MustCompile(`@[^@]*@`)
	// Function to replace each expression match with its evaluated version
	result = reExpression.ReplaceAllStringFunc(result, func(match string) string {
		evaluated, err := evaluateExpression(match)
		if err != nil {
			// If there's an error, return the original match
			return match
		}
		return evaluated
	})

	return result, nil
}
