package rule

import (
	"strconv"
)

type Rule struct {
	strValue     string
	numericValue float64
	operator     string
}

func (rule *Rule) IsValid(rowValue string, castAsNumber bool) bool {
	if castAsNumber {
		rowValueFloat, err := strconv.ParseFloat(rowValue, 64)
		if err == nil {
			return compareValues(rowValueFloat, rule.numericValue, rule.operator)
		}

	}
	return compareValues(rowValue, rule.strValue, rule.operator)
}
