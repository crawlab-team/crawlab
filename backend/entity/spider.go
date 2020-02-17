package entity

type SpiderType struct {
	Type  string `json:"type" bson:"_id"`
	Count int    `json:"count" bson:"count"`
}

type ScrapySettingParam struct {
	Key   string
	Value interface{}
	Type  string
}
