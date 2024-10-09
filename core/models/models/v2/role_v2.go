package models

type RoleV2 struct {
	any                 `collection:"roles"`
	BaseModelV2[RoleV2] `bson:",inline"`
	Key                 string `json:"key" bson:"key"`
	Name                string `json:"name" bson:"name"`
	Description         string `json:"description" bson:"description"`
}
