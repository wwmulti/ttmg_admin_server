package controllers

import "api/services"

type ServesController struct {
	BaseController
}

// @Summary	刷新内存mysql表
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /refreshTable [get]
func (c *ServesController) RefreshTable() {
	taskId, _ := (&services.TaskService{}).AddTask(services.TaskTypeRefresh, map[string]interface{}{})
	c.JSONSuccess(map[string]interface{}{
		"taskId": taskId,
	})
}

// @Summary	批量获取任务状态
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /getMultiTaskProcess [get]
func (c *ServesController) GetMultiTaskProcess() {
	taskIds := c.GetString("taskIds", "")
	data := (&services.TaskService{}).GetMultiTaskProcess(taskIds)
	c.JSONSuccess(data)
}
