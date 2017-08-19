package goscore_test

import (
	"testing"
	"encoding/xml"
	"github.com/asafschers/goscore"
)

var simpleSetPredicateTests = []struct {
	predicate	[]byte
	features    map[string]string
	expected	bool
}{
	{[]byte(`<SimpleSetPredicate field="f1" booleanOperator="isIn">
                          <Array n="6" type="string">v1 v2 v3</Array>
                          </SimpleSetPredicate>`),
		map[string]string{"f1": "v3"},
		true},
	{[]byte(`<SimpleSetPredicate field="f1" booleanOperator="isIn">
                          <Array n="6" type="string">v1 v2 v3</Array>
                          </SimpleSetPredicate>`),
		map[string]string{"f1": "v4"},
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
