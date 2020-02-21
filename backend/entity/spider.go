package entity

type SpiderType struct {
	Type  string `json:"type" bson:"_id"`
	Count int    `json:"count" bson:"count"`
}

type ScrapySettingParam struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type ScrapyItem struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields"`
}
