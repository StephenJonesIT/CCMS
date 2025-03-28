/*
 * @File: common.login.go
 * @Description: Defines Login information request from clients
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */

package common

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
