package goscore

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_parsinglogreg(t *testing.T) {
	bb, err := ioutil.ReadFile("./fixtures/logistic_regression.xml")
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}
	lr, err := NewLogisticRegression(bb)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}

	lr.SetupNumbericPredictorMap()

	features1 := map[string]float64{}
	features1["x0"] = 0.1
	features1["x1"] = 0.1
	features1["x2"] = 0.1
	features1["x3"] = 0.1

	label, confidence, err := lr.Score(features1)
	fmt.Println(label)
	fmt.Println(confidence)
	fmt.Println(err)

	label, confidence, err = lr.Score(features1, false)
	fmt.Println(label)
	fmt.Println(confidence)
	fmt.Println(err)

	features0 := map[string]float64{}
	label, confidence, err = lr.Score(features0)
	fmt.Println(label)
	fmt.Println(confidence)
	fmt.Println(err)

	label, confidence, err = lr.Score(features0, false)
	fmt.Println(label)
	fmt.Println(confidence)
	fmt.Println(err)

}
