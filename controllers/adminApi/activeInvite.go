package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type ActiveInviteController struct {
	BaseController
}

type ActiveInviteRuleResponse struct {
	JSONResponseTpl
	Data models.ActiveShareRuleDTO `json:"data"`
}

// @Summary	获取邀请活动规则列表
// @Param active_id query int true "活动ID"
// @Success 200 {object} controllers.ActiveInviteRuleResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /activeInviteRule [get]
func (c *ActiveInviteController) ActiveInviteRule() {
	request := services.ActiveInviteRuleRequest{
		ActiveId: c.GetQueryInt("active_id", 0),
	}
	data, err := (&services.ActiveInviteService{}).ActiveInviteRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(data)
}

type AllActiveInviteTypeResponse struct {
	Code int          `json:"code" example:"200"`    // 状态码
	Msg  string       `json:"msg" example:"success"` // 提示信息
	Data []InviteType `json:"data"`
}

type InviteType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	所有扫描范围
// @Success 200 {object} controllers.AllActiveInviteTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveInviteRuleScopeType [get]
func (c *ActiveInviteController) ActiveInviteRuleScopeType() {
	ActiveInviteService := new(services.ActiveInviteService)
	lang := c.Lang
	data := ActiveInviteService.ActiveInviteRuleScopeType(lang)
	c.JSONSuccess(data)
}

// @Summary	所有条件类型范围
// @Success 200 {object} controllers.AllActiveInviteTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveInviteRuleConditionType [get]
func (c *ActiveInviteController) ActiveInviteRuleConditionType() {
	ActiveInviteService := new(services.ActiveInviteService)
	lang := c.Lang
	data := ActiveInviteService.ActiveInviteRuleConditionType(lang)
	c.JSONSuccess(data)
}

// @Summary	所有奖励领取方式
// @Success 200 {object} controllers.AllActiveInviteTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveInviteRuleRewardType [get]
func (c *ActiveInviteController) ActiveInviteRuleRewardType() {
	ActiveInviteService := new(services.ActiveInviteService)
	lang := c.Lang
	data := ActiveInviteService.ActiveInviteRuleRewardType(lang)
	c.JSONSuccess(data)
}

// @Summary	所有过期策略
// @Success 200 {object} controllers.AllActiveInviteTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveInviteRuleExpireType [get]
func (c *ActiveInviteController) ActiveInviteRuleExpireType() {
	ActiveInviteService := new(services.ActiveInviteService)
	lang := c.Lang
	data := ActiveInviteService.ActiveInviteRuleExpireType(lang)
	c.JSONSuccess(data)
}

// @Summary	添加邀请活动规则
// @Param	request body services.AddActiveInviteRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveInviteRule [post]
func (c *ActiveInviteController) AddActiveInviteRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	request := services.AddActiveInviteRuleRequest{
		ActiveId:       getInt(params, "active_id"),
		ActiveTypeId:   getInt(params, "active_type_id"),
		ShareUrl:       getString(params, "share_url"),
		ScanInterval:   getInt(params, "scan_interval"),
		Scope:          getInt(params, "scope"),
		MiniTotalPays:  getFloat(params, "mini_total_pays"),
		MiniTotalWater: getFloat(params, "mini_total_water"),
		Condition:      getInt(params, "condition"),
		RewardType:     getInt(params, "reward_type"),
		ExpireType:     getInt(params, "expire_type"),
	}
	err = (&services.ActiveInviteService{}).AddActiveInviteRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type EditActiveInviteRuleRequest struct {
	services.AddActiveInviteRuleRequest
	Id int `json:"id"`
}

// @Summary	编辑邀请活动规则
// @Param	request body controllers.EditActiveInviteRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveInviteRule [post]
func (c *ActiveInviteController) EditActiveInviteRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	id := getInt(params, "id")
	request := services.AddActiveInviteRuleRequest{
		ActiveId:       getInt(params, "active_id"),
		ActiveTypeId:   getInt(params, "active_type_id"),
		ShareUrl:       getString(params, "share_url"),
		ScanInterval:   getInt(params, "scan_interval"),
		Scope:          getInt(params, "scope"),
		MiniTotalPays:  getFloat(params, "mini_total_pays"),
		MiniTotalWater: getFloat(params, "mini_total_water"),
		Condition:      getInt(params, "condition"),
		RewardType:     getInt(params, "reward_type"),
		ExpireType:     getInt(params, "expire_type"),
	}
	err = (&services.ActiveInviteService{}).EditActiveInviteRule(id, request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type AllActiveInviteRewardResponse struct {
	JSONResponseTpl
	Data models.ActiveShareRewardsDTO `json:"data"`
}

// @Summary	获取邀请活动奖励列表
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.AllActiveInviteRewardResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveInviteReward [get]
func (c *ActiveInviteController) AllActiveInviteReward() {
	request := services.AllActiveInviteRewardRequest{
		ActiveId: c.GetQueryInt("active_id", 0),
	}
	data, err := (&services.ActiveInviteService{}).AllActiveInviteReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(data)
}

// @Summary	添加邀请活动奖励
// @Param	request body services.AddActiveInviteRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveInviteReward [post]
func (c *ActiveInviteController) AddActiveInviteReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	request := services.AddActiveInviteRewardRequest{
		ActiveId: getInt(params, "active_id"),
		Mens:     getInt(params, "mens"),
		Rewards:  getFloat(params, "rewards"),
		IconOn:   getString(params, "icon_on"),
		IconOff:  getString(params, "icon_off"),
		Status:   getInt(params, "status"),
	}
	err = (&services.ActiveInviteService{}).AddActiveInviteReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type EditActiveInviteRewardRequest struct {
	services.AddActiveInviteRewardRequest
	Id int `json:"id"`
}

// @Summary	编辑邀请活动奖励
// @Param	request body controllers.EditActiveInviteRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveInviteReward [post]
func (c *ActiveInviteController) EditActiveInviteReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	id := getInt(params, "id")
	request := services.AddActiveInviteRewardRequest{
		ActiveId: getInt(params, "active_id"),
		RuleId:   getInt(params, "rule_id"),
		Mens:     getInt(params, "mens"),
		Rewards:  getFloat(params, "rewards"),
		IconOn:   getString(params, "icon_on"),
		IconOff:  getString(params, "icon_off"),
		Status:   getInt(params, "status"),
	}
	err = (&services.ActiveInviteService{}).EditActiveInviteReward(id, request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteActiveInviteRewardRequest struct {
	Id int `json:"id"`
}

// @Summary	删除邀请活动奖励
// @Param	request body controllers.DeleteActiveInviteRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteActiveInviteReward [post]
func (c *ActiveInviteController) DeleteActiveInviteReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	id := getInt(params, "id")

	err = (&services.ActiveInviteService{}).DeleteActiveInviteReward(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改邀请奖励状态
// @Param	request body services.ChangeActiveInviteRewardsStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeActiveInviteRewardStatus [post]
func (c *ActiveInviteController) ChangeActiveInviteRewardStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	ActiveInviteService := new(services.ActiveInviteService)
	request := services.ChangeActiveInviteRewardsStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = ActiveInviteService.ChangeActiveInviteRewardsStatus(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
