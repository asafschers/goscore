package goscore

import (
	"encoding/xml"
	"math"
)

// GradientBoostedModel - GradientBoostedModel PMML
type GradientBoostedModel struct {
	XMLName xml.Name
	Trees   []Node `xml:"MiningModel>Segmentation>Segment>MiningModel>Segmentation>Segment>TreeModel"`
	Target  target `xml:"MiningModel>Segmentation>Segment>MiningModel>Targets>Target"`
}

type target struct {
	XMLName         xml.Name
	RescaleConstant float64 `xml:"rescaleConstant,attr"`
}

// Score - traverses all trees in GradientBoostedModel with features and returns exp(sum) / (1 + exp(sum))
// where sum is sum of trees + rescale constant
func (gbm GradientBoostedModel) Score(features map[string]interface{}) (float64, error) {
	sum := gbm.Target.RescaleConstant
	for _, tree := range gbm.Trees {
		score, err := tree.TraverseTree(features)
		if err != nil {
			return -1, err
		}
		sum += score
	}
	return math.Exp(sum) / (1 + math.Exp(sum)), nil
}
