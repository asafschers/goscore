package goscore

import (
	"encoding/xml"
	"strconv"
)

// SimplePredicate - PMML simple predicate
type SimplePredicate struct {
	Field    string      `xml:"field,attr"`
	Operator string      `xml:"operator,attr"`
	Value    customValue `xml:"value,attr"`
}

type customValue struct {
	NumValue    float64
	StringValue string
}

func (cv *customValue) UnmarshalXMLAttr(attr xml.Attr) error {
	var err error
	cv.NumValue, err = strconv.ParseFloat(attr.Value, 64)
	if err != nil {
		cv.StringValue = attr.Value
	}
	return nil
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
			return p.Value.StringValue == features[p.Field]
		}
		numericFeatureValue, err := strconv.ParseFloat(featureValue, 64)
		if err == nil {
			return numericTrue(p, numericFeatureValue)
		}
	case bool:
		if p.Operator == "equal" {
			predicateBool, _ := strconv.ParseBool(p.Value.StringValue)
			return predicateBool == features[p.Field]
		}
	}

	return false
}

func numericTrue(p SimplePredicate, featureValue float64) bool {
	predicateValue := p.Value.NumValue

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
