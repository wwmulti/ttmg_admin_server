package controllers

import (
	"api/config"
	"api/models"
	"api/services"
	"encoding/json"
	"fmt"
)

type UserController struct {
	BaseController
}

// @Summary	平台用户管理
// @Param vip_level query int false "会员等级"
// @Param role_id query int false "角色Id"
// @Param username query string false "用户名"
// @Param vip_level query int false "会员等级"
// @Param status query int false "用户状态"
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Success 200 {object} services.UserListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /userList [get]
func (c *UserController) UserList() {
	page, pageSize := c.GetPagination()
	packageId := c.GetQueryInt("package_id", 0)
	roleId := c.GetQueryInt("role_id", 0)
	username := c.GetString("username", "")
	vipLevel := c.GetQueryInt("vip_level", 0)
	status := c.GetQueryInt("status", 0)
	needReload := c.GetQueryInt("need_reload", 0)
	lange := c.Lang
	params := services.UserRequestParams{
		VipLevel:  vipLevel,
		PackageId: packageId,
		RoleId:    roleId,
		Username:  username,
		Status:    status,
		Page:      page,
		PageSize:  pageSize,
	}
	result, err := (&services.UserService{}).UserList(params, needReload, lange, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(result)
}

// @Summary	获取用户信息
// @Param user_id query int true "用户ID"
// @Success 200 {object} services.GetUserResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /getUser [get]
func (c *UserController) GetUser() {
	userId := c.GetQueryInt("user_id", 0)
	result, err := (&services.UserService{}).GetUser(userId, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(result)
	}
}

// @Summary	创建用户
// @Param	request body services.UserRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /addUser [post]
func (c *UserController) AddUser() {
	var params services.AddUserRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	params.Language = c.Lang
	roleErr, _, _ := (&services.UserService{}).AddUser(params)
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

type EditUserRequestParams struct {
	UserId   int    `json:"user_id"`
	Password string `json:"password"`
}

// @Summary	修改用户
// @Param	request body controllers.EditUserRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /editUser [post]
func (c *UserController) EditUser() {
	var params EditUserRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	roleErr := (&services.UserService{}).EditPassword(params.UserId, params.Password, c.UserId)
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

type ClearUserBalanceRequest struct {
	UserId int `json:"user_id"`
}

// @Summary	清除用户余额
// @Param	request body controllers.ClearUserBalanceRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /clearUserBalance [post]
func (c *UserController) ClearUserBalance() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	err = (&services.UserService{}).SetUserBalance(services.SetUserBalanceRequest{
		UserId:  userId,
		Balance: 0,
	}, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	设置用户余额
// @Param	request body controllers.SetUserBalanceRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /setUserBalance [post]
func (c *UserController) SetUserBalance() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request := services.SetUserBalanceRequest{
		UserId:  getInt(params, "user_id"),
		Balance: getFloat(params, "balance"),
	}
	err = (&services.UserService{}).SetUserBalance(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	设置用户等级
// @Param	request body services.SetUserLvRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /setUserLv [post]
func (c *UserController) SetUserLv() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	err = (&services.UserService{}).SetUserLv(services.SetUserLvRequest{
		UserId: params["user_id"].(int),
		Lv:     params["lv"].(int),
	}, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改用户封禁状态
// @Param	request body services.ChangeUserStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeUserStatus [post]
func (c *UserController) ChangeUserStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	UserService := new(services.UserService)
	request := services.ChangeUserStatusRequest{
		UserId: userId,
		Status: getInt(params, "status"),
	}
	err = UserService.ChangeUserStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	设置用户打码量
// @Param	request body services.SetUserEffectiveBetRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /setUserEffectiveBet [post]
func (c *UserController) SetUserEffectiveBet() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request := services.SetUserEffectiveBetRequest{
		UserId:       getInt(params, "user_id"),
		EffectiveBet: getFloat(params, "effective_bet"),
	}
	err = (&services.UserService{}).SetUserEffectiveBet(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	清除用户打码量
// @Param	request body services.SetUserEffectiveBetRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /clearUserEffectiveBet [post]
func (c *UserController) ClearUserEffectiveBet() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request := services.SetUserEffectiveBetRequest{
		UserId:       getInt(params, "user_id"),
		EffectiveBet: 0,
	}
	err = (&services.UserService{}).SetUserEffectiveBet(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	设置用户上级ID
// @Param	request body services.SetUserPidRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /setUserPid [post]
func (c *UserController) SetUserPid() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request := services.SetUserPidRequest{
		UserId:  getInt(params, "user_id"),
		PRoleId: getInt(params, "p_role_id"),
	}
	err = (&services.UserService{}).SetUserPid(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	设置用户备注
// @Param	request body services.SetUserRemarkRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /setUserRemark [post]
func (c *UserController) SetUserRemark() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request := services.SetUserRemarkRequest{
		UserId: getInt(params, "user_id"),
		Remark: getString(params, "remark"),
	}
	err = (&services.UserService{}).SetUserRemark(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	设置用户博主
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /setUserBlogger [post]
func (c *UserController) SetUserBlogger() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	err = (&services.UserService{}).SetUserBlogger(userId, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	设置用户经纪人
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /setUserBroker [post]
func (c *UserController) SetUserBroker() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	err = (&services.UserService{}).SetUserBroker(userId, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	取消用户博主
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /cancelUserBlogger [post]
func (c *UserController) CancelUserBlogger() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	err = (&services.UserService{}).CancelUserBlogger(userId, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	取消用户经纪人
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /cancelUserBroker [post]
func (c *UserController) CancelUserBroker() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	err = (&services.UserService{}).CancelUserBroker(userId, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改用户游戏封禁状态
// @Param	request body services.ChangeUserGameStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeUserGameStatus [post]
func (c *UserController) ChangeUserGameStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	UserService := new(services.UserService)
	request := services.ChangeUserGameStatusRequest{
		UserId:    userId,
		IsBanGame: getInt(params, "is_ban_game"),
	}
	err = UserService.ChangeUserGameStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改用户提现封禁状态
// @Param	request body services.ChangeUserWithdrawStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeUserWithdrawStatus [post]
func (c *UserController) ChangeUserWithdrawStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	UserService := new(services.UserService)
	request := services.ChangeUserWithdrawStatusRequest{
		UserId:        userId,
		IsBanWithdraw: getInt(params, "is_ban_withdraw"),
	}
	err = UserService.ChangeUserWithdrawStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改用户邀请奖励封禁状态
// @Param	request body services.ChangeUserInviteRewardStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeUserInviteRewardStatus [post]
func (c *UserController) ChangeUserInviteRewardStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	UserService := new(services.UserService)
	request := services.ChangeUserInviteRewardStatusRequest{
		UserId:            userId,
		IsBanInviteReward: getInt(params, "is_ban_invite_reward"),
	}
	err = UserService.ChangeUserInviteRewardStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改用户邀请奖励封禁状态
// @Param	request body services.ChangeUserChildBetCommissionStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeUserChildBetCommissionStatus [post]
func (c *UserController) ChangeUserChildBetCommissionStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	UserService := new(services.UserService)
	request := services.ChangeUserChildBetCommissionStatusRequest{
		UserId:                  userId,
		IsBanChildBetCommission: getInt(params, "is_ban_child_bet_commission"),
	}
	err = UserService.ChangeUserChildBetCommissionStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改用户邀请奖励封禁状态
// @Param	request body services.ChangeUserOnlyAllowOfficialGameStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeUserOnlyAllowOfficialGameStatus [post]
func (c *UserController) ChangeUserOnlyAllowOfficialGameStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	UserService := new(services.UserService)
	request := services.ChangeUserOnlyAllowOfficialGameStatusRequest{
		UserId:                  userId,
		IsOnlyAllowOfficialGame: getInt(params, "is_only_allow_official_game"),
	}
	err = UserService.ChangeUserOnlyAllowOfficialGameStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改用户邀请奖励封禁状态
// @Param	request body services.ChangeUserMockStatusRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeUserMockStatus [post]
func (c *UserController) ChangeUserMockStatus() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	userId := getInt(params, "user_id")
	UserService := new(services.UserService)
	request := services.ChangeUserMockStatusRequest{
		UserId: userId,
		IsMock: getInt(params, "is_mock"),
	}
	err = UserService.ChangeUserMockStatus(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修改用户密码
// @Param	request body services.ChangeUserPasswordRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /changeUserPassword [post]
func (c *UserController) ChangeUserPassword() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	UserService := new(services.UserService)
	request := services.ChangeUserPasswordRequest{
		UserId:           getInt(params, "user_id"),
		Password:         getString(params, "password"),
		WithdrawPassword: getString(params, "withdraw_password"),
	}
	err = UserService.ChangeUserPassword(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type UserBetStatisticsListResponse struct {
	Code                  int    `json:"code" example:"200"`    // 状态码
	Msg                   string `json:"msg" example:"success"` // 提示信息
	UserBetStatisticsList `json:"data"`
}

type UserBetStatisticsList struct {
	List        []models.UserBetStatistic `json:"list"`
	Total       int64                     `json:"total"`
	CurrentPage int                       `json:"current_page"`
}

// @Summary	用户打码统计
// @Param user_id query int true "用户ID"
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Success 200 {object} controllers.UserBetStatisticsListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /userBetStatisticList [get]
func (c *UserController) UserBetStatisticList() {
	request := services.UserBetStatisticListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	userService := new(services.UserService)
	data, err := userService.UserBetStatisticList(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type UserRechargeOrderListResponse struct {
	Code                  int    `json:"code" example:"200"`    // 状态码
	Msg                   string `json:"msg" example:"success"` // 提示信息
	UserRechargeOrderList `json:"data"`
}

type UserRechargeOrderList struct {
	List           []models.RechargeOrder   `json:"list"`
	Total          int64                    `json:"total"`
	Packages       []models.Package         `json:"packages"`
	Status         []map[string]interface{} `json:"status"`
	PaymentChannel []models.PaymentChannel  `json:"payment_channel"`
	Payment        []models.Payment         `json:"payment"`
	CurrentPage    int                      `json:"current_page"`
}

// @Summary	用户充值订单记录
// @Param user_id query int true "用户ID"
// @Param need_reload query int true "是否刷新数据"
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Success 200 {object} controllers.UserRechargeOrderListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /userRechargeOrderList [get]
func (c *UserController) UserRechargeOrderList() {
	request := services.RechargeOrderListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		fmt.Println(err.Error())
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	rechargeOrderService := new(services.RechargeOrderService)
	needReload := c.GetQueryInt("need_reload", 0)
	lang := c.Lang
	data, err := rechargeOrderService.RechargeOrderList(request, needReload, lang, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type UserWithdrawOrderListResponse struct {
	Code                  int    `json:"code" example:"200"`    // 状态码
	Msg                   string `json:"msg" example:"success"` // 提示信息
	UserWithdrawOrderList `json:"data"`
}

type UserWithdrawOrderList struct {
	List           []models.WithdrawOrder   `json:"list"`
	Total          int64                    `json:"total"`
	Packages       []models.Package         `json:"packages"`
	Status         []map[string]interface{} `json:"status"`
	PaymentChannel []models.PaymentChannel  `json:"payment_channel"`
	Payment        []models.Payment         `json:"payment"`
	CurrentPage    int                      `json:"current_page"`
}

// @Summary	用户提现订单记录
// @Param user_id query int true "用户ID"
// @Param need_reload query int true "是否刷新数据"
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Success 200 {object} controllers.UserWithdrawOrderListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /userWithdrawOrderList [get]
func (c *UserController) UserWithdrawOrderList() {
	request := services.WithdrawOrderListRequest{}
	err := c.ParseForm(&request)
	if err != nil {
		fmt.Println(err.Error())
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request.RawQuery = c.Ctx.Request.URL.Query()
	withdrawOrderService := new(services.WithdrawOrderService)
	needReload := c.GetQueryInt("need_reload", 0)
	lang := c.Lang
	data, err := withdrawOrderService.WithdrawOrderList(request, needReload, lang, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type UserSameIpLoginLogListResponse struct {
	Code                   int    `json:"code" example:"200"`    // 状态码
	Msg                    string `json:"msg" example:"success"` // 提示信息
	UserSameIpLoginLogList `json:"data"`
}

type UserSameIpLoginLogList struct {
	List        []models.LoginLog `json:"list"`
	Total       int64             `json:"total"`
	CurrentPage int               `json:"current_page"`
}

// @Summary	用户同IP登录记录
// @Param user_id query int true "用户ID"
// @Param begin_time query int true "起始时间"
// @Param end_time query int true "结束时间"
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Success 200 {object} controllers.UserSameIpLoginLogListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /userSameIpLoginLogList [get]
func (c *UserController) UserSameIpLoginLogList() {
	request := services.UserSameIpLoginLogRequest{
		UserId:    c.GetQueryInt("user_id", 0),
		BeginTime: c.GetQueryInt64("begin_time", 0),
		EndTime:   c.GetQueryInt64("end_time", 0),
		Page:      c.GetQueryInt("page", 1),
		PageSize:  c.GetQueryInt("page_size", 20),
	}
	userService := new(services.UserService)
	data, err := userService.UserSameIpLoginLogList(request, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type UserRetentionLoginListResponse struct {
	Code                   int    `json:"code" example:"200"`    // 状态码
	Msg                    string `json:"msg" example:"success"` // 提示信息
	UserRetentionLoginList `json:"data"`
}

type UserRetentionLoginList struct {
	List        []models.UserRetentionDaily `json:"list"`
	Packages    []models.Package            `json:"packages"`
	Total       int64                       `json:"total"`
	CurrentPage int                         `json:"current_page"`
}

// @Summary	用户留存记录
// @Param package_id query int true "包ID"
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Success 200 {object} controllers.UserRetentionLoginListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /userLoginRetentionList [get]
func (c *UserController) UserRetentionLoginList() {
	request := services.UserRetentionRequest{
		PackageId: c.GetQueryInt("package_id", -1),
		Page:      c.GetQueryInt("page", 1),
		PageSize:  c.GetQueryInt("page_size", 20),
	}
	userService := new(services.UserService)
	needReload := c.GetQueryInt("need_reload", 0)
	retentionType := int(config.RetentionTypeLogin)
	data, err := userService.UserRetentionList(request, retentionType, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

type UserRetentionRechargeListResponse struct {
	Code                      int    `json:"code" example:"200"`    // 状态码
	Msg                       string `json:"msg" example:"success"` // 提示信息
	UserRetentionRechargeList `json:"data"`
}

type UserRetentionRechargeList struct {
	List        []models.UserRetentionDaily `json:"list"`
	Packages    []models.Package            `json:"packages"`
	Total       int64                       `json:"total"`
	CurrentPage int                         `json:"current_page"`
}

// @Summary	用户留存记录
// @Param package_id query int true "包ID"
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Success 200 {object} controllers.UserRetentionRechargeListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /userRechargeRetentionList [get]
func (c *UserController) UserRetentionRechargeList() {
	request := services.UserRetentionRequest{
		PackageId: c.GetQueryInt("package_id", -1),
		Page:      c.GetQueryInt("page", 1),
		PageSize:  c.GetQueryInt("page_size", 20),
	}
	userService := new(services.UserService)
	needReload := c.GetQueryInt("need_reload", 0)
	retentionType := int(config.RetentionTypeRecharge)
	data, err := userService.UserRetentionList(request, retentionType, needReload, c.UserId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	修复用户登录留存
// @Param	request body services.FixUserLoginRetentionRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /fixUserLoginRetention [post]
func (c *UserController) FixUserLoginRetention() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	UserService := new(services.UserService)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request := services.FixUserLoginRetentionRequest{
		StartDate:  getString(params, "start_date"),
		EndDate:    getString(params, "end_date"),
		PackageIds: getIntSlice(params, "package_ids"),
	}
	err = UserService.CalculatePackagesUserLoginRetention(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	修复用户充值留存
// @Param	request body services.FixUserRechargeRetentionRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /fixUserRechargeRetention [post]
func (c *UserController) FixUserRechargeRetention() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	UserService := new(services.UserService)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	request := services.FixUserRechargeRetentionRequest{
		StartDate:  getString(params, "start_date"),
		EndDate:    getString(params, "end_date"),
		PackageIds: getIntSlice(params, "package_ids"),
	}
	err = UserService.CalculatePackagesUserRechargeRetention(request)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
