package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

type ListResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Total   int         `json:"total"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

type ListRequestData struct {
	PageNum  int    `form:"page_num" json:"page_num"`
	PageSize int    `form:"page_size" json:"page_size"`
	SortKey  string `form:"sort_key" json:"sort_key"`
	Status   string `form:"status" json:"status"`
	Keyword  string `form:"keyword" json:"keyword"`
}

type BatchRequestPayload struct {
	Ids []primitive.ObjectID `form:"ids" json:"ids"`
}

type BatchRequestPayloadWithStringData struct {
	Ids    []primitive.ObjectID `form:"ids" json:"ids"`
	Data   string               `form:"data" json:"data"`
	Fields []string             `form:"fields" json:"fields"`
}

type FileRequestPayload struct {
	Path    string `json:"path" form:"path"`
	NewPath string `json:"new_path" form:"new_path"`
	Data    string `json:"data" form:"data"`
}
