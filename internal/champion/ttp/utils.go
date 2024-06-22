package ttp

import (
	"fmt"
	"regexp"
	"strconv"
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
