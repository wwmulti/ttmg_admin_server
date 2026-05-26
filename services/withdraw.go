package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"fmt"
	"net/url"

	"github.com/beego/beego/v2/core/logs"
)

type WithdrawService struct{}

// AutoOutMoneyConfig 自动出款配置
func (s *WithdrawService) AutoOutMoneyConfig(userId int64, packageId int) (map[string]interface{}, error) {
	result, _ := (&ConfigService{}).AllSystemConfig(AllSystemConfigRequest{
		TypeId: int(config.SystemConfigTypeAutoOutMoney),
	}, 0, userId)
	configsRaw := result["configs"]

	configs, ok := configsRaw.(map[string][]models.Config)
	if !ok {
		logs.Error("配置内容获取失败")
		return nil, fmt.Errorf("WeiZhiDeCuoWu")
	}
	configsMap := make(map[string]models.Config)
	for key, items := range configs {
		for _, item := range items {
			if item.PackageId == packageId {
				configsMap[key] = item
			}
		}
	}

	ret := make(map[string]interface{})
	ret["configs"] = configsMap

	// 包配置
	packageList, _ := (&PackageService{}).AllPackage()
	ret["packages"] = packageList

	// 支付渠道配置
	payChannelList := make([]models.PaymentChannel, 0)
	_, err := models.CreatePaymentChannelModel().
		QueryTable(new(models.PaymentChannel)).
		Filter("status", 1).
		Filter("is_deleted", 0).
		All(&payChannelList)
	if err != nil {
		return nil, err
	}
	fmt.Println(payChannelList)
	ret["pay_channels"] = payChannelList
	return ret, nil
}

type WithdrawChannelListRequest struct {
	PackageId     int        `form:"package_id"`      // 包ID
	PaymentId     int        `form:"payment_id"`      // 支付ID
	PaymentTypeId int        `form:"payment_type_id"` // 支付类型ID
	Name          string     `form:"name" op:"like"`  // 渠道名称
	Status        int        `form:"status"`          // 状态
	Page          int        `form:"page"`            // 页码
	PageSize      int        `form:"page_size"`       // 每页数量
	RawQuery      url.Values `form:"-"`
}

// WithdrawChannelList 提现渠道列表
func (s *WithdrawService) WithdrawChannelList(request WithdrawChannelListRequest, needReload int, adminId int64) (map[string]interface{}, error) {
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
	condition["channel_type"] = 2
	data, total, err := paymentTypeModel.GetPageList(&models.PaymentChannel{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}
	var payment []models.Payment
	var paymentType []models.PaymentType
	var packages []models.Package
	if needReload == 1 {
		paymentService := PaymentService{}
		payment, err = paymentService.AllPayment()
		if err != nil {
			logs.Error("获取所有支付报错:%v:", err)
			return nil, fmt.Errorf("WeiZhiDeCuoWu")
		}
		paymentType, err = paymentService.AllPaymentType()
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

type AddWithdrawChannelRequest struct {
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

// AddWithdrawChannel 添加提现渠道
func (s *WithdrawService) AddWithdrawChannel(request AddWithdrawChannelRequest, adminId int64) error {
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(request.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	paymentModel := models.CreatePaymentChannelModel()
	err := s.validateAddWithdrawChannel(request)
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
		ChannelType:   2,
		Tag:           request.Tag,
		Status:        request.Status,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加支付参数验证
func (s *WithdrawService) validateAddWithdrawChannel(request AddWithdrawChannelRequest) error {
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

type EditWithdrawChannelRequest struct {
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

// EditWithdrawChannel 编辑提现渠道
func (s *WithdrawService) EditWithdrawChannel(request EditWithdrawChannelRequest, adminId int64) error {
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(request.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	err := s.validateEditWithdrawChannel(request)
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
func (s *WithdrawService) validateEditWithdrawChannel(request EditWithdrawChannelRequest) error {
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

type ChangeWithdrawChannelStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangeWithdrawChannelStatus 修改提现渠道状态
func (s *WithdrawService) ChangeWithdrawChannelStatus(request ChangeWithdrawChannelStatusRequest, adminId int64) error {
	withdrawChannel, err := s.GetWithdrawChannelById(request.Id)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(withdrawChannel.PackageId, adminId)
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

// GetWithdrawChannelById 获取提现渠道信息
func (s *WithdrawService) GetWithdrawChannelById(id int) (*models.PaymentChannel, error) {
	withdrawChannel := &models.PaymentChannel{}
	withdrawChannelModel := models.CreatePaymentChannelModel()
	err := withdrawChannelModel.QueryTable(new(models.PaymentChannel)).
		Filter("id", id).
		One(withdrawChannel)
	if err != nil {
		return nil, err
	}
	return withdrawChannel, nil
}
