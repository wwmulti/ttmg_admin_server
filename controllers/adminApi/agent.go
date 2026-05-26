package controllers

import (
	"api/services"
	"encoding/json"
)

type AgentController struct {
	BaseController
}

// @Summary	平台管理
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Param title query string false "名称"
// @Param domain query int false "域名"
// @Success 200 {object} services.PackageListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /packageList [get]
func (c *AgentController) PackageList() {
	page, pageSize := c.GetPagination()

	title := c.GetString("title", "")
	domain := c.GetString("domain", "")

	params := services.PackageRequestParams{
		Page:     page,
		PageSize: pageSize,
		Title:    title,
		Domain:   domain,
	}

	result := (&services.PackageService{}).GetList(params, int(c.UserId))
	c.JSONSuccess(result)
}

// @Summary	所有平台
// @Success 200 {object} services.AllPackageResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /allPackage [get]
func (c *AgentController) AllPackage() {
	result, err := (&services.PackageService{}).AllPackage()
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(result)
}

// @Summary	创建平台
// @Param	request body services.AddPlatformRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /addPackage [post]
func (c *AgentController) AddPackage() {
	var params services.AddPlatformRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	addErr := (&services.PackageService{}).AddPackage(params)
	if addErr != nil {
		c.JSONError(500, addErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	编辑平台
// @Param	request body services.AddPlatformRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /editPackage [post]
func (c *AgentController) EditPackage() {
	var params services.AddPlatformRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	editErr := (&services.PackageService{}).EditPackage(params)
	if editErr != nil {
		c.JSONError(500, editErr.Error())
	}
	c.JSONSuccess(nil)
}

type DeletePackageRequestParams struct {
	Ids string `json:"ids"`
}

// @Summary	删除平台
// @Param	request body controllers.DeletePackageRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /delPackage [post]
func (c *AgentController) DelPackage() {
	var params DeletePackageRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	roleErr := (&services.PackageService{}).DelPackage(params.Ids)
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	获取平台配置
// @Param package_id query int false "package_id"
// @Success 200 {object} services.PackageConfigListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /configList [get]
func (c *AgentController) ConfigList() {
	id := c.GetQueryInt("package_id", 0)
	result := (&services.PackageService{}).ConfigList(id)
	c.JSONSuccess(result)
}

// @Summary	修改平台配置
// @Param	request body services.EditPackageConfigRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /editConfig [post]
func (c *AgentController) EditConfig() {
	var params services.EditPackageConfigRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	roleErr := (&services.PackageService{}).EditConfig(params)
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	平台分组
// @Param package_id query int false "package_id"
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Param status query int false "状态"
// @Param title query string false "名称"
// @Success 200 {object} services.PackageConfigListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /groupList [get]
func (c *AgentController) GroupList() {
	page, pageSize := c.GetPagination()

	packageId := c.GetQueryInt("package_id", 0)
	status := c.GetQueryInt("status", -1)
	title := c.GetString("title", "")

	params := services.PackageGroupListRequestParams{
		Page:      page,
		PageSize:  pageSize,
		PackageId: packageId,
		Status:    status,
		Title:     title,
	}

	result := (&services.PackageGroupService{}).GetList(params, int(c.UserId))
	c.JSONSuccess(result)
}

// @Summary	创建平台分组
// @Param	request body services.AddPackageGroupListRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /addGroup [post]
func (c *AgentController) AddGroup() {
	var params services.AddPackageGroupListRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	addErr := (&services.PackageGroupService{}).AddGroup(params, c.Lang)
	if addErr != nil {
		c.JSONError(500, addErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	修改平台分组
// @Param	request body services.AddPackageGroupListRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /editGroup [post]
func (c *AgentController) EditGroup() {
	var params services.AddPackageGroupListRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	addErr := (&services.PackageGroupService{}).EditGroup(params, c.Lang)
	if addErr != nil {
		c.JSONError(500, addErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	删除平台分组
// @Param	request body controllers.DeletePackageRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /delGroup [post]
func (c *AgentController) DelGroup() {
	var params DeletePackageRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	roleErr := (&services.PackageGroupService{}).DelGroup(params.Ids)
	if roleErr != nil {
		c.JSONError(500, roleErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	平台用户管理
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Param title query string false "名称"
// @Param domain query int false "域名"
// @Success 200 {object} services.UserListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /userList [get]
func (c *AgentController) UserList() {
	page, pageSize := c.GetPagination()

	packageId := c.GetQueryInt("package_id", 0)
	roleId := c.GetQueryInt("role_id", 0)
	username := c.GetString("username", "")

	params := services.UserRequestParams{
		Page:      page,
		PageSize:  pageSize,
		PackageId: packageId,
		RoleId:    roleId,
		Username:  username,
	}

	result := (&services.PackageGroupService{}).GetUserList(params)
	c.JSONSuccess(result)
}
