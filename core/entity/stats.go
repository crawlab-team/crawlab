package entity

type StatsDailyItem struct {
	Date    string `json:"date" bson:"_id"`
	Tasks   int64  `json:"tasks" bson:"tasks"`
	Results int64  `json:"results" bson:"results"`
}

type StatsTasksByStatusItem struct {
	Status string `json:"status" bson:"_id"`
	Tasks  int64  `json:"tasks" bson:"tasks"`
}
