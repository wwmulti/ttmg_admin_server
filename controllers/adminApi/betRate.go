package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type BetRateController struct {
	BaseController
}

// DTO 定义
type BetRateListDTO struct {
	List        []models.BetRate `json:"list"`         // 列表数据
	Total       int64            `json:"total"`        // 总条数
	CurrentPage int              `json:"current_page"` // 当前页码
}

type BetRateListResponse struct {
	Code int            `json:"code" example:"200"`    // 状态码
	Msg  string         `json:"msg" example:"success"` // 提示信息
	Data BetRateListDTO `json:"data"`
}

// @Summary 打码倍数列表
// @Param name query string false "名称"
// @Param type query int false "打码类型"
// @Param pro_rate query float64 false "打码倍率"
// @Param dev_rate query float64 false "测试倍率"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.BetRateListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /betRateList [get]
func (c *BetRateController) BetRateList() {
	request := services.BetRateListRequest{}
	c.ParseForm(&request)
	request.Lang = c.Lang
	request.RawQuery = c.Ctx.Request.URL.Query()

	if request.PackageId > 0 {
		if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
			c.JSONError(500, "ShuJuYiChang")
		}
	} else {
		request.PackageIds = c.GetPackageIds(int(c.UserId))
	}

	betRateService := new(services.BetRateService)
	data, err := betRateService.BetRateList(request, int(c.UserId))
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary 添加打码倍数
// @Param   request body services.AddBetRateRequest true "json请求参数"
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addBetRate [post]
func (c *BetRateController) AddBetRate() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
		return
	}

	betRateService := new(services.BetRateService)
	request := services.AddBetRateRequest{
		Name:      getString(params, "name"),
		Type:      getInt(params, "type"),
		ProRate:   getFloat(params, "pro_rate"),
		DevRate:   getFloat(params, "dev_rate"),
		PackageId: getInt(params, "package_id"),
		Lang:      c.Lang,
	}

	if request.PackageId > 0 {
		if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
			c.JSONError(500, "ShuJuYiChang")
		}
	}

	err = betRateService.AddBetRate(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary 编辑打码倍数
// @Param   request body services.EditBetRateRequest true "json请求参数"
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editBetRate [post]
func (c *BetRateController) EditBetRate() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
		return
	}

	betRateService := new(services.BetRateService)
	request := services.EditBetRateRequest{
		Id:      getInt(params, "id"),
		Name:    getString(params, "name"),
		ProRate: getFloat(params, "pro_rate"),
		DevRate: getFloat(params, "dev_rate"),
	}

	err = betRateService.EditBetRate(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
