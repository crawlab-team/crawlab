package models

type VariableV2 struct {
	any                     `collection:"variables"`
	BaseModelV2[VariableV2] `bson:",inline"`
	Key                     string `json:"key" bson:"key"`
	Value                   string `json:"value" bson:"value"`
	Remark                  string `json:"remark" bson:"remark"`
}
