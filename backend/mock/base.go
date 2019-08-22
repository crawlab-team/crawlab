package mock

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
