package services

import (
	"api/models"
	"api/utils"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"net/url"
)

type PaymentService struct{}

type PaymentTypeListRequest struct {
	Name     string     `form:"name" op:"like"` // 支付类型名称
	Status   int        `form:"status"`         // 状态
	Page     int        `form:"page"`           // 页码
	PageSize int        `form:"page_size"`      // 每页数量
	RawQuery url.Values `form:"-"`
}

// PaymentTypeList 支付类型列表
func (s PaymentService) PaymentTypeList(request PaymentTypeListRequest) (map[string]interface{}, error) {
	paymentTypeModel := models.CreatePaymentTypeModel()
	condition, sort := paymentTypeModel.BuildCondition(request, "-id")
	condition["is_deleted"] = 0
	data, total, err := paymentTypeModel.GetPageList(&models.PaymentType{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}

	return map[string]interface{}{
		"list":         data,
		"total":        total,
		"current_page": request.Page,
	}, nil
}

// AllPaymentType 获取所有支付类型
func (s PaymentService) AllPaymentType() ([]models.PaymentType, error) {
	var paymentType []models.PaymentType
	paymentTypeModel := models.CreatePaymentTypeModel()
	_, err := paymentTypeModel.QueryTable(new(models.PaymentType)).
		Filter("is_deleted", 0).
		All(&paymentType)
	if err != nil {
		logs.Error("获取所有支付类型报错:%v:", err)
		return nil, fmt.Errorf("WeiZhiDeCuoWu")
	}
	return paymentType, nil
}

type AddPaymentTypeRequest struct {
	Name   string `json:"name"`   // 支付类型名称
	Sort   int    `json:"sort"`   // 排序
	Status int    `json:"status"` // 状态
}

// AddPaymentType 添加支付类型
func (s PaymentService) AddPaymentType(request AddPaymentTypeRequest) error {
	err := s.validateAddPaymentType(request)
	if err != nil {
		return err
	}
	paymentTypeModel := models.CreatePaymentTypeModel()
	_, err = paymentTypeModel.Insert(&models.PaymentType{
		Name:   request.Name,
		Sort:   request.Sort,
		Status: request.Status,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加支付类型参数验证
func (s PaymentService) validateAddPaymentType(request AddPaymentTypeRequest) error {
	if request.Name == "" {
		return fmt.Errorf("ZhiFuLeiXingMingChengBiTian")
	}
	return nil
}

type EditPaymentTypeRequest struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`   // 支付类型名称
	Sort   int    `json:"sort"`   // 排序
	Status int    `json:"status"` // 状态
}

// EditPaymentType 编辑支付类型
func (s PaymentService) EditPaymentType(request EditPaymentTypeRequest) error {
	err := s.validateEditPaymentType(request)
	if err != nil {
		return err
	}
	paymentTypeModel := models.CreatePaymentTypeModel()
	fields := []string{
		"name",
		"sort",
		"status",
	}
	_, err = paymentTypeModel.Update(&models.PaymentType{
		Id:     request.Id,
		Name:   request.Name,
		Sort:   request.Sort,
		Status: request.Status,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑支付类型参数验证
func (s PaymentService) validateEditPaymentType(request EditPaymentTypeRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("ZhiFuLeiXingIDBiTian")
	}
	if request.Name == "" {
		return fmt.Errorf("ZhiFuLeiXingMingChengBiTian")
	}
	return nil
}

type ChangePaymentTypeStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangePaymentTypeStatus 修改支付类型状态
func (s PaymentService) ChangePaymentTypeStatus(request ChangePaymentTypeStatusRequest) error {
	model := models.CreatePaymentTypeModel()
	if request.Id == 0 {
		return fmt.Errorf("ZhiFuLeiXingIDBiTian")
	}
	_, err := model.Update(&models.PaymentType{
		Id:     request.Id,
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}
	return nil
}

type PaymentListRequest struct {
	PayCode  string     `form:"pay_code" op:"like"` // 支付编码
	Status   int        `form:"status"`             // 状态
	Page     int        `form:"page"`               // 页码
	PageSize int        `form:"page_size"`          // 每页数量
	RawQuery url.Values `form:"-"`
}

// PaymentList 支付列表
func (s PaymentService) PaymentList(request PaymentListRequest) (map[string]interface{}, error) {
	paymentTypeModel := models.CreatePaymentModel()
	condition, sort := paymentTypeModel.BuildCondition(request, "-id")
	condition["is_deleted"] = 0
	data, total, err := paymentTypeModel.GetPageList(&models.Payment{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}

	return map[string]interface{}{
		"list":         data,
		"total":        total,
		"current_page": request.Page,
	}, nil
}

// AllPayment 获取所有支付
func (s PaymentService) AllPayment() ([]models.Payment, error) {
	var payment []models.Payment
	paymentModel := models.CreatePaymentModel()
	_, err := paymentModel.QueryTable(new(models.Payment)).
		Filter("is_deleted", 0).
		All(&payment)
	if err != nil {
		logs.Error("获取所有支付报错:%v:", err)
		return nil, fmt.Errorf("WeiZhiDeCuoWu")
	}
	return payment, nil
}

type AddPaymentRequest struct {
	PayCode        string `json:"pay_code"`        // 支付编码
	MerchantConfig string `json:"merchant_config"` // 商户配置
	Logo           string `json:"logo"`            // 支付图标
	Remark         string `json:"remark"`          // 备注
	Status         int    `json:"status"`          // 状态
}

// AddPayment 添加支付
func (s PaymentService) AddPayment(request AddPaymentRequest) error {
	paymentModel := models.CreatePaymentModel()
	err := s.validateAddPayment(request)
	if err != nil {
		return err
	}
	_, err = paymentModel.Insert(&models.Payment{
		PayCode:        request.PayCode,
		MerchantConfig: request.MerchantConfig,
		Logo:           request.Logo,
		Remark:         request.Remark,
		Status:         request.Status,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加支付参数验证
func (s PaymentService) validateAddPayment(request AddPaymentRequest) error {
	if request.PayCode == "" {
		return fmt.Errorf("ZhiFuBianMaBiTian")
	}
	if request.Logo == "" {
		return fmt.Errorf("ZhiFuTuBiaoBiTian")
	}
	return nil
}

type EditPaymentRequest struct {
	Id             int    `json:"id"`
	PayCode        string `json:"pay_code"`        // 支付编码
	MerchantConfig string `json:"merchant_config"` // 商户配置
	Logo           string `json:"logo"`            // 支付图标
	Remark         string `json:"remark"`          // 备注
	Status         int    `json:"status"`          // 状态
}

// EditPayment 编辑支付
func (s PaymentService) EditPayment(request EditPaymentRequest) error {
	err := s.validateEditPayment(request)
	if err != nil {
		return err
	}
	paymentModel := models.CreatePaymentModel()
	fields := []string{
		"pay_code",
		"merchant_config",
		"logo",
		"remark",
		"status",
	}
	_, err = paymentModel.Update(&models.Payment{
		Id:             request.Id,
		PayCode:        request.PayCode,
		MerchantConfig: request.MerchantConfig,
		Logo:           request.Logo,
		Remark:         request.Remark,
		Status:         request.Status,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑支付参数验证
func (s PaymentService) validateEditPayment(request EditPaymentRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("ZhiFuIDBiTian")
	}
	if request.PayCode == "" {
		return fmt.Errorf("ZhiFuBianMaBiTian")
	}
	if request.Logo == "" {
		return fmt.Errorf("ZhiFuTuBiaoBiTian")
	}
	return nil
}

type ChangePaymentStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangePaymentStatus 修改支付状态
func (s PaymentService) ChangePaymentStatus(request ChangePaymentStatusRequest) error {
	model := models.CreatePaymentModel()
	if request.Id == 0 {
		return fmt.Errorf("ZhiFuIDBiTian")
	}
	_, err := model.Update(&models.Payment{
		Id:     request.Id,
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}
	return nil
}

type PaymentChannelListRequest struct {
	PackageId     int        `form:"package_id"`      // 包ID
	PaymentId     int        `form:"payment_id"`      // 支付ID
	PaymentTypeId int        `form:"payment_type_id"` // 支付类型ID
	Name          string     `form:"name" op:"like"`  // 渠道名称
	Status        int        `form:"status"`          // 状态
	Page          int        `form:"page"`            // 页码
	PageSize      int        `form:"page_size"`       // 每页数量
	RawQuery      url.Values `form:"-"`
}

// PaymentChannelList 支付列表
func (s PaymentService) PaymentChannelList(request PaymentChannelListRequest, needReload int, adminId int64) (map[string]interface{}, error) {
	paymentTypeModel := models.CreatePaymentChannelModel()
	condition, sort := paymentTypeModel.BuildCondition(request, "-id")
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
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
	condition["channel_type"] = 1
	data, total, err := paymentTypeModel.GetPageList(&models.PaymentChannel{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}
	var payment []models.Payment
	var paymentType []models.PaymentType
	var packages []models.Package
	if needReload == 1 {
		payment, err = s.AllPayment()
		if err != nil {
			logs.Error("获取所有支付报错:%v:", err)
			return nil, fmt.Errorf("WeiZhiDeCuoWu")
		}
		paymentType, err = s.AllPaymentType()
		if err != nil {
			logs.Error("获取所有支付类型报错:%v:", err)
			return nil, fmt.Errorf("WeiZhiDeCuoWu")
		}
		packageService := PackageService{}
		packages = packageService.GetMyAllPackageList(int(adminId))
	}
	return map[string]interface{}{
		"list":         data,
		"total":        total,
		"payment":      payment,
		"payment_type": paymentType,
		"packages":     packages,
		"current_page": request.Page,
	}, nil
}

// GetPaymentChannelById 获取支付渠道信息
func (s PaymentService) GetPaymentChannelById(id int) (*models.PaymentChannel, error) {
	paymentChannel := &models.PaymentChannel{}
	paymentChannelModel := models.CreatePaymentChannelModel()
	err := paymentChannelModel.QueryTable(new(models.PaymentChannel)).
		Filter("id", id).
		One(paymentChannel)
	if err != nil {
		return nil, err
	}
	return paymentChannel, nil
}

// AllPaymentChannel 获取所有支付
func (s PaymentService) AllPaymentChannel() ([]models.PaymentChannel, error) {
	var payment []models.PaymentChannel
	paymentModel := models.CreatePaymentChannelModel()
	_, err := paymentModel.QueryTable(new(models.PaymentChannel)).
		Filter("is_deleted", 0).
		All(&payment)
	if err != nil {
		logs.Error("获取所有支付子渠道报错:%v:", err)
		return nil, fmt.Errorf("WeiZhiDeCuoWu")
	}
	return payment, nil
}

type AddPaymentChannelRequest struct {
	PackageId     int     `json:"package_id"`      // 包ID
	PaymentId     int     `json:"payment_id"`      // 支付ID
	PaymentTypeId int     `json:"payment_type_id"` // 支付类型ID
	Name          string  `json:"name"`            // 渠道名称
	PrizePercent  int     `json:"prize_percent"`   // 赠送比例
	LimitSmall    float64 `json:"limit_small"`     // 最小支付金额
	LimitBig      float64 `json:"limit_big"`       // 最大支付金额
	VipLevel      int     `json:"vip_level"`       // VIP等级
	Tag           string  `json:"tag"`             // 标签文案
	ChannelConfig string  `json:"channel_config"`  // 渠道配置
	Status        int     `json:"status"`          // 状态
}

// AddPaymentChannel 添加支付
func (s PaymentService) AddPaymentChannel(request AddPaymentChannelRequest, adminId int64) error {
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(request.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	paymentModel := models.CreatePaymentChannelModel()
	err := s.validateAddPaymentChannel(request)
	if err != nil {
		return err
	}
	_, err = paymentModel.Insert(&models.PaymentChannel{
		PackageId:     request.PackageId,
		PaymentId:     request.PaymentId,
		PaymentTypeId: request.PaymentTypeId,
		Name:          request.Name,
		PrizePercent:  request.PrizePercent,
		LimitSmall:    request.LimitSmall,
		LimitBig:      request.LimitBig,
		VipLevel:      request.VipLevel,
		ChannelConfig: request.ChannelConfig,
		ChannelType:   1,
		Tag:           request.Tag,
		Status:        request.Status,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加支付参数验证
func (s PaymentService) validateAddPaymentChannel(request AddPaymentChannelRequest) error {
	if request.PaymentId == 0 {
		return fmt.Errorf("ZhiFuIDBiTian")
	}
	if request.PackageId == 0 {
		return fmt.Errorf("PingTaiIDBiTian")
	}
	if request.PaymentTypeId == 0 {
		return fmt.Errorf("ZhiFuLeiXingIDBiTian")
	}
	if request.Name == "" {
		return fmt.Errorf("ZhiFuQuDaoMingChengBiTian")
	}
	return nil
}

type EditPaymentChannelRequest struct {
	Id            int     `json:"id"`
	PackageId     int     `json:"package_id"`      // 包ID
	PaymentId     int     `json:"payment_id"`      // 支付ID
	PaymentTypeId int     `json:"payment_type_id"` // 支付类型ID
	Name          string  `json:"name"`            // 渠道名称
	PrizePercent  int     `json:"prize_percent"`   // 赠送比例
	LimitSmall    float64 `json:"limit_small"`     // 最小支付金额
	LimitBig      float64 `json:"limit_big"`       // 最大支付金额
	VipLevel      int     `json:"vip_level"`       // VIP等级
	Tag           string  `json:"tag"`             // 标签文案
	ChannelConfig string  `json:"channel_config"`  // 渠道配置
	Status        int     `json:"status"`          // 状态
}

// EditPaymentChannel 编辑支付
func (s PaymentService) EditPaymentChannel(request EditPaymentChannelRequest, adminId int64) error {
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(request.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	err := s.validateEditPaymentChannel(request)
	if err != nil {
		return err
	}
	paymentModel := models.CreatePaymentChannelModel()
	fields := []string{
		"package_id",
		"payment_id",
		"payment_type_id",
		"name",
		"prize_percent",
		"limit_small",
		"limit_big",
		"vip_level",
		"channel_config",
		"tag",
		"status",
	}
	_, err = paymentModel.Update(&models.PaymentChannel{
		Id:            request.Id,
		PackageId:     request.PackageId,
		PaymentId:     request.PaymentId,
		PaymentTypeId: request.PaymentTypeId,
		Name:          request.Name,
		PrizePercent:  request.PrizePercent,
		LimitSmall:    request.LimitSmall,
		LimitBig:      request.LimitBig,
		VipLevel:      request.VipLevel,
		ChannelConfig: request.ChannelConfig,
		Tag:           request.Tag,
		Status:        request.Status,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑支付参数验证
func (s PaymentService) validateEditPaymentChannel(request EditPaymentChannelRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("ZhiFuQuDaoIDBiTian")
	}
	if request.PaymentId == 0 {
		return fmt.Errorf("ZhiFuIDBiTian")
	}
	if request.PackageId == 0 {
		return fmt.Errorf("PingTaiIDBiTian")
	}
	if request.PaymentTypeId == 0 {
		return fmt.Errorf("ZhiFuLeiXingIDBiTian")
	}
	if request.Name == "" {
		return fmt.Errorf("ZhiFuQuDaoMingChengBiTian")
	}
	return nil
}

type ChangePaymentChannelStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangePaymentChannelStatus 修改支付状态
func (s PaymentService) ChangePaymentChannelStatus(request ChangePaymentChannelStatusRequest, adminId int64) error {
	paymentChannel, err := s.GetPaymentChannelById(request.Id)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(paymentChannel.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	model := models.CreatePaymentChannelModel()
	if request.Id == 0 {
		return fmt.Errorf("ZhiFuQuDaoIDBiTian")
	}
	_, err = model.Update(&models.PaymentChannel{
		Id:     request.Id,
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}
	return nil
}

// AllPaymentChannelQuick 快捷支付列表
func (s PaymentService) AllPaymentChannelQuick(id int, adminId int64) ([]models.PaymentChannelQuick, error) {
	paymentChannel, err := s.GetPaymentChannelById(id)
	if err != nil {
		return nil, err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(paymentChannel.PackageId, adminId)
	if !hasAuthority {
		return nil, fmt.Errorf("QuanXianBuZu")
	}
	paymentTypeModel := models.CreatePaymentChannelQuickModel()
	var data []models.PaymentChannelQuick
	_, err = paymentTypeModel.QueryTable(new(models.PaymentChannelQuick)).
		Filter("payment_channel_id", id).
		All(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetPaymentChannelQuickById 获取支付渠道快捷信息
func (s PaymentService) GetPaymentChannelQuickById(id int) (*models.PaymentChannelQuick, error) {
	paymentChannel := &models.PaymentChannelQuick{}
	paymentChannelModel := models.CreatePaymentChannelQuickModel()
	err := paymentChannelModel.QueryTable(new(models.PaymentChannelQuick)).
		Filter("id", id).
		One(paymentChannel)
	if err != nil {
		return nil, err
	}
	return paymentChannel, nil
}

type AddPaymentChannelQuickRequest struct {
	PaymentChannelId int     `json:"payment_channel_id"` // 渠道ID
	Amount           float64 `json:"amount"`             // 金额
	IsRecommend      int     `json:"is_recommend"`       // 是否推荐
	Sort             int     `json:"sort"`               // 排序
}

// AddPaymentChannelQuick 添加快捷支付
func (s PaymentService) AddPaymentChannelQuick(request AddPaymentChannelQuickRequest, adminId int64) error {
	paymentChannel, err := s.GetPaymentChannelById(request.PaymentChannelId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(paymentChannel.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	paymentModel := models.CreatePaymentChannelQuickModel()
	err = s.validateAddPaymentChannelQuick(request)
	if err != nil {
		return err
	}
	_, err = paymentModel.Insert(&models.PaymentChannelQuick{
		PaymentChannelId: request.PaymentChannelId,
		Amount:           request.Amount,
		IsRecommend:      request.IsRecommend,
		Sort:             request.Sort,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加快捷支付参数验证
func (s PaymentService) validateAddPaymentChannelQuick(request AddPaymentChannelQuickRequest) error {
	if request.PaymentChannelId == 0 {
		return fmt.Errorf("ZhiFuIDBiTian")
	}
	if request.Amount <= 0 {
		return fmt.Errorf("JinEBuNengXiaoYuDengYuLing")
	}
	return nil
}

type EditPaymentChannelQuickRequest struct {
	Id               int     `json:"id"`
	PaymentChannelId int     `json:"payment_channel_id"` // 渠道ID
	Amount           float64 `json:"amount"`             // 金额
	IsRecommend      int     `json:"is_recommend"`       // 是否推荐
	Sort             int     `json:"sort"`               // 排序
}

// EditPaymentChannelQuick 编辑快捷支付
func (s PaymentService) EditPaymentChannelQuick(request EditPaymentChannelQuickRequest, adminId int64) error {
	paymentChannel, err := s.GetPaymentChannelById(request.PaymentChannelId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(paymentChannel.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	err = s.validateEditPaymentChannelQuick(request)
	if err != nil {
		return err
	}
	paymentModel := models.CreatePaymentChannelQuickModel()
	fields := []string{
		"payment_channel_id",
		"amount",
		"is_recommend",
		"sort",
	}
	_, err = paymentModel.Update(&models.PaymentChannelQuick{
		Id:               request.Id,
		PaymentChannelId: request.PaymentChannelId,
		Amount:           request.Amount,
		IsRecommend:      request.IsRecommend,
		Sort:             request.Sort,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑快捷支付参数验证
func (s PaymentService) validateEditPaymentChannelQuick(request EditPaymentChannelQuickRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("IDBiTian")
	}
	if request.PaymentChannelId == 0 {
		return fmt.Errorf("ZhiFuIDBiTian")
	}
	if request.Amount <= 0 {
		return fmt.Errorf("JinEBuNengXiaoYuDengYuLing")
	}
	return nil
}

// DeletePaymentChannelQuick 删除快捷支付
func (s PaymentService) DeletePaymentChannelQuick(id int, adminId int64) error {
	channelQuick, err := s.GetPaymentChannelQuickById(id)
	if err != nil {
		return err
	}
	paymentChannel, err := s.GetPaymentChannelById(channelQuick.PaymentChannelId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(paymentChannel.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	model := models.CreatePaymentChannelQuickModel()
	if id == 0 {
		return fmt.Errorf("LunBoIDBiTian")
	}
	_, err = model.Delete(&models.PaymentChannelQuick{
		Id: id,
	})
	if err != nil {
		return err
	}
	return nil
}

type ChangePaymentChannelQuickRecommendRequest struct {
	Id          int `json:"id"`
	IsRecommend int `json:"is_recommend"`
}

// ChangePaymentChannelQuickRecommend 修改快捷支付推荐
func (s PaymentService) ChangePaymentChannelQuickRecommend(request ChangePaymentChannelQuickRecommendRequest, adminId int64) error {
	channelQuick, err := s.GetPaymentChannelQuickById(request.Id)
	if err != nil {
		return err
	}
	paymentChannel, err := s.GetPaymentChannelById(channelQuick.PaymentChannelId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(paymentChannel.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	model := models.CreatePaymentChannelQuickModel()
	if request.Id == 0 {
		return fmt.Errorf("IDBiTian")
	}
	_, err = model.Update(&models.PaymentChannelQuick{
		Id:          request.Id,
		IsRecommend: request.IsRecommend,
	}, "is_recommend")
	if err != nil {
		return err
	}
	return nil
}
