package common


type Response struct {
	Message string   `json:"message,omitempty"`
	Data interface{} `json:"data"`
}

func NewResponse(data interface{}) *Response {
	return &Response{
		Data: data,
	}
}
func NewDetailResponse(message string,data interface{}) *Response{
	return &Response{
		Message: message,
		Data: data,
	}
}
