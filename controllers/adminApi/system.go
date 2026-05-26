package controllers

import (
	"api/services"
	"encoding/json"
)

type SystemController struct {
	BaseController
}

// @Summary	获取菜单
// @Success 200 {object} services.Menu
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /menu [get]
func (c *SystemController) Menu() {
	menuList := (&services.AuthRuleService{}).GetMenusTreeByCache(int(c.UserId))
	c.JSONSuccess(menuList)
}

// @Summary	获取菜单规则
// @Success 200 {object} services.Menu
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /auth [get]
func (c *SystemController) Auth() {
	menuList := (&services.AuthRuleService{}).GetAuthRuleList()
	c.JSONSuccess(menuList)
}

// @Summary	创建菜单规则
// @Param	request body services.AuthAddRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /authadd [post]
func (c *SystemController) Authadd() {
	var params services.AuthAddRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	authErr := (&services.AuthRuleService{}).Authadd(params)
	if authErr != nil {
		c.JSONError(500, authErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	编辑菜单规则
// @Param	request body services.AuthAddRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /authedit [post]
func (c *SystemController) Authedit() {
	var params services.AuthAddRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	authErr := (&services.AuthRuleService{}).Authedit(params)
	if authErr != nil {
		c.JSONError(500, authErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	删除菜单规则
// @Param	request body services.AuthDelRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /authdel [post]
func (c *SystemController) Authdel() {
	var params services.AuthDelRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	authErr := (&services.AuthRuleService{}).Authdel(params)
	if authErr != nil {
		c.JSONError(500, authErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	角色列表
// @Success 200 {object} services.AuthGroupListResult
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /roleList [get]
func (c *SystemController) RoleList() {
	menuList := (&services.AuthGroupService{}).RoleList(int(c.UserId))
	c.JSONSuccess(menuList)
}

// @Summary	创建角色
// @Param	request body services.AddRoleRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /roleadd [post]
func (c *SystemController) Roleadd() {
	var params services.AddRoleRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	roleErr := (&services.AuthGroupService{}).AddRole(params, c.UserId)
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	修改角色
// @Param	request body services.AddRoleRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /roleedit [post]
func (c *SystemController) Roleedit() {
	var params services.AddRoleRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	roleErr := (&services.AuthGroupService{}).EditRole(params, c.UserId)
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

type DeleteRequestParams struct {
	Id int `json:"id"`
}

// @Summary	删除角色
// @Param	request body controllers.DeleteRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /roledel [post]
func (c *SystemController) Roledel() {
	var params DeleteRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	roleErr := (&services.AuthGroupService{}).DelRole(params.Id, int(c.UserId))
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	获取权限菜单
// @Param id query int false "ID"
// @Success 200 {object} services.RolerulesResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /rolerules [get]
func (c *SystemController) Rolerules() {
	id := c.GetQueryInt("id", 0)

	data := (&services.AuthGroupService{}).Rolerules(id)
	c.JSONSuccess(data)
}

// @Summary	授权
// @Param	request body services.RoleruleseditParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /rolerulesedit [post]
func (c *SystemController) Rolerulesedit() {
	var params services.RoleruleseditParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	roleErr := (&services.AuthGroupService{}).Rolerulesedit(params, int(c.UserId))
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	用户管理
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Param username query string false "用户名"
// @Param status query int false "状态"
// @Success 200 {object} services.AccountListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /adminList [get]
func (c *SystemController) AdminList() {
	page, pageSize := c.GetPagination()

	username := c.GetString("username", "")
	status := c.GetQueryInt("status", -1)

	params := services.UserListRequestParams{
		Page:     page,
		PageSize: pageSize,
		Username: username,
		Status:   status,
	}

	result := (&services.AccountService{}).GetUserList(params, int(c.UserId))
	c.JSONSuccess(result)
}

type AdmintatusRequestParams struct {
	Id int `json:"id"`
}

// @Summary	管理员状态修改
// @Param	request body controllers.AdmintatusRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /admintatus [post]
func (c *SystemController) Admintatus() {
	var params AdmintatusRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	adminErr := (&services.AccountService{}).EditAccountStatus(params.Id, int(c.UserId))
	if adminErr != nil {
		c.JSONError(500, adminErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	管理员添加
// @Param	request body services.AccountRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /adminAdd [post]
func (c *SystemController) AdminAdd() {
	var params services.AccountRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	adminErr := (&services.AccountService{}).AddAccount(params, int(c.UserId))
	if adminErr != nil {
		c.JSONError(500, adminErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	管理员删除
// @Param	request body controllers.DeleteRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /adminDel [post]
func (c *SystemController) AdminDel() {
	var params DeleteRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	adminErr := (&services.AccountService{}).DelAccount(params.Id, int(c.UserId))
	if adminErr != nil {
		c.JSONError(500, adminErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	管理员修改
// @Param	request body services.AccountRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /adminEdit [post]
func (c *SystemController) AdminEdit() {
	var params services.AccountRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	adminErr := (&services.AccountService{}).EditAccount(params, int(c.UserId))
	if adminErr != nil {
		c.JSONError(500, adminErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	管理员解绑google
// @Param	request body controllers.DeleteRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /unbindgoogle [post]
func (c *SystemController) Unbindgoogle() {
	var params DeleteRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	adminErr := (&services.AccountService{}).Unbindgoogle(params.Id, int(c.UserId))
	if adminErr != nil {
		c.JSONError(500, adminErr.Error())
	}
	c.JSONSuccess(nil)
}

type RequestGoogleKeyParams struct {
	Uid      int    `json:"uid"`      // 用户id
	Username string `json:"username"` // 用户
}

// @Summary	获取google key信息
// @Param	request body controllers.RequestGoogleKeyParams true "json请求参数"
// @Success 200 {object} services.GooleSecretResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /googleKey [get]
func (c *SystemController) GoogleKey() {
	uid := c.GetQueryInt("uid", -1)
	username := c.GetString("username", "")
	err, result := (&services.AccountService{}).CreateGoogleSecrect(uid, username)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(result)
}

type RequestBindGoogleKeyParams struct {
	Uid  int    `json:"uid"`  // 用户id
	Code string `json:"Code"` // 验证码
}

// @Summary	绑定google key信息
// @Param	request body controllers.RequestBindGoogleKeyParams true "json请求参数"
// @Success 200 {object} services.GooleSecretResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /bindGooleKey [post]
func (c *SystemController) BindGooleKey() {
	var params RequestBindGoogleKeyParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	bindErr := (&services.AccountService{}).Bindgoogle(params.Uid, params.Code)
	if bindErr != nil {
		c.JSONError(500, bindErr.Error())
	}
	c.JSONSuccess(nil)
}
