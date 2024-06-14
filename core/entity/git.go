package entity

type GitPayload struct {
	Paths         []string `json:"paths"`
	CommitMessage string   `json:"commit_message"`
	Branch        string   `json:"branch"`
	Tag           string   `json:"tag"`
}

type GitConfig struct {
	Url string `json:"url" bson:"url"`
}
