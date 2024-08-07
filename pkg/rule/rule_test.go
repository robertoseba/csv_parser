package rule

import (
	"testing"
)

type testTable struct {
	testName   string
	inputParam string
	rule       Rule
	expected   bool
}

var floatValue = 1.0

func TestEquals(t *testing.T) {
	tests := []testTable{
		{
			inputParam: "1",
			testName:   "eq_rule-with_floats",
			rule:       Rule{"1", &floatValue, EQ_RULE},
			expected:   true,
		},
		{
			inputParam: "2",
			testName:   "not_eq_rule-with_floats",
			rule:       Rule{"1", nil, EQ_RULE},
			expected:   false,
		},
		{
			inputParam: "hello",
			testName:   "eq_rule-with_strings",
			rule:       Rule{"hello", nil, EQ_RULE},
			expected:   true,
		},
		{
			inputParam: "not_hello",
			testName:   "not_eq_rule-with_strings",
			rule:       Rule{"hello", nil, EQ_RULE},
			expected:   false,
		},
	}

	for _, test := range tests {
		runTest(test, t)
	}
}

func TestNotEquals(t *testing.T) {
	tests := []testTable{
		{
			inputParam: "1",
			testName:   "ne_rule-with_floats",
			rule:       Rule{"1", &floatValue, NE_RULE},
			expected:   false,
		},
		{
			inputParam: "2",
			testName:   "ne_rule-with_floats",
			rule:       Rule{"1", &floatValue, NE_RULE},
			expected:   true,
		},
		{
			inputParam: "hello",
			testName:   "ne_rule-with_strings",
			rule:       Rule{"hello", nil, NE_RULE},
			expected:   false,
		},
		{
			inputParam: "not_hello",
			testName:   "ne_rule-with_strings",
			rule:       Rule{"hello", nil, NE_RULE},
			expected:   true,
		},
	}

	for _, test := range tests {
		runTest(test, t)
	}
}

func TestGreaterThan(t *testing.T) {
	tests := []testTable{
		{
			inputParam: "1.02",
			testName:   "true_with_float",
			rule:       Rule{"1", &floatValue, GT_RULE},
			expected:   true,
		},
		{
			inputParam: "2",
			testName:   "false_with_floats",
			rule:       Rule{"3", nil, GT_RULE},
			expected:   false,
		},
		{
			inputParam: "b",
			testName:   "true_with_string",
			rule:       Rule{"a", nil, GT_RULE},
			expected:   true,
		},
		{
			inputParam: "c",
			testName:   "false_with_strings",
			rule:       Rule{"d", nil, GT_RULE},
			expected:   false,
		},
	}

	for _, test := range tests {
		runTest(test, t)
	}
}

func TestLessThan(t *testing.T) {
	tests := []testTable{
		{
			inputParam: "0.99",
			testName:   "true_with_float",
			rule:       Rule{"1", &floatValue, LT_RULE},
			expected:   true,
		},
		{
			inputParam: "3",
			testName:   "false_with_float",
			rule:       Rule{"1", &floatValue, LT_RULE},
			expected:   false,
		},
		{
			inputParam: "a",
			testName:   "true_with_string",
			rule:       Rule{"b", nil, LT_RULE},
			expected:   true,
		},
		{
			inputParam: "a",
			testName:   "false_with_strings",
			rule:       Rule{"a", nil, LT_RULE},
			expected:   false,
		},
	}

	for _, test := range tests {
		runTest(test, t)
	}
}

func TestGreaterThanOrEqual(t *testing.T) {
	tests := []testTable{
		{
			inputParam: "1.02",
			testName:   "true_with_float",
			rule:       Rule{"1", &floatValue, GTE_RULE},
			expected:   true,
		},
		{
			inputParam: "1.00",
			testName:   "true_with_float",
			rule:       Rule{"1", &floatValue, GTE_RULE},
			expected:   true,
		},
		{
			inputParam: "0.45",
			testName:   "false-with_floats",
			rule:       Rule{"1", &floatValue, GTE_RULE},
			expected:   false,
		},
		{
			inputParam: "a",
			testName:   "true_with_string",
			rule:       Rule{"a", nil, GTE_RULE},
			expected:   true,
		},
		{
			inputParam: "c",
			testName:   "false_with_strings",
			rule:       Rule{"d", nil, GTE_RULE},
			expected:   false,
		},
	}

	for _, test := range tests {
		runTest(test, t)
	}
}

func TestLessThanOrEqual(t *testing.T) {
	tests := []testTable{
		{
			inputParam: "0.99",
			testName:   "true_with_float",
			rule:       Rule{"1", &floatValue, LTE_RULE},
			expected:   true,
		},
		{
			inputParam: "1.00",
			testName:   "true_with_float",
			rule:       Rule{"1", &floatValue, LTE_RULE},
			expected:   true,
		},
		{
			inputParam: "1.001",
			testName:   "false-with_floats",
			rule:       Rule{"1", nil, LTE_RULE},
			expected:   false,
		},
		{
			inputParam: "a",
			testName:   "true_with_string",
			rule:       Rule{"a", nil, LTE_RULE},
			expected:   true,
		},
		{
			inputParam: "c",
			testName:   "false_with_strings",
			rule:       Rule{"b", nil, LTE_RULE},
			expected:   false,
		},
	}

	for _, test := range tests {
		runTest(test, t)
	}
}

func runTest(test testTable, t *testing.T) {
	t.Run(test.testName, func(t *testing.T) {
		if r := test.rule.IsValid(test.inputParam); r != test.expected {
			t.Errorf("Expected %v but got %v", test.expected, r)
		}
	})
}
