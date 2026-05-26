package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"fmt"
	"net/url"

	"github.com/beego/i18n"
)

type RechargeOrderService struct{}

type RechargeOrderListRequest struct {
	PackageId        int        `form:"package_id"`
	UserId           int        `form:"user_id"`
	PaymentId        int        `form:"payment_id"`
	PaymentChannelId int        `form:"payment_channel_id"`
	OrderNumber      string     `form:"order_number" op:"like"`
	OutOrderNumber   string     `form:"out_order_number" op:"like"`
	Status           int        `form:"status"`
	PayAt            []string   `form:"pay_at[]" op:"time_range" orm:"pay_at"`
	CreatedAt        []string   `form:"created_at[]" op:"time_range" orm:"created_at"`
	Page             int        `form:"page"`
	PageSize         int        `form:"pageSize"`
	RawQuery         url.Values `form:"-"`
}

// RechargeOrderList 获取充值订单列表
func (s RechargeOrderService) RechargeOrderList(request RechargeOrderListRequest, needReload int, language string, adminId int64) (map[string]interface{}, error) {
	rechargeOrderModel := models.CreateRechargeOrderModel()
	condition, sort := rechargeOrderModel.BuildCondition(request, "-id") // 默认按类型排序
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
	data, total, err := rechargeOrderModel.GetPageList(&models.RechargeOrder{}, condition, request.Page, request.PageSize, sort)
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
		status = s.AllRechargeOrderStatus(language)
	}

	list := data.([]models.RechargeOrder)
	return map[string]interface{}{
		"list":            list,
		"total":           total,
		"packages":        packages,
		"status":          status,
		"payment_channel": paymentChannel,
		"payment":         payment,
	}, nil
}

// AllRechargeOrderStatus 所有支付订单状态
func (s RechargeOrderService) AllRechargeOrderStatus(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.RechargeOrderStatusUnpaid, "name": i18n.Tr(language, "WeiZhiFu")},
		{"id": config.RechargeOrderStatusPaid, "name": i18n.Tr(language, "YiZhiFu")},
		{"id": config.RechargeOrderStatusNotified, "name": i18n.Tr(language, "Notified")},
	}
}

type GetRechargeDailyStatSum struct {
	TotalAmount float64
	Men         int
	PackageId   int
}

type GetRechargeDailyStatResponse struct {
	TotalChargeAmount float64
	TotalChargeMen    int
	FirstChargeAmount float64
	FirstChargeMen    int
}

// GetRechargeDailyStat 按包分组取充值总额，人数,首充总额和首充人数
func (s RechargeOrderService) GetRechargeDailyStat(condition map[string]interface{}) (map[int]GetRechargeDailyStatResponse, error) {
	result := make(map[int]GetRechargeDailyStatResponse)

	model := models.CreateRechargeOrderModel()
	condition["status"] = config.RechargeOrderStatusPaid
	qs := model.Where(model.QueryTable(new(models.RechargeOrder)), condition)

	// 总充值
	var totalAmounts []GetRechargeDailyStatSum
	_, err := qs.GroupBy("package_id").Aggregate("sum(amount) as total_amount, count(distinct user_id) as men, package_id").All(&totalAmounts)
	if nil == err {
		for _, item := range totalAmounts {
			res := result[item.PackageId]
			res.TotalChargeAmount = item.TotalAmount
			res.TotalChargeMen = item.Men
			result[item.PackageId] = res
		}
	} else {
		return nil, err
	}

	// 首充
	var firstStats []GetRechargeDailyStatSum
	_, err = qs.Filter("is_first", 1).GroupBy("package_id").Aggregate("sum(amount) as total_amount, count(distinct user_id) as men, package_id").All(&firstStats)
	if err == nil {
		for _, item := range firstStats {
			res := result[item.PackageId]
			res.FirstChargeAmount = item.TotalAmount
			res.FirstChargeMen = item.Men
			result[item.PackageId] = res
		}
	} else {
		return nil, err
	}
	return result, nil
}

type GetRechargeDailyMenStatResponse struct {
	TotalMen int
}

// GetRechargeDailyMenStat 日、月报表统计实时人数
func (s RechargeOrderService) GetRechargeDailyMenStat(condition map[string]interface{}) (int, error) {

	model := models.CreateRechargeOrderModel()
	qs := model.Where(model.QueryTable(new(models.RechargeOrder)), condition)
	condition["status"] = config.RechargeOrderStatusPaid

	// 总提现金额
	var totalAmounts []GetRechargeDailyMenStatResponse
	_, err := qs.Aggregate("count(distinct user_id) as total_men").All(&totalAmounts)
	if nil == err {
		return totalAmounts[0].TotalMen, nil
	} else {
		return 0, err
	}
}
