package goscore_test

import (
	"encoding/xml"
	"github.com/asafschers/goscore"
	"testing"
	"io/ioutil"
)

var TreeTests = []struct {
	features map[string]string
	score    float64
}{
	{map[string]string{},
		4.3463944950723456E-4,
	},
	{
		map[string]string{"f2": "f2v1"},
		-1.8361380219689046E-4,
	},
	{
		map[string]string{"f2": "f2v1", "f1": "f1v3"},
		-6.237581139073701E-4,
	},
	{
		map[string]string{"f2": "f2v1", "f1": "f1v3", "f4": "0.08"},
		0.0021968294712358194,
	},
	{
		map[string]string{"f2": "f2v1", "f1": "f1v3", "f4": "0.09"},
		-9.198573460887271E-4,
	},
	{
		map[string]string{"f2": "f2v1", "f1": "f1v3", "f4": "0.09", "f3": "f3v2"},
		-0.0021187239505556523,
	},
	{
		map[string]string{"f2": "f2v1", "f1": "f1v3", "f4": "0.09", "f3": "f3v4"},
		-3.3516227414227926E-4,
	},
	{
		map[string]string{"f2": "f2v1", "f1": "f1v4"},
		0.0011015286521365208,
	},
	{
		map[string]string{"f2": "f2v4"},
		0.0022726641744997256,
	},
}

// TODO: test score distribution
// TODO: restore mining schema to pmml
// TODO: errors on unknown value

func TestTree(t *testing.T) {
	treeXml, err := ioutil.ReadFile("fixtures/tree.pmml")
	if err != nil {
		panic(err)
	}

	tree := []byte(treeXml)
	var n goscore.Node
	xml.Unmarshal(tree, &n)

	for _, tt := range TreeTests {
		actual := n.TraverseTree(tt.features)
		if actual != tt.score {
			t.Errorf("expected %f, actual %f",
				tt.score,
				actual)
		}
	}
}
