package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"fmt"
	"net/url"

	"github.com/beego/i18n"
)

type WithdrawOrderService struct{}

type WithdrawOrderListRequest struct {
	PackageId        int        `form:"package_id"`
	UserId           int        `form:"user_id"`
	RoleId           int64      `form:"role_id"`
	PaymentId        int        `form:"payment_id"`
	PaymentChannelId int        `form:"payment_channel_id"`
	OrderNo          string     `form:"order_no" op:"like"`
	TransactionNo    string     `form:"transaction_no" op:"like"`
	Status           int        `form:"status"`
	Page             int        `form:"page"`
	PageSize         int        `form:"pageSize"`
	RawQuery         url.Values `form:"-"`
}

// WithdrawOrderList 获取提现订单列表
func (s WithdrawOrderService) WithdrawOrderList(request WithdrawOrderListRequest, needReload int, language string, adminId int64) (map[string]interface{}, error) {
	withdrawOrderModel := models.CreateWithdrawOrderModel()
	condition, sort := withdrawOrderModel.BuildCondition(request, "-id") // 默认按类型排序
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
	data, total, err := withdrawOrderModel.GetPageList(&models.WithdrawOrder{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}
	var payment []models.Payment
	var paymentChannel []models.PaymentChannel
	var packages []models.Package
	var status []map[string]interface{}
	if needReload == 1 {
		paymentService := PaymentService{}
		payment, err = paymentService.AllPayment()
		if err != nil {
			return nil, err
		}
		paymentChannel, err = paymentService.AllPaymentChannel()
		if err != nil {
			return nil, err
		}
		packageService := PackageService{}
		packages = packageService.GetMyAllPackageList(int(adminId))
		status = s.AllWithdrawOrderStatus(language)
	}

	list := data.([]models.WithdrawOrder)
	return map[string]interface{}{
		"list":            list,
		"total":           total,
		"packages":        packages,
		"status":          status,
		"payment_channel": paymentChannel,
		"payment":         payment,
	}, nil
}

// AllWithdrawOrderStatus 所有提现订单状态
func (s WithdrawOrderService) AllWithdrawOrderStatus(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.WithdrawOrderStatusCreated, "name": i18n.Tr(language, "YiChuangJian")},
		{"id": config.WithdrawOrderStatusPass, "name": i18n.Tr(language, "TongGuo")},
		{"id": config.WithdrawOrderStatusRejectedBack, "name": i18n.Tr(language, "JuJueTuiHui")},
		{"id": config.WithdrawOrderStatusRejectedFrozen, "name": i18n.Tr(language, "JuJueDongJie")},
		{"id": config.WithdrawOrderStatusProcessing, "name": i18n.Tr(language, "ChuLiZhong")},
		{"id": config.WithdrawOrderStatusFailBack, "name": i18n.Tr(language, "ChuLiShiBaiHuiTui")},
		{"id": config.WithdrawOrderStatusAddWorkReview, "name": i18n.Tr(language, "JiaRuGongDanShenHe")},
		{"id": config.WithdrawOrderStatusWarning, "name": i18n.Tr(language, "FengXianGuaQi")},
		{"id": config.WithdrawOrderStatusSuccess, "name": i18n.Tr(language, "ChengGong")},
	}
}

type GetWithdrawDailyStatResponse struct {
	TotalWithdraw     float64
	TotalTax          float64
	RealAmount        float64
	TotalMen          int
	TotalFakeWithdraw float64
	ApplyWithdraw     float64
}

type GetWithdrawDailyStatSum struct {
	TotalAmount float64
	TotalTax    float64
	RealAmount  float64
	TotalMen    int
	PackageId   int
}

type GetWithdrawDailyStatFakeSum struct {
	TotalAmount float64
	PackageId   int
}

// GetWithdrawDailyStat 按分包日况统计
func (s WithdrawOrderService) GetWithdrawDailyStat(condition map[string]interface{}) (map[int]GetWithdrawDailyStatResponse, error) {
	result := make(map[int]GetWithdrawDailyStatResponse)

	model := models.CreateWithdrawOrderModel()
	qs := model.Where(model.QueryTable(new(models.WithdrawOrder)), condition)

	// 提现申请金额
	var applyWithdraws []GetWithdrawDailyStatFakeSum
	_, err := qs.Filter("status", config.WithdrawOrderStatusCreated).GroupBy("package_id").Aggregate("sum(i_money) as total_amount, package_id").All(&applyWithdraws)
	if nil == err {
		for _, item := range applyWithdraws {
			res := result[item.PackageId]
			res.ApplyWithdraw = item.TotalAmount
			result[item.PackageId] = res
		}
	} else {
		return nil, err
	}

	// 总提现金额
	var totalAmounts []GetWithdrawDailyStatSum
	_, err = qs.Filter("status", config.WithdrawOrderStatusSuccess).Filter("is_fake", 0).GroupBy("package_id").Aggregate("sum(i_money) as total_amount, sum(tax) as total_tax, sum(i_money - tax) as real_amount, count(distinct user_id) as total_men, package_id").All(&totalAmounts)
	if nil == err {
		for _, item := range totalAmounts {
			res := result[item.PackageId]
			res.TotalWithdraw = item.TotalAmount
			res.TotalTax = item.TotalTax
			res.RealAmount = item.RealAmount
			res.TotalMen = item.TotalMen
			result[item.PackageId] = res
		}
	} else {
		return nil, err
	}

	// 假提现金额
	var firstStats []GetWithdrawDailyStatFakeSum
	_, err = qs.Filter("status", config.WithdrawOrderStatusSuccess).Filter("is_fake", 1).GroupBy("package_id").Aggregate("sum(i_money) as total_amount, package_id").
		All(&firstStats)
	if err == nil {
		for _, item := range firstStats {
			res := result[item.PackageId]
			res.TotalFakeWithdraw = item.TotalAmount
			result[item.PackageId] = res
		}
	} else {
		return nil, err
	}
	return result, nil
}

type GetWithdrawDailyMenStatResponse struct {
	TotalMen int
}

// GetWithdrawDailyMenStat 日、月报表统计实时人数
func (s WithdrawOrderService) GetWithdrawDailyMenStat(condition map[string]interface{}) (int, error) {

	model := models.CreateWithdrawOrderModel()
	qs := model.Where(model.QueryTable(new(models.WithdrawOrder)), condition)
	condition["status"] = config.WithdrawOrderStatusSuccess
	condition["is_fake"] = 0

	// 总提现金额
	var totalAmounts []GetWithdrawDailyMenStatResponse
	_, err := qs.Aggregate("count(distinct user_id) as total_men").All(&totalAmounts)
	if nil == err {
		return totalAmounts[0].TotalMen, nil
	} else {
		return 0, err
	}
}
