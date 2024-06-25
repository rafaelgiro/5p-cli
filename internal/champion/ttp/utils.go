package ttp

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func initialCleanup(ttp *Tooltip) {
	reSpecial := regexp.MustCompile(`@[^@]*?(?:Postfix|Prefix)@`)
	result := reSpecial.ReplaceAllString(ttp.ToString(), "")
	*ttp = Tooltip(result)
}

func finalCleanup(ttp *Tooltip) error {
	re := regexp.MustCompile(`@(\d+(\.\d+)?(?:/\d+(\.\d+)?)*)@`)
	result := re.ReplaceAllStringFunc(ttp.ToString(), func(match string) string {
		numStr := match[1 : len(match)-1]

		return numStr
	})

	r, err := calculate(result)

	if err != nil {
		return fmt.Errorf("error while final formating of tooltip: %v", err)
	}

	*ttp = Tooltip(r)

	return nil
}

// formatDecimal formats a floating point number string to two decimal places.
func formatDecimal(number string) (string, error) {
	value, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.2f", value), nil
}

// evaluateExpression evaluates a mathematical expression within `@` and returns the result as a string.
func evaluateExpression(expression string) (string, error) {
	// Remove the enclosing `@` symbols
	expression = strings.Trim(expression, "@")

	// If the expression contains a final multiplier, handle it
	if strings.Contains(expression, "*") {
		parts := strings.Split(expression, "*")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid expression format")
		}
		numbersPart := parts[0]
		multiplierPart := parts[1]

		// Parse the multiplier
		multiplier, err := strconv.ParseFloat(multiplierPart, 64)
		if err != nil {
			return "", err
		}

		// Split the numbers part by `/` and evaluate each one
		numbers := strings.Split(numbersPart, "/")
		results := make([]string, len(numbers))
		for i, numberStr := range numbers {
			number, err := strconv.ParseFloat(numberStr, 64)
			if err != nil {
				return "", err
			}
			results[i] = fmt.Sprintf("%.2f", number*multiplier)
		}

		// Join the results with `/` and return the final string
		return strings.Join(results, "/"), nil
	}

	// If there is no multiplier, evaluate the simple expression (assuming simple arithmetic)
	value, err := strconv.ParseFloat(expression, 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.2f", value), nil
}

// calculate takes the entire input string, formats all decimal numbers, evaluates expressions between `@`, and replaces them with their results.
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
