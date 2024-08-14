package rule

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const RULE_SEPARATOR = ";"
const COL_RULE_SEPARATOR = ":"

// Rule formats examples: eq(5) or neq(3) or ||eq(5) or &&eq(5)
var STR_RULE_FORMAT = `\s*(\|\||&&)?(` + strings.Join(ALL_RULES, "|") + `)\s*\((\w+)\)\s*`

var ErrInvalidRule = errors.New("invalid rule format")

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

	//TODO: Instead of creating multiple strings, consider walking through the string and parsing it
	rulesByCols := make([]*ColRules, 0)

	for {
		col := ""
		logicalOperator := AND_OPERATOR

		ruleEndPos := strings.Index(rulesInput, RULE_SEPARATOR)
		if ruleEndPos == -1 {
			ruleEndPos = len(rulesInput)

		}

		colRulesString := rulesInput[:ruleEndPos]

		//Parses the column name
		colEndPos := strings.Index(colRulesString, COL_RULE_SEPARATOR)
		if colEndPos == -1 {
			return nil, ErrInvalidRule
		}
		col = colRulesString[:colEndPos]
		colRule := newColRules(col, 0)

		// Update colRulesString to remove the already parsed column
		colRulesString = colRulesString[colEndPos+1:]

		// Retrieves each rule for the column (one column can have multiple rules with a logical operator)
		for {
			ruleValue := ""
			var ruleType allowedRules

			// Retrieve rule type and logical operator
			ruleTypeEndPos := strings.Index(colRulesString, "(")
			if ruleTypeEndPos == -1 {
				return nil, ErrInvalidRule
			}
			ruleType = allowedRules(colRulesString[:ruleTypeEndPos])
			if string(ruleType[:len(AND_OPERATOR)]) == string(AND_OPERATOR) ||
				string(ruleType[0:len(OR_OPERATOR)]) == string(OR_OPERATOR) {

				// Only one logical operator can be set for each column rules
				// If OR logical operator has been set during parsing than it can´t have AND in the same column rules
				if logicalOperator == OR_OPERATOR && string(ruleType[:2]) == string(AND_OPERATOR) {
					return nil, ErrInvalidRule
				}

				logicalOperator = logicalOperatorType(ruleType[:2])
				ruleType = allowedRules(ruleType[2:])
			}

			if i := slices.Index(ALL_RULES, string(ruleType)); i == -1 {
				return nil, ErrInvalidRule
			}

			colRulesString = colRulesString[ruleTypeEndPos+1:]

			// Retrieves value for rule
			valueEndPos := strings.Index(colRulesString, ")")
			if valueEndPos == -1 {
				return nil, ErrInvalidRule
			}
			ruleValue = colRulesString[:valueEndPos]

			colRule.addRule(fmt.Sprintf("%s%s(%s)", logicalOperator, ruleType, ruleValue))

			if valueEndPos+1 >= len(colRulesString) {
				break
			}
			colRulesString = colRulesString[valueEndPos+1:]
		}

		rulesByCols = append(rulesByCols, colRule)

		if ruleEndPos+1 >= len(rulesInput) {
			break
		}

		// Update rulesInput to remove the already parsed column rules
		rulesInput = rulesInput[ruleEndPos+1:]
	}
	return rulesByCols, nil
}
