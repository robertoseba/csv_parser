package rule

import (
	"strconv"
	"strings"

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
	result := true

	for _, rule := range r.rules {
		if rule.IsValid(row.GetColumn(r.column)) {
			if r.logicalOperator == OR_OPERATOR {
				break
			}
		} else {
			if r.logicalOperator == AND_OPERATOR {
				result = false
				break
			}
		}
	}

	return result
}

func NewColRules(column string, initNumRules int) *ColRules {
	rules := make([]Rule, 0, initNumRules)

	return &ColRules{
		column:          column,
		logicalOperator: AND_OPERATOR,
		isNumber:        false,
		rules:           rules,
	}
}

func (r *ColRules) AddRule(strRule string) {
	switch strRule[0:2] {
	case "&&":
		strRule = strRule[2:]
	case "||":
		r.logicalOperator = OR_OPERATOR
		strRule = strRule[2:]
	}

	ruleType, ruleValue, _ := strings.Cut(strRule, "(")
	ruleValue = strings.Trim(ruleValue, ")")

	ruleAsFloat, err := strconv.ParseFloat(ruleValue, 64)

	if err == nil {
		r.rules = append(r.rules, Rule{value: ruleValue, ruleType: allowedRules(ruleType), floatValue: &ruleAsFloat})
		r.isNumber = true
		return
	}

	r.isNumber = false

	r.rules = append(r.rules, Rule{value: ruleValue, ruleType: allowedRules(ruleType), floatValue: nil})
}
