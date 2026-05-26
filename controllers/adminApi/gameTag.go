package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type GameTagController struct {
	BaseController
}

// DTO 用于 Swagger 文档定义
type GameTagListDTO struct {
	List        []models.GameTag `json:"list"`         // 列表数据
	Total       int64            `json:"total"`        // 总条数
	CurrentPage int              `json:"current_page"` // 当前页
}

type GameTagListResponse struct {
	Code int            `json:"code" example:"200"`    // 状态码
	Msg  string         `json:"msg" example:"success"` // 提示信息
	Data GameTagListDTO `json:"data"`
}

// @Summary 游戏标签列表
// @Param id query int false "标签ID"
// @Param name query string false "标签名称"
// @Param por_name query string false "葡语名称"
// @Param status query int false "状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.GameTagListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /gameTagList [get]
func (c *GameTagController) GameTagList() {
	request := services.GameTagListRequest{}
	c.ParseForm(&request)
	request.RawQuery = c.Ctx.Request.URL.Query()

	if request.PackageId > 0 {
		if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
			c.JSONError(500, "ShuJuYiChang")
		}
	} else {
		request.PackageIds = c.GetPackageIds(int(c.UserId))
	}

	gameTagService := new(services.GameTagService)
	data, err := gameTagService.GameTagList(request, int(c.UserId))
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary 添加游戏标签
// @Param   request body services.AddGameTagRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addGameTag [post]
func (c *GameTagController) AddGameTag() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
		return
	}

	gameTagService := new(services.GameTagService)
	request := services.AddGameTagRequest{
		Name:      getString(params, "name"),
		PtName:    getString(params, "pt_name"),
		Status:    getInt(params, "status"),
		PackageId: getInt(params, "package_id"),
	}

	if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
		c.JSONError(500, "ShuJuYiChang")
	}

	err = gameTagService.AddGameTag(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary 编辑游戏标签
// @Param   request body services.EditGameTagRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editGameTag [post]
func (c *GameTagController) EditGameTag() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
		return
	}

	gameTagService := new(services.GameTagService)
	request := services.EditGameTagRequest{
		Id:     getInt(params, "id"),
		Name:   getString(params, "name"),
		PtName: getString(params, "pt_name"),
		Status: getInt(params, "status"),
	}

	err = gameTagService.EditGameTag(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteGameTagRequest struct {
	Id int `json:"id"`
}

// @Summary 删除游戏标签
// @Param   request body controllers.DeleteGameTagRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteGameTag [post]
func (c *GameTagController) DeleteGameTag() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
		return
	}

	id := getInt(params, "id")
	gameTagService := new(services.GameTagService)
	err = gameTagService.DeleteGameTag(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
