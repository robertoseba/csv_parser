package rule

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type allowedRules string

const (
	EQ_RULE  allowedRules = "eq"
	GT_RULE  allowedRules = "gt"
	LT_RULE  allowedRules = "lt"
	GTE_RULE allowedRules = "gte"
	LTE_RULE allowedRules = "lte"
	NE_RULE  allowedRules = "!eq"
)

var ALL_RULES = []string{
	string(EQ_RULE),
	string(GT_RULE),
	string(LT_RULE),
	string(GTE_RULE),
	string(LTE_RULE),
	string(NE_RULE),
}

const RULE_SEPARATOR = ";"
const COL_RULE_SEPARATOR = ":"

// Rule formats examples: eq(5) or !eq(3) or ||eq(5) or &&eq(5)
var RULE_FORMAT = `\s*(\|\||&&)?(` + strings.Join(ALL_RULES, "|") + `)\s*\((\w+)\)\s*`

var ErrInvalidRule = errors.New("invalid rule format")

/**
* Returns a collection of rules grouped by column's name
* Each column can have multiple rules and have a logical operator
* that defines how the rules should be evaluated
 */
func RulesFromStr(ruleStr string) (map[string]*ColRules, error) {
	if strings.Trim(ruleStr, " ") == "" {
		return nil, nil
	}

	colRules := strings.Split(ruleStr, RULE_SEPARATOR)

	rulesByCols := make(map[string]*ColRules, len(colRules))

	regexRuleFormat := regexp.MustCompile(RULE_FORMAT)

	for _, strRule := range colRules {
		if strings.Trim(strRule, " ") == "" {
			continue
		}

		column, rules, ok := strings.Cut(strRule, COL_RULE_SEPARATOR)
		if !ok {
			return nil, ErrInvalidRule
		}
		if strings.Contains(rules, string(OR_OPERATOR)) && strings.Contains(rules, string(AND_OPERATOR)) {
			return nil, ErrInvalidRule
		}

		logicalOperator := AND_OPERATOR
		splittedStringRules := regexRuleFormat.FindAllString(rules, -1)

		rulesByCols[column] = &ColRules{
			column: column,
			rules:  make([]Rule, len(splittedStringRules)),
		}

		for idx, strRule := range splittedStringRules {
			switch strRule[0:2] {
			case "&&":
				strRule = strRule[2:]
			case "||":
				logicalOperator = OR_OPERATOR
				strRule = strRule[2:]
			}

			ruleOperator, ruleValue, _ := strings.Cut(strRule, "(")
			ruleValue = strings.Trim(ruleValue, ")")

			rule := Rule{
				operator: ruleOperator,
				value:    ruleValue,
			}

			ruleNumber, err := strconv.ParseFloat(ruleValue, 64)

			if err == nil {
				rulesByCols[column].isNumber = true
				rule.floatValue = &ruleNumber
			} else {
				rulesByCols[column].isNumber = false
			}

			rulesByCols[column].rules[idx] = rule
		}

		rulesByCols[column].logicalOperator = logicalOperator

	}
	return rulesByCols, nil
}
