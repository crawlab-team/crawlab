package controllers

type Response[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Error   string `json:"error"`
}

type ListResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Total   int    `json:"total"`
	Data    []T    `json:"data"`
	Error   string `json:"error"`
}
