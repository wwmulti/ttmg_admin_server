package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type DailyStatisticsService struct {
	BaseService
}

type LstRequest struct {
	Time                []string   `form:"time[]" op:"time_range" orm:"time"`
	MergeByDayFlag      int        `form:"merge_by_day_flag" op:"-"`
	MergeByPlatformFlag int        `form:"merge_by_plamform_flag" op:"-"`
	NeedReload          int        `form:"need_reload" op:"-"`
	PackageId           int        `form:"package_id"`
	PackageIds          []int      `form:"package_ids[]" op:"-"`
	Page                int        `form:"page"`
	PageSize            int        `form:"page_size"`
	RawQuery            url.Values `form:"-"`
}

// Lst 日况统计
func (s *DailyStatisticsService) Lst(req LstRequest, userId int) (map[string]interface{}, error) {
	model := models.CreateDailyStatModel()
	condition, sort := model.BuildCondition(req, "-time")
	condition = s.LimitPackageId(condition, req.PackageIds)

	var (
		list    []models.DailyStat
		allList []models.DailyStat
		groupBy []string
	)

	if req.MergeByDayFlag == 0 && req.MergeByPlatformFlag == 0 {
		_, err := model.Where(model.QueryTable(new(models.DailyStat)), condition).OrderBy(sort).All(&allList)
		if nil != err {
			return nil, err
		}
	} else {
		if req.MergeByDayFlag == 1 && req.MergeByPlatformFlag == 1 {
			groupBy = []string{"package_id", "time"}
		} else if req.MergeByPlatformFlag == 1 {
			groupBy = []string{"package_id"}
		} else if req.MergeByDayFlag == 1 {
			groupBy = []string{"time"}
		}

		sumStr := s.GetSumFields()
		_, err := model.Where(model.QueryTable(new(models.DailyStat)), condition).GroupBy(groupBy...).Aggregate(sumStr).OrderBy(sort).All(&allList)
		if nil != err {
			return nil, err
		}
	}
	total := len(allList)
	if total == 0 {
		return nil, nil
	}

	list = utils.SlicePage(req.Page, req.PageSize, allList)

	// 计算总充值、总提现、总利润
	var (
		totalRecharge, totalWithdrawal, totalProfit int64
	)
	for _, v := range allList {
		totalRecharge += v.TotalChargeAmount
		totalWithdrawal += v.TotalWithdraw
		totalProfit += v.TotalProfit
	}

	// 下拉选项
	var packages interface{}
	if req.NeedReload == 1 {
		packages = (&PackageService{}).GetMyAllPackageList(userId)
	}

	return map[string]interface{}{
		"list":           list,
		"total":          total,
		"packages":       packages,
		"total_recharge": totalRecharge,
		"total_withdraw": totalWithdrawal,
		"total_profit":   totalProfit,
	}, nil
}

// GetSumFields 构建SUM字段
func (s *DailyStatisticsService) GetSumFields() string {
	t := reflect.TypeOf(models.DailyStat{})
	var sums []string
	ignoreFields := []string{"id", "time", "update_time", "date", "package_id"}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		column := field.Tag.Get("orm")
		if column == "" || !strings.Contains(column, "column(") {
			continue
		}

		colName := column[strings.Index(column, "(")+1 : strings.Index(column, ")")]
		if slices.Contains(ignoreFields, colName) {
			sums = append(sums, fmt.Sprintf("MAX(%s) as %s", colName, colName))
			continue
		}

		sums = append(sums, fmt.Sprintf("SUM(%s) as %s", colName, colName))
	}
	return strings.Join(sums, ", ")
}

type TotalDailyRequest struct {
	Time       []string   `form:"time[]" op:"time_range" orm:"time"`
	NeedReload int        `form:"need_reload" op:"-"`
	PackageId  int        `form:"package_id"`
	PackageIds []int      `form:"package_ids[]" op:"-"`
	Page       int        `form:"page"`
	PageSize   int        `form:"page_size"`
	RawQuery   url.Values `form:"-"`
}

