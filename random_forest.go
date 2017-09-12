package goscore

import "encoding/xml"

// RandomForest - PMML Random Forest
type RandomForest struct {
	XMLName            xml.Name
	Trees              []Node	`xml:"MiningModel>Segmentation>Segment>TreeModel"`
}


func (rf RandomForest) Score(features map[string]string) float64 {
	return 3
}
