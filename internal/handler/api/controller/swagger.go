package controller

type jsonResponse struct {
	Data    any            `json:"data,omitempty"`
	Code    int            `json:"code,omitempty"`
	Message string         `json:"message,omitempty"`
	Error   *errorResponse `json:"error,omitempty"`
}

type errorResponse struct {
	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
