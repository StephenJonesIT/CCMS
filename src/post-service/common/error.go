package common

type ErrorResponse struct {
	Code int `json:"code,omitempty"`
	Message string `json:"message"`
}

func NewErrorResponse(message string) *ErrorResponse{
	return &ErrorResponse{
		Message: message,
	}
}

func NewDetailErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code: code,
		Message: message,
	}
}