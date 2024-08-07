package rule

import (
	"errors"
	"regexp"
	"strings"
)

const RULE_SEPARATOR = ";"
const COL_RULE_SEPARATOR = ":"

// Rule formats examples: eq(5) or !eq(3) or ||eq(5) or &&eq(5)
var STR_RULE_FORMAT = `\s*(\|\||&&)?(` + strings.Join(ALL_RULES, "|") + `)\s*\((\w+)\)\s*`

var ErrInvalidRule = errors.New("invalid rule format")

/**
* Returns a collection of rules grouped by column's name
* Each column can have multiple rules and have a logical operator
* that defines how the rules should be evaluated.
* Also, based on the rules, the column can be marked as a number column.
 */
func NewFrom(ruleInput string) ([]*ColRules, error) {
	if strings.Trim(ruleInput, " ") == "" {
		return nil, nil
	}

	splittedRulesByCols := strings.Split(ruleInput, RULE_SEPARATOR)

	rulesByCols := make([]*ColRules, 0, len(splittedRulesByCols))

	regexRuleFormat := regexp.MustCompile(STR_RULE_FORMAT)

	for _, colStrRules := range splittedRulesByCols {
		if strings.Trim(colStrRules, " ") == "" {
			continue
		}

		column, rules, ok := strings.Cut(colStrRules, COL_RULE_SEPARATOR)
		if !ok {
			return nil, ErrInvalidRule
		}
		if strings.Contains(rules, string(OR_OPERATOR)) && strings.Contains(rules, string(AND_OPERATOR)) {
			return nil, ErrInvalidRule
		}

		strRules := regexRuleFormat.FindAllString(rules, -1)

		colRule := newColRules(column, len(strRules))

		for _, strRule := range strRules {
			colRule.addRule(strRule)
		}

		rulesByCols = append(rulesByCols, colRule)

	}
	return rulesByCols, nil
}
