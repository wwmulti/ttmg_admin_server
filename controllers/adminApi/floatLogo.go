package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type FloatLogoController struct {
	BaseController
}

type FloatLogoListDTO struct {
	List        []models.FloatLogo `json:"list"`         // 列表数据
	Total       int64              `json:"total"`        // 总条数
	CurrentPage int                `json:"current_page"` // 当前页
}
type FloatLogoListResponse struct {
	Code int              `json:"code" example:"200"`    // 状态码
	Msg  string           `json:"msg" example:"success"` // 提示信息
	Data FloatLogoListDTO `json:"data"`
}

// @Summary	浮动图标列表
// @Param id query string false "ID"
// @Param package_id query string false "包ID"
// @Param need_reload query int false "是否刷新"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.FloatLogoListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /floatLogoList [get]
func (c *FloatLogoController) FloatLogoList() {
	request := services.FloatLogoListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	needReload := c.GetQueryInt("need_reload", 0)
	FloatLogoService := new(services.FloatLogoService)
	data, err := FloatLogoService.FloatLogoList(request, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type AllFloatLogoTypeResponse struct {
	Code int             `json:"code" example:"200"`    // 状态码
	Msg  string          `json:"msg" example:"success"` // 提示信息
	Data []FloatLogoType `json:"data"`
}

type FloatLogoType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	添加浮动图标
// @Param	request body services.AddFloatLogoRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addFloatLogo [post]
func (c *FloatLogoController) AddFloatLogo() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	FloatLogoService := new(services.FloatLogoService)
	request := services.AddFloatLogoRequest{
		PackageIds: getIntSlice(params, "package_ids"),
		Link:       getString(params, "link"),
		Logo:       getString(params, "logo"),
		Status:     getInt(params, "status"),
	}
	err = FloatLogoService.AddFloatLogo(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑浮动图标
// @Param	request body services.EditFloatLogoRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editFloatLogo [post]
func (c *FloatLogoController) EditFloatLogo() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	FloatLogoService := new(services.FloatLogoService)
	request := services.EditFloatLogoRequest{
		Id:     getInt(params, "id"),
		Link:   getString(params, "link"),
		Logo:   getString(params, "logo"),
		Status: getInt(params, "status"),
	}
	err = FloatLogoService.EditFloatLogo(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteFloatLogoRequest struct {
	Id string `json:"id" example:"123456"` // 浮动图标id
}

// @Summary	删除浮动图标
// @Param	request body controllers.DeleteFloatLogoRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteFloatLogo [post]
func (c *FloatLogoController) DeleteFloatLogo() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	FloatLogoService := new(services.FloatLogoService)
	err = FloatLogoService.DeleteFloatLogo(id, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改浮动图标状态
// @Param	request body services.ChangeFloatLogoStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeFloatLogoStatus [post]
func (c *FloatLogoController) ChangeFloatLogoStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	FloatLogoService := new(services.FloatLogoService)
	request := services.ChangeFloatLogoStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = FloatLogoService.ChangeFloatLogoStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type RefreshPackageFloatLogoRequest struct {
	PackageId int `json:"package_id"`
}

// @Summary	修改浮动图标状态
// @Param	request body controllers.RefreshPackageFloatLogoRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /refreshPackageFloatLogo [post]
func (c *FloatLogoController) RefreshPackageFloatLogo() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	packageId := getInt(params, "package_id")
	err = services.FloatLogoService{}.RefreshPackageFloatLogo(packageId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
