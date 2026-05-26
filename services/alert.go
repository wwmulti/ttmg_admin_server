package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
	"net/url"
	"strconv"
	"time"
)

type AlertService struct{}

type AlertListRequest struct {
	Id        int        `form:"id"`             // ID
	PackageId int        `form:"package_id"`     // 包ID
	Name      string     `form:"name" op:"like"` // 弹窗名称
	Status    int        `form:"status"`         // 状态
	Page      int        `form:"page"`           // 页码
	PageSize  int        `form:"page_size"`      // 每页数量
	RawQuery  url.Values `form:"-"`
}

// AlertList 弹窗列表
func (s AlertService) AlertList(request AlertListRequest, needReload int, lang string, adminId int64) (map[string]interface{}, error) {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	alertModel := models.CreateAlertModel()
	condition, sort := alertModel.BuildCondition(request, "-id")
	if !isRootAdmin {
		if request.PackageId != 0 {
			inArray := utils.InArray(request.PackageId, packageSlice)
			if !inArray {
				return nil, fmt.Errorf("QuanXianBuZu")
			}
		} else {
			if len(packageSlice) == 0 {
				return nil, fmt.Errorf("QuanXianBuZu")
			}
			delete(condition, "package_id")
			condition["package_id__in"] = packageSlice
		}
	}
	condition["is_deleted"] = 0
	data, total, err := alertModel.GetPageList(&models.Alert{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}
	var alertType []map[string]interface{}
	var AlertRuleType []map[string]interface{}
	var packages []models.Package
	if needReload == 1 {
		alertType = s.AllAlertType(lang)
		AlertRuleType = s.AllAlertRuleType(lang)
		packageService := PackageService{}
		packages = packageService.GetMyAllPackageList(int(adminId))
	}
	return map[string]interface{}{
		"list":            data,
		"total":           total,
		"current_page":    request.Page,
		"alert_type":      alertType,
		"alert_rule_type": AlertRuleType,
		"packages":        packages,
	}, nil
}

// GetAlertById 获取弹窗信息
func (s AlertService) GetAlertById(id int) (*models.Alert, error) {
	alert := &models.Alert{}
	alertModel := models.CreateAlertModel()
	err := alertModel.QueryTable(new(models.Alert)).
		Filter("id", id).
		One(alert)
	if err != nil {
		return nil, err
	}
	return alert, nil
}

// AllAlertType 所有弹窗类型
func (s AlertService) AllAlertType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.DownloadAlert, "name": i18n.Tr(language, "XiaZaiTanChuang")},
		{"id": config.MultiTagAlert, "name": i18n.Tr(language, "DuoBiaoQianTanChuang")},
		{"id": config.FirstRechargeAlert, "name": i18n.Tr(language, "ShouChongTanChuang")},
	}
}

// AllAlertRuleType 所有弹窗规则类型
func (s AlertService) AllAlertRuleType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.Always, "name": i18n.Tr(language, "MeiCiTanChuang")},
		{"id": config.DayOneTime, "name": i18n.Tr(language, "MeiRiZhiTanYiCi")},
		{"id": config.DeviceOnTime, "name": i18n.Tr(language, "ZhiTanYiCi")},
		{"id": config.PeriodOneTime, "name": i18n.Tr(language, "DingShiShiTanYiCi")},
	}
}

type AddAlertRequest struct {
	PackageIds  []int  `json:"package_ids"`  // 包ID
	Name        string `json:"name"`         // 弹窗名称
	Type        int    `json:"type"`         // 弹窗类型
	ContentType int    `json:"content_type"` // 内容类型
	EnTitle     string `json:"en_title"`     // 英语标题
	PtTitle     string `json:"pt_title"`     // 葡萄牙语标题
	EnContent   string `json:"en_content"`   // 英语内容
	PtContent   string `json:"pt_content"`   // 葡萄牙语内容
	Image       string `json:"image"`        // 图片
	Sort        int    `json:"sort"`         // 排序
	AlertRule   int    `json:"alert_rule"`   // 弹窗规则
	AlertHours  int    `json:"alert_hours"`  // 提醒时间
	Status      int    `json:"status"`       // 状态
}

