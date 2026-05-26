package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/beego/i18n"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type ActiveService struct {
	BaseService
}

// GetActiveTypeList 获取活动类型
func (s ActiveService) GetActiveTypeList() []config.ActiveTypeItem {
	return config.ActiveTypeList
}

type ActiveListRequest struct {
	Id           int `form:"id"`
	PackageId    int `form:"package_id"`
	Status       int `form:"status"`
	ActiveTypeId int `form:"active_type_id"`
	/* IsNewcomer     int        `form:"is_newcomer"`
	IsPay          int        `form:"is_pay"`
	IsHideComplete int        `form:"is_hide_complete"` */
	Page       int        `form:"page"`
	PageSize   int        `form:"page_size"`
	RawQuery   url.Values `form:"-"`
	PackageIds []int      // 自身授权的分包
}

// ActiveList 活动列表
func (s ActiveService) ActiveList(request ActiveListRequest, needReload int, language string, adminId int64) (map[string]interface{}, error) {
	activeModel := models.CreateActivesModel()
	condition, sort := activeModel.BuildCondition(request, "-package_id")

	condition = s.LimitPackageId(condition, request.PackageIds)
	condition["is_deleted"] = 0
	data, total, err := activeModel.GetPageList(&models.Actives{}, condition, request.Page, request.PageSize, sort)

	if nil != err {
		return nil, err
	}
	var activityType []map[string]interface{}
	var packages []models.Package
	if needReload == 1 {
		activityType = s.AllActivityType(language)
		packages = (&PackageService{}).GetMyAllPackageList(int(adminId))
	}
	list := data.([]models.Actives)
	return map[string]interface{}{
		"list":          list,
		"total":         total,
		"activity_type": activityType,
		"packages":      packages,
	}, nil
}

// AllActivityType 所有活动类型
func (s ActiveService) AllActivityType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.ActiveTypeInvite, "name": i18n.Tr(language, "YaoQingHuoDong")},
		{"id": config.ActiveTypeSign, "name": i18n.Tr(language, "QianDaoHuoDong")},
		{"id": config.ActiveTypeFirstRecharge, "name": i18n.Tr(language, "ShouChongHuoDong")},
		{"id": config.ActiveTypeRelief, "name": i18n.Tr(language, "JiuJiJinHuoDong")},
		{"id": config.ActiveTypeInterest, "name": i18n.Tr(language, "LiXiBaoHuoDong")},
		{"id": config.ActiveTypeLuckyWheel, "name": i18n.Tr(language, "XinYunZhuanPanHuoDong")},
		{"id": config.ActiveTypeVipLv, "name": i18n.Tr(language, "VipDengJiHuoDong")},
	}
}

type AddActiveRequest struct {
	Name         string `json:"name"`           // 活动名
	PtName       string `json:"pt_name"`        // 葡语名称
	EnName       string `json:"en_name"`        // 英语名称
	Icon         string `json:"icon"`           // 活动图标
	PcIcon       string `json:"pc_icon"`        // pc端图片
	ActiveTypeId int    `json:"active_type_id"` // 活动类型id
	RedirectLink int    `json:"redirect_link"`  // 跳转链接
	PtDesc       string `json:"pt_desc"`        // 葡语描述
	EnDesc       string `json:"en_desc"`        // 英语描述
	StartTime    int64  `json:"start_time"`     // 活动开始时间
	EndTime      int64  `json:"end_time"`       // 活动结束时间
	CollecTime   int64  `json:"collec_time"`    // 领取时间
	Sort         int    `json:"sort"`           // 排序
	Status       int    `json:"status"`         // 活动状态
	PackageIds   string `json:"package_ids"`    // 开放平台id(1,2,3逗号分隔)
}

// 活动类型唯一
var uniqueActiveTypes = []config.ActiveType{
	config.ActiveTypeInvite,
	config.ActiveTypeVipLv,
}

