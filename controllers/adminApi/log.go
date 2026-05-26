package controllers

import (
	"api/models"
	"api/services"
)

type LogController struct {
	BaseController
}

type UserBalanceLogListDTO struct {
	List        []models.UserBalanceLog `json:"list"`         // 列表数据
	Total       int64                   `json:"total"`        // 总条数
	CurrentPage int                     `json:"current_page"` // 当前页
}
type UserBalanceLogListResponse struct {
	Code int                   `json:"code" example:"200"`    // 状态码
	Msg  string                `json:"msg" example:"success"` // 提示信息
	Data UserBalanceLogListDTO `json:"data"`
}

// @Summary	用户余额变动日志列表
// @Param user_id query int false "用户ID"
// @Param package_id query int false "包ID"
// @Param type query int false "日志类型"
// @Param activity_type_id query int false "活动类型ID"
// @Param activity_id query int false "活动ID"
// @Param game_type_id query int false "游戏类型ID"
// @Param game_id query int false "游戏ID"
// @Param platform_id query int false "平台ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.UserBalanceLogListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /userBalanceLogList [get]
func (c *LogController) UserBalanceLogList() {
	page, pageSize := c.GetPagination()
	beginTime, endTime := c.GetSearchTime()

	gameTypeName := c.GetString("game_type_name", "")
	packageId := c.GetQueryInt("package_id", 0)
	packageIds := c.GetPackageIds(int(c.UserId))
	gameTypeId := (services.GameTypeService{}).GetIds(packageIds, gameTypeName)

	request := services.UserBalanceLogRequestParams{
		BeginTime:      beginTime,
		EndTime:        endTime + 1, // 多表查询里时间是小于
		Id:             c.GetQueryInt("id", 0),
		RoleId:         c.GetQueryInt("role_id", 0),
		UserId:         c.GetQueryInt("user_id", 0),
		Type:           c.GetQueryInt("type", 0),
		ActivityTypeId: c.GetQueryInt("activity_type_id", 0),
		ActivityId:     c.GetQueryInt("activity_id", 0),
		GameTypeId:     c.GetQueryInt("game_type_id", 0),
		GameTypeIds:    gameTypeId,
		GameId:         c.GetQueryInt("game_id", 0),
		PlatformId:     c.GetQueryInt("platform_id", 0),
		Page:           page,
		PageSize:       pageSize,
		OrderBy:        "created_time DESC",
	}

	if packageId > 0 { // 优先
		request.PackageId = packageId
	} else if len(packageIds) > 0 {
		request.PackageIds = packageIds
	}

	needReload := c.GetQueryInt("need_reload", 0)
	lang := c.Lang

	if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
		c.JSONError(500, "ShuJuYiChang")
	}

	LogService := new(services.LogService)
	data, err := LogService.UserBalanceLogList(request, needReload, lang, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type AllUserBalanceLogTypeResponse struct {
	Code int                  `json:"code" example:"200"`    // 状态码
	Msg  string               `json:"msg" example:"success"` // 提示信息
	Data []UserBalanceLogType `json:"data"`
}

type UserBalanceLogType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @Summary	所有用户余额变动日志类型
// @Success 200 {object} controllers.AllUserBalanceLogTypeResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /allUserBalanceLogType [get]
func (c *LogController) AllUserBalanceLogType() {
	LogService := new(services.LogService)
	lang := c.Lang
	data := LogService.AllUserBalanceLogType(lang)
	c.JSONSuccess(data)
}

type GameWaterLogListDTO struct {
	List        []models.GameWaterLog `json:"list"`         // 列表数据
	Total       int64                 `json:"total"`        // 总条数
	CurrentPage int                   `json:"current_page"` // 当前页
}
type GameWaterLogListResponse struct {
	Code int                 `json:"code" example:"200"`    // 状态码
	Msg  string              `json:"msg" example:"success"` // 提示信息
	Data GameWaterLogListDTO `json:"data"`
}

// @Summary	用户游戏押注日志
// @Param user_id query int false "用户ID"
// @Param order_id query int false "订单ID"
// @Param parent_order_id query int false "父订单ID"
// @Param game_cat_id query int false "游戏类型ID"
// @Param game_id query int false "游戏ID"
// @Param game_platform_id query int false "平台ID"
// @Param createTime query []string false "创建时间"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.GameWaterLogListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /gameWaterLogList [get]
func (c *LogController) GameWaterLogList() {
	page, pageSize := c.GetPagination()
	beginTime, endTime := c.GetSearchTime()

	gameCateName := c.GetString("game_cat_name", "")
	gamePlatformName := c.GetString("game_platform_name", "")
	packageIds := c.GetPackageIds(int(c.UserId))
	request := services.GameWaterLogRequestParams{
		UserId:         c.GetQueryInt("user_id", 0),
		RoleId:         c.GetQueryInt("role_id", 0),
		GameId:         c.GetQueryInt("game_id", 0),
		OrderId:        c.GetString("order_id", ""),
		PackageId:      c.GetQueryInt("package_id", 0),
		ParentOrderId:  c.GetString("parent_order_id", ""),
		Page:           page,
		PageSize:       pageSize,
		BeginTime:      beginTime,
		EndTime:        endTime,
		OrderBy:        "c_time DESC",
		GameCateId:     (services.GameTypeService{}).GetIds(packageIds, gameCateName),
		GamePlatformId: (services.PlatformService{}).GetIds(packageIds, gamePlatformName),
		PackageIds:     packageIds,
	}
	needReload := c.GetQueryInt("need_reload", 0)

	if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
		c.JSONError(500, "ShuJuYiChang")
	}

	LogService := new(services.LogService)
	data, err := LogService.GameWaterLogList(request, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	系统日志
// @Param id query int false "ID"
// @Param admin_id query int false "管理员ID"
// @Param role_id query int false "角色ID"
// @Param path query string false "操作方法"
// @Param body query string false "操作数据"
// @Param c_time query []string false "操作时间"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.GameListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /systemLog [get]
func (c *LogController) SystemLog() {
	request := services.GetSystemLogRequest{}
	c.ParseForm(&request)
	request.RawQuery = c.Ctx.Request.URL.Query()

	if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
		c.JSONError(500, "ShuJuYiChang")
	}

	if request.PackageId > 0 {
		if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
			c.JSONError(500, "ShuJuYiChang")
		}
	} else {
		request.PackageIds = c.GetPackageIds(int(c.UserId))
	}

	data, err := (services.LogService{}).GetSystemLog(request, int(c.UserId), c.Lang)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type UserLoginLogListResponse struct {
	Code             int    `json:"code" example:"200"`    // 状态码
	Msg              string `json:"msg" example:"success"` // 提示信息
	UserLoginLogList `json:"data"`
}

type UserLoginLogList struct {
	List        []models.LoginLog `json:"list"`
	Total       int64             `json:"total"`
	Packages    []models.Package  `json:"packages"`
	CurrentPage int               `json:"current_page"`
}

// @Summary	用户登录记录
// @Param user_id query int true "用户ID"
// @Param begin_time query int true "起始时间"
// @Param end_time query int true "结束时间"
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Success 200 {object} controllers.UserLoginLogListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /userLoginLogList [get]
func (c *LogController) UserLoginLogList() {
	beginTime, endTime := c.GetSearchTime()
	request := services.LoginLogRequestParams{
		UserId:    c.GetQueryInt("user_id", 0),
		BeginTime: beginTime,
		EndTime:   endTime,
		Page:      c.GetQueryInt("page", 1),
		PageSize:  c.GetQueryInt("page_size", 20),
	}
	needReload := c.GetQueryInt("need_reload", 0)
	userService := new(services.UserService)
	data, err := userService.UserLoginLogList(request, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type UserOperateLogListResponse struct {
	Code               int    `json:"code" example:"200"`    // 状态码
	Msg                string `json:"msg" example:"success"` // 提示信息
	UserOperateLogList `json:"data"`
}

type UserOperateLogList struct {
	List        []models.UserOperateLog `json:"list"`
	Packages    []models.Package        `json:"packages"`
	Total       int64                   `json:"total"`
	CurrentPage int                     `json:"current_page"`
}

// @Summary	用户操作记录
// @Param user_id query int true "用户ID"
// @Param role_id query int true "用户RoleID"
// @Param start_time query int true "起始时间"
// @Param end_time query int true "结束时间"
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Success 200 {object} controllers.UserOperateLogListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /userOperateLogList [get]
func (c *LogController) UserOperateLogList() {
	beginTime, endTime := c.GetSearchTime()
	request := services.UserOperateLogRequestParams{
		UserId:    c.GetQueryInt("user_id", 0),
		RoleId:    c.GetQueryInt("role_id", 0),
		BeginTime: beginTime,
		EndTime:   endTime,
		Page:      c.GetQueryInt("page", 1),
		PageSize:  c.GetQueryInt("page_size", 20),
	}
	needReload := c.GetQueryInt("need_reload", 0)
	userService := new(services.UserService)
	data, err := userService.UserOperateLogList(request, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}
