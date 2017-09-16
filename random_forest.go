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


func (rf RandomForest) LabelScores(features map[string]string) map[string]float64 {
	scores := map[string]float64 {}
	for _, tree := range rf.Trees {
		score := strconv.FormatFloat(tree.TraverseTree(features), 'f', -1, 64)
		scores[score] += 1
	}
	return scores
}

func (rf RandomForest) Score(features map[string]string, label string) float64 {
	labelScores := rf.LabelScores(features)

	all_count := 0.0
	for  _, value := range labelScores {
		all_count += value
	}

	return labelScores[label] / all_count
}
