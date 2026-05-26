package controllers

import "api/services"

type DailyStatisticsController struct {
	BaseController
}

// @Summary	日况统计
// @Param page query int false "当前页"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /lst [get]
func (c *DailyStatisticsController) Lst() {
	request := services.LstRequest{}
	c.ParseForm(&request)
	request.RawQuery = c.Ctx.Request.URL.Query()

	if request.PackageId > 0 {
		if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
			c.JSONError(500, "ShuJuYiChang")
		}
	} else {
		request.PackageIds = c.GetPackageIds(int(c.UserId))
	}

	data, err := (&services.DailyStatisticsService{}).Lst(request, int(c.UserId))
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	日报表
// @Param page query int false "当前页"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /totalDaily [get]
func (c *DailyStatisticsController) TotalDaily() {
	request := services.TotalDailyRequest{}
	c.ParseForm(&request)
	request.RawQuery = c.Ctx.Request.URL.Query()

	if request.PackageId > 0 {
		if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
			c.JSONError(500, "ShuJuYiChang")
		}
	} else {
		request.PackageIds = c.GetPackageIds(int(c.UserId))
	}

	data, err := (&services.DailyStatisticsService{}).TotalDaily(request, int(c.UserId))
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}

// @Summary	月报表
// @Param page query int false "当前页"
// @Param page_size query int false "每页数量"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {object} {"code": 0,"msg": "error message"}
// @router /totalMonth [get]
func (c *DailyStatisticsController) TotalMonth() {
	request := services.TotalMonthRequest{}
	c.ParseForm(&request)
	request.RawQuery = c.Ctx.Request.URL.Query()

	if request.PackageId > 0 {
		if !c.IsMyPackageId(request.PackageId, int(c.UserId)) {
			c.JSONError(500, "ShuJuYiChang")
		}
	} else {
		request.PackageIds = c.GetPackageIds(int(c.UserId))
	}

	data, err := (&services.DailyStatisticsService{}).TotalMonth(request, int(c.UserId))
	if err != nil {
		c.JSONError(500, err.Error())
	} else {
		c.JSONSuccess(data)
	}
}
