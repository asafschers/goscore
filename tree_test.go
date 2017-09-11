package goscore_test

import (
	"github.com/asafschers/goscore"
	"testing"
	"encoding/xml"
)

const treeXml = `
<TreeModel functionName="regression" splitCharacteristic="multiSplit">
    <Node id="1">
        <True/>
        <Node id="13" score="4.3463944950723456E-4">
            <SimplePredicate field="f2" operator="isMissing"/>
        </Node>
        <Node id="2">
            <SimpleSetPredicate field="f2" booleanOperator="isIn">
                <Array type="string">f2v1 f2v2 f2v3</Array>
            </SimpleSetPredicate>
            <Node id="11" score="-1.8361380219689046E-4">
                <SimplePredicate field="f1" operator="isMissing"/>
            </Node>
            <Node id="3">
                <SimpleSetPredicate field="f1" booleanOperator="isIn">
                    <Array type="string">f1v1 f1v2 f1v3</Array>
                </SimpleSetPredicate>
                <Node id="9" score="-6.237581139073701E-4">
                    <SimplePredicate field="f4" operator="isMissing"/>
                </Node>
                <Node id="4" score="0.0021968294712358194">
                    <SimplePredicate field="f4" operator="lessThan" value="0.08086312118570185"/>
                </Node>
                <Node id="5">
                    <SimplePredicate field="f4" operator="greaterOrEqual" value="0.08086312118570185"/>
                    <Node id="8" score="-9.198573460887271E-4">
                        <SimplePredicate field="f3" operator="isMissing"/>
                    </Node>
                    <Node id="6" score="-0.0021187239505556523">
                        <SimpleSetPredicate field="f3" booleanOperator="isIn">
                            <Array type="string">f3v1 f3v2 f3v3</Array>
                        </SimpleSetPredicate>
                    </Node>
                    <Node id="7" score="-3.3516227414227926E-4">
                        <SimpleSetPredicate field="f3" booleanOperator="isIn">
                            <Array type="string">f3v4 f3v5 f3v6</Array>
                        </SimpleSetPredicate>
                    </Node>
                </Node>
            </Node>
            <Node id="10" score="0.0011015286521365208">
                <SimpleSetPredicate field="f1" booleanOperator="isIn">
                    <Array type="string">f1v4 f1v5 f1v6</Array>
                </SimpleSetPredicate>
            </Node>
        </Node>
        <Node id="12" score="0.0022726641744997256">
            <SimplePredicate field="f2" operator="equal" value="f2v4"/>
        </Node>
    </Node>
</TreeModel>`

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
	tree := []byte(treeXml)
	var n goscore.Node
	xml.Unmarshal(tree, &n)

	for _, tt := range TreeTests {
		actual := goscore.TraverseTree(n, tt.features)
		if actual != tt.score {
			t.Errorf("expected %f, actual %f",
				tt.score,
				actual)
		}
	}
}
