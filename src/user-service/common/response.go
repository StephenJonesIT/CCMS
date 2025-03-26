/*
 * @File: common.response.go
 * @Description: Defines Response information will be returned to the clients
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package common

type Response struct {
	Data 	interface{} `json:"data"`
	Paging 	interface{}	`json:"paging,omitempty"` 
}

func NewResponse(data interface{}) *Response {
	return &Response{
		Data: data,
	}
}

func NewDetailResponse(data, page interface{}) *Response {
	return &Response{
		Data: data,
		Paging: page,
	}
}