// AddActive 添加活动
func (s ActiveService) AddActive(request AddActiveRequest, adminId int64) error {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	err := s.validateAddActive(request)
	if err != nil {
		return err
	}

	uniqueCheck := false
	if slices.Contains(uniqueActiveTypes, config.ActiveType(request.ActiveTypeId)) {
		uniqueCheck = true
	}

	tx, err := models.CreateActivesModel().Begin()
	if err != nil {
		logs.Error("添加活动事务开启失败，失败原因为：%v", err)
		return err
	}
	packageIds := strings.Split(request.PackageIds, ",")
	for _, packageId := range packageIds {
		pkgId, _ := strconv.Atoi(packageId)
		if !isRootAdmin && !utils.InArray(pkgId, packageSlice) {
			return fmt.Errorf("QuanXianBuZu")
		}
		if uniqueCheck {
			isExist := s.CheckActiveExist(pkgId, request.ActiveTypeId)
			if isExist {
				tx.Rollback()
				logs.Info("活动类型已存在 包id: %v 类型: %v ", pkgId, request.ActiveTypeId)
				continue
			}
		}
		_, err = tx.Insert(&models.Actives{
			Name:         request.Name,
			PtName:       request.PtName,
			EnName:       request.EnName,
			Icon:         request.Icon,
			PcIcon:       request.PcIcon,
			ActiveTypeId: request.ActiveTypeId,
			RedirectLink: request.RedirectLink,
			PtDesc:       request.PtDesc,
			EnDesc:       request.EnDesc,
			StartTime:    request.StartTime,
			EndTime:      request.EndTime,
			CollecTime:   request.CollecTime,
			Sort:         request.Sort,
			Status:       request.Status,
			PackageId:    pkgId,
			Ctime:        time.Now().Unix(),
		})
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("添加活动事务提交失败，失败原因为：%v", err)
		return err
	}
	return nil
}

// CheckActiveExist 检查活动是否存在
func (s ActiveService) CheckActiveExist(packageId int, activeTypeId int) bool {
	activeModel := models.CreateActivesModel()
	return activeModel.QueryTable(new(models.Actives)).
		Filter("package_id", packageId).
		Filter("active_type_id", activeTypeId).
		Exist()
}

// 添加活动参数验证
func (s ActiveService) validateAddActive(request AddActiveRequest) error {
	if request.Name == "" {
		return fmt.Errorf("HuoDongMingChengBiTian")
	}
	if request.PackageIds == "" {
		return fmt.Errorf("PingTaiBiTian")
	}
	//isExist := models.CreateActivesModel().QueryTable(new(models.Actives)).
	//	Filter("name", request.Name).
	//	Filter("is_deleted", 0).
	//	Exist()
	//if isExist {
	//	return fmt.Errorf("HuoDongMingChengYiCunZai")
	//}
	/* if request.Icon == "" {
		return fmt.Errorf("HuoDongLogoBiTian")
	}
	if request.StartTime == 0 {
		return fmt.Errorf("HuoDongKaiShiShiJianBiTian")
	}
	if request.EndTime == 0 {
		return fmt.Errorf("HuoDongJieShiShiJianBiTian")
	}
	if request.CollecTime == 0 {
		return fmt.Errorf("HuoDongLingQuShiJianBiTian")
	} */
	if request.ActiveTypeId == 0 {
		return fmt.Errorf("HuoDongLeiXingBiTian")
	}
	/* if request.RedirectLink == 0 {
		return fmt.Errorf("HuoDongTiaoZhuanDiZhiBiTian")
	} */
	return nil
}

type EditActiveRequest struct {
	Id           int    `json:"id"`             // 活动ID
	Name         string `json:"name"`           // 活动名
	PtName       string `json:"pt_name"`        // 葡语名称
	EnName       string `json:"en_name"`        // 英语名称
	Icon         string `json:"icon"`           // 活动图标
	PcIcon       string `json:"pc_icon"`        // pc端图片
	ActiveTypeId int    `json:"active_type_id"` // 活动类型id
	RedirectLink int    `json:"redirect_link"`  // 跳转链接
	PtDesc       string `json:"pt_desc"`        // 葡语描述
	EnDesc       string `json:"en_desc"`        // 英语描述
	StartTime    int64  `json:"start_time"`     // 活动开始时间
	EndTime      int64  `json:"end_time"`       // 活动结束时间
	CollecTime   int64  `json:"collec_time"`    // 领取时间
	Sort         int    `json:"sort"`           // 排序
	Status       int    `json:"status"`         // 活动状态
}

