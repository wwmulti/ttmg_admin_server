package controllers

import (
	"api/config"
	"api/models"
	"api/services"
	"encoding/json"
)

type ActiveController struct {
	BaseController
}

type ActiveListDTO struct {
	List        []models.Actives `json:"list"`         // 列表数据
	Total       int64            `json:"total"`        // 总条数
	CurrentPage int              `json:"current_page"` // 当前页
}
type ActiveListResponse struct {
	Code int           `json:"code" example:"200"`    // 状态码
	Msg  string        `json:"msg" example:"success"` // 提示信息
	Data ActiveListDTO `json:"data"`
}

// @Summary	活动列表
// @Param package_id query int false "开放平台"
// @Param status query int false "活动状态"
// @Param is_newcomer query int false "是否新人"
// @Param is_pay query int false "是否充值"
// @Param is_hide_complete query int false "完成后是否隐藏"
// @Param need_reload query int false "是否刷新"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.ActiveListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /activeList [get]
func (c *ActiveController) ActiveList() {
	request := services.ActiveListRequest{}
	c.ParseForm(&request)
	request.RawQuery = c.Ctx.Request.URL.Query()

	activityService := new(services.ActiveService)
	needReload := c.GetQueryInt("need_reload", 0)
	lang := c.Lang

	if request.PackageId > 0 {
		if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
			c.JSONError(500, "ShuJuYiChang")
		}
	} else {
		request.PackageIds = c.GetPackageIds(int(c.UserId))
	}
	data, err := activityService.ActiveList(request, needReload, lang, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加活动
// @Param	request body services.AddActiveRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addActive [post]
func (c *ActiveController) AddActive() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activityService := new(services.ActiveService)
	request := services.AddActiveRequest{
		Name:         getString(params, "name"),
		PtName:       getString(params, "pt_name"),
		EnName:       getString(params, "en_name"),
		Icon:         getString(params, "icon"),
		PcIcon:       getString(params, "pc_icon"),
		ActiveTypeId: getInt(params, "active_type_id"),
		RedirectLink: getInt(params, "redirect_link"),
		PtDesc:       getString(params, "pt_desc"),
		EnDesc:       getString(params, "en_desc"),
		StartTime:    getInt64(params, "start_time"),
		EndTime:      getInt64(params, "end_time"),
		CollecTime:   getInt64(params, "collec_time"),
		Sort:         getInt(params, "sort"),
		Status:       getInt(params, "status"),
		PackageIds:   getString(params, "package_ids"),
	}
	err = activityService.AddActive(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑活动
// @Param	request body services.EditActiveRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActive [post]
func (c *ActiveController) EditActive() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activityService := new(services.ActiveService)
	request := services.EditActiveRequest{
		Id:           getInt(params, "id"),
		Name:         getString(params, "name"),
		PtName:       getString(params, "pt_name"),
		EnName:       getString(params, "en_name"),
		Icon:         getString(params, "icon"),
		PcIcon:       getString(params, "pc_icon"),
		ActiveTypeId: getInt(params, "active_type_id"),
		RedirectLink: getInt(params, "redirect_link"),
		PtDesc:       getString(params, "pt_desc"),
		EnDesc:       getString(params, "en_desc"),
		StartTime:    getInt64(params, "start_time"),
		EndTime:      getInt64(params, "end_time"),
		CollecTime:   getInt64(params, "collec_time"),
		Sort:         getInt(params, "sort"),
		Status:       getInt(params, "status"),
	}
	err = activityService.EditActive(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type EditActiveAttrRequest struct {
	Id    int    `json:"id"`
	Field string `json:"field"`
}

// @Summary	编辑活动属性
// @Param	request body controllers.EditActiveAttrRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editActiveAttr [post]
func (c *ActiveController) EditActiveAttr() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	activityService := new(services.ActiveService)
	id := getInt(params, "id")
	field := getString(params, "field")

	err = activityService.EditActiveAttr(id, field, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteActiveRequest struct {
	Id int `json:"id" example:"123456"` // 活动id
}

// @Summary	删除活动
// @Param	request body controllers.DeleteActiveRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteActive [post]
func (c *ActiveController) DeleteActive() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	activityService := new(services.ActiveService)
	err = activityService.DeleteActive(id, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type ActiveTypeListResponse struct {
	JSONResponseTpl
	Data []config.ActiveTypeItem `json:"data"`
}

// @Summary	活动类型列表
// @Success 200 {object} controllers.ActiveTypeListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /activeTypeList [get]
func (c *ActiveController) ActiveTypeList() {
	activityService := new(services.ActiveService)
	data := activityService.GetActiveTypeList()
	c.JSONSuccess(data)
}

type UserActivityRewardLogListResponse struct {
	Code int                          `json:"code" example:"200"`    // 状态码
	Msg  string                       `json:"msg" example:"success"` // 提示信息
	Data models.UserActivityRewardLog `json:"data"`
}

// @Summary	用户活动奖励日志
// @Param activity_id query int false "活动ID"
// @Param user_id query int false "用户ID"
// @Param activity_type_id query int false "活动类型ID"
// @Param receive_type_id query int false "领取类型ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.UserActivityRewardLogListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /userActivityRewardLogList [get]
func (c *ActiveController) UserActivityRewardLogList() {
	activityId := c.GetQueryInt("activity_id", 0)
	userId := c.GetQueryInt("user_id", 0)
	activityTypeId := c.GetQueryInt("activity_type_id", 0)
	receiveTypeId := c.GetQueryInt("receive_type_id", 0)
	page, pageSize := c.GetPagination()
	request := services.UserActivityRewardLogRequest{
		ActivityId:     activityId,
		UserId:         userId,
		ActivityTypeId: activityTypeId,
		ReceiveTypeId:  receiveTypeId,
		Page:           page,
		PageSize:       pageSize,
	}
	activityService := new(services.ActiveService)
	data, err := activityService.UserActivityRewardLogList(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}
