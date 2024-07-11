package models

type TokenV2 struct {
	any                  `collection:"tokens"`
	BaseModelV2[TokenV2] `bson:",inline"`
	Name                 string `json:"name" bson:"name"`
	Token                string `json:"token" bson:"token"`
}
