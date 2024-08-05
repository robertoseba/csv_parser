package rule

import (
	"errors"
	"testing"
)

func TestRuleFromStr(t *testing.T) {
	expectedFloat1 := 5.0
	expectedFloat2 := 23.0
	expectedFloat3 := 3.0
	expectedFloat4 := 10.0

	tests := []struct {
		name             string
		inputParams      string
		expectedColRules map[string]*ColRules
		expectedError    error
	}{
		{name: "no-rules", inputParams: "", expectedColRules: nil, expectedError: nil},
		{name: "invalid-rule-col-separator", inputParams: "col-eq(5)", expectedColRules: nil, expectedError: ErrInvalidRule},
		{name: "invalid-rule-more-than-one-logical-operator", inputParams: "col:eq(5)||lte(10)&&eq(10)", expectedColRules: nil, expectedError: ErrInvalidRule},

		{name: "two-rules-2-cols", inputParams: "col1:eq(5)||eq(23);col2:!eq(3)&&lt(10)",
			expectedColRules: map[string]*ColRules{
				"col1": {
					logicalOperator: "||",
					column:          "col1",
					isNumber:        true,
					rules: []Rule{
						{value: "5", operator: "eq", floatValue: &expectedFloat1},
						{value: "23", operator: "eq", floatValue: &expectedFloat2},
					},
				},
				"col2": {
					logicalOperator: "&&",
					column:          "col2",
					isNumber:        true,
					rules: []Rule{
						{value: "3", operator: "!eq", floatValue: &expectedFloat3},
						{value: "10", operator: "lt", floatValue: &expectedFloat4},
					},
				},
			},
		},
		{name: "simple-rule-1-col", inputParams: "col1:eq(5)",
			expectedColRules: map[string]*ColRules{
				"col1": {
					logicalOperator: "&&",
					column:          "col1",
					isNumber:        true,
					rules: []Rule{
						{value: "5", operator: "eq", floatValue: &expectedFloat1},
					},
				},
			},
		},
		{name: "rule_with_strings", inputParams: "email:eq(test@email.com);",
			expectedColRules: map[string]*ColRules{
				"email": {
					logicalOperator: "&&",
					column:          "email",
					isNumber:        false,
					rules: []Rule{
						{value: "test@email.com", operator: "eq", floatValue: nil},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rules, err := RulesFromStr(test.inputParams)

			if err != nil && test.expectedError == nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrong error: %v, Got: %v", test.expectedError, err)
			}

			for column, colRules := range rules {
				if column != test.expectedColRules[column].Column() {
					t.Errorf("Wrong column: %v, Got: %v", test.expectedColRules[column], colRules)
				}

				if colRules.IsNumber() != test.expectedColRules[column].IsNumber() {
					t.Errorf("Wrong isNumber: %v, Got: %v", test.expectedColRules[column].IsNumber(), colRules.IsNumber())
				}

				if colRules.logicalOperator != test.expectedColRules[column].logicalOperator {
					t.Errorf("Wrong logicalOperator: %v, Got: %v", test.expectedColRules[column].logicalOperator, colRules.logicalOperator)
				}

				for idx, rule := range colRules.rules {
					if rule.operator != test.expectedColRules[column].rules[idx].operator {
						t.Errorf("Wrong operator: %v, Got: %v", test.expectedColRules[column].rules[idx].operator, rule.operator)
					}

					if rule.value != test.expectedColRules[column].rules[idx].value {
						t.Errorf("Wrong value: %v, Got: %v", test.expectedColRules[column].rules[idx].value, rule.value)
					}

					if *rule.floatValue != *test.expectedColRules[column].rules[idx].floatValue {
						t.Errorf("Wrong floatValue: %v, Got: %v", test.expectedColRules[column].rules[idx].floatValue, rule.floatValue)
					}
				}
			}
		})
	}
}
