package controllers

import (
	"api/services"
	"encoding/json"
)

type WithdrawController struct {
	BaseController
}

// @Summary	自动出款配置
// @Param package_id query int false "开放平台"
// @Success 200 {object} controllers.ActiveListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /autoOutMoneyConfig [get]
func (c *WithdrawController) AutoOutMoneyConfig() {
	userIdInterface := c.Data["UserId"]
	userId := userIdInterface.(int64)

	packageId := c.GetQueryInt("package_id", 0)
	data, err := (&services.WithdrawService{}).AutoOutMoneyConfig(userId, packageId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	提现子渠道列表
// @Param name query string false "提现子渠道名称"
// @Param status query int false "提现子渠道状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /withdrawChannelList [get]
func (c *WithdrawController) WithdrawChannelList() {
	request := services.WithdrawChannelListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	withDrawService := new(services.WithdrawService)
	needReload := c.GetQueryInt("need_reload", 0)
	data, err := withDrawService.WithdrawChannelList(request, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加提现子渠道
// @Param	request body services.AddWithdrawChannelRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addWithdrawChannel [post]
func (c *WithdrawController) AddWithdrawChannel() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	withDrawService := new(services.WithdrawService)
	request := services.AddWithdrawChannelRequest{
		PackageId:     getInt(params, "package_id"),
		PaymentId:     getInt(params, "payment_id"),
		PaymentTypeId: getInt(params, "payment_type_id"),
		Tag:           getString(params, "tag"),
		VipLevel:      getInt(params, "vip_level"),
		LimitSmall:    getFloat(params, "limit_small"),
		LimitBig:      getFloat(params, "limit_big"),
		PrizePercent:  getInt(params, "prize_percent"),
		Name:          getString(params, "name"),
		ChannelConfig: getJSONString(params, "channel_config"),
		Status:        getInt(params, "status"),
	}
	err = withDrawService.AddWithdrawChannel(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑提现子渠道
// @Param	request body services.EditWithdrawChannelRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editWithdrawChannel [post]
func (c *WithdrawController) EditWithdrawChannel() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	withDrawService := new(services.WithdrawService)
	request := services.EditWithdrawChannelRequest{
		Id:            getInt(params, "id"),
		PackageId:     getInt(params, "package_id"),
		PaymentId:     getInt(params, "payment_id"),
		PaymentTypeId: getInt(params, "payment_type_id"),
		Tag:           getString(params, "tag"),
		VipLevel:      getInt(params, "vip_level"),
		LimitSmall:    getFloat(params, "limit_small"),
		LimitBig:      getFloat(params, "limit_big"),
		PrizePercent:  getInt(params, "prize_percent"),
		Name:          getString(params, "name"),
		ChannelConfig: getJSONString(params, "channel_config"),
		Status:        getInt(params, "status"),
	}
	err = withDrawService.EditWithdrawChannel(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改提现子渠道状态
// @Param	request body services.ChangeWithdrawChannelStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeWithdrawChannelStatus [post]
func (c *WithdrawController) ChangeWithdrawChannelStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	withDrawService := new(services.WithdrawService)
	request := services.ChangeWithdrawChannelStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = withDrawService.ChangeWithdrawChannelStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
