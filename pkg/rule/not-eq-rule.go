package rule

import (
	"github.com/robertoseba/csv_parser/pkg/row"
)

type NotEqRule struct {
	Rule
}

func (rule *NotEqRule) Validate(row *row.Row) bool {
	ruleValue, colValue, err := rule.preProcessNumbers(row.GetColumn(rule.column))
	if err == nil {
		return ruleValue != colValue
	}

	return rule.value != row.GetColumn(rule.column)
}
