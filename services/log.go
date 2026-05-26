package services

import (
	"api/config"
	"api/models"
	"net/url"

	"github.com/beego/i18n"
)

type LogService struct {
	BaseService
}

// UserBalanceLogList 余额变动日志列表
func (s LogService) UserBalanceLogList(request UserBalanceLogRequestParams, needReload int, lang string, adminId int64) (map[string]interface{}, error) {
	data, err := (&UserBalanceLogService{}).GetList(request)
	if nil != err {
		return nil, err
	}
	var activityType []config.ActiveTypeItem
	var gameType []models.GameType
	var packages []models.Package
	var userBalanceLogType []map[string]interface{}
	if needReload == 1 {
		activityService := ActiveService{}
		activityType = activityService.GetActiveTypeList()
		gameType, err = (GameTypeService{}).AllGameType(request.PackageIds)
		if err != nil {
			return nil, err
		}
		packages = (&PackageService{}).GetMyAllPackageList(int(adminId))

		userBalanceLogType = s.AllUserBalanceLogType(lang)
	}
	return map[string]interface{}{
		"list":                  data.List,
		"total":                 data.Total,
		"current_page":          request.Page,
		"activity_type":         activityType,
		"game_type":             gameType,
		"packages":              packages,
		"user_balance_log_type": userBalanceLogType,
	}, nil
}

// AllUserBalanceLogType 所有余额变动日志类型
func (s LogService) AllUserBalanceLogType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.BalanceLogTypeRecharge, "name": i18n.Tr(language, "ChongZhi")},
		{"id": config.BalanceLogTypeRechargeExtra, "name": i18n.Tr(language, "ChongZhiEWaiZengSong")},
		{"id": config.BalanceLogTypeWithdraw, "name": i18n.Tr(language, "TiXian")},
		{"id": config.BalanceLogTypeActivityInvite, "name": i18n.Tr(language, "YaoQingHuoDong")},
		{"id": config.BalanceLogTypeActivitySign, "name": i18n.Tr(language, "QianDaoHuoDong")},
		{"id": config.BalanceLogTypeActivityFirstRecharge, "name": i18n.Tr(language, "ShouChongHuoDong")},
		{"id": config.BalanceLogTypeActivityReliefDaily, "name": i18n.Tr(language, "JiuJiJinHuoDongRiJiangLi")},
		{"id": config.BalanceLogTypeActivityReliefWeekly, "name": i18n.Tr(language, "JiuJiJinHuoDongZhouJiangLi")},
		{"id": config.BalanceLogTypeActivityReliefMonthly, "name": i18n.Tr(language, "JiuJiJinHuoDongYueJiangLi")},
		{"id": config.BalanceLogTypeActivityInterest, "name": i18n.Tr(language, "LiXiBaoHuoDong")},
		{"id": config.BalanceLogTypeActivityLuckyWheel, "name": i18n.Tr(language, "ZhuanPanHuoDong")},
		{"id": config.BalanceLogTypeActivityVipLv, "name": i18n.Tr(language, "VipHuoDongJinJiJiangLi")},
		{"id": config.BalanceLogTypeActivityVipLvDaily, "name": i18n.Tr(language, "VipHuoDongRenWuRiJiangLi")},
		{"id": config.BalanceLogTypeActivityVipLvWeekly, "name": i18n.Tr(language, "VipHuoDongRenWuZhouJiangLi")},
		{"id": config.BalanceLogTypeActivityVipLvMonthly, "name": i18n.Tr(language, "VipHuoDongRenWuYueJiangLi")},
		{"id": config.BalanceLogTypeRegister, "name": i18n.Tr(language, "ZhuCe")},
		{"id": config.BalanceLogTypeSettlement, "name": i18n.Tr(language, "JieSuan")},
		{"id": config.BalanceLogTypeAgent, "name": i18n.Tr(language, "DaiLiYongJinJiangLi")},
		{"id": config.BalanceLogTypeSetUserBalance, "name": i18n.Tr(language, "XiuGaiYongHuYuE")},
		{"id": config.BalanceLogTypeGmAdd, "name": i18n.Tr(language, "GMShangFen")},
		{"id": config.BalanceLogTypeGmSubtract, "name": i18n.Tr(language, "GMXiaFen")},
		{"id": config.BalanceLogTypeAddInviteReward, "name": i18n.Tr(language, "GMZengJiaDaiLiJieMianYaoQingJiangLi")},
		{"id": config.BalanceLogTypeSubInviteReward, "name": i18n.Tr(language, "GMJianShaoDaiLiJieMianYaoQingJiangLi")},
		{"id": config.BalanceLogTypeAddAdReward, "name": i18n.Tr(language, "GMGuangGaoFeiXiaFa")},
		{"id": config.BalanceLogTypeSubAdReward, "name": i18n.Tr(language, "GMGuangGaoFeiKouChu")},
		{"id": config.BalanceLogTypeAddAccount, "name": i18n.Tr(language, "GMMoNiZhangHuJiaKuan")},
		{"id": config.BalanceLogTypeSubAccount, "name": i18n.Tr(language, "GMMoNiZhangHuJianKuan")},
	}
}

