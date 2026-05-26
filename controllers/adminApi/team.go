package controllers

import (
	"api/services"
	"encoding/json"
)

type TeamController struct {
	BaseController
}

// @Summary	团队管理
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Param username query string false "用户名"
// @Success 200 {object} services.AccountListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /getList [get]
func (c *TeamController) GetList() {
	page, pageSize := c.GetPagination()

	title := c.GetString("title", "")
	packageId := c.GetQueryInt("package_id", -1)
	pid := c.GetQueryInt("pid", -1)
	needReload := c.GetQueryInt("need_reload", 0)

	params := services.TeamListRequestParams{
		Page:       page,
		PageSize:   pageSize,
		Title:      title,
		PackageId:  packageId,
		Pid:        pid,
		NeedReload: needReload,
	}

	result := (&services.TeamService{}).GetList(params)
	c.JSONSuccess(result)
}

// @Summary	业务员团队管理
// @Param page query int false "页数"
// @Param page_size query int false "条数"
// @Param username query string false "用户名"
// @Success 200 {object} services.AccountListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /getSalespersonList [get]
func (c *TeamController) GetSalespersonList() {
	page, pageSize := c.GetPagination()

	title := c.GetString("title", "")
	packageId := c.GetQueryInt("package_id", -1)
	pid := c.GetQueryInt("pid", -1)
	needReload := c.GetQueryInt("need_reload", 0)

	params := services.TeamListRequestParams{
		Page:       page,
		PageSize:   pageSize,
		Title:      title,
		PackageId:  packageId,
		Pid:        pid,
		NeedReload: needReload,
	}

	result := (&services.TeamService{}).GetSalespersonList(params)
	c.JSONSuccess(result)
}

// @Summary	团队添加
// @Param	request body services.TeamRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /addTeam [post]
func (c *TeamController) AddTeam() {
	var params services.TeamRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	params.Language = c.Lang
	addErr := (&services.TeamService{}).AddTeam(params)
	if addErr != nil {
		c.JSONError(500, addErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	业务员团队添加
// @Param	request body services.TeamRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /addSalespersonTeam [post]
func (c *TeamController) AddSalespersonTeam() {
	var params services.TeamRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}

	params.Language = c.Lang
	addErr := (&services.TeamService{}).AddTeam(params)
	if addErr != nil {
		c.JSONError(500, addErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	团队修改
// @Param	request body services.TeamRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /editTeam [post]
func (c *TeamController) EditTeam() {
	var params services.TeamRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	addErr := (&services.TeamService{}).EditTeam(params)
	if addErr != nil {
		c.JSONError(500, addErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	业务员团队修改
// @Param	request body services.TeamRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /editSalespersonTeam [post]
func (c *TeamController) EditSalespersonTeam() {
	var params services.TeamRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	addErr := (&services.TeamService{}).EditTeam(params)
	if addErr != nil {
		c.JSONError(500, addErr.Error())
	}
	c.JSONSuccess(nil)
}

type TeamDelRequestParams struct {
	Id int `json:"ids"`
}

// @Summary	团队删除
// @Param	request body controllers.TeamDelRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /delTeam [post]
func (c *TeamController) DelTeam() {
	var params TeamDelRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	addErr := (&services.TeamService{}).DelTeam(params.Id)
	if addErr != nil {
		c.JSONError(500, addErr.Error())
	}
	c.JSONSuccess(nil)
}

// @Summary	业务员团队删除
// @Param	request body controllers.TeamDelRequestParams true "json请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /delSalespersonTeam [post]
func (c *TeamController) DelSalespersonTeam() {
	var params TeamDelRequestParams
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "CanShuHuoQuShiBai")
	}
	addErr := (&services.TeamService{}).DelTeam(params.Id)
	if addErr != nil {
		c.JSONError(500, addErr.Error())
	}
	c.JSONSuccess(nil)
}
