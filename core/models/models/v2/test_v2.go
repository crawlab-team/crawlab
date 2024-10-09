package models

type TestModelV2 struct {
	any                      `collection:"testmodels"`
	BaseModelV2[TestModelV2] `bson:",inline"`
	Name                     string `json:"name" bson:"name"`
}
