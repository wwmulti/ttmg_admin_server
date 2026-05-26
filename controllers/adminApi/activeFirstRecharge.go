package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type ActiveFirstRechargeRuleController struct {
	BaseController
}

type ActiveFirstRechargeRuleListResponse struct {
	Code int                      `json:"code" example:"200"`    // 状态码
	Msg  string                   `json:"msg" example:"success"` // 提示信息
	Data models.FirstRechargeRule `json:"data"`
}

// @Summary	首充活动规则
// @Param activity_id query string false "活动ID"
// @Success 200 {object} controllers.ActiveFirstRechargeRuleListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /activeFirstRechargeRule [get]
func (c *ActiveFirstRechargeRuleController) ActiveFirstRechargeRule() {
	ActivityId := c.GetQueryInt("activity_id", 0)
	request := services.ActiveFirstRechargeRuleListRequest{
		ActivityId: ActivityId,
	}
	activeFirstRechargeService := new(services.ActiveFirstRechargeService)
	data, err := activeFirstRechargeService.ActiveFirstRechargeRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type AllTypeResponse struct {
	Code int                    `json:"code" example:"200"`    // 状态码
	Msg  string                 `json:"msg" example:"success"` // 提示信息
	Data []ActiveSignRewardType `json:"data"`
}
type Type struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	首充统计方式
// @Success 200 {object} controllers.AllTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveFirstRechargeRuleBillingType [get]
func (c *ActiveFirstRechargeRuleController) AllActiveFirstRechargeRuleBillingType() {
	lang := c.Lang
	activeFirstRechargeService := new(services.ActiveFirstRechargeService)
	data := activeFirstRechargeService.AllActiveFirstRechargeRuleBillingType(lang)
	c.JSONSuccess(data)
}

// @Summary	首充领取方式
// @Success 200 {object} controllers.AllTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveFirstRechargeRuleReceiveType [get]
func (c *ActiveFirstRechargeRuleController) AllActiveFirstRechargeRuleReceiveType() {
	lang := c.Lang
	activeFirstRechargeService := new(services.ActiveFirstRechargeService)
	data := activeFirstRechargeService.AllActiveFirstRechargeRuleReceiveType(lang)
	c.JSONSuccess(data)
}

// @Summary	首充任务类型
// @Success 200 {object} controllers.AllTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveFirstRechargeRuleTaskType [get]
func (c *ActiveFirstRechargeRuleController) AllActiveFirstRechargeRuleTaskType() {
	lang := c.Lang
	activeFirstRechargeService := new(services.ActiveFirstRechargeService)
	data := activeFirstRechargeService.AllActiveFirstRechargeRuleTaskType(lang)
	c.JSONSuccess(data)
}

// @Summary	添加首充活动规则
// @Param	request body services.AddActiveFirstRechargeRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveFirstRechargeRule [post]
func (c *ActiveFirstRechargeRuleController) AddActiveFirstRechargeRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activeFirstRechargeService := new(services.ActiveFirstRechargeService)
	request := services.AddActiveFirstRechargeRuleRequest{
		ActivityId:      getInt(params, "activity_id"),
		BetAmount:       getFloat(params, "bet_amount"),
		BetNumber:       getInt(params, "bet_number"),
		BillingType:     getInt(params, "billing_type"),
		RepeatActive:    getInt(params, "repeat_active"),
		TaskType:        getInt(params, "task_type"),
		UpdateFrequency: getInt(params, "update_frequency"),
		ReceiveType:     getInt(params, "receive_type"),
	}
	err = activeFirstRechargeService.AddActiveFirstRechargeRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑首充活动规则
// @Param	request body services.EditActiveFirstRechargeRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveFirstRechargeRule [post]
func (c *ActiveFirstRechargeRuleController) EditActiveFirstRechargeRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activeFirstRechargeService := new(services.ActiveFirstRechargeService)
	request := services.EditActiveFirstRechargeRuleRequest{
		Id:              getInt(params, "id"),
		ActivityId:      getInt(params, "activity_id"),
		BetAmount:       getFloat(params, "bet_amount"),
		BetNumber:       getInt(params, "bet_number"),
		BillingType:     getInt(params, "billing_type"),
		RepeatActive:    getInt(params, "repeat_active"),
		TaskType:        getInt(params, "task_type"),
		UpdateFrequency: getInt(params, "update_frequency"),
		ReceiveType:     getInt(params, "receive_type"),
	}
	err = activeFirstRechargeService.EditActiveFirstRechargeRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteActiveFirstRechargeRuleRequest struct {
	Id string `json:"id" example:"123456"` // 首充活动规则id
}

// @Summary	删除首充活动规则
// @Param	request body controllers.DeleteActiveFirstRechargeRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteActiveFirstRechargeRule [post]
func (c *ActiveFirstRechargeRuleController) DeleteActiveFirstRechargeRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	activeFirstRechargeService := new(services.ActiveFirstRechargeService)
	err = activeFirstRechargeService.DeleteActiveFirstRechargeRule(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type AllActiveFirstRechargeRuleRewardResponse struct {
	Code int                              `json:"code" example:"200"`    // 状态码
	Msg  string                           `json:"msg" example:"success"` // 提示信息
	Data []models.FirstRechargeRuleReward `json:"data"`                  // 列表数据
}

// @Summary	首充活动奖励列表
// @Param activity_id query string false "活动ID"
// @Success 200 {object} controllers.AllActiveFirstRechargeRuleRewardResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveFirstRechargeReward [get]
func (c *ActiveFirstRechargeRuleController) AllActiveFirstRechargeReward() {
	ActivityId := c.GetQueryInt("activity_id", 0)
	request := services.AllActiveFirstRechargeRuleRewardRequest{
		ActivityId: ActivityId,
	}
	activeFirstRechargeService := new(services.ActiveFirstRechargeService)
	data, err := activeFirstRechargeService.AllActiveFirstRechargeReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加首充活动奖励
// @Param	request body services.AddActiveFirstRechargeRuleRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveFirstRechargeReward [post]
func (c *ActiveFirstRechargeRuleController) AddActiveFirstRechargeReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activeFirstRechargeService := new(services.ActiveFirstRechargeService)
	request := services.AddActiveFirstRechargeRuleRewardRequest{
		ActivityId:          getInt(params, "activity_id"),
		SerialNumber:        getInt(params, "serial_number"),
		TotalRechargeAmount: getFloat(params, "total_recharge_amount"),
		RewardAmount:        getFloat(params, "reward_amount"),
	}
	err = activeFirstRechargeService.AddActiveFirstRechargeReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑首充活动奖励
// @Param	request body services.EditActiveFirstRechargeRuleRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveFirstRechargeReward [post]
func (c *ActiveFirstRechargeRuleController) EditActiveFirstRechargeReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activeFirstRechargeService := new(services.ActiveFirstRechargeService)
	request := services.EditActiveFirstRechargeRuleRewardRequest{
		Id:                  getInt(params, "id"),
		ActivityId:          getInt(params, "activity_id"),
		SerialNumber:        getInt(params, "serial_number"),
		TotalRechargeAmount: getFloat(params, "total_recharge_amount"),
		RewardAmount:        getFloat(params, "reward_amount"),
	}
	err = activeFirstRechargeService.EditActiveFirstRechargeReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteActiveFirstRechargeRuleRewardRequest struct {
	Id string `json:"id" example:"123456"` // 首充活动规则id
}

// @Summary	删除首充活动奖励
// @Param	request body controllers.DeleteActiveFirstRechargeRuleRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteActiveFirstRechargeRuleReward [post]
func (c *ActiveFirstRechargeRuleController) DeleteActiveFirstRechargeReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	activeFirstRechargeService := new(services.ActiveFirstRechargeService)
	err = activeFirstRechargeService.DeleteActiveFirstRechargeReward(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
