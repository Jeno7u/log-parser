package dto

type ErrorResponse struct {
	Error string `json:"error"`
}

type ParseResponse struct {
	LogID string `json:"log_id"`
}
