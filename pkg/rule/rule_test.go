package rule

import (
	"reflect"
	"testing"
)

func TestRuleFromStr(t *testing.T) {
	tests := []struct {
		name        string
		inputParams string
		expected    map[string]*ColRules
	}{
		{name: "no-rules", inputParams: "", expected: nil},
		{name: "simple-rule-2-cols", inputParams: "col1:eq(5)||eq(23);col2:!eq(3)&&lt(10)",
			expected: make(map[string]*ColRules, 2)}, {
			"col1": &ColRules{
				logicalOperator: "||",
				column:          "col1",
				Rules: []*Rule{
					&Rule{value: "5", operator: "eq", floatValue: 5},
					&Rule{value: "23", operator: "eq", floatValue: 23},
				},
			},
			"col2": &ColRules{
				logicalOperator: "&&",
				column:          "col2",
				Rules: []*Rule{
					&Rule{value: "3", operator: "!eq", floatValue: 3},
					&Rule{value: "10", operator: "lt", floatValue: 10},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rules, err := RulesFromStr(test.inputParams)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			for i, rule := range rules {
				if rule.Column() != test.expected[i].Column() {
					t.Errorf("Wrong column: %v, Got: %v", test.expected[i], rule)
				}

				if rule.IsNumber() != test.expected[i].IsNumber() {
					t.Errorf("Wrong isNumber: %v, Got: %v", test.expected[i], rule)
				}

				if reflect.TypeOf(rule) != reflect.TypeOf(test.expected[i]) {
					t.Errorf("Wrong rule type: %T, Got: %T", test.expected[i], rule)
				}
			}
		})
	}
}
