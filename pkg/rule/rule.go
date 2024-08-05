package rule

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

type allowedRules string

const (
	EQ_RULE  allowedRules = "eq"
	GT_RULE  allowedRules = "gt"
	LT_RULE  allowedRules = "lt"
	GTE_RULE allowedRules = "gte"
	LTE_RULE allowedRules = "lte"
	NE_RULE  allowedRules = "!eq"
)

var ALL_RULES = []string{
	string(EQ_RULE),
	string(GT_RULE),
	string(LT_RULE),
	string(GTE_RULE),
	string(LTE_RULE),
	string(NE_RULE),
}

type Rule struct {
	value      string
	floatValue *float64
	operator   allowedRules
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

func compareValues[T constraints.Ordered](first, second T, operator allowedRules) bool {
	switch operator {
	case EQ_RULE:
		return first == second
	case GT_RULE:
		return first > second
	case LT_RULE:
		return first < second
	case GTE_RULE:
		return first >= second
	case LTE_RULE:
		return first <= second
	case NE_RULE:
		return first != second
	default:
		panic("invalid operator")
	}
}