// EditActive 编辑活动
func (s ActiveService) EditActive(request EditActiveRequest, adminId int64) error {
	activity, err := s.GetActiveById(request.Id)
	if err != nil {
		logs.Error("活动获取失败: %v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(activity.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	err = s.validateEditActive(request)
	if err != nil {
		return err
	}
	fields := []string{
		"name",
		"pt_name",
		"en_name",
		"icon",
		"pc_icon",
		"redirect_link",
		"pt_desc",
		"en_desc",
		"start_time",
		"end_time",
		"collec_time",
		"sort",
		"status",
	}
	activeModel := models.CreateActivesModel()
	_, err = activeModel.Update(&models.Actives{
		Id:           request.Id,
		Name:         request.Name,
		PtName:       request.PtName,
		EnName:       request.EnName,
		Icon:         request.Icon,
		PcIcon:       request.PcIcon,
		RedirectLink: request.RedirectLink,
		PtDesc:       request.PtDesc,
		EnDesc:       request.EnDesc,
		StartTime:    request.StartTime,
		EndTime:      request.EndTime,
		CollecTime:   request.CollecTime,
		Sort:         request.Sort,
		Status:       request.Status,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑活动参数验证
func (s ActiveService) validateEditActive(request EditActiveRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	if request.Name == "" {
		return fmt.Errorf("HuoDongMingChengBiTian")
	}
	//isExist := models.CreateActivesModel().QueryTable(new(models.Actives)).
	//	Filter("id__ne", request.Id).
	//	Filter("name", request.Name).
	//	Filter("is_deleted", 0).
	//	Exist()
	//if isExist {
	//	return fmt.Errorf("HuoDongMingChengYiCunZai")
	//}
	if request.Icon == "" {
		return fmt.Errorf("HuoDongLogoBiTian")
	}
	/* if request.StartTime == 0 {
		return fmt.Errorf("HuoDongKaiShiShiJianBiTian")
	}
	if request.EndTime == 0 {
		return fmt.Errorf("HuoDongJieShiShiJianBiTian")
	}
	if request.CollecTime == 0 {
		return fmt.Errorf("HuoDongLingQuShiJianBiTian")
	}
	if request.ActiveTypeId == 0 {
		return fmt.Errorf("HuoDongLeiXingBiTian")
	}
	if request.RedirectLink == 0 {
		return fmt.Errorf("HuoDongTiaoZhuanDiZhiBiTian")
	} */
	return nil
}

// DeleteActive 删除活动
func (s ActiveService) DeleteActive(id int, adminId int64) error {
	activity, err := s.GetActiveById(id)
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return fmt.Errorf("HuoDongBuCunZai")
		} else {
			return err
		}
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(activity.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	if activity.PackageId == 0 {
		return fmt.Errorf("WuPingTaiHuoDongWuFaShanChu")
	}
	activeModel := models.CreateActivesModel()
	if id == 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	_, err = activeModel.Update(&models.Actives{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	return nil
}

// EditActiveAttr 编辑活动属性
func (s ActiveService) EditActiveAttr(id int, field string, adminId int64) error {
	activity, err := s.GetActiveById(id)
	if err != nil {
		return fmt.Errorf("HuoDongBuCunZai")
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(activity.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	activityModel := models.CreateActivesModel()
	switch field {
	case "status":
		if activity.Status == 1 {
			activity.Status = 0
		} else {
			activity.Status = 1
		}
	case "is_newcomer":
		if activity.IsNewcomer == 1 {
			activity.IsNewcomer = 0
		} else {
			activity.IsNewcomer = 1
		}
	case "is_pay":
		if activity.IsPay == 1 {
			activity.IsPay = 0
		} else {
			activity.IsPay = 1
		}
	case "is_hide_complete":
		if activity.IsHideComplete == 1 {
			activity.IsHideComplete = 0
		} else {
			activity.IsHideComplete = 1
		}
	default:
		return fmt.Errorf("GengXinNeiRongBuCunZai")
	}
	_, err = activityModel.Update(&activity, field)
	if err != nil {
		return err
	}
	return nil
}

type UserActivityRewardLogRequest struct {
	UserId         int `json:"user_id"`          // 用户ID
	ActivityId     int `json:"activity_id"`      // 活动ID
	ActivityTypeId int `json:"activity_type_id"` // 活动类型ID
	ReceiveTypeId  int `json:"receive_type_id"`  // 领取类型ID
	Page           int `json:"page"`             // 页码
	PageSize       int `json:"page_size"`        // 每页数量
}

// UserActivityRewardLogList 用户活动奖励列表
func (s ActiveService) UserActivityRewardLogList(request UserActivityRewardLogRequest) (map[string]interface{}, error) {
	signRuleModel := models.CreateUserActivityRewardLogModel()
	var list []models.UserActivityRewardLog
	qs := signRuleModel.QueryTable(&models.UserActivityRewardLog{})
	if request.UserId > 0 {
		qs = qs.Filter("user_id", request.UserId)
	}
	if request.ActivityId != 0 {
		qs = qs.Filter("activity_id", request.ActivityId)
	}
	if request.ActivityTypeId > 0 {
		qs = qs.Filter("activity_type_id", request.ActivityTypeId)
	}
	if request.ReceiveTypeId > 0 {
		qs = qs.Filter("receive_type_id", request.ReceiveTypeId)
	}
	_, err := qs.Limit(request.PageSize, (request.Page-1)*request.PageSize).
		OrderBy("-id").
		All(&list)
	if nil != err {
		return nil, err
	}
	total, err := qs.Count()
	if nil != err {
		return nil, err
	}

	return map[string]interface{}{
		"list":         list,
		"total":        total,
		"current_page": request.Page,
	}, nil
}

// GetActiveById 根据ID获取活动
func (s ActiveService) GetActiveById(activeId int) (*models.Actives, error) {
	var active models.Actives
	err := models.CreateActivesModel().QueryTable(&models.Actives{}).
		Filter("id", activeId).
		Filter("is_deleted", 0).
		One(&active)
	if err != nil {
		return nil, err
	}
	return &active, nil
}

type CreateActivityByTypeRequest struct {
	TypeIds       []int `json:"type_ids"`        // 类型ID数组
	PackageId     int   `json:"package_id"`      // 包ID
	StatusList    []int `json:"status"`          // 状态数组
	CopyPackageId int   `json:"copy_package_id"` // 复制的包ID
}

// CopyActivity 复制活动
func (s ActiveService) CopyActivity(params CreateActivityByTypeRequest, tx orm.TxOrmer) error {
	var err error
	for _, item := range config.ActiveTypeList {
		params := CreateActivityParams{
			PackageId:     params.PackageId,
			TypeId:        int(item.ID),
			IsCopy:        true,
			CopyPackageId: params.CopyPackageId,
		}

		typeId := int(item.ID)

		switch typeId {
		case int(config.ActiveTypeInvite):
			err = s.createInviteActivity(params, tx)
		case int(config.ActiveTypeSign):
			err = s.createSignActivity(params, tx)
		case int(config.ActiveTypeFirstRecharge):
			err = s.createFirstRechargeActivity(params, tx)
		case int(config.ActiveTypeRelief):
			err = s.createReliefActivity(params, tx)
		case int(config.ActiveTypeInterest):
			err = s.createInterestActivity(params, tx)
		case int(config.ActiveTypeLuckyWheel):
			err = s.createLuckWheelActivity(params, tx)
		case int(config.ActiveTypeVipLv):
			err = s.createVipLvActivity(params, tx)
		default:
			logs.Error("活动类型不存在：%v", typeId)
			return fmt.Errorf("HuoDongLeiXingBuCunZai")
		}
	}

	return err
}

// CreateActivityByType 根据活动类型批量创建活动
func (s ActiveService) CreateActivityByType(params CreateActivityByTypeRequest, tx orm.TxOrmer) error {
	var err error
	for index, typeId := range params.TypeIds {
		params := CreateActivityParams{
			PackageId: params.PackageId,
			Status:    params.StatusList[index],
			TypeId:    typeId,
		}

		switch typeId {
		case int(config.ActiveTypeInvite):
			err = s.createInviteActivity(params, tx)
		case int(config.ActiveTypeSign):
			err = s.createSignActivity(params, tx)
		case int(config.ActiveTypeFirstRecharge):
			err = s.createFirstRechargeActivity(params, tx)
		case int(config.ActiveTypeRelief):
			err = s.createReliefActivity(params, tx)
		case int(config.ActiveTypeInterest):
			err = s.createInterestActivity(params, tx)
		case int(config.ActiveTypeLuckyWheel):
			err = s.createLuckWheelActivity(params, tx)
		case int(config.ActiveTypeVipLv):
			err = s.createVipLvActivity(params, tx)
		default:
			logs.Error("活动类型不存在：%v", typeId)
			return fmt.Errorf("HuoDongLeiXingBuCunZai")
		}
	}
	if err != nil {
		return err
	} else {
		return nil
	}
}

type CreateActivityParams struct {
	PackageId     int  `json:"package_id"` // 包ID
	Status        int  `json:"status"`
	TypeId        int  `json:"type_id"`
	IsCopy        bool `json:"is_copy"`         // 是否复制指定分包的配置
	CopyPackageId int  `json:"copy_package_id"` // 复制包ID
}

// CreateActivity 创建活动
func (s ActiveService) createActivity(params CreateActivityParams, tx orm.TxOrmer) (*models.Actives, int, bool, error) {
	now := time.Now().Unix()
	tomorrow := utils.GetTomorrowTimestamp()
	var template models.Actives

	// 默认分包id为0的模板，复制指定分包id
	templateId := 0
	activityStatus := params.Status
	if params.IsCopy {
		templateId = params.CopyPackageId
	}

	err := models.CreateActivesModel().QueryTable(&models.Actives{}).
		Filter("active_type_id", params.TypeId).
		Filter("package_id", templateId).
		One(&template)
	if err != nil {
		logs.Error("模板不存在：%v", params.TypeId)
		return nil, 0, false, fmt.Errorf("HuoDongMuBanBuCunZai")
	}

	if params.IsCopy {
		activityStatus = template.Status
	}

	var activity models.Actives
	err = models.CreateActivesModel().QueryTable(&models.Actives{}).
		Filter("active_type_id", params.TypeId).
		Filter("package_id", params.PackageId).
		One(&activity)
	if err != nil {
		if !errors.Is(err, orm.ErrNoRows) {
			return nil, 0, false, fmt.Errorf("WeiZhiDeCuoWu")
		} else {
			activity := models.Actives{
				ActiveTypeId:    template.ActiveTypeId,
				PackageId:       params.PackageId,
				Name:            template.Name,
				PtName:          template.PtName,
				EnName:          template.EnName,
				Icon:            template.Icon,
				RedirectLink:    template.RedirectLink,
				PtDesc:          template.PtDesc,
				EnDesc:          template.EnDesc,
				StartTime:       now,
				EndTime:         tomorrow,
				CollecTime:      tomorrow,
				TimerCheckAward: template.TimerCheckAward,
				TimerCheckTime:  template.TimerCheckTime,
				ValidBets:       template.ValidBets,
				IsNewcomer:      template.IsNewcomer,
				IsPay:           template.IsPay,
				IsHideComplete:  template.IsHideComplete,
				Sort:            template.Sort,
				Ctime:           now,
				Status:          activityStatus,
			}
			id, err := tx.Insert(&activity)
			if err != nil {
				return nil, 0, false, err
			}
			activity.Id = int(id)

			return &activity, template.Id, false, nil
		}
	} else {
		data := models.Actives{
			Id:     activity.Id,
			Status: activityStatus,
		}
		fields := []string{
			"status",
		}
		_, err := tx.Update(&data, fields...)
		if err != nil {
			return nil, 0, true, fmt.Errorf("GengXinHuoDongZhuangTaiShiBai")
		}
		return &activity, template.Id, true, nil
	}
}

// createInviteActivity 创建邀请活动
func (s ActiveService) createInviteActivity(params CreateActivityParams, tx orm.TxOrmer) error {
	activity, templateActivityId, isUpdate, err := s.createActivity(params, tx)
	if err != nil {
		return err
	}
	if isUpdate {
		return nil
	}
	var templateRule models.ActiveShareRule
	ruleModel := models.CreateActiveShareRuleModel()
	err = ruleModel.QueryTable(&models.ActiveShareRule{}).
		Filter("active_id", templateActivityId).
		One(&templateRule)
	if err != nil {
		return fmt.Errorf("YaoQingHuoDongMuBanGuiZeBuCunZai")
	}
	rule := models.ActiveShareRule{
		ActiveId:       activity.Id,
		ActiveTypeId:   templateRule.ActiveTypeId,
		ShareUrl:       templateRule.ShareUrl,
		ScanInterval:   templateRule.ScanInterval,
		Scope:          templateRule.Scope,
		MiniTotalPays:  templateRule.MiniTotalPays,
		MiniTotalWater: templateRule.MiniTotalWater,
		Condition:      templateRule.Condition,
		RewardType:     templateRule.RewardType,
		ExpireType:     templateRule.ExpireType,
		Ctime:          time.Now().Unix(),
		Status:         1,
	}
	ruleId, err := tx.Insert(&rule)
	if err != nil {
		return err
	}
	var list []models.ActiveShareRewards
	ActiveShareRewardsModel := models.CreateActiveShareRewardsModel()
	_, err = ActiveShareRewardsModel.QueryTable(&models.ActiveShareRewards{}).
		Filter("active_id", templateActivityId).
		All(&list)
	if nil != err {
		return err
	}
	for _, item := range list {
		rewards := models.ActiveShareRewards{
			ActiveId: activity.Id,
			RuleId:   int(ruleId),
			Mens:     item.Mens,
			Rewards:  item.Rewards,
			IconOff:  item.IconOff,
			IconOn:   item.IconOn,
			Ctime:    time.Now().Unix(),
			Status:   1,
		}
		_, err = tx.Insert(&rewards)
		if err != nil {
			return err
		}
	}
	return nil
}

// createSignActivity 创建签到活动
func (s ActiveService) createSignActivity(params CreateActivityParams, tx orm.TxOrmer) error {
	activity, templateActivityId, isUpdate, err := s.createActivity(params, tx)
	if err != nil {
		return err
	}
	if isUpdate {
		return nil
	}
	var templateRule models.SignRule
	ruleModel := models.CreateSignRuleModel()
	err = ruleModel.QueryTable(&models.SignRule{}).
		Filter("activity_id", templateActivityId).
		One(&templateRule)
	if err != nil {
		return fmt.Errorf("QianDaoHuoDongMuBanGuiZeBuCunZai")
	}
	rule := models.SignRule{
		ActivityId:       activity.Id,
		Days:             templateRule.Days,
		IsLoop:           templateRule.IsLoop,
		IsInterruptReset: templateRule.IsInterruptReset,
		Status:           1,
	}
	ruleId, err := tx.Insert(&rule)
	if err != nil {
		return err
	}
	var list []models.SignRuleReward
	SignRewardsModel := models.CreateSignRuleRewardModel()
	_, err = SignRewardsModel.QueryTable(&models.SignRuleReward{}).
		Filter("sign_rule_id", templateRule.Id).
		All(&list)
	if nil != err {
		return err
	}
	for _, item := range list {
		rewards := models.SignRuleReward{
			SignRuleId:      int(ruleId),
			Day:             item.Day,
			RewardTypeId:    item.RewardTypeId,
			RewardAmount:    item.RewardAmount,
			Icon:            item.Icon,
			DayRechargeLine: item.DayRechargeLine,
			DayRunningLine:  item.DayRunningLine,
		}
		_, err = tx.Insert(&rewards)
		if err != nil {
			return err
		}
	}
	return nil
}

// createFirstRechargeActivity 创建首充活动
func (s ActiveService) createFirstRechargeActivity(params CreateActivityParams, tx orm.TxOrmer) error {
	activity, templateActivityId, isUpdate, err := s.createActivity(params, tx)
	if err != nil {
		return err
	}
	if isUpdate {
		return nil
	}
	var templateRule models.FirstRechargeRule
	ruleModel := models.CreateFirstRechargeRuleModel()
	err = ruleModel.QueryTable(&models.FirstRechargeRule{}).
		Filter("activity_id", templateActivityId).
		One(&templateRule)
	if err != nil {
		return fmt.Errorf("ShouChongHuoDongMuBanGuiZeBuCunZai")
	}
	rule := models.FirstRechargeRule{
		ActivityId:      activity.Id,
		BillingType:     templateRule.BillingType,
		ReceiveType:     templateRule.ReceiveType,
		TaskType:        templateRule.TaskType,
		BetNumber:       templateRule.BetNumber,
		BetAmount:       templateRule.BetAmount,
		UpdateFrequency: templateRule.UpdateFrequency,
		RepeatActive:    templateRule.RepeatActive,
		Status:          1,
	}
	ruleId, err := tx.Insert(&rule)
	if err != nil {
		return err
	}
	var list []models.FirstRechargeRuleReward
	FirstRechargeRewardsModel := models.CreateFirstRechargeRuleRewardModel()
	_, err = FirstRechargeRewardsModel.QueryTable(&models.FirstRechargeRuleReward{}).
		Filter("first_recharge_rule_id", templateRule.Id).
		All(&list)
	if nil != err {
		return err
	}
	for _, item := range list {
		rewards := models.FirstRechargeRuleReward{
			FirstRechargeRuleId: int(ruleId),
			SerialNumber:        item.SerialNumber,
			TotalRechargeAmount: item.TotalRechargeAmount,
			RewardAmount:        item.RewardAmount,
		}
		_, err = tx.Insert(&rewards)
		if err != nil {
			return err
		}
	}
	return nil
}

// createReliefActivity 创建救济金活动
func (s ActiveService) createReliefActivity(params CreateActivityParams, tx orm.TxOrmer) error {
	activity, templateActivityId, isUpdate, err := s.createActivity(params, tx)
	if err != nil {
		return err
	}
	if isUpdate {
		return nil
	}
	var templateRule models.ActiveReliefRule
	ruleModel := models.CreateActiveReliefRuleModel()
	err = ruleModel.QueryTable(&models.ActiveReliefRule{}).
		Filter("active_id", templateActivityId).
		One(&templateRule)
	if err != nil {
		return fmt.Errorf("JiuJiHuoDongMuBanGuiZeBuCunZai")
	}
	rule := models.ActiveReliefRule{
		ActiveId: activity.Id,
		Cycle:    templateRule.Cycle,
		OpenDay:  templateRule.OpenDay,
		OpenTime: templateRule.OpenTime,
		IsRepeat: templateRule.IsRepeat,
		Ctime:    time.Now().Unix(),
		Status:   1,
	}
	ruleId, err := tx.Insert(&rule)
	if err != nil {
		return err
	}
	var list []models.ActiveReliefRewards
	ActiveReliefRewardsModel := models.CreateActiveReliefRewardsModel()
	_, err = ActiveReliefRewardsModel.QueryTable(&models.ActiveReliefRewards{}).
		Filter("active_id", templateActivityId).
		All(&list)
	if nil != err {
		return err
	}
	for _, item := range list {
		rewards := models.ActiveReliefRewards{
			ActiveId:      activity.Id,
			RuleId:        int(ruleId),
			Amount:        item.Amount,
			RebatePercent: item.RebatePercent,
			Ctime:         time.Now().Unix(),
			Status:        1,
		}
		_, err = tx.Insert(&rewards)
		if err != nil {
			return err
		}
	}
	return nil
}

// createInterestActivity 创建利息宝活动
func (s ActiveService) createInterestActivity(params CreateActivityParams, tx orm.TxOrmer) error {
	activity, templateActivityId, isUpdate, err := s.createActivity(params, tx)
	if err != nil {
		return err
	}
	if isUpdate {
		return nil
	}
	var templateRule models.InterestRule
	ruleModel := models.CreateInterestRuleModel()
	err = ruleModel.QueryTable(&models.InterestRule{}).
		Filter("activity_id", templateActivityId).
		One(&templateRule)
	if err != nil {
		return fmt.Errorf("LiXiBaoHuoDongMuBanGuiZeBuCunZai")
	}
	rule := models.InterestRule{
		ActivityId:          activity.Id,
		InterestRate:        templateRule.InterestRate,
		DepositAmount:       templateRule.DepositAmount,
		Interval:            templateRule.Interval,
		ReceiveType:         templateRule.ReceiveType,
		InterestLimitType:   templateRule.InterestLimitType,
		InterestLimitAmount: templateRule.InterestLimitAmount,
		Status:              1,
	}
	_, err = tx.Insert(&rule)
	if err != nil {
		return err
	}
	return nil
}

// createLuckWheelActivity 创建幸运转盘活动
func (s ActiveService) createLuckWheelActivity(params CreateActivityParams, tx orm.TxOrmer) error {
	activity, templateActivityId, isUpdate, err := s.createActivity(params, tx)
	if err != nil {
		return err
	}
	if isUpdate {
		return nil
	}
	var list []models.ActiveLuckyReward
	SignRewardsModel := models.CreateActiveLuckyRewardModel()
	_, err = SignRewardsModel.QueryTable(&models.ActiveLuckyReward{}).
		Filter("activity_id", templateActivityId).
		All(&list)
	if nil != err {
		return err
	}
	for _, item := range list {
		rewards := models.ActiveLuckyReward{
			ActivityId: activity.Id,
			WheelType:  item.WheelType,
			Reward:     item.Reward,
			Weight:     item.Weight,
			CreateTime: time.Now().Unix(),
		}
		_, err = tx.Insert(&rewards)
		if err != nil {
			return err
		}
	}
	return nil
}

// createVipLvActivity 创建VIP等级活动
func (s ActiveService) createVipLvActivity(params CreateActivityParams, tx orm.TxOrmer) error {
	activity, templateActivityId, isUpdate, err := s.createActivity(params, tx)
	if err != nil {
		return err
	}
	if isUpdate {
		return nil
	}
	var ruleList []models.ActiveVipRule
	ActiveVipRuleModel := models.CreateActiveVipRuleModel()
	_, err = ActiveVipRuleModel.QueryTable(&models.ActiveVipRule{}).
		Filter("active_id", templateActivityId).
		All(&ruleList)
	if nil != err {
		return err
	}
	for _, item := range ruleList {
		rewards := models.ActiveVipRule{
			ActiveId:            activity.Id,
			Lv:                  item.Lv,
			TotalBets:           item.TotalBets,
			TotalPays:           item.TotalPays,
			ConAnd:              item.ConAnd,
			Rewards:             item.Rewards,
			WithdrawNumLimit:    item.WithdrawNumLimit,
			WithdrawAmountLimit: item.WithdrawAmountLimit,
			WithdrawFreeNum:     item.WithdrawFreeNum,
			WithdrawFee:         item.WithdrawFee,
			Status:              1,
		}
		_, err = tx.Insert(&rewards)
		if err != nil {
			return err
		}
	}
	var welfareList []models.ActiveVipWelfare
	ActiveVipWelfareModel := models.CreateActiveVipWelfareModel()
	_, err = ActiveVipWelfareModel.QueryTable(&models.ActiveVipWelfare{}).
		Filter("active_id", templateActivityId).
		All(&welfareList)
	if nil != err {
		return err
	}
	for _, item := range welfareList {
		rewards := models.ActiveVipWelfare{
			ActiveId:  activity.Id,
			Cycle:     item.Cycle,
			Lv:        item.Lv,
			TotalBets: item.TotalBets,
			TotalPays: item.TotalPays,
			ConAnd:    item.ConAnd,
			Rewards:   item.Rewards,
			Ctime:     time.Now().Unix(),
			Status:    1,
		}
		_, err = tx.Insert(&rewards)
		if err != nil {
			return err
		}
	}
	return nil
}
