package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type ActiveSignController struct {
	BaseController
}

type ActiveSignRuleResponse struct {
	Code int             `json:"code" example:"200"`    // 状态码
	Msg  string          `json:"msg" example:"success"` // 提示信息
	Data models.SignRule `json:"data"`
}

// @Summary	签到活动规则
// @Param activity_id query string false "活动ID"
// @Success 200 {object} controllers.ActiveSignRuleResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /activeSignRule [get]
func (c *ActiveSignController) ActiveSignRule() {
	ActivityId := c.GetQueryInt("activity_id", 0)
	request := services.ActiveSignRule{
		ActivityId: ActivityId,
	}
	activityService := new(services.ActiveSignService)
	data, err := activityService.ActiveSignRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加签到活动规则
// @Param	request body services.AddActiveSignRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveSignRule [post]
func (c *ActiveSignController) AddActiveSignRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activityService := new(services.ActiveSignService)
	request := services.AddActiveSignRuleRequest{
		ActivityId:       getInt(params, "activity_id"),
		Days:             getInt(params, "days"),
		IsLoop:           getInt(params, "is_loop"),
		IsInterruptReset: getInt(params, "is_interrupt_reset"),
	}
	err = activityService.AddActiveSignRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑签到活动规则
// @Param	request body services.EditActiveSignRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveSignRule [post]
func (c *ActiveSignController) EditActiveSignRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activityService := new(services.ActiveSignService)
	request := services.EditActiveSignRuleRequest{
		Id:               getInt(params, "id"),
		ActivityId:       getInt(params, "activity_id"),
		Days:             getInt(params, "days"),
		IsLoop:           getInt(params, "is_loop"),
		IsInterruptReset: getInt(params, "is_interrupt_reset"),
	}
	err = activityService.EditActiveSignRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteActiveSignRuleRequest struct {
	Id string `json:"id" example:"123456"` // 签到活动规则id
}

type AllActiveSignRewardResponse struct {
	Code int               `json:"code" example:"200"`    // 状态码
	Msg  string            `json:"msg" example:"success"` // 提示信息
	Data []models.SignRule `json:"data"`                  // 列表数据
}

// @Summary	签到活动奖励列表
// @Param activity_id query string false "活动ID"
// @Success 200 {object} controllers.AllActiveSignRewardResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveSignReward [get]
func (c *ActiveSignController) AllActiveSignReward() {
	ActivityId := c.GetQueryInt("activity_id", 0)
	request := services.AllActiveSignRewardRequest{
		ActivityId: ActivityId,
	}
	activityService := new(services.ActiveSignService)
	data, err := activityService.AllActiveSignReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type AllActiveSignRewardTypeResponse struct {
	Code int                    `json:"code" example:"200"`    // 状态码
	Msg  string                 `json:"msg" example:"success"` // 提示信息
	Data []ActiveSignRewardType `json:"data"`
}
type ActiveSignRewardType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	所有签到奖励类型
// @Success 200 {object} controllers.AllActiveSignRewardTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveSignRewardType [get]
func (c *ActiveSignController) AllActiveSignRewardType() {
	ActiveSignService := new(services.ActiveSignService)
	lang := c.Lang
	data := ActiveSignService.AllActiveSignRewardType(lang)
	c.JSONSuccess(data)
}

// @Summary	添加签到活动奖励
// @Param	request body services.AddActiveSignRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveSignReward [post]
func (c *ActiveSignController) AddActiveSignReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activityService := new(services.ActiveSignService)
	request := services.AddActiveSignRewardRequest{
		ActivityId:      getInt(params, "activity_id"),
		Day:             getInt(params, "day"),
		Icon:            getString(params, "icon"),
		RewardAmount:    getFloat(params, "reward_amount"),
		RewardTypeId:    getInt(params, "reward_type_id"),
		DayRechargeLine: getFloat(params, "day_recharge_line"),
		DayRunningLine:  getFloat(params, "day_running_line"),
	}
	err = activityService.AddActiveSignReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑签到活动奖励
// @Param	request body services.EditActiveSignRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveSignReward [post]
func (c *ActiveSignController) EditActiveSignReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activityService := new(services.ActiveSignService)
	request := services.EditActiveSignRewardRequest{
		Id:              getInt(params, "id"),
		ActivityId:      getInt(params, "activity_id"),
		Day:             getInt(params, "day"),
		Icon:            getString(params, "icon"),
		RewardAmount:    getFloat(params, "reward_amount"),
		RewardTypeId:    getInt(params, "reward_type_id"),
		DayRechargeLine: getFloat(params, "day_recharge_line"),
		DayRunningLine:  getFloat(params, "day_running_line"),
	}
	err = activityService.EditActiveSignReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteActiveSignRewardRequest struct {
	Id string `json:"id" example:"123456"` // 签到活动规则id
}

// @Summary	删除签到活动奖励
// @Param	request body controllers.DeleteActiveSignRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteActiveSignReward [post]
func (c *ActiveSignController) DeleteActiveSignReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	activityService := new(services.ActiveSignService)
	err = activityService.DeleteActiveSignReward(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
