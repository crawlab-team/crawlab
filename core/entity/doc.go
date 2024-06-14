package entity

type DocItem struct {
	Title    string    `json:"title"`
	Url      string    `json:"url"`
	Path     string    `json:"path"`
	Children []DocItem `json:"children"`
}
