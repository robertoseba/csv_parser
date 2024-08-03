package rule

// type ComparatorRule struct {
// 	Rule
// 	compare func(any, any) bool
// }
//
// func (rule *ComparatorRule) preProcessNumbers(rowValue string) (ruleValue, colValue float64, err error) {
// 	if rule.isValueNumber {
// 		rowValueFloat, err := strconv.ParseFloat(rowValue, 64)
// 		if err == nil {
// 			ruleValueFloat, _ := strconv.ParseFloat(rule.value, 64)
// 			return ruleValueFloat, rowValueFloat, nil
// 		}
// 	}
//
// 	return 0, 0, fmt.Errorf("cannot convert values to numbers")
// }
//
// func (rule *ComparatorRule) Validate(row *parser.Row) bool {
// 	ruleValue, colValue, err := rule.preProcessNumbers(row.GetColumn(rule.column))
// 	if err == nil {
// 		return rule.compare(ruleValue, colValue)
// 	}
//
// 	return rule.compare(row.GetColumn(rule.column), rule.value)
// }
//
// func NewComparatorRule(ruleStr string) IRule {
// 	if ruleStr == "eq" {
// 		return &ComparatorRule{
// 			Rule: Rule{
// 				column:        "",
// 				value:         "",
// 				isValueNumber: false,
// 			},
// 			compare: func(ruleValue, colValue any) bool { return ruleValue == colValue },
// 		}
// 	}
// }
