package goscore

import (
	"encoding/xml"
	"errors"
	"strconv"
)

type truePredicate struct{}
type dummyMiningSchema struct{}

// Node - PMML tree node
type Node struct {
	XMLName            xml.Name
	Attrs              []xml.Attr         `xml:",any,attr"`
	Content            []byte             `xml:",innerxml"`
	Nodes              []Node             `xml:",any"`
	True               truePredicate      `xml:"True"`
	DummyMiningSchema  dummyMiningSchema  `xml:"MiningSchema"`
	SimplePredicate    SimplePredicate    `xml:"SimplePredicate"`
	SimpleSetPredicate SimpleSetPredicate `xml:"SimpleSetPredicate"`
}

// TraverseTree - traverses Node predicates with features and returns score by terminal node
func (n Node) TraverseTree(features map[string]interface{}) (score float64, err error) {
	curr := n.Nodes[0]
	for len(curr.Nodes) > 0 {
		prevID := curr.Attrs[0].Value
		curr = step(curr, features)
		if prevID == curr.Attrs[0].Value {
			break
		}
	}

	if len(curr.Attrs) < 2 {
		return -1, errors.New("Terminal node without score, Node id: " + curr.Attrs[0].Value)
	}
	return strconv.ParseFloat(curr.Attrs[1].Value, 64)
}

func step(curr Node, features map[string]interface{}) Node {
	for _, node := range curr.Nodes {
		if node.XMLName.Local == "True" || node.SimplePredicate.True(features) || node.SimpleSetPredicate.True(features) {
			curr = node
			break
		}
	}
	return curr
}
