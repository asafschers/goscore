package goscore

import (
	"encoding/xml"
	"fmt"
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
	nn := NeuralNetwork{}
	err = xml.Unmarshal(bb, &nn)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}
	fmt.Println(nn.NeuralOutputs)
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
func TestScore(t *testing.T) {
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
	nn := NeuralNetwork{}
	err = xml.Unmarshal(bb, &nn)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
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
		t.Log("Misclassification")
		t.Fail()
		return
	}
	t.Log(score0, score1, score2)

}
