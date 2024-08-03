package rule

import (
	"github.com/robertoseba/csv_parser/pkg/row"
)

type EqRule struct {
	Rule
}

func (rule *EqRule) Validate(row *row.Row) bool {
	ruleValue, colValue, err := rule.preProcessNumbers(row.GetColumn(rule.column))
	if err == nil {
		return ruleValue == colValue
	}

	return rule.value == row.GetColumn(rule.column)
}
