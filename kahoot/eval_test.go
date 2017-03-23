package kahoot

import "testing"

func TestEval(t *testing.T) {
	exprs := map[string]int64{
		"(23 + 64 + 35 * 35)":                    1312,
		"88 * 94 * 9 * 48":                       3573504,
		"59 * 93 * (89 *\t 9) * 60 * (4 + 47)":   13448966220,
		"(7 + 80 + ((23 * 35) + 32))":            924,
		"(58 + 8 + ((72 * 46) + 56 * 13  + 49))": 4155,
	}
	for expr, expected := range exprs {
		actual, err := eval(expr)
		if err != nil {
			t.Error(expr+": ", err)
		} else if actual != expected {
			t.Errorf("%s: expected %d got %d", expr, expected, actual)
		}
	}
}
