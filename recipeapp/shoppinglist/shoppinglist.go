package shoppinglist

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// - The numeric value is returned as a decimal (float64).
// - The remainder (text part) is trimmed of leading/trailing whitespace.
// - If the text part contains a "/", the text part is cut where the first "/" appears (the "/" and anything after it is removed).
// Returns an error if there is no leading numeric token or if parsing fails (including zero denominator).
func SplitLeadingNumberDecimal(s string) (float64, string, error) {
	// Trim only leading whitespace so we still preserve mid-string spacing for the text part.
	s = strings.TrimLeft(s, " \t\r\n")

	// Regex to capture a leading numeric token (mixed number, fraction, or integer) and the rest.
	// Order matters: try mixed number first, then fraction, then integer.
	// Token examples matched:
	//  - "123"
	//  - "-42"
	//  - "3/4"
	//  - "-3/4"
	//  - "1 1/2"
	//  - "+1 2/3"
	re := regexp.MustCompile(`^([+-]?(?:\d+\s+\d+/\d+|\d+/\d+|\d+))(.*)$`)
	m := re.FindStringSubmatch(s)
	if m == nil {
		return 0, "", fmt.Errorf("no leading numeric token found")
	}

	numToken := strings.TrimSpace(m[1])
	rest := m[2]

	// Parse numeric token into decimal
	value, err := parseNumericToken(numToken)
	if err != nil {
		return 0, "", err
	}

	// Trim leading whitespace of rest, then cut at first "/" if present, finally trim surrounding whitespace.
	rest = strings.TrimLeft(rest, " \t\r\n")
	if i := strings.Index(rest, "/"); i != -1 {
		rest = rest[:i]
	}
	rest = strings.TrimSpace(rest)

	return value, rest, nil
}

// parseNumericToken parses tokens like:
//
//	"123" -> 123.0
//	"-3/4" -> -0.75
//	"1 1/2" -> 1.5
//	"+2" -> 2.0
func parseNumericToken(token string) (float64, error) {
	if token == "" {
		return 0, errors.New("empty numeric token")
	}

	// Extract sign
	sign := 1.0
	if token[0] == '+' || token[0] == '-' {
		if token[0] == '-' {
			sign = -1.0
		}
		token = token[1:]
		token = strings.TrimSpace(token)
	}

	// Mixed number: "1 1/2"
	if strings.Contains(token, " ") {
		parts := strings.Fields(token)
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid mixed number format: %q", token)
		}
		wholeStr, fracStr := parts[0], parts[1]
		whole, err := strconv.ParseInt(wholeStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid whole number in mixed number: %v", err)
		}
		numer, denom, err := splitNumerDenom(fracStr)
		if err != nil {
			return 0, err
		}
		if denom == 0 {
			return 0, errors.New("zero denominator in fraction")
		}
		fracVal := float64(numer) / float64(denom)
		return sign * (float64(whole) + fracVal), nil
	}

	// Simple fraction: "3/4"
	if strings.Contains(token, "/") {
		numer, denom, err := splitNumerDenom(token)
		if err != nil {
			return 0, err
		}
		if denom == 0 {
			return 0, errors.New("zero denominator in fraction")
		}
		return sign * (float64(numer) / float64(denom)), nil
	}

	// Integer
	i, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer token: %v", err)
	}
	return sign * float64(i), nil
}

// splitNumerDenom splits "a/b" into (a, b)
func splitNumerDenom(frac string) (int64, int64, error) {
	parts := strings.SplitN(frac, "/", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid fraction: %q", frac)
	}
	numStr := strings.TrimSpace(parts[0])
	denStr := strings.TrimSpace(parts[1])
	numer, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid numerator: %v", err)
	}
	denom, err := strconv.ParseInt(denStr, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid denominator: %v", err)
	}
	return numer, denom, nil
}
