package entity

type FilterSelectOption struct {
	Value interface{} `json:"value" bson:"value"`
	Label string      `json:"label" bson:"label"`
}