// TotalDaily 日报表
func (s *DailyStatisticsService) TotalDaily(req TotalDailyRequest, userId int) (map[string]interface{}, error) {
	model := models.CreateDailyStatModel()
	condition, sort := model.BuildCondition(req, "-time")
	condition = s.LimitPackageId(condition, req.PackageIds)
	packageIds := req.PackageIds
	if req.PackageId > 0 {
		packageIds = []int{req.PackageId}
	}

	var data models.DailyStat

	// 单个包单天统计
	if req.Time[0] != "" && req.Time[0] == req.Time[1] && req.PackageId > 0 {
		err := model.Where(model.QueryTable(new(models.DailyStat)), condition).OrderBy(sort).One(&data)
		if nil != err {
			return nil, err
		}
	} else {
		var list []models.DailyStat
		sumStr := s.GetSumFields()
		_, err := model.Where(model.QueryTable(new(models.DailyStat)), condition).Aggregate(sumStr).OrderBy(sort).All(&list)
		if nil != err {
			return nil, err
		}
		data = list[0]

		// 计算活动人数去重
		uniqueMens := s.CalUniqueActiveMen(condition["time__gte"].(int64), condition["time__lt"].(int64), packageIds)
		data.TotalChargeMen = uniqueMens.TotalChargeMen
		data.TotalWithdrawMen = uniqueMens.TotalWithdrawMen
		data.PlatformSendMen = uniqueMens.PlatformSendMen
		data.TotalRechargeSendMen = uniqueMens.TotalRechargeSendMen
		data.TotalSignSendMen = uniqueMens.TotalSignSendMen
		data.TotalVipSendMen = uniqueMens.TotalVipSendMen
		data.TotalSpinSendMen = uniqueMens.TotalSpinSendMen
		data.TotalHelpSendMen = uniqueMens.TotalHelpSendMen
		data.TotalInterestSendMen = uniqueMens.TotalInterestSendMen
	}

	// 计算用户总数
	totalUserMens := s.CalTotalUser(packageIds)
	data.TotalNormalUserMen = totalUserMens.TotalNormalUserMen
	data.TotalDisabledUserMen = totalUserMens.TotalDisabledUserMen

	// 下拉选项
	var packages interface{}
	if req.NeedReload == 1 {
		packages, _ = (&PackageService{}).AllPackage()
	}

	return map[string]interface{}{
		"packages": packages,
		"data":     data,
	}, nil
}

type TotalMonthRequest struct {
	Time       []string   `form:"time[]" op:"-"`
	NeedReload int        `form:"need_reload" op:"-"`
	PackageId  int        `form:"package_id" op:"-"`
	PackageIds []int      `form:"package_ids[]" op:"-"`
	Page       int        `form:"page"`
	PageSize   int        `form:"page_size"`
	RawQuery   url.Values `form:"-"`
}

// TotalMonth 月报表
func (s *DailyStatisticsService) TotalMonth(req TotalMonthRequest, userId int) (map[string]interface{}, error) {
	model := models.CreateDailyStatModel()
	condition, sort := model.BuildCondition(req, "-time")
	condition = s.LimitPackageId(condition, req.PackageIds)
	packageIds := req.PackageIds
	if req.PackageId > 0 {
		packageIds = []int{req.PackageId}
	}

	// 计算月份时间范围
	startT, err := time.ParseInLocation("2006-01", req.Time[0], time.Local)
	if err != nil {
		return nil, fmt.Errorf("start month format error: %v", err)
	}
	startTime := startT.Unix()

	endT, err := time.ParseInLocation("2006-01", req.Time[1], time.Local)
	if err != nil {
		return nil, fmt.Errorf("end month format error: %v", err)
	}
	endTime := endT.AddDate(0, 1, 0).Unix()
	condition["time__gte"] = startTime
	condition["time__lt"] = endTime

	var data models.DailyStat
	var list []models.DailyStat
	sumStr := s.GetSumFields()
	_, err = model.Where(model.QueryTable(new(models.DailyStat)), condition).Aggregate(sumStr).OrderBy(sort).All(&list)
	if nil != err {
		return nil, err
	}
	data = list[0]

	// 计算活动人数去重
	uniqueMens := s.CalUniqueActiveMen(startTime, endTime, packageIds)
	data.TotalChargeMen = uniqueMens.TotalChargeMen
	data.TotalWithdrawMen = uniqueMens.TotalWithdrawMen
	data.PlatformSendMen = uniqueMens.PlatformSendMen
	data.TotalRechargeSendMen = uniqueMens.TotalRechargeSendMen
	data.TotalSignSendMen = uniqueMens.TotalSignSendMen
	data.TotalVipSendMen = uniqueMens.TotalVipSendMen
	data.TotalSpinSendMen = uniqueMens.TotalSpinSendMen
	data.TotalHelpSendMen = uniqueMens.TotalHelpSendMen
	data.TotalInterestSendMen = uniqueMens.TotalInterestSendMen

	// 计算用户总数
	totalUserMens := s.CalTotalUser(packageIds)
	data.TotalNormalUserMen = totalUserMens.TotalNormalUserMen
	data.TotalDisabledUserMen = totalUserMens.TotalDisabledUserMen

	// 下拉选项
	var packages interface{}
	if req.NeedReload == 1 {
		packages, _ = (&PackageService{}).AllPackage()
	}

	return map[string]interface{}{
		"packages": packages,
		"data":     data,
	}, nil
}

