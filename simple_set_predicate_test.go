package goscore_test

import (
	"encoding/xml"
	"github.com/asafschers/goscore"
	"testing"
)

const simpleSetPredicate1 = `<SimpleSetPredicate field="f1" booleanOperator="isIn">
                            <Array n="6" type="string">v1 v2 v3</Array>
                            </SimpleSetPredicate>`
const simpleSetPredicate2 = `<SimpleSetPredicate field="f2" booleanOperator="isIn">
							  <Array n="6" type="string">"Missing"   "No Match"</Array>
							  </SimpleSetPredicate>`

var simpleSetPredicateTests = []struct {
	predicate []byte
	features  map[string]interface{}
	expected  bool
}{
	{[]byte(simpleSetPredicate1),
		map[string]interface{}{"f1": "v3"},
		true},
	{[]byte(simpleSetPredicate1),
		map[string]interface{}{"f1": "v4"},
		false},
	{[]byte(simpleSetPredicate2),
		map[string]interface{}{"f2": "No Match"},
		true},
	{[]byte(simpleSetPredicate2),
		map[string]interface{}{"f2": "Match"},
		false},
}

func TestSimpleSetPredicate(t *testing.T) {

	for _, tt := range simpleSetPredicateTests {
		var predicate goscore.SimpleSetPredicate
		xml.Unmarshal(tt.predicate, &predicate)

		actual := predicate.True(tt.features)
		if actual != tt.expected {
			t.Errorf("Predicate: %s %s %s, Feature value : %s, expected %t, actual %t",
				predicate.Field,
				predicate.Operator,
				predicate.Values,
				tt.features[predicate.Field],
				tt.expected,
				actual)
		}
	}
}
