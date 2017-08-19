package goscore

import  (
	"strings"
	"github.com/mattn/go-shellwords"
)
type SimpleSetPredicate struct {
	Field     string `xml:"field,attr"`
	Operator  string `xml:"booleanOperator,attr"`
	Values    string `xml:"Array"`
}

func (p SimpleSetPredicate) True(features map[string]string) bool {
	values := setValues(p)

	if p.Operator == "isIn" {
		for _, value := range values {

			if value == features[p.Field] {
				return true
			}
		}
	}
	return false
}
func setValues(p SimpleSetPredicate) []string {
	var values []string
	if strings.Contains(p.Values, `"`) {
		values, _ = shellwords.Parse(p.Values)
	} else {
		values = strings.Split(p.Values, " ")
	}
	return values
}