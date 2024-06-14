package entity

type DataField struct {
	Key  string `json:"key" bson:"key"`
	Type string `json:"type" bson:"type"`
}
