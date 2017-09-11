package goscore

import (
	"bytes"
	"encoding/xml"
	"strconv"
)

type truePredicate struct{}

// Node - PMML tree after unmarshal
type Node struct {
	XMLName            xml.Name
	Attrs              []xml.Attr         `xml:"-"`
	Content            []byte             `xml:",innerxml"`
	Nodes              []Node             `xml:",any"`
	True               truePredicate      `xml:"True"`
	SimplePredicate    SimplePredicate    `xml:"SimplePredicate"`
	SimpleSetPredicate SimpleSetPredicate `xml:"SimpleSetPredicate"`
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node Node
	return d.DecodeElement((*node)(n), &start)
}

func GenerateTree(data []byte, n *Node) {
	buffer := bytes.NewBuffer(data)
	decoder := xml.NewDecoder(buffer)

	err := decoder.Decode(&n)
	if err != nil {
		panic(err)
	}
}

func TraverseTree(n Node, features map[string]string) (score float64) {
	curr := n.Nodes[0]
	for len(curr.Nodes) > 0 {
		prevID := curr.Attrs[0].Value
		curr = step(curr, features)
		if prevID == curr.Attrs[0].Value {
			break
		}
	}
	score, _ = strconv.ParseFloat(curr.Attrs[1].Value, 64)
	return score
}

func step(curr Node, features map[string]string) Node {
	for _, node := range curr.Nodes {
		if node.XMLName.Local == "True" || node.SimplePredicate.True(features) || node.SimpleSetPredicate.True(features) {
			curr = node
			break
		}
	}
	return curr
}
