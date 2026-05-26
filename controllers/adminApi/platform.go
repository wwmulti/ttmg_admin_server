package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type PlatformController struct {
	BaseController
}

type PlatformListDTO struct {
	List        []models.PlatformDTO `json:"list"`         // 列表数据
	Total       int64                `json:"total"`        // 总条数
	CurrentPage int                  `json:"current_page"` // 当前页
}
type PlatformListResponse struct {
	Code int             `json:"code" example:"200"`    // 状态码
	Msg  string          `json:"msg" example:"success"` // 提示信息
	Data PlatformListDTO `json:"data"`
}

// @Summary	平台列表
// @Param id query int false "平台ID"
// @Param name query string false "平台名称"
// @Param status query int false "平台状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.PlatformListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /platformList [get]
func (c *PlatformController) PlatformList() {
	request := services.PlatformListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	platformService := new(services.PlatformService)
	data, err := platformService.PlatformList(request, int(c.UserId))
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type AllPlatformResponse struct {
	Code int               `json:"code" example:"200"`    // 状态码
	Msg  string            `json:"msg" example:"success"` // 提示信息
	Data []models.Platform `json:"data"`
}

// @Summary	所有平台
// @Success 200 {object} controllers.AllPlatformResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allPlatform [get]
func (c *PlatformController) AllPlatform() {
	platformService := new(services.PlatformService)
	data, err := platformService.AllPlatformByPackageIds(c.GetPackageIds(int(c.UserId)))
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加平台
// @Param	request body services.AddPlatformRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addPlatform [post]
func (c *PlatformController) AddPlatform() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	platformService := new(services.PlatformService)
	request := services.AddPlatformRequest{
		Name:          getString(params, "name"),
		Alias:         getString(params, "alias"),
		Logo:          getString(params, "logo"),
		ClickLogo:     getString(params, "click_logo"),
		Image:         getString(params, "image"),
		FrontColor:    getString(params, "front_color"),
		MiniMoney:     getFloat(params, "mini_money"),
		Status:        getInt(params, "status"),
		Sort:          getInt(params, "sort"),
		GameShowCount: getInt(params, "game_show_count"),
		GameShowMore:  getInt(params, "game_show_more"),
		ApiRate:       getFloat(params, "api_rate"),
		PackageId:     1,
	}

	if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
		c.JSONError(500, "ShuJuYiChang")
	}
	err = platformService.AddPlatform(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑平台
// @Param	request body services.EditPlatformRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editPlatform [post]
func (c *PlatformController) EditPlatform() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	platformService := new(services.PlatformService)
	request := services.EditPlatformRequest{
		Id:            getInt(params, "id"),
		Name:          getString(params, "name"),
		Alias:         getString(params, "alias"),
		Logo:          getString(params, "logo"),
		ClickLogo:     getString(params, "click_logo"),
		Image:         getString(params, "image"),
		FrontColor:    getString(params, "front_color"),
		MiniMoney:     getFloat(params, "mini_money"),
		Status:        getInt(params, "status"),
		Sort:          getInt(params, "sort"),
		GameShowCount: getInt(params, "game_show_count"),
		GameShowMore:  getInt(params, "game_show_more"),
		ApiRate:       getFloat(params, "api_rate"),
	}
	err = platformService.EditPlatform(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type EditPlatformAttrRequest struct {
	Id    int    `json:"id"`
	Field string `json:"field"`
}

// @Summary	编辑平台属性
// @Param	request body controllers.EditPlatformAttrRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editPlatformAttr [post]
func (c *PlatformController) EditPlatformAttr() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	platformService := new(services.PlatformService)
	id := getInt(params, "id")
	field := getString(params, "field")

	err = platformService.EditPlatformAttr(id, field)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeletePlatformRequest struct {
	Id string `json:"id" example:"123456"` // 平台id
}

// @Summary	删除平台
// @Param	request body controllers.DeletePlatformRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deletePlatform [post]
func (c *PlatformController) DeletePlatform() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	platformService := new(services.PlatformService)
	err = platformService.DeletePlatform(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
