package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type ActiveVipController struct {
	BaseController
}

type AllActiveVipRuleResponse struct {
	Code int                    `json:"code" example:"200"`    // 状态码
	Msg  string                 `json:"msg" example:"success"` // 提示信息
	Data []models.ActiveVipRule `json:"data"`
}

// @Summary 获取所有VIP活动规则
// @Param active_id query int true "活动ID"
// @Success 200 {object} controllers.AllActiveVipRuleResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveVipRule [get]
func (c *ActiveVipController) AllActiveVipRule() {
	ActiveId := c.GetQueryInt("active_id", 0)
	request := services.AllActiveVipRuleRequest{
		ActiveId: ActiveId,
	}
	data, err := (&services.ActiveVipService{}).AllActiveVipRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(data)
}

type AllActiveVipTypeResponse struct {
	Code int       `json:"code" example:"200"`    // 状态码
	Msg  string    `json:"msg" example:"success"` // 提示信息
	Data []VipType `json:"data"`
}

type VipType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	Vip规则所有条件类型
// @Success 200 {object} controllers.AllActiveVipTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveVipRuleConditionType [get]
func (c *ActiveVipController) AllActiveVipRuleConditionType() {
	ActivityVipService := new(services.ActiveVipService)
	lang := c.Lang
	data := ActivityVipService.AllActiveVipRuleConditionType(lang)
	c.JSONSuccess(data)
}

