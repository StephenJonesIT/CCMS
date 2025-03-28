/*
 * @File: common.error.go
 * @Description: Defines Error information will be returned to the clients
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package common

type ErrorResponse struct {
	Error 	string 	 `json:"error"`
	Code 	int		 `json:"code,omitempty"`
}

func NewErrorResponse(err string)  *ErrorResponse{
	return &ErrorResponse{
		Error: err,
	}
}

func NewDetailErrorResponse(err string, code int) *ErrorResponse{
	return &ErrorResponse{
		Error: err,
		Code: code,
	}
}
