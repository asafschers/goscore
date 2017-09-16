package goscore

import (
	"encoding/xml"
	"strconv"
)

// RandomForest - PMML Random Forest
type RandomForest struct {
	XMLName            xml.Name
	Trees              []Node	`xml:"MiningModel>Segmentation>Segment>TreeModel"`
}


func (rf RandomForest) Score(features map[string]string) map[string]int {
	scores := map[string]int {}
	for _, tree := range rf.Trees {
		score := strconv.FormatFloat(tree.TraverseTree(features), 'f', -1, 64)
		scores[score] += 1
	}
	return scores
}
