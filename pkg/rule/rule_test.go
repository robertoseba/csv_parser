package rule

import (
	"reflect"
	"testing"
)

func TestRuleFromStr(t *testing.T) {
	tests := []struct {
		name        string
		inputParams string
		expected    []IRule
	}{
		{name: "no-rules", inputParams: "", expected: nil},
		{name: "simple-rule-2-cols", inputParams: "col1:eq(5)||eq(23);col2:!eq(3)&&lt(10)", expected: []IRule{
			&EqRule{Rule: Rule{column: "col1", value: "5", isValueNumber: true}},
			&NotEqRule{Rule: Rule{column: "col2", value: "3", isValueNumber: true}},
		}},
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
