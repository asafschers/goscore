package goscore

import (
	"encoding/xml"
	"math"
)

// GradientBoostedModel - GradientBoostedModel PMML
type GradientBoostedModel struct {
	Version  string  `xml:"version,attr"`
	Trees    []Node  `xml:"MiningModel>Segmentation>Segment>MiningModel>Segmentation>Segment>TreeModel"`
	Target   target  `xml:"MiningModel>Segmentation>Segment>MiningModel>Targets>Target"`
	Constant float64 `xml:"MiningModel>Segmentation>Segment>MiningModel>Output>OutputField>Apply>Constant"`
}

type target struct {
	XMLName         xml.Name
	RescaleConstant float64 `xml:"rescaleConstant,attr"`
}

// Score - traverses all trees in GradientBoostedModel with features and returns exp(sum) / (1 + exp(sum))
// where sum is sum of trees + rescale constant
func (gbm GradientBoostedModel) Score(features map[string]interface{}) (float64, error) {

	sum := fetchConst(gbm)

	for _, tree := range gbm.Trees {
		score, err := tree.TraverseTree(features)
		if err != nil {
			return -1, err
		}
		sum += score
	}
	return math.Exp(sum) / (1 + math.Exp(sum)), nil
}
func fetchConst(gbm GradientBoostedModel) (sum float64) {
	if gbm.Version == "4.2" {
		sum = gbm.Constant
	} else if gbm.Version == "4.3" {
		sum = gbm.Target.RescaleConstant
	}
	return sum
}
