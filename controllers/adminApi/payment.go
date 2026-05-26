package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type PaymentController struct {
	BaseController
}

type PaymentTypeListDTO struct {
	List        []models.PaymentType `json:"list"`         // 列表数据
	Total       int64                `json:"total"`        // 总条数
	CurrentPage int                  `json:"current_page"` // 当前页
}
type PaymentTypeListResponse struct {
	Code int                `json:"code" example:"200"`    // 状态码
	Msg  string             `json:"msg" example:"success"` // 提示信息
	Data PaymentTypeListDTO `json:"data"`
}

// @Summary	支付类型列表
// @Param name query int false "支付类型名称"
// @Param status query int false "支付类型状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.PaymentTypeListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /paymentTypeList [get]
func (c *PaymentController) PaymentTypeList() {
	request := services.PaymentTypeListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	PaymentService := new(services.PaymentService)
	data, err := PaymentService.PaymentTypeList(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加支付类型
// @Param	request body services.AddPaymentTypeRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addPaymentType [post]
func (c *PaymentController) AddPaymentType() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	PaymentService := new(services.PaymentService)
	request := services.AddPaymentTypeRequest{
		Name:   getString(params, "name"),
		Sort:   getInt(params, "sort"),
		Status: getInt(params, "status"),
	}
	err = PaymentService.AddPaymentType(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑支付类型
// @Param	request body services.EditPaymentTypeRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editPaymentType [post]
func (c *PaymentController) EditPaymentType() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	PaymentService := new(services.PaymentService)
	request := services.EditPaymentTypeRequest{
		Id:     getInt(params, "id"),
		Name:   getString(params, "name"),
		Sort:   getInt(params, "sort"),
		Status: getInt(params, "status"),
	}
	err = PaymentService.EditPaymentType(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改支付类型状态
// @Param	request body services.ChangePaymentTypeStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changePaymentTypeStatus [post]
func (c *PaymentController) ChangePaymentTypeStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	PaymentService := new(services.PaymentService)
	request := services.ChangePaymentTypeStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = PaymentService.ChangePaymentTypeStatus(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type PaymentListDTO struct {
	List        []models.Payment `json:"list"`         // 列表数据
	Total       int64            `json:"total"`        // 总条数
	CurrentPage int              `json:"current_page"` // 当前页
}
type PaymentListResponse struct {
	Code int            `json:"code" example:"200"`    // 状态码
	Msg  string         `json:"msg" example:"success"` // 提示信息
	Data PaymentListDTO `json:"data"`
}

// @Summary	支付渠道列表
// @Param name query int false "支付渠道名称"
// @Param status query int false "支付渠道状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.PaymentListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /paymentList [get]
func (c *PaymentController) PaymentList() {
	request := services.PaymentListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	PaymentService := new(services.PaymentService)
	data, err := PaymentService.PaymentList(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加支付渠道
// @Param	request body services.AddPaymentRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addPayment [post]
func (c *PaymentController) AddPayment() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	PaymentService := new(services.PaymentService)
	request := services.AddPaymentRequest{
		PayCode:        getString(params, "pay_code"),
		MerchantConfig: getJSONString(params, "merchant_config"),
		Logo:           getString(params, "logo"),
		Remark:         getString(params, "remark"),
		Status:         getInt(params, "status"),
	}
	err = PaymentService.AddPayment(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑支付渠道
// @Param	request body services.EditPaymentRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editPayment [post]
func (c *PaymentController) EditPayment() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	PaymentService := new(services.PaymentService)
	request := services.EditPaymentRequest{
		Id:             getInt(params, "id"),
		PayCode:        getString(params, "pay_code"),
		MerchantConfig: getJSONString(params, "merchant_config"),
		Logo:           getString(params, "logo"),
		Remark:         getString(params, "remark"),
		Status:         getInt(params, "status"),
	}
	err = PaymentService.EditPayment(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改支付类型状态
// @Param	request body services.ChangePaymentStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changePaymentStatus [post]
func (c *PaymentController) ChangePaymentStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	PaymentService := new(services.PaymentService)
	request := services.ChangePaymentStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = PaymentService.ChangePaymentStatus(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	支付子渠道列表
// @Param name query string false "支付子渠道名称"
// @Param status query int false "支付子渠道状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.PaymentChannelListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /paymentChannelList [get]
func (c *PaymentController) PaymentChannelList() {
	request := services.PaymentChannelListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	PaymentService := new(services.PaymentService)
	needReload := c.GetQueryInt("need_reload", 0)
	data, err := PaymentService.PaymentChannelList(request, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加支付子渠道
// @Param	request body services.AddPaymentChannelRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addPaymentChannel [post]
func (c *PaymentController) AddPaymentChannel() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	PaymentService := new(services.PaymentService)
	request := services.AddPaymentChannelRequest{
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
	err = PaymentService.AddPaymentChannel(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑支付子渠道
// @Param	request body services.EditPaymentChannelRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editPaymentChannel [post]
func (c *PaymentController) EditPaymentChannel() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	PaymentService := new(services.PaymentService)
	request := services.EditPaymentChannelRequest{
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
	err = PaymentService.EditPaymentChannel(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改支付类型状态
// @Param	request body services.ChangePaymentChannelStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changePaymentChannelStatus [post]
func (c *PaymentController) ChangePaymentChannelStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	PaymentService := new(services.PaymentService)
	request := services.ChangePaymentChannelStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = PaymentService.ChangePaymentChannelStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type PaymentChannelQuickListResponse struct {
	Code int                          `json:"code" example:"200"`
	Msg  string                       `json:"msg" example:"success"`
	Data []models.PaymentChannelQuick `json:"data"`
}

// @Summary	支付子渠道快捷列表
// @Param payment_channel_id query int false "支付子渠道ID"
// @Success 200 {object} controllers.PaymentChannelQuickListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allPaymentChannelQuick [get]
func (c *PaymentController) AllPaymentChannelQuick() {
	paymentChannelId := c.GetQueryInt("payment_channel_id", 0)
	PaymentService := new(services.PaymentService)
	data, err := PaymentService.AllPaymentChannelQuick(paymentChannelId, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加支付子渠道快捷
// @Param	request body services.AddPaymentChannelQuickRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addPaymentChannelQuick [post]
func (c *PaymentController) AddPaymentChannelQuick() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	PaymentService := new(services.PaymentService)
	request := services.AddPaymentChannelQuickRequest{
		PaymentChannelId: getInt(params, "payment_channel_id"),
		Amount:           getFloat(params, "amount"),
		IsRecommend:      getInt(params, "is_recommend"),
		Sort:             getInt(params, "sort"),
	}
	err = PaymentService.AddPaymentChannelQuick(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑支付子渠道快捷
// @Param	request body services.EditPaymentChannelQuickRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editPaymentChannelQuick [post]
func (c *PaymentController) EditPaymentChannelQuick() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	PaymentService := new(services.PaymentService)
	request := services.EditPaymentChannelQuickRequest{
		Id:               getInt(params, "id"),
		PaymentChannelId: getInt(params, "payment_channel_id"),
		Amount:           getFloat(params, "amount"),
		IsRecommend:      getInt(params, "is_recommend"),
		Sort:             getInt(params, "sort"),
	}
	err = PaymentService.EditPaymentChannelQuick(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeletePaymentChannelQuickRequest struct {
	Id string `json:"id" example:"123456"` // 轮播图id
}

// @Summary	删除快捷支付
// @Param	request body controllers.DeletePaymentChannelQuickRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deletePaymentChannelQuick [post]
func (c *PaymentController) DeletePaymentChannelQuick() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	PaymentService := new(services.PaymentService)
	err = PaymentService.DeletePaymentChannelQuick(id, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改快捷支付推荐
// @Param	request body services.ChangePaymentChannelQuickRecommendRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changePaymentChannelQuickRecommend [post]
func (c *PaymentController) ChangePaymentChannelQuickRecommend() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	PaymentService := new(services.PaymentService)
	request := services.ChangePaymentChannelQuickRecommendRequest{
		Id:          id,
		IsRecommend: getInt(params, "is_recommend"),
	}
	err = PaymentService.ChangePaymentChannelQuickRecommend(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
