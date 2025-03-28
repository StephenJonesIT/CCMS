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

type CreateOrUpdate struct {
	Message  interface{} `json:"message"`
	Data 	 interface{} `json:"data"`
}

func NewCreateOrUpdate (message, data interface{}) *CreateOrUpdate{
	return &CreateOrUpdate{
		Message: message,
		Data: data,
	}
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
