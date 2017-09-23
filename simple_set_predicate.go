package goscore

import (
	"github.com/mattn/go-shellwords"
)

// SimpleSetPredicate - PMML simple set predicate
type SimpleSetPredicate struct {
	Field    string `xml:"field,attr"`
	Operator string `xml:"booleanOperator,attr"`
	Values   string `xml:"Array"`
}

// True - Evaluates to true if features input is true for SimpleSetPredicate
func (p SimpleSetPredicate) True(features map[string]interface{}) bool {
	values, _ := shellwords.Parse(p.Values)

	if p.Operator == "isIn" {
		for _, value := range values {
			if value == features[p.Field] {
				return true
			}
		}
	}
	return false
}
