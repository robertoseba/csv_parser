package csv_parser

import (
	"fmt"
	"strings"
)

type Rule struct{
	column 		string
	operator 	func(a, b string) bool	
	value 		string
}

var operators = map[string]func(a, b string) bool{
	"!=": func(a, b string) bool { return a != b },
	">=": func(a, b string) bool { return a >= b },
	"<=": func(a, b string) bool { return a <= b },
	"=": func(a, b string) bool { return a == b },
	">": func(a, b string) bool { return a > b },
	"<": func(a, b string) bool { return a < b },
}

func NewRule(strRule string, headers *Row) (*Rule, error){
	var operator func (a, b string) bool
	var col, value string

	for operatorKey, operatorFunc := range operators{
		if strings.Contains(strRule, operatorKey){
			operator = operatorFunc
			parts := strings.Split(strRule, operatorKey)
			col = parts[0]
			value = parts[1]
			break	
		}
	}

	if operator == nil{
		return nil, fmt.Errorf("invalid rule: %s", strRule)
	}

	if !headers.Contains(col) {
		return nil, fmt.Errorf("invalid column for rule: %s", col)
	}

	return &Rule{
		column: col,
		operator: operator,
		value: value,
	},nil
}

func (r *Rule) Validate(row *Row) bool{
	return r.operator(row.GetColumn(r.column), r.value)
}
