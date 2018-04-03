package goscore

import (
	"encoding/xml"
	"errors"
	"math"
	"sort"
	"strings"
)

func init() {
	ActivationFunctions = map[string]ActivationFunction{}

	ActivationFunctions["identity"] = IdentityActivationFunction
	NormalizationMethods = map[string]NormalizationMethod{}
	NormalizationMethods["softmax"] = SoftmaxNormalizationMethod
}

var ActivationFunctions map[string]ActivationFunction
var NormalizationMethods map[string]NormalizationMethod

type MiningField struct {
	Name string `xml:"name,attr"`
}
type NeuralNetWorkStructure struct {
	//DummyMiningSchema dummyMiningSchema `xml:"MiningSchema"`
	//Fields MiningField `xml:"PMML>NeuralNetwork>MiningSchema>MiningField"`
}
type FieldRef struct {
	Field string `xml:"field,attr"`
}

type DerivedField struct {
	DataType     string `xml:"dataType,attr"`
	FieldRef     FieldRef
	NormDiscrete NormDiscrete
}
type NormDiscrete struct {
	Value string `xml:"value,attr"`
}

func (d *DerivedField) GetInputName() string {
	if strings.HasPrefix(d.FieldRef.Field, d.DataType+"(") {
		return d.FieldRef.Field[len(d.DataType)+1 : len(d.FieldRef.Field)-1]
	}
	return d.FieldRef.Field
}

type Contribution struct {
	From   string  `xml:"from,attr"`
	Weight float64 `xml:"weight,attr"`
}
type Neuron struct {
	Id               string         `xml:"id,attr"`
	DerivedFieldType DerivedField   `xml:"DerivedField"`
	Bias             float64        `xml:"bias,attr"`
	Contribution     []Contribution `xml:"Con"`
}
type NeuralInputs struct {
	Input []Neuron `xml:"NeuralInput"`
}
type OutputField struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}
type Outputs struct {
	OutputField []OutputField `xml:"OutputField"`
}
type NeuralOutput struct {
	OutputNeuron string       `xml:"outputNeuron,attr"`
	DerivedField DerivedField `xml:"DerivedField"`
}
type NeuralOutputs struct {
	NeuralOutput []NeuralOutput `xml:"NeuralOutput"`
}
type NeuralLayer struct {
	Neuron              []Neuron `xml:"Neuron"`
	ActivationFunction  string   `xml:"activationFunction,attr"`
	NormalizationMethod string   `xml:"normalizationMethod,attr"`
}
type NeuralNetwork struct {
	XMLName xml.Name
	//Struct  NeuralNetWorkStructure `xml:"PMML>NeuralNetwork"`
	InputLayer          NeuralInputs  `xml:"NeuralNetwork>NeuralInputs"`
	NeuralOutputs       NeuralOutputs `xml:"NeuralNetwork>NeuralOutputs"`
	OutputLayer         Outputs       `xml:"NeuralNetwork>Output"`
	Fields              []MiningField `xml:"NeuralNetwork>MiningSchema>MiningField"`
	Layers              []NeuralLayer `xml:"NeuralNetwork>NeuralLayer"`
	ActivationFunction  string        `xml:"activationFunction,attr"`
	NormalizationMethod string        `xml:"normalizationMethod,attr"`
}

func (nn *NeuralNetwork) Score(feature map[string]interface{}, outputName string) (float64, error) {
	tempMap := map[string]float64{}
	outputMap := map[string]string{}
	for _, np := range nn.InputLayer.Input {
		tempMap[np.Id] = feature[np.DerivedFieldType.GetInputName()].(float64)
	}
	for _, np := range nn.NeuralOutputs.NeuralOutput {
		outputMap[np.DerivedField.NormDiscrete.Value] = np.OutputNeuron
	}
	//initialize actFunc and normMethod
	var actFunc ActivationFunction
	var normMethod NormalizationMethod
	if nn.ActivationFunction == "" {
		actFunc = ActivationFunctions["identity"]
	} else {
		if _, ok := ActivationFunctions[nn.ActivationFunction]; !ok {
			return 0, ActivationFunctionNotImplemented
		} else {
			actFunc = ActivationFunctions["identity"]
		}
	}
	if nn.NormalizationMethod != "" {
		if _, ok := NormalizationMethods[nn.NormalizationMethod]; !ok {
			return 0, NormalizationMethodNotImplemented
		} else {
			normMethod = NormalizationMethods[nn.NormalizationMethod]
		}
	}
	for _, layer := range nn.Layers {
		var actFuncL ActivationFunction
		var normMethodL NormalizationMethod
		//init actFuncL and normMethodL if the layer have different activation function or normalization method
		if layer.ActivationFunction == "" {
			actFuncL = actFunc
		} else {
			if _, ok := ActivationFunctions[layer.ActivationFunction]; !ok {
				return 0, ActivationFunctionNotImplemented
			} else {
				actFuncL = ActivationFunctions["identity"]
			}
		}
		if layer.NormalizationMethod != "" {
			if _, ok := NormalizationMethods[layer.NormalizationMethod]; !ok {
				return 0, NormalizationMethodNotImplemented
			} else {
				normMethodL = NormalizationMethods[layer.NormalizationMethod]
			}
		} else {
			normMethodL = normMethod
		}

		newTemp := map[string]float64{}
		orderedPair := []pair{}

		for _, neuron := range layer.Neuron {
			neuronValue := 0.0
			for _, con := range neuron.Contribution {
				neuronValue += con.Weight * tempMap[con.From]
			}
			neuronValue += neuron.Bias
			newTemp[neuron.Id] = neuronValue
			if actFuncL != nil {
				newTemp[neuron.Id] = actFuncL(neuronValue)
				newPair := pair{neuron.Id, newTemp[neuron.Id]}
				orderedPair = append(orderedPair, newPair)
			}
		}
		sort.Slice(orderedPair, func(i, j int) bool {
			return -1 == strings.Compare(orderedPair[i].name, orderedPair[j].name)
		})
		rawFloat := []float64{}
		for _, p := range orderedPair {
			rawFloat = append(rawFloat, p.value)
		}
		if normMethodL != nil {
			normalizedValue := normMethodL(rawFloat...)
			for i, _ := range orderedPair {
				newTemp[orderedPair[i].name] = normalizedValue[i]
			}

		}
		tempMap = newTemp
	}
	return tempMap[outputMap[outputName]], nil
}

type pair struct {
	name  string
	value float64
}
type ActivationFunction func(float64) float64
type NormalizationMethod func(...float64) []float64

var ActivationFunctionNotImplemented = errors.New("Activation Function Not Implemented Yet")
var NormalizationMethodNotImplemented = errors.New("Normalization Method Not Implemented Yet")

func NewThresHoldFunction(a float64) ActivationFunction {
	return func(b float64) float64 {
		if b < a {
			return 1.0
		} else {
			return 0.0
		}
	}
}
func IdentityActivationFunction(b float64) float64 {
	return b
}
func SoftmaxNormalizationMethod(input ...float64) []float64 {
	hasil := []float64{}
	tempExp := []float64{}
	for _, i := range input {
		tempExp = append(tempExp, math.Exp(i))
	}
	sum := 0.0
	for _, j := range tempExp {
		sum += j
	}
	for _, i := range tempExp {
		hasil = append(hasil, i/sum)
	}
	return hasil
}
