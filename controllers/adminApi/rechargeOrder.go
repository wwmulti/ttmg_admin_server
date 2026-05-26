package controllers

import (
	"api/models"
	"api/services"
)

type RechargeOrderController struct {
	BaseController
}

type RechargeOrderListDTO struct {
	List        []models.RechargeOrder `json:"list"`         // 列表数据
	Total       int64                  `json:"total"`        // 总条数
	CurrentPage int                    `json:"current_page"` // 当前页码
}

type RechargeOrderListResponse struct {
	Code int                  `json:"code" example:"200"`    // 状态码
	Msg  string               `json:"msg" example:"success"` // 提示信息
	Data RechargeOrderListDTO `json:"data"`
}

// @Summary 充值订单列表
// @Param name query string false "名称"
// @Param type query int false "打码类型"
// @Param pro_rate query float64 false "打码倍率"
// @Param dev_rate query float64 false "测试倍率"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.RechargeOrderListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /rechargeOrderList [get]
func (c *RechargeOrderController) RechargeOrderList() {
	request := services.RechargeOrderListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	needReload := c.GetQueryInt("need_reload", 0)
	rechargeOrderService := new(services.RechargeOrderService)
	lang := c.Lang
	data, err := rechargeOrderService.RechargeOrderList(request, needReload, lang, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}
