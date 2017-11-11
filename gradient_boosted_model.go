package goscore

import (
	"encoding/xml"
	"io/ioutil"
	"math"
	"sync"
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

// LoadGradientBoostedModel - Load Gradient Boosted Model PMML file to GradientBoostedModel struct
func LoadGradientBoostedModel(fileName string) (gbm GradientBoostedModel, err error) {
	GradientBoostedModelXml, err := ioutil.ReadFile(fileName)
	if err != nil {
		return gbm, err
	}

	err = xml.Unmarshal([]byte(GradientBoostedModelXml), &gbm)
	if err != nil {
		return gbm, err
	}
	return gbm, nil
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

// ScoreConcurrently - same as Score but concurrent
func (gbm GradientBoostedModel) ScoreConcurrently(features map[string]interface{}) (float64, error) {
	scores := gbm.traverseConcurrently(features)
	sum, err := sumScores(scores, len(gbm.Trees))
	if err != nil {
		return -1, err
	}
	sum += fetchConst(gbm)
	return math.Exp(sum) / (1 + math.Exp(sum)), nil
}

type result struct {
	ErrorName error
	Score     float64
}

func (gbm GradientBoostedModel) traverseConcurrently(features map[string]interface{}) chan result {
	scores := make(chan result, len(gbm.Trees))
	var wg sync.WaitGroup
	wg.Add(len(gbm.Trees))
	for _, tree := range gbm.Trees {
		go func(tree Node, features map[string]interface{}) {
			treeScore, err := tree.TraverseTree(features)
			scores <- result{ErrorName: err, Score: treeScore}
			wg.Done()
		}(tree, features)
	}
	wg.Wait()
	return scores
}

func sumScores(messages chan result, treeCount int) (float64, error) {
	sum := 0.0
	for i := 0; i < treeCount; i++ {
		res := <-messages
		if res.ErrorName != nil {
			return -1, res.ErrorName
		}
		sum += res.Score
	}
	return sum, nil
}

func fetchConst(gbm GradientBoostedModel) (sum float64) {
	if gbm.Version == "4.2" {
		sum = gbm.Constant
	} else if gbm.Version == "4.3" {
		sum = gbm.Target.RescaleConstant
	}
	return sum
}