// GameWaterLogList 流水日志列表
func (s LogService) GameWaterLogList(waterLogParams GameWaterLogRequestParams, needReload int, adminId int64) (map[string]interface{}, error) {
	result, err := (&GameWaterLogService{}).GetList(waterLogParams, int(adminId))
	if nil != err {
		return nil, err
	}
	var gameType []models.GameType
	var platform []models.Platform
	var packages []models.Package
	if needReload == 1 {
		gameTypeService := GameTypeService{}
		gameType, err = gameTypeService.AllGameType(waterLogParams.PackageIds)
		if err != nil {
			return nil, err
		}
		platformService := PlatformService{}
		platform, err = platformService.AllPlatformByPackageIds(waterLogParams.PackageIds)
		if err != nil {
			return nil, err
		}
		packages = (&PackageService{}).GetMyAllPackageList(int(adminId))

	}
	return map[string]interface{}{
		"list":         result.List,
		"total":        result.Total,
		"current_page": waterLogParams.Page,
		"game_type":    gameType,
		"packages":     packages,
		"platform":     platform,
	}, nil
}

type GetSystemLogRequest struct {
	Id         int        `form:"id"`
	AdminId    int        `form:"admin_id"`
	PlayerId   int        `form:"player_id"`
	PackageId  int        `form:"package_id"`
	Path       string     `form:"path"`
	Body       string     `form:"body" op:"like"`
	LogType    int        `form:"log_type"`
	CTime      []string   `form:"c_time[]" op:"time_range" orm:"c_time"`
	NeedReload int        `form:"need_reload" op:"-"`
	Page       int        `form:"page"`
	PageSize   int        `form:"page_size"`
	RawQuery   url.Values `form:"-"`
	PackageIds []int
}

// GetSystemLog 系统日志
func (s LogService) GetSystemLog(request GetSystemLogRequest, adminId int, lang string) (map[string]interface{}, error) {
	adminLog := models.CreateAdminLogModel()
	condition, sort := adminLog.BuildCondition(request, "-id")

	condition = s.LimitPackageId(condition, request.PackageIds)
	data, total, err := adminLog.GetPageList(&models.AdminLog{}, condition, request.Page, request.PageSize, sort)

	if nil != err {
		return nil, err
	}

	// 取操作方法
	var methods []models.AuthRule
	var typeList []map[string]interface{}
	packages := make([]models.Package, 0)
	if request.NeedReload == 1 {
		models.CreateAuthRuleModel().QueryTable(new(models.AuthRule)).Filter("tag", 0).All(&methods, "title", "name")
		typeList = s.GetAllSystemLogTypes(lang)

		packages = (&PackageService{}).GetMyAllPackageList(adminId)
	}

	list := data.([]models.AdminLog)
	return map[string]interface{}{
		"router_list":  methods,
		"type_list":    typeList,
		"list":         list,
		"total":        total,
		"package_list": packages,
	}, nil
}

// GetAllSystemLogTypes 获取所有系统日志类型
func (s LogService) GetAllSystemLogTypes(lang string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.AdminLogTypeCommon, "title": i18n.Tr(lang, "ChangGuiCaoZuo")},
		{"id": config.AdminLogTypeRefreshTable, "title": i18n.Tr(lang, "ShuaXinBiaoShuJu")},
		{"id": config.AdminLogTypeGmOperateAdd, "title": i18n.Tr(lang, "ShangXiaFenCaoZuo")},
	}
}
