package goscore

import (
	"strings"
)
type SimpleSetPredicate struct {
	Field     string `xml:"field,attr"`
	Operator  string `xml:"booleanOperator,attr"`
	Values    string `xml:"Array"`
}

func (p SimpleSetPredicate) True(features map[string]string) bool {
	if p.Operator == "isIn" {
		return strings.Contains(p.Values, features[p.Field])
	}
	return true
}