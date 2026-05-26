package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type GameTypeController struct {
	BaseController
}

type GameTypeListDTO struct {
	List        []models.GameType `json:"list"`         // 列表数据
	Total       int64             `json:"total"`        // 总条数
	CurrentPage int               `json:"current_page"` // 当前页
}
type GameTypeListResponse struct {
	Code int             `json:"code" example:"200"`    // 状态码
	Msg  string          `json:"msg" example:"success"` // 提示信息
	Data GameTypeListDTO `json:"data"`
}

// @Summary	游戏类型列表
// @Param id query int false "游戏类型ID"
// @Param name query string false "游戏类型名称"
// @Param alias query string false "游戏类型别名"
// @Param status query int false "游戏类型状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.GameTypeListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /gameTypeList [get]
func (c *GameTypeController) GameTypeList() {
	request := services.GameTypeListRequest{}
	c.ParseForm(&request)
	request.RawQuery = c.Ctx.Request.URL.Query()
	gameTypeService := new(services.GameTypeService)
	data, err := gameTypeService.GameTypeList(request, int(c.UserId))
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type AllGameTypeResponse struct {
	Code int               `json:"code" example:"200"`    // 状态码
	Msg  string            `json:"msg" example:"success"` // 提示信息
	Data []models.GameType `json:"data"`
}

// @Summary	所有游戏类型
// @Success 200 {object} controllers.AllGameTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allGameType [get]
func (c *GameTypeController) AllGameType() {
	gameTypeService := new(services.GameTypeService)
	data, err := gameTypeService.AllGameType(c.GetPackageIds(int(c.UserId)))
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加游戏类型
// @Param	request body services.AddGameTypeRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addGameType [post]
func (c *GameTypeController) AddGameType() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	gameTypeService := new(services.GameTypeService)
	request := services.AddGameTypeRequest{
		Name:           getString(params, "name"),
		Alias:          getString(params, "alias"),
		Logo:           getString(params, "logo"),
		Sort:           getInt(params, "sort"),
		Status:         getInt(params, "status"),
		PageGameNumber: getInt(params, "page_game_number"),
		BetRate:        getInt(params, "bet_rate"),
		PackageId:      1,
		PlatformIds:    getString(params, "platform_ids"),
	}

	if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
		c.JSONError(500, "ShuJuYiChang")
	}

	err = gameTypeService.AddGameType(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑游戏类型
// @Param	request body services.EditGameTypeRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editGameType [post]
func (c *GameTypeController) EditGameType() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	gameTypeService := new(services.GameTypeService)
	request := services.EditGameTypeRequest{
		Id:             getInt(params, "id"),
		Name:           getString(params, "name"),
		Alias:          getString(params, "alias"),
		Logo:           getString(params, "logo"),
		Sort:           getInt(params, "sort"),
		Status:         getInt(params, "status"),
		PageGameNumber: getInt(params, "page_game_number"),
		BetRate:        getInt(params, "bet_rate"),
		PlatformIds:    getString(params, "platform_ids"),
	}
	err = gameTypeService.EditGameType(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type EditGameTypeAttrRequest struct {
	Id    int    `json:"id"`
	Field string `json:"field"`
}

// @Summary	编辑游戏类型属性
// @Param	request body controllers.EditGameTypeAttrRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editGameTypeAttr [post]
func (c *GameTypeController) EditGameTypeAttr() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	gameTypeService := new(services.GameTypeService)
	id := getInt(params, "id")
	field := getString(params, "field")

	err = gameTypeService.EditGameTypeAttr(id, field)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteGameTypeRequest struct {
	Id string `json:"id" example:"123456"` // 游戏类型id
}

// @Summary	删除游戏类型
// @Param	request body controllers.DeleteGameTypeRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteGameType [post]
func (c *GameTypeController) DeleteGameType() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	gameTypeService := new(services.GameTypeService)
	err = gameTypeService.DeleteGameType(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type AgentConfigRequest struct {
	Id          int                           `json:"id"`
	AgentConfig []services.AgentConfigRequest `json:"agent_config"`
}

// @Summary	代理配置
// @Param	request body controllers.AgentConfigRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /agentConfig [post]
func (c *GameTypeController) AgentConfig() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")

	configBytes, _ := json.Marshal(params["agent_config"])
	var configs []map[string]interface{}
	err = json.Unmarshal(configBytes, &configs)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	var agentConfigs []services.AgentConfigRequest
	for _, v := range configs {
		agentConfigs = append(agentConfigs, services.AgentConfigRequest{
			Level:      getInt(v, "level"),
			Members:    getInt(v, "members"),
			Commission: getFloat(v, "commission"),
			Rate:       getFloat(v, "rate"),
		})
	}

	err = (services.GameTypeService{}).AgentConfig(id, agentConfigs)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type BetRateConfigRequest struct {
	Id      int    `json:"id"`             // 游戏类型id
	BetRate string `json:"games_bet_rate"` // 游戏打码比例（格式：id,rate;）
}

// @Summary	游戏打码比例配置
// @Param	request body controllers.BetRateConfigRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /betRateConfig [post]
func (c *GameTypeController) BetRateConfig() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	err = (services.GameTypeService{}).BetRateConfig(getInt(params, "id"), getString(params, "games_bet_rate"))
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
