package controllers

import (
	"api/config"
	"api/models"
	"api/services"
	"encoding/json"
)

type ConfigController struct {
	BaseController
}

type AllSystemConfigResponse struct {
	Code int             `json:"code" example:"200"`    // 状态码
	Msg  string          `json:"msg" example:"success"` // 提示信息
	Data []models.Config `json:"data"`
}

// @Summary	系统配置列表
// @Param type_id query string false "配置类型ID"
// @Success 200 {object} controllers.AllSystemConfigResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allSystemConfig [get]
func (c *ConfigController) AllSystemConfig() {
	configService := new(services.ConfigService)
	request := services.AllSystemConfigRequest{
		PackageId: c.GetQueryInt("package_id", 0),
		TypeId:    c.GetQueryInt("type_id", 0),
	}
	needReload := c.GetQueryInt("need_reload", 0)
	data, err := configService.AllSystemConfig(request, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	编辑系统配置
// @Param	request body []services.EditConfigRequest true "json请求参数"
// @Success 200 {object} controllers.AllSystemConfigResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editSystemConfig [post]
func (c *ConfigController) EditSystemConfig() {
	var params []services.EditConfigRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	configService := new(services.ConfigService)
	err = configService.EditSystemConfig(params, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	转盘配置列表
// @Success 200 {object} controllers.AllSystemConfigResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allLuckyWheelConfig [get]
func (c *ConfigController) AllLuckyWheelConfig() {
	configService := new(services.ConfigService)
	request := services.AllSystemConfigRequest{
		PackageId: c.GetQueryInt("package_id", 0),
		TypeId:    int(config.SystemConfigTypeLuckWheel),
	}
	userIdInterface := c.Data["UserId"]
	userId := userIdInterface.(int64)
	needReload := c.GetQueryInt("need_reload", 0)
	data, err := configService.AllSystemConfig(request, needReload, userId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	编辑转盘配置
// @Param	request body []services.EditConfigRequest true "json请求参数"
// @Success 200 {object} controllers.AllSystemConfigResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editLuckyWheelConfig [post]
func (c *ConfigController) EditLuckyWheelConfig() {
	var params []services.EditConfigRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	configService := new(services.ConfigService)
	err = configService.EditSystemConfig(params, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	奖池配置列表
// @Success 200 {object} controllers.AllSystemConfigResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allPrizeConfig [get]
func (c *ConfigController) AllPrizeConfig() {
	configService := new(services.ConfigService)
	request := services.AllSystemConfigRequest{
		PackageId: c.GetQueryInt("package_id", 0),
		TypeId:    int(config.SystemConfigTypePrize),
	}
	needReload := c.GetQueryInt("need_reload", 0)
	data, err := configService.AllSystemConfig(request, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	编辑奖池配置
// @Param	request body []services.EditConfigRequest true "json请求参数"
// @Success 200 {object} controllers.AllSystemConfigResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editPrizeConfig [post]
func (c *ConfigController) EditPrizeConfig() {
	var params []services.EditConfigRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	configService := new(services.ConfigService)
	err = configService.EditSystemConfig(params, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	代理配置列表
// @Success 200 {object} controllers.AllSystemConfigResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allAgentConfig [get]
func (c *ConfigController) AllAgentConfig() {
	configService := new(services.ConfigService)
	request := services.AllSystemConfigRequest{
		PackageId: c.GetQueryInt("package_id", 0),
		TypeId:    int(config.SystemConfigTypeAgent),
	}
	needReload := c.GetQueryInt("need_reload", 0)
	data, err := configService.AllSystemConfig(request, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	编辑代理配置
// @Param	request body []services.EditConfigRequest true "json请求参数"
// @Success 200 {object} controllers.AllSystemConfigResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editAgentConfig [post]
func (c *ConfigController) EditAgentConfig() {
	var params []services.EditConfigRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	configService := new(services.ConfigService)
	err = configService.EditSystemConfig(params, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑自动出款配置
// @Param	request body []services.EditConfigRequest true "json请求参数"
// @Success 200 {object} controllers.AllSystemConfigResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editAutoOutMoneyConfig [post]
func (c *ConfigController) EditAutoOutMoneyConfig() {
	var params []services.EditConfigRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	configService := new(services.ConfigService)
	err = configService.EditSystemConfig(params, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
