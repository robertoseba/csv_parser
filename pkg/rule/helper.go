package rule

import "golang.org/x/exp/constraints"

func compareValues[T constraints.Ordered](first, second T, operator string) bool {
	switch operator {
	case "eq":
		return first == second
	case "gt":
		return first > second
	case "lt":
		return first < second
	case "gte":
		return first >= second
	case "lte":
		return first <= second
	case "!eq":
		return first != second
	default:
		panic("invalid operator")
	}
}
