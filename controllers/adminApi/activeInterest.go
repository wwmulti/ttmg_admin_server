package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type ActiveInterestController struct {
	BaseController
}

type ActiveInterestRuleResponse struct {
	Code int                 `json:"code" example:"200"`    // 状态码
	Msg  string              `json:"msg" example:"success"` // 提示信息
	Data models.InterestRule `json:"data"`
}

// @Summary	利息宝活动规则
// @Param activity_id query string false "活动ID"
// @Success 200 {object} controllers.ActiveInterestRuleResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /activeInterestRule [get]
func (c *ActiveInterestController) ActiveInterestRule() {
	ActivityId := c.GetQueryInt("activity_id", 0)
	request := services.ActiveInterestRule{
		ActivityId: ActivityId,
	}
	activityService := new(services.ActiveInterestService)
	data, err := activityService.ActiveInterestRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type AllActiveInterestTypeResponse struct {
	Code int                  `json:"code" example:"200"`    // 状态码
	Msg  string               `json:"msg" example:"success"` // 提示信息
	Data []ActiveInterestType `json:"data"`
}
type ActiveInterestType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	所有利息宝领取时间类型
// @Success 200 {object} controllers.AllActiveInterestTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveInterestReceiveType [get]
func (c *ActiveInterestController) AllActiveInterestReceiveType() {
	ActiveInterestService := new(services.ActiveInterestService)
	lang := c.Lang
	data := ActiveInterestService.AllActiveInterestReceiveType(lang)
	c.JSONSuccess(data)
}

// @Summary	所有利息宝利息最大限制类型
// @Success 200 {object} controllers.AllActiveInterestTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveInterestInterestLimitType [get]
func (c *ActiveInterestController) AllActiveInterestInterestLimitType() {
	ActiveInterestService := new(services.ActiveInterestService)
	lang := c.Lang
	data := ActiveInterestService.AllActiveInterestInterestLimitType(lang)
	c.JSONSuccess(data)
}

// @Summary	添加利息宝活动规则
// @Param	request body services.AddActiveInterestRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveInterestRule [post]
func (c *ActiveInterestController) AddActiveInterestRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activityService := new(services.ActiveInterestService)
	request := services.AddActiveInterestRuleRequest{
		ActivityId:          getInt(params, "activity_id"),
		InterestRate:        getInt(params, "interest_rate"),
		DepositAmount:       getInt64(params, "deposit_amount"),
		Interval:            getInt(params, "interval"),
		ReceiveType:         getInt(params, "receive_type"),
		InterestLimitType:   getInt(params, "interest_limit_type"),
		InterestLimitAmount: getInt64(params, "interest_limit_amount"),
	}
	err = activityService.AddActiveInterestRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑利息宝活动规则
// @Param	request body services.EditActiveInterestRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveInterestRule [post]
func (c *ActiveInterestController) EditActiveInterestRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activityService := new(services.ActiveInterestService)
	request := services.EditActiveInterestRuleRequest{
		Id:                  getInt(params, "id"),
		ActivityId:          getInt(params, "activity_id"),
		InterestRate:        getInt(params, "interest_rate"),
		DepositAmount:       getInt64(params, "deposit_amount"),
		Interval:            getInt(params, "interval"),
		ReceiveType:         getInt(params, "receive_type"),
		InterestLimitType:   getInt(params, "interest_limit_type"),
		InterestLimitAmount: getInt64(params, "interest_limit_amount"),
	}
	err = activityService.EditActiveInterestRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
