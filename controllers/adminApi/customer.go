package controllers

import (
	"api/services"
	"encoding/json"
)

type CustomerController struct {
	BaseController
}

// @Summary	GMs上下分操作列表
// @Param role_id query int false "角色ID"
// @Param money query int false "金额"
// @Param operate_type query int false "操作类型ID"
// @Param if_inner_proxy query int false "是否模拟号"
// @Param status query int false "审核状态"
// @Param type query int false "操作来源"
// @Param channel_remark query string false "渠道备注"
// @Param channel_role_id query int false "渠道ID"
// @Param check_user query int false "审核人ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /gmSendMoneyList [get]
func (c *CustomerController) GmSendMoneyList() {
	request := services.GmSendMoneyListRequest{}
	c.ParseForm(&request)
	request.RawQuery = c.Ctx.Request.URL.Query()
	request.PackageIds = c.GetPackageIds(int(c.UserId))

	customerService := new(services.CustomerService)
	data, err := customerService.GmSendMoneyList(request, c.Lang)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	上下分操作
// @Param	request body services.GmOperateAddRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /gmOperateAdd [post]
func (c *CustomerController) GmOperateAdd() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	customerService := new(services.CustomerService)
	request := services.GmOperateAddRequest{
		Descript:    getString(params, "descript"),
		Money:       getString(params, "money"),
		OperateType: getInt(params, "operate_type"),
		Password:    getString(params, "password"),
		RoleId:      getString(params, "role_id"),
		WageMul:     getFloat(params, "wage_mul"),
		AdminId:     int(c.UserId),
		Ip:          c.Ctx.Input.IP(),
		PackageIds:  c.GetPackageIds(int(c.UserId)),
	}
	err = customerService.GmOperateAdd(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type GmOperateCheckRequest struct {
	Id string `json:"id"`
}

// @Summary	上下分审核
// @Param	request body services.GmOperateCheckRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /gmOperateCheck [post]
func (c *CustomerController) GmOperateCheck() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	customerService := new(services.CustomerService)
	idStr := getString(params, "id")
	err = customerService.GmOperateCheck(idStr, int(c.UserId), c.UserName)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type GmOperateRefuseRequest struct {
	Id int `json:"id"`
}

// @Summary	上下分拒绝
// @Param	request body services.GmOperateRefuseRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /gmOperateRefuse [post]
func (c *CustomerController) GmOperateRefuse() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	customerService := new(services.CustomerService)
	err = customerService.GmOperateRefuse(getInt(params, "id"), int(c.UserId), c.UserName)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
