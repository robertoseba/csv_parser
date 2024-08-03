package rule

import "github.com/robertoseba/csv_parser/pkg/row"

type Validator struct {
	rulesByColumn map[string][]IRule
}

func NewValidator(rules []IRule) *Validator {
	rulesByColumn := make(map[string][]IRule)

	for _, rule := range rules {
		rulesByColumn[rule.Column()] = append(rulesByColumn[rule.Column()], rule)
	}

	return &Validator{
		rulesByColumn: rulesByColumn,
	}
}

// TODO: Implement multiple rules for the same column
func (v *Validator) IsValid(row *row.Row) bool {
	for _, rules := range v.rulesByColumn {
		for _, rule := range rules {
			if !rule.Validate(row) {
				return false
			}
		}
	}

	return true
}

func (v *Validator) Columns() []string {
	columns := make([]string, 0, len(v.rulesByColumn))

	for col := range v.rulesByColumn {
		columns = append(columns, col)
	}

	return columns
}
