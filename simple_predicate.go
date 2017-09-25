package goscore

import (
	"strconv"
)

// SimplePredicate - PMML simple predicate
type SimplePredicate struct {
	Field    string `xml:"field,attr"`
	Operator string `xml:"operator,attr"`
	Value    string `xml:"value,attr"`
}

// True - Evaluates to true if features input is true for SimplePredicate
func (p SimplePredicate) True(features map[string]interface{}) bool {

	if p.Operator == "isMissing" {
		featureValue, exists := features[p.Field]
		return featureValue == "" || featureValue == nil || !exists
	}

	switch featureValue := features[p.Field].(type) {
	case int:
		return numericTrue(p, float64(featureValue))
	case float64:
		return numericTrue(p, featureValue)
	case string:
		if p.Operator == "equal" {
			return p.Value == features[p.Field]
		}
		numericFeatureValue, err := strconv.ParseFloat(featureValue, 64)
		if err == nil {
			return numericTrue(p, numericFeatureValue)
		}
	case bool:
		if p.Operator == "equal" {
			predicateBool, _ := strconv.ParseBool(p.Value)
			return predicateBool == features[p.Field]
		}
	}

	return false
}

func numericTrue(p SimplePredicate, featureValue float64) bool {
	predicateValue, _ := strconv.ParseFloat(p.Value, 64)

	if p.Operator == "equal" {
		return featureValue == predicateValue
	} else if p.Operator == "lessThan" {
		return featureValue < predicateValue
	} else if p.Operator == "lessOrEqual" {
		return featureValue <= predicateValue
	} else if p.Operator == "greaterThan" {
		return featureValue > predicateValue
	} else if p.Operator == "greaterOrEqual" {
		return featureValue >= predicateValue
	}
	return false
}
