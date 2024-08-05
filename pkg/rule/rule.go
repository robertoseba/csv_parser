package rule

import (
	"strconv"
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
