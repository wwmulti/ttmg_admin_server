package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type AlertController struct {
	BaseController
}

type AlertListDTO struct {
	List        []models.Alert `json:"list"`         // 列表数据
	Total       int64          `json:"total"`        // 总条数
	CurrentPage int            `json:"current_page"` // 当前页
}
type AlertListResponse struct {
	Code int          `json:"code" example:"200"`    // 状态码
	Msg  string       `json:"msg" example:"success"` // 提示信息
	Data AlertListDTO `json:"data"`
}

// @Summary	弹窗列表
// @Param id query string false "ID"
// @Param name query int false "弹窗名称"
// @Param status query int false "弹窗状态"
// @Param need_reload query int false "是否刷新"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.AlertListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /alertList [get]
func (c *AlertController) AlertList() {
	request := services.AlertListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	AlertService := new(services.AlertService)
	lang := c.Lang
	needReload := c.GetQueryInt("need_reload", 0)
	data, err := AlertService.AlertList(request, needReload, lang, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type AllAlertTypeResponse struct {
	Code int         `json:"code" example:"200"`    // 状态码
	Msg  string      `json:"msg" example:"success"` // 提示信息
	Data []AlertType `json:"data"`
}

type AlertType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	所有弹窗类型
// @Success 200 {object} controllers.AllAlertTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allAlertType [get]
func (c *AlertController) AllAlertType() {
	AlertService := new(services.AlertService)
	lang := c.Lang
	data := AlertService.AllAlertType(lang)
	c.JSONSuccess(data)
}

type AllAlertRuleTypeResponse struct {
	Code int             `json:"code" example:"200"`    // 状态码
	Msg  string          `json:"msg" example:"success"` // 提示信息
	Data []AlertRuleType `json:"data"`
}

type AlertRuleType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	所有弹窗规则类型
// @Success 200 {object} controllers.AllAlertRuleTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allAlertRuleType [get]
func (c *AlertController) AllAlertRuleType() {
	AlertService := new(services.AlertService)
	lang := c.Lang
	data := AlertService.AllAlertRuleType(lang)
	c.JSONSuccess(data)
}

// @Summary	添加弹窗
// @Param	request body services.AddAlertRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addAlert [post]
func (c *AlertController) AddAlert() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	AlertService := new(services.AlertService)
	request := services.AddAlertRequest{
		PackageIds:  getIntSlice(params, "package_ids"),
		Name:        getString(params, "name"),
		Type:        getInt(params, "type"),
		ContentType: getInt(params, "content_type"),
		EnTitle:     getString(params, "en_title"),
		PtTitle:     getString(params, "pt_title"),
		EnContent:   getString(params, "en_content"),
		PtContent:   getString(params, "pt_content"),
		Image:       getString(params, "image"),
		Sort:        getInt(params, "sort"),
		AlertRule:   getInt(params, "alert_rule"),
		AlertHours:  getInt(params, "alert_hours"),
		Status:      getInt(params, "status"),
	}
	err = AlertService.AddAlert(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑弹窗
// @Param	request body services.EditAlertRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editAlert [post]
func (c *AlertController) EditAlert() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	AlertService := new(services.AlertService)
	request := services.EditAlertRequest{
		Id:          getInt(params, "id"),
		Name:        getString(params, "name"),
		Type:        getInt(params, "type"),
		ContentType: getInt(params, "content_type"),
		EnTitle:     getString(params, "en_title"),
		PtTitle:     getString(params, "pt_title"),
		EnContent:   getString(params, "en_content"),
		PtContent:   getString(params, "pt_content"),
		Image:       getString(params, "image"),
		Sort:        getInt(params, "sort"),
		AlertRule:   getInt(params, "alert_rule"),
		AlertHours:  getInt(params, "alert_hours"),
		Status:      getInt(params, "status"),
	}
	err = AlertService.EditAlert(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteAlertRequest struct {
	Id string `json:"id" example:"123456"` // 弹窗id
}

// @Summary	删除弹窗
// @Param	request body controllers.DeleteAlertRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteAlert [post]
func (c *AlertController) DeleteAlert() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	AlertService := new(services.AlertService)
	err = AlertService.DeleteAlert(id, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改弹窗状态
// @Param	request body services.ChangeAlertStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeAlertStatus [post]
func (c *AlertController) ChangeAlertStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	AlertService := new(services.AlertService)
	request := services.ChangeAlertStatusRequest{
		Id:     id,
		Status: getInt(params, "status"),
	}
	err = AlertService.ChangeAlertStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type RefreshPackageAlertRequest struct {
	PackageId int `json:"package_id"`
}

// @Summary	修改浮动图标状态
// @Param	request body controllers.RefreshPackageAlertRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /refreshPackageAlert [post]
func (c *AlertController) RefreshPackageAlert() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	packageId := getInt(params, "package_id")
	err = services.AlertService{}.RefreshPackageAlert(packageId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
