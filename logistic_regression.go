package goscore

import (
	"encoding/xml"
	"errors"
	"math"
)

type PMMLLR struct {
	// struct xml:PMML
	LogisticRegression LogisticRegression `xml:"RegressionModel"`
}

type LogisticRegression struct {
	// struct xml:PMML>RegressionModel
	NormalizationMethod string            `xml:"normalizationMethod,attr"`
	Fields              []MiningField     `xml:"MiningSchema>MiningField"`
	RegressionTable     []RegressionTable `xml:"RegressionTable"`
}

// type MiningField struct {
// 	// struct xml:PMML>RegressionModel>MiningSchema>MiningField
// 	Name string `xml:"name,attr"`
// }

type RegressionTable struct {
	// struct xml:PMML>RegressionModel>RegressionTable
	Intercept           float64            `xml:"intercept,attr"`
	TargetCategory      string             `xml:"targetCategory,attr"`
	NumericPredictor    []NumericPredictor `xml:"NumericPredictor"`
	NumericPredictorMap *map[string]float64
}

type NumericPredictor struct {
	// struct xml:PMML>RegressionModel>RegressionTable>NumbericPredictor
	Name        string  `xml:"name,attr"`
	Coefficient float64 `xml:"coefficient,attr"`
}

type NormalizationMethodMap func(map[string]float64) (map[string]float64, error)

var NormalizationMethodMaps map[string]NormalizationMethodMap

//var NormalizationMethodNotImplemented = errors.New("Normalization Method Not Implemented Yet")

func init() {
	NormalizationMethodMaps = map[string]NormalizationMethodMap{}
	NormalizationMethodMaps["softmax"] = SoftmaxNormalizationMethods
}

// function for compute confidence value
// into probability using softMax function
// input  : map of confidence value with float64 type
// output : map of probability each class with float64 type
func SoftmaxNormalizationMethods(confidence map[string]float64) (map[string]float64, error) {
	if confidence != nil {
		result := map[string]float64{}
		tempExp := []float64{}
		for _, v := range confidence {
			tempExp = append(tempExp, math.Exp(v))
		}
		sum := 0.0
		for _, j := range tempExp {
			sum += j
		}

		i := 0
		for k, _ := range confidence {
			result[k] = tempExp[i]
			i += 1
		}
		return result, nil
	}
	return nil, errors.New("feature is empty")
}

func NewLogisticRegression(source []byte) (*LogisticRegression, error) {
	pmml := PMMLLR{}
	err := xml.Unmarshal(source, &pmml)
	if err != nil {
		return nil, err
	}
	return &pmml.LogisticRegression, nil
}

// func NewLogisticRegressionFromReader(source io.Reader) (*LogisticRegression, error) {
// 	pmml := PMMLLR{}
// 	err := xml.NewDecoder(source).Decode(&pmml)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pmml.LogisticRegression, nil
// }

// method for score test data
// input : 	independent variable with map["var name"]value
//			voting with boolean type
//				true (default)  -> using normalization
//				false			-> without normalization
// return : -label with string type
//			-confident/prob with map type
//			-errors
func (lr *LogisticRegression) Score(args ...interface{}) (string, map[string]float64, error) {

	features := map[string]float64{}
	voting := true

	for _, arg := range args {
		switch t := arg.(type) {
		case map[string]float64:
			features = t
		case bool:
			voting = t
		default:
			return "", nil, errors.New("Unknown argument")
		}
	}
	// calculate confident value using log reg function
	confident := lr.RegressionFunctionContinuous(features)

	if !voting {
		return getMaxMap(confident), confident, nil
	}

	// calculate confident value with normalization method
	var normMethod NormalizationMethodMap
	if lr.NormalizationMethod != "" {
		if _, ok := NormalizationMethods[lr.NormalizationMethod]; !ok {
			return "", nil, NormalizationMethodNotImplemented
		} else {
			normMethod = NormalizationMethodMaps[lr.NormalizationMethod]
		}
	}

	prob, err := normMethod(confident)
	if err != nil {
		return "", nil, err
	}
	return getMaxMap(prob), prob, nil
}

// create map for containing numeric predictor
func (lr *LogisticRegression) SetupNumbericPredictorMap() {
	for i, rt := range lr.RegressionTable {
		m := make(map[string]float64)
		for _, np := range rt.NumericPredictor {
			m[np.Name] = np.Coefficient
		}
		lr.RegressionTable[i].NumericPredictorMap = &m
		//fmt.Println(rt.NumericPredictorMap)
	}
	//fmt.Println(lr.RegressionTable[0].NumericPredictorMap)
}

// method for calculate feature using logistic regression
// function for countinous independent variable
func (lr *LogisticRegression) RegressionFunctionContinuous(features map[string]float64) map[string]float64 {
	confidence := map[string]float64{}

	for _, regressionTable := range lr.RegressionTable {
		var intercept float64
		if regressionTable.Intercept != 0.0 {
			intercept = regressionTable.Intercept
		}

		//fmt.Println(regressionTable.NumericPredictorMap)
		if regressionTable.NumericPredictorMap != nil {
			m := *regressionTable.NumericPredictorMap
			sum := 0.0
			for k, v := range features {
				if c, ok := m[k]; ok {
					sum += v * c
				}
			}
			confidence[regressionTable.TargetCategory] = intercept + sum
		}
	}

	return confidence
}

// method for key with search max value in map
func getMaxMap(feature map[string]float64) string {
	result := ""
	max := -999.999
	for k, v := range feature {
		if result != "" {
			if max < v {
				result = k
				max = v
			}
		} else {
			result = k
			max = v
		}
	}
	return result
}
