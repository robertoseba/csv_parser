package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type mytype interface {
	string | int64 | float64
}

type Rule struct {
	column   string
	operator func(a, b string) bool
	value    string
}

type IRule interface {
	Validate(row *Row) bool
	Column() string
}

func NewRule(strRule string) (IRule, error) {
	var operator func(a, b string) bool
	var col, value string

	for idx, operatorSymbol := range operatorsSymbols {
		if strings.Contains(strRule, operatorSymbol) {
			operator = operatorsFuncs[idx]
			parts := strings.Split(strRule, operatorSymbol)
			col = parts[0]
			value = parts[1]
			break
		}
	}

	if operator == nil {
		return nil, fmt.Errorf("invalid rule: %s", strRule)
	}

	return &Rule{
		column:   col,
		operator: operator,
		value:    value,
	}, nil
}

// TODO: Implement conversion of string to number if necessary
func (r *Rule) Validate(row *Row) bool {
	return r.operator(r.value, row.GetColumn(r.column))

}

func (r *Rule) Column() string {
	return r.column
}

var operatorsSymbols = [6]string{"!=", ">=", "<=", "=", ">", "<"}
var operatorsFuncs = [6]func(a, b string) bool{
	func(a, b string) bool { return a != b },
	gte,
	func(a, b string) bool { return a <= b },
	func(a, b string) bool { return a == b },
	gt,
	func(a, b string) bool { return a < b },
}

func gte(ruleValue string, colValue string) bool {
	rValue, err := strconv.ParseFloat(ruleValue, 64)

	if err == nil {
		cValue, err := strconv.ParseFloat(colValue, 64)
		if err != nil {
			panic("Cannot compare numbers and strings")
		}
		return cValue >= rValue
	}

	return ruleValue >= colValue
}

func gt(ruleValue string, colValue string) bool {
	rValue, err := strconv.ParseFloat(ruleValue, 64)
	if err == nil {
		cValue, err := strconv.ParseFloat(colValue, 64)
		if err != nil {
			panic("Cannot compare numbers and strings")
		}
		return cValue > rValue

	}

	return ruleValue > colValue
}
