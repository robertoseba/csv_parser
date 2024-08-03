package rule

import (
	"fmt"
	"strconv"
	"strings"
)

var suportedOperators = [...]string{"!=", ">=", "<=", "=", ">", "<"}

func NewRuleFromStr(strRule string) (IRule, error) {
	var column, value string
	isValueNumber := false

	for _, operatorSymbol := range suportedOperators {
		if strings.Contains(strRule, operatorSymbol) {
			parts := strings.Split(strRule, operatorSymbol)

			column = parts[0]
			value = parts[1]

			_, err := strconv.ParseFloat(value, 64)
			if err == nil {
				isValueNumber = true
			}
			rule := Rule{
				column:        column,
				value:         value,
				isValueNumber: isValueNumber,
			}

			switch operatorSymbol {
			case "=":
				return &EqRule{Rule: rule}, nil
			case "!=":
				return &NotEqRule{Rule: rule}, nil

			}
		}
	}

	return nil, fmt.Errorf("invalid rule: %s", strRule)
}
