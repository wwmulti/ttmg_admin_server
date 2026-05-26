package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type BroadcastController struct {
	BaseController
}

type BroadcastListDTO struct {
	List        []models.Broadcast `json:"list"`         // 列表数据
	Total       int64              `json:"total"`        // 总条数
	CurrentPage int                `json:"current_page"` // 当前页
}
type BroadcastListResponse struct {
	Code int              `json:"code" example:"200"`    // 状态码
	Msg  string           `json:"msg" example:"success"` // 提示信息
	Data BroadcastListDTO `json:"data"`
}

// @Summary	广播列表
// @Param start_hour query string false "广播开始时间"
// @Param end_hour query int false "广播结束时间"
// @Param need_reload query int false "是否刷新"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.BroadcastListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /broadcastList [get]
func (c *BroadcastController) BroadcastList() {
	request := services.BroadcastListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	needReload := c.GetQueryInt("need_reload", 0)
	BroadcastService := new(services.BroadcastService)
	data, err := BroadcastService.BroadcastList(request, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加广播
// @Param	request body services.AddBroadcastRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addBroadcast [post]
func (c *BroadcastController) AddBroadcast() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	BroadcastService := new(services.BroadcastService)
	request := services.AddBroadcastRequest{
		PackageIds: getIntSlice(params, "package_ids"),
		Name:       getString(params, "name"),
		EnContent:  getString(params, "en_content"),
		PtContent:  getString(params, "pt_content"),
		Sort:       getInt(params, "sort"),
		Status:     getInt(params, "status"),
	}
	err = BroadcastService.AddBroadcast(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑广播
// @Param	request body services.EditBroadcastRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editBroadcast [post]
func (c *BroadcastController) EditBroadcast() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	BroadcastService := new(services.BroadcastService)
	request := services.EditBroadcastRequest{
		Id:        getInt(params, "id"),
		Name:      getString(params, "name"),
		EnContent: getString(params, "en_content"),
		PtContent: getString(params, "pt_content"),
		Sort:      getInt(params, "sort"),
		Status:    getInt(params, "status"),
	}
	err = BroadcastService.EditBroadcast(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteBroadcastRequest struct {
	Id string `json:"id" example:"123456"` // 广播id
}

// @Summary	删除广播
// @Param	request body controllers.DeleteBroadcastRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteBroadcast [post]
func (c *BroadcastController) DeleteBroadcast() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	BroadcastService := new(services.BroadcastService)
	err = BroadcastService.DeleteBroadcast(id, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改广播状态
// @Param	request body services.ChangeBroadcastStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeBroadcastStatus [post]
func (c *BroadcastController) ChangeBroadcastStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	broadcastService := new(services.BroadcastService)
	request := services.ChangeBroadcastStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = broadcastService.ChangeBroadcastStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
