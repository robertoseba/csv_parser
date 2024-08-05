package rule

import "github.com/robertoseba/csv_parser/pkg/row"

type logicalOperatorType string

const AND_OPERATOR logicalOperatorType = "&&"
const OR_OPERATOR logicalOperatorType = "||"

type ColRules struct {
	column          string
	rules           []Rule
	logicalOperator logicalOperatorType
	isNumber        bool
}

func (r *ColRules) Column() string {
	return r.column
}

func (r *ColRules) IsNumber() bool {
	return r.isNumber
}

func (r *ColRules) IsValid(row *row.Row) bool {
	result := true

	for _, rule := range r.rules {
		if rule.IsValid(row.GetColumn(r.column)) {
			if r.logicalOperator == OR_OPERATOR {
				break
			}
		} else {
			if r.logicalOperator == AND_OPERATOR {
				result = false
				break
			}
		}
	}

	return result
}
