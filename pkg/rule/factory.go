package rule

import (
	"fmt"
	"regexp"
	"strings"
)

/**
* Returns a collection of rules grouped by column's name
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

			//TODO: extract the operator and value from the rule
			operator, value, _ := strings.Cut(rule, "(")
			value = strings.Trim(value, ")")
			fmt.Println(operator)
			fmt.Println(value)

		}
		fmt.Println(logicalOperator)

	}
	return nil, nil
}
