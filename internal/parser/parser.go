package parser

import (
	"errors"
	"slices"
	"strings"
)

var ErrInvalidRule = errors.New("invalid rule format")
var ErrInvalidOperator = errors.New("we currently only support one type of logical operator per column rules")
var ErrInvalidRuleType = errors.New("invalid rule type")

/**
* Returns a slice of items with grouped rules by their columns.
* Each column can have multiple rules and have a logical operator
* that defines how the rules should be evaluated.
* Also, based on the rule value, the column can be set to be a number type.
 */
func ParseRules(ruleInput string) ([]ColRules, error) {
	if strings.Trim(ruleInput, " ") == "" {
		return nil, nil
	}
	return parse(ruleInput)

}

func parse(rulesInput string) ([]ColRules, error) {
	const ruleSeparator = ";"

	rulesCount := strings.Count(rulesInput, ruleSeparator)
	rulesByCols := make([]ColRules, 0, rulesCount)

	for {
		logicalOperator := andOperator

		colRuleEndPos := strings.Index(rulesInput, ruleSeparator)
		if colRuleEndPos == -1 {
			colRuleEndPos = len(rulesInput)
		}

		colName, remaining, ok := parseColName(rulesInput[:colRuleEndPos])
		if !ok {
			return nil, ErrInvalidRule
		}

		colRule := newColRules(colName, 0)
		// Retrieves each rule for the column (one column can have multiple rules with a logical operator)
		for {

			ruleTypeEndPos := strings.IndexByte(remaining, '(')
			if ruleTypeEndPos == -1 {
				return nil, ErrInvalidRule
			}

			var rule allowedRules
			var err error

			rule, logicalOperator, err = parseRuleTypeAndOperator(remaining[:ruleTypeEndPos], logicalOperator)

			if err != nil {
				return nil, err
			}

			remaining = remaining[ruleTypeEndPos+1:]

			valueEndPos := strings.IndexByte(remaining, ')')
			if valueEndPos == -1 {
				return nil, ErrInvalidRule
			}
			ruleValue := remaining[:valueEndPos]

			colRule.addRule(allowedRules(rule), ruleValue)

			if valueEndPos+1 >= len(remaining) {
				break
			}
			remaining = remaining[valueEndPos+1:]
		}

		colRule.logicalOperator = logicalOperator
		rulesByCols = append(rulesByCols, *colRule)

		if colRuleEndPos+1 >= len(rulesInput) {
			break
		}

		// Update rulesInput to remove the already parsed column rules
		rulesInput = rulesInput[colRuleEndPos+1:]
	}
	return rulesByCols, nil
}

func parseColName(ruleInput string) (string, string, bool) {
	colEndPos := strings.IndexByte(ruleInput, ':')
	if colEndPos == -1 {
		return "", ruleInput, false
	}
	return ruleInput[:colEndPos], ruleInput[colEndPos+1:], true
}

func parseRuleTypeAndOperator(ruleInput string, previousOperator logicalOperatorType) (allowedRules, logicalOperatorType, error) {

	var rule string
	var logicalOperator logicalOperatorType

	switch {
	case ruleInput[:len(andOperator)] == string(andOperator):
		logicalOperator = andOperator
		rule = ruleInput[len(andOperator):]

	case ruleInput[:len(orOperator)] == string(orOperator):
		logicalOperator = orOperator
		rule = ruleInput[len(orOperator):]
	default:
		logicalOperator = previousOperator
		rule = ruleInput
	}
	// Only one logical operator can be set for each column rules
	// Since we start with AND_OPERATOR, If OR_OPERATOR  has been set during
	// parsing than it can´t be set again to AND_OPERATOR as it would signify
	// multiple operators for the same column rules
	if previousOperator == orOperator && logicalOperator == andOperator {
		return "", previousOperator, ErrInvalidOperator
	}

	if !isRuleTypeValid(rule) {
		return "", logicalOperator, ErrInvalidRuleType
	}

	return allowedRules(rule), logicalOperator, nil
}

func isRuleTypeValid(rule string) bool {
	if i := slices.Index(all_rules, rule); i == -1 {
		return false
	}
	return true
}
