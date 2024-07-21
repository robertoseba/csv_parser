package parser

import (
	"fmt"
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

var operators = map[string]func(a, b string) bool{
	"!=": func(a, b string) bool { return a != b },
	">=": func(a, b string) bool { return a >= b },
	"<=": func(a, b string) bool { return a <= b },
	"=":  func(a, b string) bool { return a == b },
	">":  func(a, b string) bool { return a > b },
	"<":  func(a, b string) bool { return a < b },
}

func NewRule(strRule string) (IRule, error) {
	var operator func(a, b string) bool
	var col, value string

	for operatorKey, operatorFunc := range operators {
		if strings.Contains(strRule, operatorKey) {
			operator = operatorFunc
			parts := strings.Split(strRule, operatorKey)
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

func NewRulesFrom(strRules []string) ([]IRule, error) {
	if len(strRules) == 0 {
		return nil, nil
	}

	rules := make([]IRule, 0, len(strRules))

	for _, strRule := range strRules {
		r, err := NewRule(strRule)

		if err != nil {
			return nil, err
		}

		rules = append(rules, r)
	}

	return rules, nil
}

// TODO: Implement conversion of string to number if necessary
func (r *Rule) Validate(row *Row) bool {
	return r.operator(row.GetColumn(r.column), r.value)
}

func (r *Rule) Column() string {
	return r.column
}
