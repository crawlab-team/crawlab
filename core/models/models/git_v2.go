package models

type GitV2 struct {
	any                `collection:"gits"`
	BaseModelV2[GitV2] `bson:",inline"`
	Url                string     `json:"url" bson:"url"`
	Name               string     `json:"name" bson:"name"`
	AuthType           string     `json:"auth_type" bson:"auth_type"`
	Username           string     `json:"username" bson:"username"`
	Password           string     `json:"password" bson:"password"`
	CurrentBranch      string     `json:"current_branch" bson:"current_branch"`
	Status             string     `json:"status" bson:"status"`
	Error              string     `json:"error" bson:"error"`
	Spiders            []SpiderV2 `json:"spiders,omitempty" bson:"-"`

	// settings
	AutoPull bool `json:"auto_pull" bson:"auto_pull"`
}
