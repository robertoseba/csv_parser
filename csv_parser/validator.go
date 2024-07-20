package csv_parser

type Validator struct {
	rules []Rule
	headers *Row
}

func NewValidator(rules []string, headers *Row) (*Validator, error) {
	validator := &Validator{
		headers: headers,
	}

	for _, rule := range rules {
		r, err := NewRule(rule, headers)
		if err != nil {
			return nil, err
		}

		validator.rules = append(validator.rules, *r)
	}

	return validator, nil
}

func (v *Validator) IsValid(row *Row) bool {
	for _, rule := range v.rules {
		if !rule.Validate(row) {
			return false
		}
	}

	return true
}

 