// @Summary 添加VIP活动规则
// @Param   request body services.AddActiveVipRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveVipRule [post]
func (c *ActiveVipController) AddActiveVipRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	request := services.AddActiveVipRuleRequest{
		ActiveId:            getInt(params, "active_id"),
		Lv:                  getInt(params, "lv"),
		TotalPays:           getFloat(params, "total_pays"),
		TotalBets:           getFloat(params, "total_bets"),
		ConAnd:              getInt(params, "con_and"),
		Rewards:             getFloat(params, "rewards"),
		WithdrawNumLimit:    getInt(params, "withdraw_num_limit"),
		WithdrawAmountLimit: getFloat(params, "withdraw_amount_limit"),
		WithdrawFreeNum:     getInt(params, "withdraw_free_num"),
		WithdrawFee:         getInt(params, "withdraw_fee"),
		Status:              getInt(params, "status"),
	}
	err = (&services.ActiveVipService{}).AddActiveVipRule(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type EditActiveVipRuleRequest struct {
	services.AddActiveVipRuleRequest
	Id int `json:"id"`
}

// @Summary 编辑VIP活动规则
// @Param   request body controllers.EditActiveVipRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveVipRule [post]
func (c *ActiveVipController) EditActiveVipRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	id := getInt(params, "id")
	request := services.AddActiveVipRuleRequest{
		ActiveId:            getInt(params, "active_id"),
		Lv:                  getInt(params, "lv"),
		TotalPays:           getFloat(params, "total_pays"),
		TotalBets:           getFloat(params, "total_bets"),
		ConAnd:              getInt(params, "con_and"),
		Rewards:             getFloat(params, "rewards"),
		WithdrawNumLimit:    getInt(params, "withdraw_num_limit"),
		WithdrawAmountLimit: getFloat(params, "withdraw_amount_limit"),
		WithdrawFreeNum:     getInt(params, "withdraw_free_num"),
		WithdrawFee:         getInt(params, "withdraw_fee"),
		Status:              getInt(params, "status"),
	}
	err = (&services.ActiveVipService{}).EditActiveVipRule(id, request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteActiveVipRuleRequest struct {
	Id int `json:"id"`
}

// @Summary 删除VIP活动规则
// @Param   request body controllers.DeleteActiveVipRuleRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteActiveVipRule [post]
func (c *ActiveVipController) DeleteActiveVipRule() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	id := getInt(params, "id")
	err = (&services.ActiveVipService{}).DeleteActiveVipRule(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改Vip规则状态
// @Param	request body services.ChangeActiveVipRuleStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeActiveVipRuleStatus [post]
func (c *ActiveVipController) ChangeActiveVipRuleStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	ActiveVipService := new(services.ActiveVipService)
	request := services.ChangeActiveVipRuleStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = ActiveVipService.ChangeActiveVipRuleStatus(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type AllActiveVipTaskResponse struct {
	Code int                       `json:"code" example:"200"`    // 状态码
	Msg  string                    `json:"msg" example:"success"` // 提示信息
	Data []models.ActiveVipWelfare `json:"data"`
}

// @Summary 获取所有VIP活动福利
// @Success 200 {object} controllers.AllActiveVipTaskResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveVipTask [get]
func (c *ActiveVipController) AllActiveVipTask() {
	ActiveId := c.GetQueryInt("active_id", 0)
	request := services.AllActiveVipTaskRequest{
		ActiveId: ActiveId,
	}
	data, err := (&services.ActiveVipService{}).AllActiveVipTask(request)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(data)
}

// @Summary	Vip任务所有条件类型
// @Success 200 {object} controllers.AllActiveVipTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveVipTaskConditionType [get]
func (c *ActiveVipController) AllActiveVipTaskConditionType() {
	ActivityVipService := new(services.ActiveVipService)
	lang := c.Lang
	data := ActivityVipService.AllActiveVipTaskConditionType(lang)
	c.JSONSuccess(data)
}

// @Summary	Vip任务所有奖励类型
// @Success 200 {object} controllers.AllActiveVipTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allActiveVipTaskCycleType [get]
func (c *ActiveVipController) AllActiveVipTaskCycleType() {
	ActivityVipService := new(services.ActiveVipService)
	lang := c.Lang
	data := ActivityVipService.AllActiveVipTaskCycleType(lang)
	c.JSONSuccess(data)
}

// @Summary 添加VIP活动福利
// @Param   request body services.AddActiveVipTaskRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActiveVipTask [post]
func (c *ActiveVipController) AddActiveVipTask() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	request := services.AddActiveVipTaskRequest{
		ActiveId:  getInt(params, "active_id"),
		Cycle:     getInt(params, "cycle"),
		Lv:        getInt(params, "lv"),
		TotalPays: getFloat(params, "total_pays"),
		TotalBets: getFloat(params, "total_bets"),
		ConAnd:    getInt(params, "con_and"),
		Rewards:   getFloat(params, "rewards"),
		Status:    getInt(params, "status"),
	}
	err = (&services.ActiveVipService{}).AddActiveVipTask(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type EditActiveVipTaskRequest struct {
	services.AddActiveVipTaskRequest
	Id int `json:"id"`
}

// @Summary 编辑VIP活动福利
// @Param   request body controllers.EditActiveVipTaskRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveVipTask [post]
func (c *ActiveVipController) EditActiveVipTask() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	id := getInt(params, "id")
	request := services.AddActiveVipTaskRequest{
		ActiveId:  getInt(params, "active_id"),
		Cycle:     getInt(params, "cycle"),
		Lv:        getInt(params, "lv"),
		TotalPays: getFloat(params, "total_pays"),
		TotalBets: getFloat(params, "total_bets"),
		ConAnd:    getInt(params, "con_and"),
		Rewards:   getFloat(params, "rewards"),
		Status:    getInt(params, "status"),
	}
	err = (&services.ActiveVipService{}).EditActiveVipTask(id, request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteActiveVipTaskRequest struct {
	Id int `json:"id"`
}

// @Summary 删除VIP活动福利
// @Param   request body controllers.DeleteActiveVipTaskRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteActiveVipTask [post]
func (c *ActiveVipController) DeleteActiveVipTask() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	id := getInt(params, "id")
	err = (&services.ActiveVipService{}).DeleteActiveVipTask(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改Vip规则状态
// @Param	request body services.ChangeActiveVipTaskStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeActiveVipTaskStatus [post]
func (c *ActiveVipController) ChangeActiveVipTaskStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	ActiveVipService := new(services.ActiveVipService)
	request := services.ChangeActiveVipTaskStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = ActiveVipService.ChangeActiveVipTaskStatus(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
