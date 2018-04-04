package goscore

import (
	"io/ioutil"
	"math"
	"testing"
)

func Test_parsing(t *testing.T) {
	bb, err := ioutil.ReadFile("./fixtures/neural_network.pmml")
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}
	nn, err := NewNeuralNetwork(bb)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}
	if nn.ActivationFunction != "rectifier" {
		t.Error("Activation Function error", nn.ActivationFunction)
		t.Fail()
		return
	}
	if len(nn.Layers) != 3 {
		t.Error("Number of Layer wrong", nn.Layers)
		t.Fail()
		return
	}
	//t.Log(nn)

}
func TestSoftmaxFunction(t *testing.T) {
	Output := SoftmaxNormalizationMethod([]float64{1.0, 2.0, 3.0, 4.0, 1.0, 2.0, 3.0}...)
	targetOutput := []float64{0.0236405, 0.0642617, 0.174681, 0.474833, 0.0236405, 0.0642617, 0.174681}
	for i, _ := range targetOutput {
		if math.Abs(Output[i]-targetOutput[i]) > 0.01 {
			t.Error(Output[i], targetOutput[i])
			t.Fail()
			return
		}
	}
}
func TestScore1(t *testing.T) {
	input := map[string]interface{}{}
	input["Sepal.Length"] = 5.1
	input["Sepal.Width"] = 3.5
	input["Petal.Length"] = 1.4
	input["Petal.Width"] = 0.2

	bb, err := ioutil.ReadFile("./fixtures/neural_network.pmml")
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}
	nn, err := NewNeuralNetwork(bb)
	if err != nil {
		t.Fail()
		t.Error(err.Error())
		return
	}
	score0, err := nn.Score(input, "0")
	score1, err := nn.Score(input, "1")
	score2, err := nn.Score(input, "2")
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}
	if !(score0 > score1 && score0 > score2) {
		t.Log("Misclassification 1")
		t.Error(score0, score1, score2)
		t.Fail()
		return
	}
	t.Log(input, score0, score1, score2)
	// 7.7,2.6,6.9,2.3
	data := [][]float64{
		[]float64{7, 3.2, 4.7, 1.4},
		[]float64{6.4, 3.2, 4.5, 1.5},
		[]float64{6.9, 3.1, 4.9, 1.5},
		[]float64{5.5, 2.3, 4, 1.3},
		[]float64{6.5, 2.8, 4.6, 1.5},
		[]float64{5.7, 2.8, 4.5, 1.3},
		[]float64{6.3, 3.3, 4.7, 1.6},
		[]float64{4.9, 2.4, 3.3, 1},
		[]float64{6.6, 2.9, 4.6, 1.3},
		[]float64{5.2, 2.7, 3.9, 1.4},
		[]float64{5, 2, 3.5, 1},
		[]float64{5.9, 3, 4.2, 1.5},
		[]float64{6, 2.2, 4, 1},
		[]float64{6.1, 2.9, 4.7, 1.4},
		[]float64{5.6, 2.9, 3.6, 1.3},
		[]float64{6.7, 3.1, 4.4, 1.4},
		[]float64{5.6, 3, 4.5, 1.5},
		[]float64{5.8, 2.7, 4.1, 1},
		[]float64{6.2, 2.2, 4.5, 1.5},
		[]float64{5.6, 2.5, 3.9, 1.1},
		[]float64{5.9, 3.2, 4.8, 1.8},
		[]float64{6.1, 2.8, 4, 1.3},
		[]float64{6.3, 2.5, 4.9, 1.5},
		[]float64{6.1, 2.8, 4.7, 1.2},
		[]float64{6.4, 2.9, 4.3, 1.3},
		[]float64{6.6, 3, 4.4, 1.4},
		[]float64{6.8, 2.8, 4.8, 1.4},
		[]float64{6.7, 3, 5, 1.7},
		[]float64{6, 2.9, 4.5, 1.5},
		[]float64{5.7, 2.6, 3.5, 1},
		[]float64{5.5, 2.4, 3.8, 1.1},
		[]float64{5.5, 2.4, 3.7, 1},
		[]float64{5.8, 2.7, 3.9, 1.2},
		[]float64{6, 2.7, 5.1, 1.6},
		[]float64{5.4, 3, 4.5, 1.5},
		[]float64{6, 3.4, 4.5, 1.6},
		[]float64{6.7, 3.1, 4.7, 1.5},
		[]float64{6.3, 2.3, 4.4, 1.3},
		[]float64{5.6, 3, 4.1, 1.3},
		[]float64{5.5, 2.5, 4, 1.3},
		[]float64{5.5, 2.6, 4.4, 1.2},
		[]float64{6.1, 3, 4.6, 1.4},
		[]float64{5.8, 2.6, 4, 1.2},
		[]float64{5, 2.3, 3.3, 1},
		[]float64{5.6, 2.7, 4.2, 1.3},
		[]float64{5.7, 3, 4.2, 1.2},
		[]float64{5.7, 2.9, 4.2, 1.3},
		[]float64{6.2, 2.9, 4.3, 1.3},
		[]float64{5.1, 2.5, 3, 1.1},
		[]float64{5.7, 2.8, 4.1, 1.3},
	}
	wrongCount := 0
	for _, i := range data {
		input["Sepal.Length"] = i[0]
		input["Sepal.Width"] = i[1]
		input["Petal.Length"] = i[2]
		input["Petal.Width"] = i[3]
		score0, err = nn.Score(input, "0")
		score1, err = nn.Score(input, "1")
		score2, err = nn.Score(input, "2")
		if !(score1 > score0 && score1 > score2) {
			t.Log("Misclassification 2")
			//t.Error(score0, score1, score2)
			wrongCount++
		}
	}
	if wrongCount > 10 {
		t.Fail()
	}
	t.Log(wrongCount)

}
