package controllers

import (
	"api/services"
	"encoding/json"
)

type ActiveLuckyWheelController struct {
	BaseController
}

// @Summary	幸运转盘活动奖励列表
// @Param activity_id query string false "活动ID"
// @Success 200 {object} services.AllActiveLuckyWheelRewardRequest
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveLuckyWheelReward [get]
func (c *ActiveLuckyWheelController) AllActiveLuckyWheelReward() {
	ActivityId := c.GetQueryInt("activity_id", 0)
	WheelType := c.GetQueryInt("wheel_type", 0)
	request := services.AllActiveLuckyWheelRewardRequest{
		ActivityId: ActivityId,
		WheelType:  WheelType,
	}
	activityService := new(services.ActiveLuckyWheelService)
	data, err := activityService.AllActiveLuckyWheelReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type AllActiveLuckyWheelRewardTypeResponse struct {
	Code int                          `json:"code" example:"200"`    // 状态码
	Msg  string                       `json:"msg" example:"success"` // 提示信息
	Data []ActiveLuckyWheelRewardType `json:"data"`
}
type ActiveLuckyWheelRewardType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	所有幸运转盘奖励类型
// @Success 200 {object} controllers.AllActiveLuckyWheelRewardTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveLuckyWheelRewardType [get]
func (c *ActiveLuckyWheelController) AllActiveLuckyWheelRewardType() {
	ActiveLuckyWheelService := new(services.ActiveLuckyWheelService)
	lang := c.Lang
	data := ActiveLuckyWheelService.AllActiveLuckyWheelRewardType(lang)
	c.JSONSuccess(data)
}

// @Summary	添加幸运转盘活动奖励
// @Param	request body services.AddActiveLuckyWheelRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveLuckyWheelReward [post]
func (c *ActiveLuckyWheelController) AddActiveLuckyWheelReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activityService := new(services.ActiveLuckyWheelService)
	request := services.AddActiveLuckyWheelRewardRequest{
		ActivityId: getInt(params, "activity_id"),
		WheelType:  getInt(params, "wheel_type"),
		Reward:     getFloat(params, "reward"),
		Weight:     getInt(params, "weight"),
	}
	err = activityService.AddActiveLuckyWheelReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑幸运转盘活动奖励
// @Param	request body services.EditActiveLuckyWheelRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveLuckyWheelReward [post]
func (c *ActiveLuckyWheelController) EditActiveLuckyWheelReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activityService := new(services.ActiveLuckyWheelService)
	request := services.EditActiveLuckyWheelRewardRequest{
		Id:         getInt(params, "id"),
		ActivityId: getInt(params, "activity_id"),
		WheelType:  getInt(params, "wheel_type"),
		Reward:     getFloat(params, "reward"),
		Weight:     getInt(params, "weight"),
	}
	err = activityService.EditActiveLuckyWheelReward(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteActiveLuckyWheelRewardRequest struct {
	Id string `json:"id" example:"123456"` // 幸运转盘活动规则ID
}

// @Summary	删除幸运转盘活动奖励
// @Param	request body controllers.DeleteActiveLuckyWheelRewardRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteActiveLuckyWheelReward [post]
func (c *ActiveLuckyWheelController) DeleteActiveLuckyWheelReward() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	activityService := new(services.ActiveLuckyWheelService)
	err = activityService.DeleteActiveLuckyWheelReward(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
