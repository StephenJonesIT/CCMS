/*
 * @File: common.token.go
 * @Description: Defines Token information will be returned to the clients
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package common

type LoginResponse struct {
    Message string `json:"message"`
    Token   string `json:"token"`
}

func NewLoginResponse (mes, token string) *LoginResponse {
	return &LoginResponse{
		Message: mes,
		Token: token,
	}
}