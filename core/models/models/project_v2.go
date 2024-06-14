package models

type ProjectV2 struct {
	any                    `collection:"projects"`
	BaseModelV2[ProjectV2] `bson:",inline"`
	Name                   string `json:"name" bson:"name"`
	Description            string `json:"description" bson:"description"`
	Spiders                int    `json:"spiders" bson:"-"`
}
