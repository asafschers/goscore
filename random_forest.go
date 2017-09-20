package goscore

import (
	"encoding/xml"
	"strconv"
)

// RandomForest - PMML Random Forest
type RandomForest struct {
	XMLName xml.Name
	Trees   []Node `xml:"MiningModel>Segmentation>Segment>TreeModel"`
}

// LabelScores - traverses all trees in RandomForest with features and maps result
// labels to how many trees returned those label
func (rf RandomForest) LabelScores(features map[string]string) (scores map[string]float64, err error) {
	scores = map[string]float64{}
	for _, tree := range rf.Trees {
		var score float64
		score, err = tree.TraverseTree(features)
		if err != nil {
			return scores, err
		}
		scoreString := strconv.FormatFloat(score, 'f', -1, 64)
		scores[scoreString] += 1
	}
	return scores, err
}

// Score - traverses all trees in RandomForest with features and returns ratio of
// given label results count to all results count
func (rf RandomForest) Score(features map[string]string, label string) (float64, error) {
	labelScores, err := rf.LabelScores(features)

	allCount := 0.0
	for _, value := range labelScores {
		allCount += value
	}

	return labelScores[label] / allCount, err
}
