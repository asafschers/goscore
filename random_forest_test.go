package goscore_test

import (
	"encoding/xml"
	"errors"
	"github.com/asafschers/goscore"
	"io/ioutil"
	"testing"
)

var RandomForestTests = []struct {
	features map[string]interface{}
	score    float64
	err      error
}{
	{map[string]interface{}{
		"Sex":      "male",
		"Parch":    0,
		"Age":      30,
		"Fare":     9.6875,
		"Pclass":   2,
		"SibSp":    0,
		"Embarked": "Q",
	},
		2.0 / 15.0,
		nil,
	},
	{map[string]interface{}{
		"Sex":      "female",
		"Parch":    0,
		"Age":      38,
		"Fare":     71.2833,
		"Pclass":   2,
		"SibSp":    1,
		"Embarked": "C",
	},
		14.0 / 15.0,
		nil,
	},
	{map[string]interface{}{
		"Sex":      "female",
		"Parch":    0,
		"Age":      38,
		"Fare":     71.2833,
		"Pclass":   2,
		"SibSp":    1,
		"Embarked": "UnknownCategory",
	},
		-1,
		errors.New("Terminal node without score"),
	},
}

func TestRandomForest(t *testing.T) {
	randomForestXml, err := ioutil.ReadFile("fixtures/random_forest.pmml")
	if err != nil {
		panic(err)
	}

	var rf goscore.RandomForest
	err = xml.Unmarshal([]byte(randomForestXml), &rf)
	if err != nil {
		panic(err)
	}

	for _, tt := range RandomForestTests {
		actual, err := rf.Score(tt.features, "1")

		if err != nil {
			if tt.err == nil {
				t.Errorf("expected no error, actual: %s",
					err)
			} else if tt.err.Error() != err.Error() {
				t.Errorf("expected error %s, actual: %s",
					tt.err.Error(),
					err)
			}
		}

		if err == nil && actual != tt.score {
			t.Errorf("expected %f, actual %f",
				tt.score,
				actual)
		}
	}
}
