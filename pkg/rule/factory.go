package rule

import (
	"errors"
	"fmt"
	"regexp"
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
	parseRules(ruleInput)

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

func parseRules(rulesInput string) ([]ColRules, error) {
	if strings.Trim(rulesInput, " ") == "" {
		return nil, nil
	}
	col := ""
	ruleValue := ""
	logicalOperator := AND_OPERATOR
	var ruleType allowedRules

	fmt.Println("input: ", rulesInput)
	for {
		//Parses end of the each rule
		ruleEndPos := strings.Index(rulesInput, RULE_SEPARATOR)
		if ruleEndPos == -1 {
			ruleEndPos = len(rulesInput)

		}

		colRulesString := rulesInput[:ruleEndPos]
		fmt.Printf("rule string: %s\n", colRulesString)

		//Parses column
		colEndPos := strings.Index(colRulesString, COL_RULE_SEPARATOR)
		if colEndPos == -1 {
			return nil, ErrInvalidRule
		}
		col = colRulesString[:colEndPos]
		fmt.Println("col: ", col)

		colRulesString = colRulesString[colEndPos+1:]
		for {
			ruleTypeEndPos := strings.Index(colRulesString, "(")
			if ruleTypeEndPos == -1 {
				return nil, ErrInvalidRule
			}
			ruleType = allowedRules(colRulesString[:ruleTypeEndPos])
			if string(ruleType[:len(AND_OPERATOR)]) == string(AND_OPERATOR) ||
				string(ruleType[0:len(OR_OPERATOR)]) == string(OR_OPERATOR) {

				// Only one logical operator can be set for each column rules
				// If OR logical operator has been set during parsing than it can´t have AND in the same column rules
				if logicalOperator == OR_OPERATOR || string(ruleType[:2]) == string(AND_OPERATOR) {
					return nil, ErrInvalidRule
				}

				logicalOperator = logicalOperatorType(ruleType[:2])
				ruleType = allowedRules(ruleType[2:])
			}
			fmt.Println("logical operator:", logicalOperator)
			//TODO: validate ruleType here
			fmt.Println("type:", ruleType)

			colRulesString = colRulesString[ruleTypeEndPos+1:]

			// retrieve value for rule
			valueEndPos := strings.Index(colRulesString, ")")
			if valueEndPos == -1 {
				return nil, ErrInvalidRule
			}
			ruleValue = colRulesString[:valueEndPos]
			fmt.Println("value", ruleValue)

			//Create rule here and append

			if valueEndPos+1 >= len(colRulesString) {
				break
			}
			colRulesString = colRulesString[valueEndPos+1:]
		}

		// Create colRule here
		if ruleEndPos+1 >= len(rulesInput) {
			break
		}
		rulesInput = rulesInput[ruleEndPos+1:]
	}
	return nil, nil
}
