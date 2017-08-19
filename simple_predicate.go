package goscore

import "strconv"

type SimplePredicate struct {
	Field     string `xml:"field,attr"`
	Operator  string `xml:"operator,attr"`
	Value     string `xml:"value,attr"`
}

func (p SimplePredicate) True(features map[string]string) bool {
	predicateValue, _ := strconv.ParseFloat(p.Value, 64)
	featureValue, _ := strconv.ParseFloat(features[p.Field], 64)

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
	} else if p.Operator == "isMissing" {
		_, exists := features[p.Field];
		return !exists
	}
	return false
}

