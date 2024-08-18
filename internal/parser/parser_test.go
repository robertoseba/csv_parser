package parser

import (
	"errors"
	"testing"
)

func TestParseRules(t *testing.T) {
	expectedFloat1 := 5.0
	expectedFloat2 := 23.0
	expectedFloat3 := 3.0
	expectedFloat4 := 10.0

	tests := []struct {
		name             string
		inputParams      string
		expectedColRules []*ColRules
		expectedError    error
	}{
		{name: "no-rules", inputParams: "", expectedColRules: nil, expectedError: nil},
		{name: "invalid-rule-col-separator", inputParams: "col-eq(5)", expectedColRules: nil, expectedError: ErrInvalidRule},
		{name: "invalid-rule-no-parenthesis", inputParams: "col:eq[5]", expectedColRules: nil, expectedError: ErrInvalidRule},

		{name: "invalid-rule-more-than-one-logical-operator", inputParams: "col:eq(5)||lte(10)&&eq(10)",
			expectedColRules: nil, expectedError: ErrInvalidOperator},

		{name: "invalid-rule-type", inputParams: "col:eq(5)||ltx(10)",
			expectedColRules: nil, expectedError: ErrInvalidRuleType},

		{name: "two-rules-2-cols", inputParams: "col1:eq(5)||eq(23);col2:neq(3)&&lt(10)",
			expectedColRules: []*ColRules{
				{
					logicalOperator: "||",
					column:          "col1",
					isNumber:        true,
					rules: []rule{
						{value: "5", ruleType: eqRule, floatValue: &expectedFloat1},
						{value: "23", ruleType: eqRule, floatValue: &expectedFloat2},
					},
				},
				{
					logicalOperator: "&&",
					column:          "col2",
					isNumber:        true,
					rules: []rule{
						{value: "3", ruleType: neRule, floatValue: &expectedFloat3},
						{value: "10", ruleType: ltRule, floatValue: &expectedFloat4},
					},
				},
			},
		},

		{name: "simple-rule-1-col", inputParams: "col1:eq(5)",
			expectedColRules: []*ColRules{
				{
					logicalOperator: "&&",
					column:          "col1",
					isNumber:        true,
					rules: []rule{
						{value: "5", ruleType: eqRule, floatValue: &expectedFloat1},
					},
				},
			},
		},
		{name: "repeated-col", inputParams: "col1:eq(5);col1:lte(10)",
			expectedColRules: []*ColRules{
				{
					logicalOperator: "&&",
					column:          "col1",
					isNumber:        true,
					rules: []rule{
						{value: "5", ruleType: eqRule, floatValue: &expectedFloat1},
					},
				},
				{
					logicalOperator: "&&",
					column:          "col1",
					isNumber:        true,
					rules: []rule{
						{value: "10", ruleType: lteRule, floatValue: &expectedFloat4},
					},
				},
			},
		},

		{name: "implict-and-logical-operator", inputParams: "col1:eq(5)lte(10)",
			expectedColRules: []*ColRules{
				{
					logicalOperator: "&&",
					column:          "col1",
					isNumber:        true,
					rules: []rule{
						{value: "5", ruleType: eqRule, floatValue: &expectedFloat1},
						{value: "10", ruleType: lteRule, floatValue: &expectedFloat4},
					},
				},
			},
		},

		{name: "rule_with_strings", inputParams: "email:eq(test@email.com);",
			expectedColRules: []*ColRules{
				{
					logicalOperator: "&&",
					column:          "email",
					isNumber:        false,
					rules: []rule{
						{value: "test@email.com", ruleType: eqRule, floatValue: nil},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rules, err := ParseRules(test.inputParams)

			if err != nil && test.expectedError == nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrong error: %v, Got: %v", test.expectedError, err)
			}
			for i, colRules := range rules {
				if rules[i].Column() != test.expectedColRules[i].Column() {
					t.Errorf("Wrong column: %v, Got: %v", test.expectedColRules[i], colRules)
				}

				if colRules.IsNumber() != test.expectedColRules[i].IsNumber() {
					t.Errorf("Wrong isNumber: %v, Got: %v", test.expectedColRules[i].IsNumber(), colRules.IsNumber())
				}
				if colRules.logicalOperator != test.expectedColRules[i].logicalOperator {
					t.Errorf("Wrong logicalOperator: %v, Got: %v", test.expectedColRules[i].logicalOperator, colRules.logicalOperator)
				}

				for idx, rule := range colRules.rules {
					if rule.ruleType != test.expectedColRules[i].rules[idx].ruleType {
						t.Errorf("Wrong operator: %v, Got: %v", test.expectedColRules[i].rules[idx].ruleType, rule.ruleType)
					}
					if rule.value != test.expectedColRules[i].rules[idx].value {
						t.Errorf("Wrong value: %v, Got: %v", test.expectedColRules[i].rules[idx].value, rule.value)
					}
					if rule.floatValue != nil && *rule.floatValue != *test.expectedColRules[i].rules[idx].floatValue {
						t.Errorf("Wrong floatValue: %v, Got: %v", test.expectedColRules[i].rules[idx].floatValue, rule.floatValue)
					}
				}
			}
		})
	}
}
