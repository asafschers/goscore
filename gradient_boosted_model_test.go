package goscore_test

import (
	"github.com/asafschers/goscore"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

var _ = Describe("GradientBoostedModel", func() {
	var (
		lowScoreFeatures, highScoreFeatures map[string]interface{}
		highScore, lowScore                 float64
		tolerance                           float64 = 0.000000001
		gbm                                 goscore.GradientBoostedModel
		first                               bool = true
		err                                 error
	)

	BeforeSuite(func() {
		lowScoreFeatures = map[string]interface{}{
			"Sex":      "male",
			"Parch":    0,
			"Age":      30,
			"Fare":     9.6875,
			"Pclass":   2,
			"SibSp":    0,
			"Embarked": "Q",
		}
		lowScore = 0.3652639329522468

		highScoreFeatures = map[string]interface{}{
			"Sex":      "female",
			"Parch":    0,
			"Age":      38,
			"Fare":     71.2833,
			"Pclass":   2,
			"SibSp":    1,
			"Embarked": "C",
		}
		highScore = 0.4178155014037758

		if first {
			gbm, err = goscore.LoadGradientBoostedModel("fixtures/gradient_boosted_model.pmml")
			if err != nil {
				panic(err)
			}
			first = false
		}
	})

	Describe("Loads model", func() {
		It("Loads all trees", func() {
			Expect(len(gbm.Trees)).To(Equal(100))
		})
	})

	Describe("Scores", func() {
		It("Scores low", func() {
			Expect(gbm.Score(lowScoreFeatures)).To(BeNumerically("~", lowScore, tolerance))
		})

		It("Scores low concurrently", func() {
			Expect(gbm.ScoreConcurrently(lowScoreFeatures)).To(BeNumerically("~", lowScore, tolerance))
		})

		It("Scores high", func() {
			Expect(gbm.Score(highScoreFeatures)).To(BeNumerically("~", highScore, tolerance))
		})

		It("Scores high concurrently", func() {
			Expect(gbm.ScoreConcurrently(highScoreFeatures)).To(BeNumerically("~", highScore, tolerance))
		})
	})
})

func TestGradientBoostedModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GradientBoostedModel Suite")
}
