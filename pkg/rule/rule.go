package rule

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/robertoseba/csv_parser/pkg/row"
)

type IRule interface {
	Validate(row *row.Row) bool
	Column() string
	IsNumber() bool
}

type Rule struct {
	column        string
	value         string
	isValueNumber bool
	// OrlogicalOpt  bool
}

func (r *Rule) Column() string {
	return r.column
}

func (r *Rule) IsNumber() bool {
	return r.isValueNumber
}

func (rule *Rule) preProcessNumbers(rowValue string) (ruleValue, colValue float64, err error) {
	if rule.isValueNumber {
		rowValueFloat, err := strconv.ParseFloat(rowValue, 64)
		if err == nil {
			ruleValueFloat, _ := strconv.ParseFloat(rule.value, 64)
			return ruleValueFloat, rowValueFloat, nil
		}
	}

	return 0, 0, fmt.Errorf("cannot convert values to numbers")
}

func RulesFromStr(ruleStr string) ([]IRule, error) {
	if strings.Trim(ruleStr, " ") == "" {
		return nil, nil
	}

	expressions := []string{"!eq", "eq", "gt", "lt", "gte", "lte"}
	// logicalOperators := []string{"&&", "||"}

	// Rule formats examples: eq(5) or !eq(3)
	regexRuleFormat := regexp.MustCompile(`\s*(\|\||&&)?(` + strings.Join(expressions, "|") + `)\s*\((\w+)\)\s*`)

	// logicalOperators := []string{"AND", "OR"}
	// regexLogialOperators := regexp.MustCompile(`\s*(` + strings.Join(logicalOperators, "|") + `)\s*`)

	splittedGroupRules := strings.Split(ruleStr, ";")

	for _, strRule := range splittedGroupRules {
		if strings.Trim(strRule, " ") == "" {
			continue
		}

		column, rules, ok := strings.Cut(strRule, ":")
		if !ok {
			return nil, fmt.Errorf("invalid rule format")
		}
		fmt.Println(column)
		for _, rule := range regexRuleFormat.FindAllString(rules, -1) {
			fmt.Println(rule)
		}

	}
	// rules := regexLogialOperators.Split(strRule, -1)

	// for idx, rule := range rules {
	// 	if strings.Trim(rule, " ") == "" {
	// 		continue
	// 	}

	// }

	return []IRule{
		&EqRule{Rule: Rule{column: "col1", value: "5", isValueNumber: true}},
		&NotEqRule{Rule: Rule{column: "col2", value: "3", isValueNumber: true}},
	}, nil
}
