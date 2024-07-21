package parser

type Validator struct {
	rulesByColumn map[string][]IRule
}

func NewValidator(rules []IRule) *Validator {
	rulesByColumn := make(map[string][]IRule)

	for _, rule := range rules {
		rulesByColumn[rule.Column()] = append(rulesByColumn[rule.Column()], rule)
	}

	validator := &Validator{
		rulesByColumn: rulesByColumn,
	}

	return validator
}

// TODO: Implement multiple rules for the same column
func (v *Validator) IsValid(row *Row) bool {
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
