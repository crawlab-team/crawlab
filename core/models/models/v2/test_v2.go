package models

type TestModel struct {
	any                    `collection:"testmodels"`
	BaseModelV2[TestModel] `bson:",inline"`
	Name                   string `json:"name" bson:"name"`
}
