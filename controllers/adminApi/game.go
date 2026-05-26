package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
)

type GameController struct {
	BaseController
}

type GameListDTO struct {
	List        []models.Game `json:"list"`         // 列表数据
	Total       int64         `json:"total"`        // 总条数
	CurrentPage int           `json:"current_page"` // 当前页
}
type GameListResponse struct {
	Code int         `json:"code" example:"200"`    // 状态码
	Msg  string      `json:"msg" example:"success"` // 提示信息
	Data GameListDTO `json:"data"`
}

// @Summary	游戏列表
// @Param id query int false "游戏ID"
// @Param name query string false "游戏名称"
// @Param status query int false "游戏状态"
// @Param game_type_id query int false "游戏类型ID"
// @Param platform_id query int false "平台ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.GameListResponse
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /gameList [get]
func (c *GameController) GameList() {
	page, pageSize := c.GetPagination()

	params := services.GameListRequest{
		Page:         page,
		PageSize:     pageSize,
		Name:         c.GetString("name", ""),
		Code:         c.GetString("code", ""),
		Recommend:    c.GetQueryInt("recommend", -1),
		Maintain:     c.GetQueryInt("maintain", -1),
		GameTypeName: c.GetString("game_type_name", ""),
		PlatformName: c.GetString("platform_name", ""),
		Status:       c.GetQueryInt("status", -1),
		Tag:          c.GetString("tag", ""),
		CodeRule:     c.GetQueryInt("code_rule", -1),
		NeedReload:   c.GetQueryInt("need_reload", 0),
		Event:        c.GetString("event", ""),
		PackageId:    c.GetQueryInt("package_id", -1),
	}

	if params.PackageId > 0 {
		if !c.IsMyPackageId(params.PackageId, int(c.UserId)) {
			c.JSONError(500, "ShuJuYiChang")
		}
	} else {
		params.PackageIds = c.GetPackageIds(int(c.UserId))
	}

	gameService := new(services.GameService)
	data, err := gameService.GameList(params, int(c.UserId))
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	添加游戏
// @Param	request body services.AddGameRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /addGame [post]
func (c *GameController) AddGame() {
	var params services.EditGameRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	params.PackageId = 1
	if !c.IsMyPackageId(params.PackageId, int(c.UserId)) {
		c.JSONError(500, "ShuJuYiChang")
	}

	params.PackageIds = c.GetPackageIds(int(c.UserId))

	err = (services.GameService{}).AddGame(params)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	编辑游戏
// @Param	request body services.EditGameRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editGame [post]
func (c *GameController) EditGame() {
	var params services.EditGameRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	params.PackageId = 1
	params.PackageIds = c.GetPackageIds(int(c.UserId))

	err = (services.GameService{}).EditGame(params)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type EditGameAttrRequest struct {
	Id    int    `json:"id"`
	Field string `json:"field"`
}

// @Summary	编辑游戏属性
// @Param	request body controllers.EditGameAttrRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /editGameAttr [post]
func (c *GameController) EditGameAttr() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	gameService := new(services.GameService)
	id := getInt(params, "id")
	field := getString(params, "field")
	v := getInt(params, "value")

	err = gameService.EditGameAttr(id, field, v)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

type DeleteGameRequest struct {
	Id string `json:"id" example:"123456"` // 游戏id
}

// @Summary	删除游戏
// @Param	request body controllers.DeleteGameRequest true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /deleteGame [post]
func (c *GameController) DeleteGame() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	id := getInt(params, "id")
	gameService := new(services.GameService)
	err = gameService.DeleteGame(id)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}

// @Summary	获取商户列表
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Param title query string false "名称"
// @Param status query int false "状态"
// @Success 200 {object} services.AccountListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /merchantList [get]
func (c *GameController) MerchantList() {
	page, pageSize := c.GetPagination()

	status := c.GetQueryInt("status", -1)
	title := c.GetString("title", "")

	params := services.MerchantRequestParams{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Title:    title,
	}

	result := (&services.MerchantService{}).GetList(params, int(c.UserId))
	c.JSONSuccess(result)
}

// @Summary	修改商户信息
// @Param	request body services.EditMerchantRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /editMerchant [post]
func (c *GameController) EditMerchant() {
	var params services.EditMerchantRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	roleErr := (&services.MerchantService{}).EditMerchant(params)
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary 创建商户信息
// @Param	request body services.EditMerchantRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /addMerchant [post]
func (c *GameController) AddMerchant() {
	var params services.EditMerchantRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	roleErr := (&services.MerchantService{}).AddMerchant(params)
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	删除商户
// @Param	request body services.DeleteMerchantRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /delMerchant [post]
func (c *GameController) DelMerchant() {
	var params services.DeleteMerchantRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	roleErr := (&services.MerchantService{}).DeleteMerchant(params)
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	拉取商户游戏
// @Param	request body services.DeleteMerchantRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /pullMerchantGame [post]
func (c *GameController) PullMerchantGame() {
	var params services.PullMerchantGameRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	if !c.IsMyPackageId(params.PackageId, int(c.UserId)) {
		c.JSONError(500, "ShuJuYiChang")
	}
	roleErr := (&services.MerchantService{}).PullMerchantGame(params)
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	Excel导入游戏
// @Param	file formData file true "Excel文件"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /importExcel [post]
func (c *GameController) ImportExcel() {
	_, fileHeader, err := c.GetFile("file")
	if err != nil {
		c.JSONError(500, "未找到上传的文件")
		return
	}

	fileExt := ""
	if len(fileHeader.Filename) > 5 {
		fileExt = fileHeader.Filename[len(fileHeader.Filename)-5:]
	}
	if fileExt != ".xlsx" && fileExt != ".xls" {
		c.JSONError(500, "只支持Excel文件格式 (.xlsx, .xls)")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSONError(500, "文件打开失败")
		return
	}
	defer file.Close()

	packageId := 1
	packageIds := c.GetPackageIds(int(c.UserId))

	stats, err := (services.GameService{}).ImportGamesFromExcel(file, packageId, packageIds)
	if err != nil {
		c.JSONError(500, err.Error())
		return
	}

	if stats.SuccessCount == 0 {
		c.JSONError(500, "Excel中没有有效数据，请检查文件格式和内容")
		return
	}

	c.JSONSuccess(stats)
}

// @Summary	创建用户会话
// @Param	request body map[string]interface{} true "json请求参数"
// @Success 200 {object} controllers.JSONResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /createUserSession [post]
func (c *GameController) CreateUserSession() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	userId := getString(params, "user_id")
	userName := getString(params, "user_name")
	rtp := getInt(params, "rtp")

	data, err := (services.GameService{}).CreateUserSession(userId, userName, rtp)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	获取游戏URL
// @Param	request body map[string]interface{} true "json请求参数"
// @Success 200 {object} controllers.JSONResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /getGameUrl [post]
func (c *GameController) GetGameUrl() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	gameCode := getString(params, "game_code")
	language := getString(params, "language")
	gameType := getString(params, "type")
	userId := getString(params, "user_id")
	userToken := getString(params, "user_token")

	data, err := (services.GameService{}).GetGameUrl(gameCode, language, gameType, userId, userToken)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	设置用户RTP
// @Param	request body map[string]interface{} true "json请求参数"
// @Success 200 {object} controllers.JSONResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /setUserRtp [post]
func (c *GameController) SetUserRtp() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	userId := getString(params, "user_id")
	rtp := getInt(params, "rtp")

	data, err := (services.GameService{}).SetUserRtp(userId, rtp)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	获取所有游戏列表
// @Param	package_id query int false "分包ID"
// @Success 200 {object} controllers.JSONResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /getAllGames [get]
func (c *GameController) GetAllGames() {
	packageId := 1
	data, err := (services.GameService{}).GetAllGamesList(packageId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	获取客服地址
// @Success 200 {object} controllers.JSONResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /getCustomer [get]
func (c *GameController) GetCustomer() {
	packageIds := c.GetPackageIds(int(c.UserId))
	packageId := 1
	if len(packageIds) > 0 {
		packageId = packageIds[0]
	}

	url, err := (services.GameService{}).GetCustomerService(packageId)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(url)
	}
}

// @Summary	设置客服地址
// @Param	request body map[string]interface{} true "json请求参数"
// @Success 200 {object} controllers.JSONResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /customer [post]
func (c *GameController) SetCustomer() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
		return
	}

	packageIds := c.GetPackageIds(int(c.UserId))
	packageId := 1
	if len(packageIds) > 0 {
		packageId = packageIds[0]
	}

	url := getString(params, "url")
	if url == "" {
		c.JSONError(400, "KeFuDiZhiBuNengWeiKong")
		return
	}

	err = (services.GameService{}).SetCustomerService(packageId, url)
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(nil)
	}
}
