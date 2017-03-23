package kahoot

import (
	"regexp"
	"strconv"
	"strings"
)

var simpleExprRegexp = regexp.MustCompile(`\(([0-9\+\*\s]*)\)`)

// eval evaluates a mathematical expression, such as:
//
//     ((76 * 21) * (((81 + 4) * 55) + 10))
//
// This is necessary to obtain session tokens.
func eval(expr string) (int64, error) {
	for {
		// Evaluate a simple sub-expression.
		match := simpleExprRegexp.FindStringSubmatch(expr)
		if match == nil {
			break
		}
		simple := match[1]
		val, err := evalSimple(simple)
		if err != nil {
			return 0, err
		}
		valStr := strconv.FormatInt(val, 10)
		expr = strings.Replace(expr, match[0], valStr, 1)
	}
	return evalSimple(expr)
}

// evalSimple evaluates an expression with no nested
// parentheses, like
//
//     23 + 64 + 35 * 35
//
func evalSimple(expr string) (int64, error) {
	var sum int64
	for _, sumTerm := range strings.Split(expr, "+") {
		product := int64(1)
		for _, numStr := range strings.Split(sumTerm, "*") {
			num, err := strconv.ParseInt(strings.TrimSpace(numStr), 10, 64)
			if err != nil {
				return 0, err
			}
			product *= num
		}
		sum += product
	}
	return sum, nil
}
