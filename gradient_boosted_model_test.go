package goscore_test

import (
	"encoding/xml"
	"github.com/asafschers/goscore"
	"io/ioutil"
	"math"
	"testing"
)

var GradientBoostedModelTests = []struct {
	features map[string]string
	score    float64
}{
	{map[string]string{
		"Sex":      "male",
		"Parch":    "0",
		"Age":      "30",
		"Fare":     "9.6875",
		"Pclass":   "2",
		"SibSp":    "0",
		"Embarked": "Q",
	},
		0.3652639329522468,
	},
	{map[string]string{
		"Sex":      "female",
		"Parch":    "0",
		"Age":      "38",
		"Fare":     "71.2833",
		"Pclass":   "2",
		"SibSp":    "1",
		"Embarked": "C",
	},
		0.4178155014037758,
	},
}

const TOLERANCE = 0.000001

func TestGradientBoostedModel(t *testing.T) {
	GradientBoostedModelXml, err := ioutil.ReadFile("fixtures/gradient_boosted_model.pmml")
	if err != nil {
		panic(err)
	}

	var gbm goscore.GradientBoostedModel
	err = xml.Unmarshal([]byte(GradientBoostedModelXml), &gbm)
	if err != nil {
		panic(err)
	}

	if len(gbm.Trees) != 100 {
		t.Errorf("expected 100 trees, actual %d", len(gbm.Trees))
	}

	for _, tt := range GradientBoostedModelTests {

		actual, _ := gbm.Score(tt.features)

		if diff := math.Abs(actual - tt.score); diff > TOLERANCE {
			t.Errorf("expected %f, actual %f",
				tt.score,
				actual)
		}
	}
}
