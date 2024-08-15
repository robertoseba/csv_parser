package rule

import (
	"strconv"

	"github.com/robertoseba/csv_parser/pkg/row"
)

type logicalOperatorType string

const AND_OPERATOR logicalOperatorType = "&&"
const OR_OPERATOR logicalOperatorType = "||"

type ColRules struct {
	column          string
	rules           []Rule
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
			if r.logicalOperator == OR_OPERATOR {
				return true
			}
		} else {
			if r.logicalOperator == AND_OPERATOR {
				return false
			}
		}
	}

	return r.logicalOperator == AND_OPERATOR
}

func newColRules(column string, initNumRules int) *ColRules {
	rules := make([]Rule, 0, initNumRules)

	return &ColRules{
		column:          column,
		logicalOperator: AND_OPERATOR,
		isNumber:        false,
		rules:           rules,
	}
}

func (r *ColRules) addRule(ruleType allowedRules, ruleValue string) {

	ruleAsFloat, err := strconv.ParseFloat(ruleValue, 64)

	if err == nil {
		r.rules = append(r.rules, Rule{value: ruleValue, ruleType: ruleType, floatValue: &ruleAsFloat})
		r.isNumber = true
		return
	}

	r.isNumber = false

	r.rules = append(r.rules, Rule{value: ruleValue, ruleType: ruleType, floatValue: nil})
}
