package controllers

import (
	"api/config"
	"api/middleware"
	"api/models"
	"api/services"
	"api/utils"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/beego/i18n"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// SkipAuthPaths 定义不需要认证的路径前缀
var SkipAuthPaths = []string{
	"/api/login/login",
	"/api/active/testCreateActive",
	"/api/serves/refreshTable",
	"/api/serves/getMultiTaskProcess",
	"/api/game/getAllGames",
	"/api/game/setUserRtp",
	"/api/game/getGameUrl",
	"/api/game/createUserSession",
}

// JSONResponse JSON响应结构
type JSONResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// JSONResponseTpl JSON响应结构模板
type JSONResponseTpl struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// BaseController 基础控制器
type BaseController struct {
	beego.Controller

	// 自定义字段
	UserId   int64
	UserName string
	IsAdmin  bool
	IsLogin  bool

	ModuleName     string // 模块名
	ControllerName string // 控制器名
	ActionName     string // 方法名

	Lang string
}

// Prepare 准备方法
func (c *BaseController) Prepare() {
	c.setI18nLanguage()
	// 获取路由信息
	c.ModuleName, c.ControllerName, c.ActionName = c.ParseRouteParams()

	if c.ShouldSkipAuth() {
		// 记录请求日志
		if beego.BConfig.RunMode == "dev" {
			c.logRequest()
		}

		return
	}

	// 检查登录状态
	middleware.CheckToken(c.Ctx)

	// 设置通用模板数据
	c.setCommonData()

	// 检测是否有权限
	c.checkAuth()

	// 记录请求日志
	c.logRequest()
}

// 设置国际化语言
func (c *BaseController) setI18nLanguage() {
	paramLang := c.GetString("lang")
	if paramLang == "" {
		al := c.Ctx.Input.Header("Accept-Language")
		if len(al) >= 5 {
			paramLang = al[:5]
		}
	}

	lang := ""
	paramLang = utils.ConvertLangCode(paramLang)
	if !i18n.IsExist(paramLang) {
		lang = "en-US"
	} else {
		lang = paramLang
	}
	c.Lang = lang
}

func (c *BaseController) ShouldSkipAuth() bool {
	// 1. 通过路径判断
	if c.skipByPath() {
		return true
	}

	// 3. 通过HTTP方法判断（如OPTIONS预检请求）
	if c.Ctx.Request.Method == "OPTIONS" {
		return true
	}

	return false
}

// skipByPath 通过路径判断是否跳过认证
func (c *BaseController) skipByPath() bool {
	currentPath := c.Ctx.Request.URL.Path

	for _, path := range SkipAuthPaths {
		if path == currentPath {
			return true
		}
		// 支持通配符 * 结尾
		if strings.HasSuffix(path, "/*") {
			prefix := strings.TrimSuffix(path, "/*")
			if strings.HasPrefix(currentPath, prefix+"/") || currentPath == prefix {
				return true
			}
		}
		// 支持前缀匹配
		if strings.HasPrefix(path, "/") && strings.HasSuffix(path, "/") {
			if strings.HasPrefix(currentPath, path) {
				return true
			}
		}
		// 简单前缀匹配
		if strings.HasPrefix(currentPath, path) {
			return true
		}
	}

	return false
}

func (c *BaseController) ParseRouteParams() (module, controller, action string) {
	// 方法1：从 URL 路径解析
	urlPath := c.Ctx.Request.URL.Path
	parts := strings.Split(strings.Trim(urlPath, "/"), "/")

	if len(parts) >= 3 {
		module = parts[0]     // api
		controller = parts[1] // system
		action = parts[2]     // roleList
	}

	return module, controller, action
}

// checkAuth 检查认证
func (c *BaseController) checkAuth() {
	groupId := (&services.AccountService{}).GetGroupId(int(c.UserId))
	path := fmt.Sprintf("%v/%v/%v", c.ModuleName, c.ControllerName, c.ActionName)
	menuId := (&services.AuthRuleService{}).GetMenuId(path)
	if !(&services.Authservice{}).Check(menuId, groupId) {
		c.JSONError(500, "MeiYouQuanXian")
	}
}

// setCommonData 设置通用数据
func (c *BaseController) setCommonData() {
	data := c.Ctx.Input.Data()

	uid, _ := data["userId"].(int64)
	c.UserId = uid
	c.UserName = data["userName"].(string)
	c.Data["UserId"] = data["userId"]
	c.Data["UserName"] = data["userName"]

	// // CSRF Token
	// c.Data["csrf_token"] = c.XSRFToken()
}

// logRequest 记录请求日志
func (c *BaseController) logRequest() {
	logs.Info("[%s][%s][%s.%s] UserId:%d IP:%s Body:%s",
		c.Ctx.Request.Method,
		c.Ctx.Request.URL.Path,
		c.ControllerName,
		c.ActionName,
		c.UserId,
		c.Ctx.Input.IP(),
		c.Ctx.Input.RequestBody,
	)

	if c.Ctx.Request.Method == "POST" {
		body := string(c.Ctx.Input.RequestBody)
		if c.Ctx.Request.URL.Path == "/api/login/login" {
			re := regexp.MustCompile(`("password"\s*:\s*")[^"]+(")`)
			body = re.ReplaceAllString(body, `${1}******${2}`)
		}

		go func() {
			var params map[string]interface{}
			err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
			var playerId, packageId int
			var roleId int64
			if err == nil {
				playerId = getInt(params, "user_id")
				packageId = getInt(params, "package_id")
			}

			if playerId > 0 {
				user, err := (&services.UserService{}).GetUserById(playerId)
				if err == nil {
					roleId = user.RoleId
					packageId = user.PackageId
				}
			}
			(&services.UserService{}).AddAdminLog(models.AdminLog{
				AdminId:    int(c.UserId),
				Path:       strings.TrimLeft(c.Ctx.Request.URL.Path, "/"),
				Controller: c.ControllerName,
				Action:     c.ActionName,
				Body:       body,
				Status:     1,
				LogType:    int(config.AdminLogTypeCommon),
				Ip:         c.Ctx.Input.IP(),
				CTime:      time.Now().Unix(),
				PlayerId:   playerId,
				RoleId:     roleId,
				PackageId:  packageId,
			})
		}()
	}
}

// ========== 响应方法 ==========

// JSONSuccess 成功响应
func (c *BaseController) JSONSuccess(data interface{}) {
	resp := JSONResponse{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
	c.Data["json"] = resp
	c.ServeJSON()
	c.StopRun()
}

// JSONError 错误响应
func (c *BaseController) JSONError(code int, msg string) {
	resp := JSONResponse{
		Code: code,
		Msg:  c.Tr(msg),
	}
	c.Data["json"] = resp
	c.ServeJSON()
	c.StopRun()
}

// JSONErrorWithData 带数据的错误响应
func (c *BaseController) JSONErrorWithData(code int, msg string, data interface{}) {
	resp := JSONResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	c.Data["json"] = resp
	c.ServeJSON()
	c.StopRun()
}

// ========== 工具方法 ==========
// GetQueryInt 获取查询参数(整数)
func (c *BaseController) GetQueryInt(key string, defaultValue int) int {
	str := c.GetString(key)
	if str == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return val
}

// GetQueryInt64 获取查询参数(int64)
func (c *BaseController) GetQueryInt64(key string, defaultValue int64) int64 {
	str := c.GetString(key)
	if str == "" {
		return defaultValue
	}
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

// getInt 获取int64 默认为0
func getInt64(params map[string]interface{}, key string) int64 {
	v, ok := params[key]
	if !ok || v == nil {
		return 0
	}

	switch val := v.(type) {
	case int:
		return int64(val)
	case int8:
		return int64(val)
	case int16:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val

	case float32:
		return int64(val)
	case float64:
		// json.Unmarshal 默认是 float64
		return int64(val)

	case string:
		if val == "" {
			return 0
		}
		if i, err := strconv.Atoi(val); err == nil {
			return int64(i)
		}
	}

	return 0
}

// getInt 获取int 默认为0
func getInt(params map[string]interface{}, key string) int {
	v, ok := params[key]
	if !ok || v == nil {
		return 0
	}

	switch val := v.(type) {
	case int:
		return val
	case int8:
		return int(val)
	case int16:
		return int(val)
	case int32:
		return int(val)
	case int64:
		return int(val)

	case float32:
		return int(val)
	case float64:
		// json.Unmarshal 默认是 float64
		return int(val)

	case string:
		if val == "" {
			return 0
		}
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}

	return 0
}

// getFloat 获取float64 默认为0
func getFloat(params map[string]interface{}, key string) float64 {
	v, ok := params[key]
	if !ok || v == nil {
		return 0
	}

	switch val := v.(type) {
	case float32:
		return float64(val)
	case float64:
		// json.Unmarshal 默认是 float64
		return val

	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)

	case string:
		if val == "" {
			return 0
		}
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}

	return 0
}

// getString 获取string 默认为""
func getString(params map[string]interface{}, key string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getJSONString 获取JSON字符串
func getJSONString(params map[string]interface{}, key string) string {
	if v, ok := params[key]; ok {
		b, _ := json.Marshal(v)
		return string(b)
	}
	return ""
}

// getIntSlice 获取int数组
func getIntSlice(params map[string]interface{}, key string) []int {
	if v, ok := params[key]; ok {
		arr, ok := v.([]interface{})
		if !ok {
			return nil
		}

		var result []int
		for _, item := range arr {
			switch val := item.(type) {
			case float64: // JSON数字默认是float64
				result = append(result, int(val))
			case int:
				result = append(result, val)
			}
		}
		return result
	}
	return nil
}

// GetPagination 获取分页参数
func (c *BaseController) GetPagination() (page, pageSize int) {
	page = c.GetQueryInt("page", 1)
	pageSize = c.GetQueryInt("page_size", 20)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return page, pageSize
}

// GetDeviceInfo 获取设备信息
func (c *BaseController) GetDeviceInfo() map[string]interface{} {
	userAgent := c.Ctx.Input.Header("user-agent")
	var deviceInfo map[string]interface{}
	deviceInfo = map[string]interface{}{
		"user_agent": userAgent,
	}
	return deviceInfo
}

// GetClientIP 获取客户端IP
func (c *BaseController) GetClientIP() string {
	// 尝试从 Header 获取真实 IP
	realIP := c.Ctx.Input.Header("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	forwardedFor := c.Ctx.Input.Header("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	return c.Ctx.Input.IP()
}

func (c *BaseController) Tr(key string, args ...interface{}) string {
	// 调用 github.com/beego/i18n 包提供的 Tr 方法
	return i18n.Tr(c.Lang, key, args...)
}

// 接口限流（N 秒内最多 M 次）
// userId: 用户ID
// action: 方法
// max: 最多次数
// window: 时间窗口
func (c *BaseController) RateLimit(ctx context.Context, userId int, action string, max int, window time.Duration) bool {
	key := fmt.Sprintf(config.RedisKeyName.RateLimit+"%s:%d", action, userId)

	RedisCache := utils.GetRedisClient()
	cnt, err := RedisCache.Incr(ctx, key).Result()
	if err != nil {
		return true
	}
	// 避免删除key失败
	RedisCache.Expire(ctx, key, window)
	return cnt <= int64(max)
}

// 是否进行中
func (c *BaseController) IsActiveRuning(uid, activeId int) bool {
	redis := utils.GetRedisClient()
	ctx := context.Background()
	cacheKey := fmt.Sprintf("is_active_runing_%v_%v", activeId, uid)
	result, err := redis.SetNX(ctx, cacheKey, "1", time.Second*2).Result()
	if err != nil {
		logs.Error("redis设置活动IsActiveRuning失败:%v", err.Error())
		return false
	}
	return result
}

// 获取查询时间
func (c *BaseController) GetSearchTime() (int64, int64) {
	createTimeArr := c.Ctx.Request.URL.Query()["createTime[]"]
	var beginTime, endTime int64
	if len(createTimeArr) == 2 {
		// 转开始时间戳
		t1, _ := time.ParseInLocation("2006-01-02 15:04:05", createTimeArr[0], time.Local)
		beginTime = t1.Unix()

		// 转结束时间戳
		t2, _ := time.ParseInLocation("2006-01-02 15:04:05", createTimeArr[1], time.Local)
		endTime = t2.Unix()
	}
	return beginTime, endTime
}

// 是否是自己的分包
func (c *BaseController) IsMyPackageId(packageId int, userId int) bool {
	if packageId == 0 {
		return true
	}
	packageList := (&services.PackageService{}).GetMyAllPackageList(userId)
	if len(packageList) == 0 {
		return false
	}
	for _, packageInfo := range packageList {
		if packageInfo.Id == packageId {
			return true
		}
	}
	return false
}

// 获取自己的分包id
func (c *BaseController) GetPackageIds(userId int, packageId ...int) []int {
	if len(packageId) > 0 && packageId[0] > 0 {
		if c.IsMyPackageId(packageId[0], userId) {
			return []int{packageId[0]}
		}
	} else {
		packageList := (&services.PackageService{}).GetMyAllPackageList(userId)
		if len(packageList) == 0 {
			return []int{-1} // 不存在数据
		}
		list := make([]int, 0)
		for _, packageInfo := range packageList {
			list = append(list, packageInfo.Id)
		}

		return list
	}
	return []int{-1} // 不存在数据
}

// ========== 文件操作 ==========

// SaveUploadFile 保存上传文件
// func (c *BaseController) SaveUploadFile(formName, savePath string) (string, error) {
// 	// file, header, err := c.GetFile(formName)
// 	// if err != nil {
// 	// 	return "", err
// 	// }
// 	// defer file.Close()

// 	// // 创建目录
// 	// // ...

// 	// // 保存文件
// 	// filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
// 	// savePath = path.Join(savePath, filename)

// 	// err = c.SaveToFile(formName, savePath)
// 	// if err != nil {
// 	// 	return "", err
// 	// }

// 	// return savePath, nil
// }

// Finish 结束请求
func (c *BaseController) Finish() {
	// 可以在这里进行资源清理
}
