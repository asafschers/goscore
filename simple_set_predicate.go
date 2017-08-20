package goscore

import  (
	"github.com/mattn/go-shellwords"
)
type SimpleSetPredicate struct {
	Field     string `xml:"field,attr"`
	Operator  string `xml:"booleanOperator,attr"`
	Values    string `xml:"Array"`
}

func (p SimpleSetPredicate) True(features map[string]string) bool {
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