// AddAlert 添加弹窗
func (s AlertService) AddAlert(request AddAlertRequest, adminId int64) error {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	err := s.validateAddAlert(request)
	if err != nil {
		return err
	}
	if len(request.PackageIds) == 0 {
		return fmt.Errorf("PingTaiBiXuan")
	}
	tx, err := models.CreateAlertModel().Begin()
	if err != nil {
		logs.Error("添加打脸事务开启失败，失败原因为：%v", err)
		return err
	}
	for _, packageId := range request.PackageIds {
		if !isRootAdmin && !utils.InArray(packageId, packageSlice) {
			return fmt.Errorf("QuanXianBuZu")
		}
		_, err = tx.Insert(&models.Alert{
			PackageId:   packageId,
			Name:        request.Name,
			Type:        request.Type,
			ContentType: request.ContentType,
			EnTitle:     request.EnTitle,
			PtTitle:     request.PtTitle,
			EnContent:   request.EnContent,
			PtContent:   request.PtContent,
			Image:       request.Image,
			Sort:        request.Sort,
			AlertRule:   request.AlertRule,
			AlertHours:  request.AlertHours,
			Status:      request.Status,
		})
		if err != nil {
			tx.Rollback()
			logs.Error("添加打脸失败，失败原因为：%v", err)
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("添加打脸事务提交失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	for _, packageId := range request.PackageIds {
		err := s.refreshAlert(packageId)
		if err != nil {
			logs.Error("刷新打脸失败，失败原因为：%v", err)
			return fmt.Errorf("ShuaXinHuanCunShiBai")
		}
	}
	return nil
}

// 添加弹窗参数验证
func (s AlertService) validateAddAlert(request AddAlertRequest) error {
	if request.Name == "" {
		return fmt.Errorf("TanChuangMingChengBiTian")
	}
	if request.Type == 0 {
		return fmt.Errorf("TanChuangLeiXingBiTian")
	}
	if request.ContentType == 0 {
		return fmt.Errorf("TanChuangNeiRongLeiXingBiTian")
	}
	if request.AlertRule == 0 {
		return fmt.Errorf("TanChuangGuiZeBiTian")
	}
	return nil
}

type EditAlertRequest struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`         // 弹窗名称
	Type        int    `json:"type"`         // 弹窗类型
	ContentType int    `json:"content_type"` // 内容类型
	EnTitle     string `json:"en_title"`     // 英语标题
	PtTitle     string `json:"pt_title"`     // 葡萄牙语标题
	EnContent   string `json:"en_content"`   // 英语内容
	PtContent   string `json:"pt_content"`   // 葡萄牙语内容
	Image       string `json:"image"`        // 图片
	Sort        int    `json:"sort"`         // 排序
	AlertRule   int    `json:"alert_rule"`   // 弹窗规则
	AlertHours  int    `json:"alert_hours"`  // 提醒时间
	Status      int    `json:"status"`       // 状态
}

// EditAlert 编辑弹窗
func (s AlertService) EditAlert(request EditAlertRequest, adminId int64) error {
	err := s.validateEditAlert(request)
	if err != nil {
		return err
	}
	alert, err := s.GetAlertById(request.Id)
	if err != nil {
		logs.Error("编辑打脸失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(alert.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	fields := []string{
		"name",
		"type",
		"content_type",
		"en_title",
		"pt_title",
		"en_content",
		"pt_content",
		"image",
		"sort",
		"alert_rule",
		"alert_hours",
		"status",
	}
	tx, err := models.CreateAlertModel().Begin()
	if err != nil {
		logs.Error("编辑打脸事务开启失败，失败原因为：%v", err)
		return err
	}
	_, err = tx.Update(&models.Alert{
		Id:          request.Id,
		Name:        request.Name,
		Type:        request.Type,
		ContentType: request.ContentType,
		EnTitle:     request.EnTitle,
		PtTitle:     request.PtTitle,
		EnContent:   request.EnContent,
		PtContent:   request.PtContent,
		Image:       request.Image,
		Sort:        request.Sort,
		AlertRule:   request.AlertRule,
		AlertHours:  request.AlertHours,
		Status:      request.Status,
	}, fields...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("添加打脸事务提交失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	err = s.refreshAlert(alert.PackageId)
	if err != nil {
		logs.Error("刷新打脸失败，失败原因为：%v", err)
		return fmt.Errorf("ShuaXinHuanCunShiBai")
	}
	return nil
}

// 编辑弹窗参数验证
func (s AlertService) validateEditAlert(request EditAlertRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("TanChuangIDBiTian")
	}
	if request.Name == "" {
		return fmt.Errorf("TanChuangMingChengBiTian")
	}
	if request.Type == 0 {
		return fmt.Errorf("TanChuangLeiXingBiTian")
	}
	if request.ContentType == 0 {
		return fmt.Errorf("TanChuangNeiRongLeiXingBiTian")
	}
	if request.AlertRule == 0 {
		return fmt.Errorf("TanChuangGuiZeBiTian")
	}
	return nil
}

// DeleteAlert 删除弹窗
func (s AlertService) DeleteAlert(id int, adminId int64) error {
	if id == 0 {
		return fmt.Errorf("TanChuangIDBiTian")
	}
	alert, err := s.GetAlertById(id)
	if err != nil {
		logs.Error("获取打脸失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(alert.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	tx, err := models.CreateAlertModel().Begin()
	if err != nil {
		logs.Error("删除打脸事务开启失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	_, err = tx.Update(&models.Alert{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("删除打脸事务提交失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	err = s.refreshAlert(alert.PackageId)
	if err != nil {
		logs.Error("刷新打脸失败，失败原因为：%v", err)
		return fmt.Errorf("ShuaXinHuanCunShiBai")
	}
	return nil
}

type ChangeAlertStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangeAlertStatus 修改弹窗状态
func (s AlertService) ChangeAlertStatus(request ChangeAlertStatusRequest, adminId int64) error {
	if request.Id == 0 {
		return fmt.Errorf("TanChuangIDBiTian")
	}
	alert, err := s.GetAlertById(request.Id)
	if err != nil {
		logs.Error("获取打脸失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(alert.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	tx, err := models.CreateAlertModel().Begin()
	if err != nil {
		logs.Error("删除打脸事务开启失败，失败原因为：%v", err)
		return err
	}
	_, err = tx.Update(&models.Alert{
		Id:     request.Id,
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("删除打脸事务提交失败，失败原因为：%v", err)
		return err
	}
	err = s.refreshAlert(alert.PackageId)
	if err != nil {
		logs.Error("刷新打脸失败，失败原因为：%v", err)
		return fmt.Errorf("ShuaXinHuanCunShiBai")
	}
	return nil
}

// 刷新弹窗缓存
func (s AlertService) refreshAlert(packageId int) error {
	now := time.Now().Unix()
	secret := config.System.CurlSecretKey.Value
	params := map[string]interface{}{
		"package_id": strconv.Itoa(packageId),
		"time":       now,
	}
	request, _ := utils.SignMap(params, secret)
	packageService := PackageService{}
	resp := RequestPackageUrlRequest{
		PackageId: packageId,
		Params:    request,
		Url:       "/data/refresh/refreshAlert",
		Type:      int(config.CurlRequestTypePost),
	}
	response, err := packageService.RequestPackageUrl(resp)
	if err != nil {
		logs.Error("刷新弹窗失败：%v", err)
		return fmt.Errorf("ShuXinHuanCunShiBai")
	}
	if !response.OK() {
		logs.Error("刷新弹窗失败：%v", response)
		return fmt.Errorf("ShuXinHuanCunShiBai")
	}
	return nil
}

// RefreshPackageAlert 刷新包打脸缓存
func (s AlertService) RefreshPackageAlert(packageId int) error {
	err := s.refreshAlert(packageId)
	if err != nil {
		logs.Error("刷新浮动图标失败：%v", err)
		return fmt.Errorf(err.Error())
	}
	return nil
}
