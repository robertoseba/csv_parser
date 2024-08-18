package parser

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

type allowedRules string

const (
	eqRule  allowedRules = "eq"
	gtRule  allowedRules = "gt"
	ltRule  allowedRules = "lt"
	gteRule allowedRules = "gte"
	lteRule allowedRules = "lte"
	neRule  allowedRules = "neq"
)

var all_rules = []string{
	string(eqRule),
	string(gtRule),
	string(ltRule),
	string(gteRule),
	string(lteRule),
	string(neRule),
}

type rule struct {
	value      string
	floatValue *float64
	ruleType   allowedRules
}

func (rule *rule) isValid(rowValue string) bool {
	if rule.floatValue != nil {
		rowValueFloat, err := strconv.ParseFloat(rowValue, 64)
		if err == nil {
			return compareValues(rowValueFloat, *rule.floatValue, rule.ruleType)
		}

	}
	return compareValues(rowValue, rule.value, rule.ruleType)
}

func compareValues[T constraints.Ordered](first, second T, operator allowedRules) bool {
	switch operator {
	case eqRule:
		return first == second
	case gtRule:
		return first > second
	case ltRule:
		return first < second
	case gteRule:
		return first >= second
	case lteRule:
		return first <= second
	case neRule:
		return first != second
	default:
		panic("invalid operator")
	}
}
