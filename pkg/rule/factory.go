package rule

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var ErrInvalidRule = errors.New("invalid rule format")
var ErrInvalidOperator = errors.New("invalid rule logical operator format")
var ErrInvalidRuleType = errors.New("invalid rule type")

/**
* Returns a slice of items with grouped rules by their columns.
* Each column can have multiple rules and have a logical operator
* that defines how the rules should be evaluated.
* Also, based on the rule value, the column can be set to be a number type.
 */
func NewFrom(ruleInput string) ([]*ColRules, error) {
	if strings.Trim(ruleInput, " ") == "" {
		return nil, nil
	}
	return parseRules(ruleInput)

}

func parseRules(rulesInput string) ([]*ColRules, error) {
	const ruleSeparator = ";"

	rulesCount := strings.Count(rulesInput, ruleSeparator)
	rulesByCols := make([]*ColRules, 0, rulesCount)

	for {
		logicalOperator := AND_OPERATOR

		endPos := parseEndPosColRules(rulesInput, ruleSeparator)
		colRulesString := rulesInput[:endPos]

		colName, remaining, ok := parseColName(colRulesString)
		if !ok {
			return nil, ErrInvalidRule
		}

		colRule := newColRules(colName, 0)

		colRulesString = remaining

		// Retrieves each rule for the column (one column can have multiple rules with a logical operator)
		for {

			ruleTypeEndPos := strings.IndexByte(colRulesString, '(')
			if ruleTypeEndPos == -1 {
				return nil, ErrInvalidRule
			}

			rule := colRulesString[:ruleTypeEndPos]
			if rule[:len(AND_OPERATOR)] == string(AND_OPERATOR) ||
				rule[0:len(OR_OPERATOR)] == string(OR_OPERATOR) {

				// Only one logical operator can be set for each column rules
				// If OR logical operator has been set during parsing
				// than it can´t have AND in the same column rules
				if logicalOperator == OR_OPERATOR && rule[:2] == string(AND_OPERATOR) {
					return nil, ErrInvalidOperator
				}

				logicalOperator = logicalOperatorType(rule[:2])
				rule = rule[2:]
			}

			if !isRuleTypeValid(rule) {
				return nil, ErrInvalidRuleType
			}

			colRulesString = colRulesString[ruleTypeEndPos+1:]

			// Retrieves value for rule
			valueEndPos := strings.IndexByte(colRulesString, ')')
			if valueEndPos == -1 {
				return nil, ErrInvalidRule
			}
			ruleValue := colRulesString[:valueEndPos]

			colRule.addRule(fmt.Sprintf("%s%s(%s)", logicalOperator, allowedRules(rule), ruleValue))

			if valueEndPos+1 >= len(colRulesString) {
				break
			}
			colRulesString = colRulesString[valueEndPos+1:]
		}

		rulesByCols = append(rulesByCols, colRule)

		if endPos+1 >= len(rulesInput) {
			break
		}

		// Update rulesInput to remove the already parsed column rules
		rulesInput = rulesInput[endPos+1:]
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

func parseEndPosColRules(ruleInput string, ruleSeparator string) int {
	ruleEndPos := strings.Index(ruleInput, ruleSeparator)
	if ruleEndPos == -1 {
		ruleEndPos = len(ruleInput)
	}

	return ruleEndPos
}

func isRuleTypeValid(rule string) bool {
	if i := slices.Index(ALL_RULES, rule); i == -1 {
		return false
	}
	return true
}

// func parseRuleType(ruleInput string, currLogicalOperator logicalOperatorType) (allowedRules, string, bool) {
// 	ruleTypeEndPos := strings.IndexByte(ruleInput, '(')
// 	if ruleTypeEndPos == -1 {
// 		return "", ruleInput, false
// 	}

// 	ruleType := allowedRules(ruleInput[:ruleTypeEndPos])
// 	if string(ruleType[:len(AND_OPERATOR)]) == string(AND_OPERATOR) ||
// 		string(ruleType[:len(OR_OPERATOR)]) == string(OR_OPERATOR) {

// 		// Only one logical operator can be set for each column rules
// 		// If OR logical operator has been set during parsing than it can´t have AND in the same column rules
// 		if currLogicalOperator == OR_OPERATOR && string(ruleType[:2]) == string(AND_OPERATOR) {
// 			return "", "", false
// 		}

// 		logicalOperator = logicalOperatorType(ruleType[:2])
// 		ruleType = allowedRules(ruleType[2:])
// 	}

// 	if i := slices.Index(ALL_RULES, string(ruleType)); i == -1 {
// 		return "", ruleInput, false
// 	}

// 	return ruleType, ruleInput[ruleTypeEndPos+1:],

// }
