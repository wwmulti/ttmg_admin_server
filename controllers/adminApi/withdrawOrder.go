package controllers

import (
	"api/models"
	"api/services"
)

type WithdrawOrderController struct {
	BaseController
}

type WithdrawOrderListDTO struct {
	List        []models.WithdrawOrder `json:"list"`         // 列表数据
	Total       int64                  `json:"total"`        // 总条数
	CurrentPage int                    `json:"current_page"` // 当前页码
}

type WithdrawOrderListResponse struct {
	Code int                  `json:"code" example:"200"`    // 状态码
	Msg  string               `json:"msg" example:"success"` // 提示信息
	Data WithdrawOrderListDTO `json:"data"`
}

// @Summary 提现订单列表
// @Param name query string false "名称"
// @Param type query int false "打码类型"
// @Param pro_rate query float64 false "打码倍率"
// @Param dev_rate query float64 false "测试倍率"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.WithdrawOrderListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /withdrawOrderList [get]
func (c *WithdrawOrderController) WithdrawOrderList() {
	request := services.WithdrawOrderListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	needReload := c.GetQueryInt("need_reload", 0)
	withdrawOrderService := new(services.WithdrawOrderService)
	lang := c.Lang
	data, err := withdrawOrderService.WithdrawOrderList(request, needReload, lang, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}