// CalUniqueActiveMen 实时计算活动人数去重
func (s *DailyStatisticsService) CalUniqueActiveMen(startTime int64, endTime int64, packageIds []int) models.DailyStat {

	var wg sync.WaitGroup
	var mu sync.Mutex
	var data models.DailyStat

	// 构建条件
	bulidCond := func(column string) map[string]interface{} {
		condition := map[string]interface{}{}

		condition[fmt.Sprintf("%s__gte", column)] = startTime
		condition[fmt.Sprintf("%s__lt", column)] = endTime
		if len(packageIds) > 0 {
			condition["package_id__in"] = packageIds
		}
		return condition
	}

	// 充值人数
	wg.Add(1)
	go func() {
		defer wg.Done()
		totalRechargeMen, err := (&RechargeOrderService{}).GetRechargeDailyMenStat(bulidCond("pay_at"))
		if err != nil {
			logs.Error("计算充值去重人数失败:%v", err)
		} else {
			mu.Lock()
			data.TotalChargeMen = totalRechargeMen
			mu.Unlock()
		}
	}()

	// 提现人数
	wg.Add(1)
	go func() {
		defer wg.Done()
		totalWithdrawMen, err := (&WithdrawOrderService{}).GetWithdrawDailyMenStat(bulidCond("update_time"))
		if err != nil {
			logs.Error("计算提现去重人数失败:%v", err)
		} else {
			mu.Lock()
			data.TotalWithdrawMen = totalWithdrawMen
			mu.Unlock()
		}
	}()

	// 注册赠送人数
	// 充值赠送人数
	// 签到人数
	// vip奖励人数
	// 转盘人数
	// 亏损返利人数
	// 利息宝人数
	wg.Add(1)
	go func() {
		defer wg.Done()
		types := []int{
			int(config.BalanceLogTypeRegister),
			int(config.BalanceLogTypeRechargeExtra),
			int(config.BalanceLogTypeActivitySign),
			int(config.BalanceLogTypeActivityVipLv),
			int(config.BalanceLogTypeActivityVipLvDaily),
			int(config.BalanceLogTypeActivityVipLvWeekly),
			int(config.BalanceLogTypeActivityVipLvMonthly),
			int(config.BalanceLogTypeActivityLuckyWheel),
			int(config.BalanceLogTypeActivityReliefDaily),
			int(config.BalanceLogTypeActivityReliefWeekly),
			int(config.BalanceLogTypeActivityReliefMonthly),
			int(config.BalanceLogTypeActivityInterest),
		}

		totalUserBalanceMen, err := (&UserBalanceLogService{}).GetUserBalanceLogDailyMenStat(UserBalanceLogGroupListRequestParams{
			BeginTime:  startTime,
			EndTime:    endTime,
			Types:      types,
			PackageIds: packageIds,
		})
		if err != nil {
			logs.Error("计算余额变动去重人数失败:%v", err)
		} else {
			mu.Lock()
			for _, item := range totalUserBalanceMen {
				switch item.Type {
				case int(config.BalanceLogTypeRegister):
					data.PlatformSendMen += item.TotalMen
				case int(config.BalanceLogTypeRechargeExtra):
					data.TotalRechargeSendMen += item.TotalMen
				case int(config.BalanceLogTypeActivitySign):
					data.TotalSignSendMen += item.TotalMen
				case int(config.BalanceLogTypeActivityVipLv), int(config.BalanceLogTypeActivityVipLvDaily), int(config.BalanceLogTypeActivityVipLvWeekly), int(config.BalanceLogTypeActivityVipLvMonthly):
					data.TotalVipSendMen += item.TotalMen
				case int(config.BalanceLogTypeActivityLuckyWheel):
					data.TotalSpinSendMen += item.TotalMen
				case int(config.BalanceLogTypeActivityReliefDaily), int(config.BalanceLogTypeActivityReliefWeekly), int(config.BalanceLogTypeActivityReliefMonthly):
					data.TotalHelpSendMen += item.TotalMen
				case int(config.BalanceLogTypeActivityInterest):
					data.TotalInterestSendMen += item.TotalMen
				}
			}
			mu.Unlock()
		}
	}()

	wg.Wait()

	return data
}

// CalTotalUser 实时计算用户总数
func (s *DailyStatisticsService) CalTotalUser(packageIds []int) models.DailyStat {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var data models.DailyStat

	// 总正常人数
	wg.Add(1)
	go func() {
		defer wg.Done()
		condNormalUserMen := map[string]interface{}{
			"status": 1,
		}
		if len(packageIds) > 0 {
			condNormalUserMen["package_id__in"] = packageIds
		}
		totalNormalUserMen, err := (&UserService{}).GetUserDailyMenStat(condNormalUserMen)
		if err != nil {
			data.TotalNormalUserMen = 0
			logs.Error("计算正常去重人数失败:%v", err)
		} else {
			mu.Lock()
			data.TotalNormalUserMen = totalNormalUserMen
			mu.Unlock()
		}
	}()

	// 总禁用人数
	wg.Add(1)
	go func() {
		defer wg.Done()
		condDisableUserMen := map[string]interface{}{
			"status": 0,
		}
		if len(packageIds) > 0 {
			condDisableUserMen["package_id__in"] = packageIds
		}
		totalDisabledUserMen, err := (&UserService{}).GetUserDailyMenStat(condDisableUserMen)
		if err != nil {
			logs.Error("计算禁用去重人数失败:%v", err)
		} else {
			mu.Lock()
			data.TotalDisabledUserMen = totalDisabledUserMen
			mu.Unlock()
		}
	}()

	wg.Wait()

	return data
}
