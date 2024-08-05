package rule

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

type Rule struct {
	value      string
	floatValue *float64
	operator   string
}

func (rule *Rule) IsValid(rowValue string) bool {
	if rule.floatValue != nil {
		rowValueFloat, err := strconv.ParseFloat(rowValue, 64)
		if err == nil {
			return compareValues(rowValueFloat, *rule.floatValue, rule.operator)
		}

	}
	return compareValues(rowValue, rule.value, rule.operator)
}

func compareValues[T constraints.Ordered](first, second T, operator string) bool {
	switch operator {
	case "eq":
		return first == second
	case "gt":
		return first > second
	case "lt":
		return first < second
	case "gte":
		return first >= second
	case "lte":
		return first <= second
	case "!eq":
		return first != second
	default:
		panic("invalid operator")
	}
}
