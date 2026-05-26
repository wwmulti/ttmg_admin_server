package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type ActiveReliefController struct {
	BaseController
}

type ActiveReliefRuleResponse struct {
	Code int                             `json:"code"`
	Msg  string                          `json:"msg"`
	Data models.ActiveReliefRuleModelDTO `json:"data"`
}

// @Summary 获取救济活动规则列表
// @Param active_id query int true "活动ID"
// @Success 200 {object} controllers.ActiveReliefRuleResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /activeReliefRule [get]
func (c *ActiveReliefController) ActiveReliefRule() {
	request := services.ActiveReliefRuleRequest{
		ActiveId: c.GetQueryInt("active_id", 0),
	}
	data, err := (&services.ActiveReliefService{}).ActiveReliefRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(data)
}

type AllActiveReliefRuleCycleTypeResponse struct {
	Code int               `json:"code" example:"200"`    // 状态码
	Msg  string            `json:"msg" example:"success"` // 提示信息
	Data []ReliefCycleType `json:"data"`
}

type ReliefCycleType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	所有救济统计周期
// @Success 200 {object} controllers.AllActiveReliefRuleCycleTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveReliefRuleCycleType [get]
func (c *ActiveReliefController) AllActiveReliefRuleCycleType() {
	ActivityReliefService := new(services.ActiveReliefService)
	lang := c.Lang
	data := ActivityReliefService.AllActiveReliefRuleCycleType(lang)
	c.JSONSuccess(data)
}

// @Summary 添加救济活动规则
// @Param   request body services.AddActiveReliefRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveReliefRule [post]
func (c *ActiveReliefController) AddActiveReliefRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	request := services.AddActiveReliefRuleRequest{
		ActiveId: getInt(params, "active_id"),
		Cycle:    getInt(params, "cycle"),
		OpenDay:  getInt(params, "open_day"),
		OpenTime: getString(params, "open_time"),
		IsRepeat: getInt(params, "is_repeat"),
	}
	err = (&services.ActiveReliefService{}).AddActiveReliefRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type EditActiveReliefRuleRequest struct {
	services.AddActiveReliefRuleRequest
	Id int `json:"id"`
}

// @Summary 编辑救济活动规则
// @Param   request body controllers.EditActiveReliefRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveReliefRule [post]
func (c *ActiveReliefController) EditActiveReliefRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	id := getInt(params, "id")
	request := services.AddActiveReliefRuleRequest{
		ActiveId: getInt(params, "active_id"),
		Cycle:    getInt(params, "cycle"),
		OpenDay:  getInt(params, "open_day"),
		OpenTime: getString(params, "open_time"),
		IsRepeat: getInt(params, "is_repeat"),
	}
	err = (&services.ActiveReliefService{}).EditActiveReliefRule(id, request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type AllActiveReliefRewardResponse struct {
	JSONResponseTpl
	Data models.ActiveReliefRewardsDTO `json:"data"`
}

// @Summary 获取救济活动奖励
// @Param active_id query int false "活动ID"
// @Success 200 {object} controllers.AllActiveReliefRewardResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveReliefReward [get]
func (c *ActiveReliefController) AllActiveReliefReward() {
	request := services.AllActiveReliefRewardRequest{
		ActiveId: c.GetQueryInt("active_id", 0),
	}
	data, err := (&services.ActiveReliefService{}).AllActiveReliefReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(data)
}

// @Summary 添加救济活动奖励
// @Param   request body services.AddActiveReliefRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveReliefReward [post]
func (c *ActiveReliefController) AddActiveReliefReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	request := services.AddActiveReliefRewardRequest{
		ActiveId:      getInt(params, "active_id"),
		Amount:        getFloat(params, "amount"),
		RebatePercent: getFloat(params, "rebate_percent"),
		Status:        getInt(params, "status"),
	}
	err = (&services.ActiveReliefService{}).AddActiveReliefReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type EditActiveReliefRewardRequest struct {
	services.AddActiveReliefRewardRequest
	Id int `json:"id"`
}

// @Summary 编辑救济活动奖励
// @Param   request body controllers.EditActiveReliefRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveReliefReward [post]
func (c *ActiveReliefController) EditActiveReliefReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	id := getInt(params, "id")
	request := services.AddActiveReliefRewardRequest{
		ActiveId:      getInt(params, "active_id"),
		Amount:        getFloat(params, "amount"),
		RebatePercent: getFloat(params, "rebate_percent"),
		Status:        getInt(params, "status"),
	}
	err = (&services.ActiveReliefService{}).EditActiveReliefReward(id, request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteActiveReliefRewardRequest struct {
	Id int `json:"id"`
}

// @Summary 删除救济活动奖励
// @Param   request body controllers.DeleteActiveReliefRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteActiveReliefReward [post]
func (c *ActiveReliefController) DeleteActiveReliefReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	id := getInt(params, "id")
	err = (&services.ActiveReliefService{}).DeleteActiveReliefReward(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改救济活动奖励状态
// @Param	request body services.ChangeActiveReliefRewardStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeActiveReliefRewardStatus [post]
func (c *ActiveReliefController) ChangeActiveReliefRewardStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	ActiveReliefService := new(services.ActiveReliefService)
	request := services.ChangeActiveReliefRewardStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = ActiveReliefService.ChangeActiveReliefRewardStatus(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
