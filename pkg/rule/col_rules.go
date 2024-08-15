package rule

import (
	"strconv"

	"github.com/robertoseba/csv_parser/pkg/row"
)

type logicalOperatorType string

const andOperator logicalOperatorType = "&&"
const orOperator logicalOperatorType = "||"

type ColRules struct {
	column          string
	rules           []rule
	logicalOperator logicalOperatorType
	isNumber        bool
}

func (r *ColRules) Column() string {
	return r.column
}

func (r *ColRules) IsNumber() bool {
	return r.isNumber
}

func (r *ColRules) IsValid(row *row.Row) bool {
	for _, rule := range r.rules {
		if rule.isValid(row.GetColumn(r.column)) {
			if r.logicalOperator == orOperator {
				return true
			}
		} else {
			if r.logicalOperator == andOperator {
				return false
			}
		}
	}

	return r.logicalOperator == andOperator
}

func newColRules(column string, initNumRules int) *ColRules {
	rules := make([]rule, 0, initNumRules)

	return &ColRules{
		column:          column,
		logicalOperator: andOperator,
		isNumber:        false,
		rules:           rules,
	}
}

func (r *ColRules) addRule(ruleType allowedRules, ruleValue string) {

	ruleAsFloat, err := strconv.ParseFloat(ruleValue, 64)

	if err == nil {
		r.rules = append(r.rules, rule{value: ruleValue, ruleType: ruleType, floatValue: &ruleAsFloat})
		r.isNumber = true
		return
	}

	r.isNumber = false

	r.rules = append(r.rules, rule{value: ruleValue, ruleType: ruleType, floatValue: nil})
}
