package rule

import (
	"testing"
)

type testTable struct {
	testName   string
	inputParam string
	rule       rule
	expected   bool
}

var floatValue = 1.0

func TestEquals(t *testing.T) {
	tests := []testTable{
		{
			inputParam: "1",
			testName:   "eq_rule-with_floats",
			rule:       rule{"1", &floatValue, eqRule},
			expected:   true,
		},
		{
			inputParam: "2",
			testName:   "not_eq_rule-with_floats",
			rule:       rule{"1", nil, eqRule},
			expected:   false,
		},
		{
			inputParam: "hello",
			testName:   "eq_rule-with_strings",
			rule:       rule{"hello", nil, eqRule},
			expected:   true,
		},
		{
			inputParam: "not_hello",
			testName:   "not_eq_rule-with_strings",
			rule:       rule{"hello", nil, eqRule},
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
			rule:       rule{"1", &floatValue, neRule},
			expected:   false,
		},
		{
			inputParam: "2",
			testName:   "ne_rule-with_floats",
			rule:       rule{"1", &floatValue, neRule},
			expected:   true,
		},
		{
			inputParam: "hello",
			testName:   "ne_rule-with_strings",
			rule:       rule{"hello", nil, neRule},
			expected:   false,
		},
		{
			inputParam: "not_hello",
			testName:   "ne_rule-with_strings",
			rule:       rule{"hello", nil, neRule},
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
			rule:       rule{"1", &floatValue, gtRule},
			expected:   true,
		},
		{
			inputParam: "2",
			testName:   "false_with_floats",
			rule:       rule{"3", nil, gtRule},
			expected:   false,
		},
		{
			inputParam: "b",
			testName:   "true_with_string",
			rule:       rule{"a", nil, gtRule},
			expected:   true,
		},
		{
			inputParam: "c",
			testName:   "false_with_strings",
			rule:       rule{"d", nil, gtRule},
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
			rule:       rule{"1", &floatValue, ltRule},
			expected:   true,
		},
		{
			inputParam: "3",
			testName:   "false_with_float",
			rule:       rule{"1", &floatValue, ltRule},
			expected:   false,
		},
		{
			inputParam: "a",
			testName:   "true_with_string",
			rule:       rule{"b", nil, ltRule},
			expected:   true,
		},
		{
			inputParam: "a",
			testName:   "false_with_strings",
			rule:       rule{"a", nil, ltRule},
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
			rule:       rule{"1", &floatValue, gteRule},
			expected:   true,
		},
		{
			inputParam: "1.00",
			testName:   "true_with_float",
			rule:       rule{"1", &floatValue, gteRule},
			expected:   true,
		},
		{
			inputParam: "0.45",
			testName:   "false-with_floats",
			rule:       rule{"1", &floatValue, gteRule},
			expected:   false,
		},
		{
			inputParam: "a",
			testName:   "true_with_string",
			rule:       rule{"a", nil, gteRule},
			expected:   true,
		},
		{
			inputParam: "c",
			testName:   "false_with_strings",
			rule:       rule{"d", nil, gteRule},
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
			rule:       rule{"1", &floatValue, lteRule},
			expected:   true,
		},
		{
			inputParam: "1.00",
			testName:   "true_with_float",
			rule:       rule{"1", &floatValue, lteRule},
			expected:   true,
		},
		{
			inputParam: "1.001",
			testName:   "false-with_floats",
			rule:       rule{"1", nil, lteRule},
			expected:   false,
		},
		{
			inputParam: "a",
			testName:   "true_with_string",
			rule:       rule{"a", nil, lteRule},
			expected:   true,
		},
		{
			inputParam: "c",
			testName:   "false_with_strings",
			rule:       rule{"b", nil, lteRule},
			expected:   false,
		},
	}

	for _, test := range tests {
		runTest(test, t)
	}
}

func runTest(test testTable, t *testing.T) {
	t.Run(test.testName, func(t *testing.T) {
		if r := test.rule.isValid(test.inputParam); r != test.expected {
			t.Errorf("Expected %v but got %v", test.expected, r)
		}
	})
}
