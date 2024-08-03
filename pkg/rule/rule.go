package rule

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/robertoseba/csv_parser/pkg/row"
)

type logicalOperatorType string

const AND_OPERATOR logicalOperatorType = "&&"
const OR_OPERATOR logicalOperatorType = "||"

type Rule struct {
	strValue     string
	numericValue float64
	operator     string
}

func (r *Rule) IsValid(rowValue string, castAsNumber bool) bool {
	if castAsNumber {
		rowValueFloat, err := strconv.ParseFloat(rowValue, 64)
		if err == nil {
			return compareValues(rowValueFloat, r.numericValue, r.operator)
		}

	}
	return compareValues(rowValue, r.strValue, r.operator)
}

type ColRules struct {
	column          string
	rules           []Rule
	logicalOperator logicalOperatorType
	castAsNumber    bool
}

func (r *ColRules) Column() string {
	return r.column
}

func (r *ColRules) IsNumber() bool {
	return r.castAsNumber
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

/**
* Returns a collection of rules ordered by column's name
* Each column can have multiple rules and have a logical operator
* that defines how the rules should be evaluated
 */
func RulesFromStr(ruleStr string) (map[string][]ColRules, error) {
	if strings.Trim(ruleStr, " ") == "" {
		return nil, nil
	}

	expressions := []string{"!eq", "eq", "gt", "lt", "gte", "lte"}

	// Rule formats examples: eq(5) or !eq(3) or ||eq(5) or &&eq(5)
	regexRuleFormat := regexp.MustCompile(`\s*(\|\||&&)?(` + strings.Join(expressions, "|") + `)\s*\((\w+)\)\s*`)

	colRules := strings.Split(ruleStr, ";")

	for _, strRule := range colRules {
		if strings.Trim(strRule, " ") == "" {
			continue
		}

		column, rules, ok := strings.Cut(strRule, ":")
		if !ok {
			return nil, fmt.Errorf("invalid rule format")
		}
		fmt.Println(column)

		if strings.Contains(rules, string(OR_OPERATOR)) && strings.Contains(rules, string(AND_OPERATOR)) {
			return nil, fmt.Errorf("invalid rule format. Rule can only contain one type of logical operator per column")
		}

		logicalOperator := AND_OPERATOR

		for _, rule := range regexRuleFormat.FindAllString(rules, -1) {
			switch rule[0:2] {
			case "&&":
				rule = rule[2:]
			case "||":
				logicalOperator = OR_OPERATOR
				rule = rule[2:]
			}

			fmt.Println(rule)

		}
		fmt.Println(logicalOperator)

	}
	return []Rule{
		&EqRule{Rule: Rule{column: "col1", strValue: "5", isValueNumber: true}},
		&NotEqRule{Rule: Rule{column: "col2", strValue: "3", isValueNumber: true}},
	}, nil
}
