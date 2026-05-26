package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type BannerController struct {
	BaseController
}

type BannerListDTO struct {
	List        []models.Banner `json:"list"`         // 列表数据
	Total       int64           `json:"total"`        // 总条数
	CurrentPage int             `json:"current_page"` // 当前页
}
type BannerListResponse struct {
	Code int           `json:"code" example:"200"`    // 状态码
	Msg  string        `json:"msg" example:"success"` // 提示信息
	Data BannerListDTO `json:"data"`
}

// @Summary	轮播图列表
// @Param id query string false "ID"
// @Param package_id query string false "包ID"
// @Param name query int false "轮播图名称"
// @Param status query int false "轮播图状态"
// @Param need_reload query string false "是否刷新"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.BannerListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /bannerList [get]
func (c *BannerController) BannerList() {
	request := services.BannerListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	BannerService := new(services.BannerService)
	needReload := c.GetQueryInt("need_reload", 0)
	lang := c.Lang
	data, err := BannerService.BannerList(request, needReload, lang, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type AllBannerTypeResponse struct {
	Code int          `json:"code" example:"200"`    // 状态码
	Msg  string       `json:"msg" example:"success"` // 提示信息
	Data []BannerType `json:"data"`
}

type BannerType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	所有轮播图类型
// @Success 200 {object} controllers.AllBannerTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allBannerType [get]
func (c *BannerController) AllBannerType() {
	BannerService := new(services.BannerService)
	lang := c.Lang
	data := BannerService.AllBannerType(lang)
	c.JSONSuccess(data)
}

// @Summary	添加轮播图
// @Param	request body services.AddBannerRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addBanner [post]
func (c *BannerController) AddBanner() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	BannerService := new(services.BannerService)
	request := services.AddBannerRequest{
		PackageIds: getIntSlice(params, "package_ids"),
		Name:       getString(params, "name"),
		Type:       getInt(params, "type"),
		Url:        getString(params, "url"),
		Sort:       getInt(params, "sort"),
		Status:     getInt(params, "status"),
	}
	err = BannerService.AddBanner(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑轮播图
// @Param	request body services.EditBannerRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editBanner [post]
func (c *BannerController) EditBanner() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	BannerService := new(services.BannerService)
	request := services.EditBannerRequest{
		Id:     getInt(params, "id"),
		Name:   getString(params, "name"),
		Type:   getInt(params, "type"),
		Url:    getString(params, "url"),
		Sort:   getInt(params, "sort"),
		Status: getInt(params, "status"),
	}
	err = BannerService.EditBanner(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteBannerRequest struct {
	Id string `json:"id" example:"123456"` // 轮播图id
}

// @Summary	删除轮播图
// @Param	request body controllers.DeleteBannerRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteBanner [post]
func (c *BannerController) DeleteBanner() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	BannerService := new(services.BannerService)
	err = BannerService.DeleteBanner(id, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改轮播图状态
// @Param	request body services.ChangeBannerStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeBannerStatus [post]
func (c *BannerController) ChangeBannerStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	BannerService := new(services.BannerService)
	request := services.ChangeBannerStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = BannerService.ChangeBannerStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type RefreshPackageBannerRequest struct {
	PackageId int `json:"package_id"`
}

// @Summary	修改轮播图状态
// @Param	request body controllers.RefreshPackageBannerRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /refreshPackageBanner [post]
func (c *BannerController) RefreshPackageBanner() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	packageId := getInt(params, "package_id")
	err = services.BannerService{}.RefreshPackageBanner(packageId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
