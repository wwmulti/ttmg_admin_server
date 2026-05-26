package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type LoginController struct {
	BaseController
}

type LoginResponse struct {
	JSONResponseTpl
	Data models.Account `json:"code"` // 用户信息
}

// @Summary	登陆
// @Param request body services.LoginRequestParams true "请求参数"
// @Success 200 {object} controllers.LoginResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /login [post]
func (c *LoginController) Login() {
	var params services.LoginRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	data, err := (&services.AccountService{}).Login(params)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(data)
}
