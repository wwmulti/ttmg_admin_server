package services

import (
	"api/config"
	"api/models"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
	"regexp"
	"time"
)

type ActiveReliefService struct{}

type ActiveReliefRuleRequest struct {
	ActiveId int `json:"active_id"` // 活动ID
}

// ActiveReliefRule 救济活动规则
func (s *ActiveReliefService) ActiveReliefRule(request ActiveReliefRuleRequest) (*models.ActiveReliefRule, error) {
	reliefRuleModel := models.CreateActiveReliefRuleModel()
	rule := &models.ActiveReliefRule{}
	err := reliefRuleModel.QueryTable(models.ActiveReliefRule{}).
		Filter("active_id", request.ActiveId).
		One(rule)
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("WeiZhiDeCuoWu")
		}

	}
	return rule, nil
}

// AllActiveReliefRuleCycleType 所有奖励类型
func (s *ActiveReliefService) AllActiveReliefRuleCycleType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.ReliefCycleTypeDay, "name": i18n.Tr(language, "Ri")},
		{"id": config.ReliefCycleTypeWeek, "name": i18n.Tr(language, "Zhou")},
		{"id": config.ReliefCycleTypeMonth, "name": i18n.Tr(language, "Yue")},
	}
}

type AddActiveReliefRuleRequest struct {
	ActiveId int    `json:"active_id" orm:"column(active_id)"` // 活动id
	Cycle    int    `json:"cycle" orm:"column(cycle)"`         // 统计周期：1-按日，2-按周，3-按月
	OpenDay  int    `json:"open_day" orm:"column(open_day)"`   // 开放领取日期（周：1-7；月：1-31；日：固定 1）
	OpenTime string `json:"open_time" orm:"column(open_time)"` // 开放领取时间点（如 "00:00:00"）
	IsRepeat int    `json:"is_repeat" orm:"column(is_repeat)"` // 是否重复：0-否，1-是
	Status   int    `json:"status" orm:"column(status)"`       // 状态：0-关闭，1-开启
}

// AddActiveReliefRule 添加救济活动规则
func (s *ActiveReliefService) AddActiveReliefRule(request AddActiveReliefRuleRequest) error {
	if err := s.VerifyActiveReliefRule(request); err != nil {
		return err
	}

	_, err := models.CreateActiveReliefRuleModel().Insert(&models.ActiveReliefRule{
		ActiveId: request.ActiveId,
		Cycle:    request.Cycle,
		OpenDay:  request.OpenDay,
		OpenTime: request.OpenTime,
		IsRepeat: request.IsRepeat,
		Ctime:    time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}

	return nil
}

// EditActiveReliefRule 编辑救济活动规则
func (s *ActiveReliefService) EditActiveReliefRule(id int, request AddActiveReliefRuleRequest) error {
	if err := s.VerifyActiveReliefRule(request); err != nil {
		return err
	}
	fields := []string{
		"cycle",
		"open_day",
		"open_time",
		"is_repeat",
	}
	affected, err := models.CreateActiveReliefRuleModel().Update(&models.ActiveReliefRule{
		Id:       id,
		Cycle:    request.Cycle,
		OpenDay:  request.OpenDay,
		OpenTime: request.OpenTime,
		IsRepeat: request.IsRepeat,
	}, fields...)
	if err != nil || affected == 0 {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// VerifyActiveReliefRule 验证救济活动规则
func (s *ActiveReliefService) VerifyActiveReliefRule(request AddActiveReliefRuleRequest) error {
	if request.ActiveId <= 0 {
		return fmt.Errorf("HuoDongBuCunZai")
	}

	if request.Cycle < 1 || request.Cycle > 3 {
		return fmt.Errorf("ZhouQiXuanZeCuoWu")
	}

	switch request.Cycle {
	case 1:
		if request.OpenDay != 1 {
			return fmt.Errorf("GuDingZhiWeiYi")
		}
	case 2:
		if request.OpenDay > 7 || request.OpenDay < 1 {
			return fmt.Errorf("ZhiBiXuZaiYiDaoQiDeQuJianNei")
		}
	case 3:
		// 获取下个月的天数
		// 比如现在是3月，下下个月是5月，5月0日就是4月30日
		now := time.Now()
		firstDayOfNextNextMonth := time.Date(now.Year(), now.Month()+2, 0, 0, 0, 0, 0, now.Location())
		day := firstDayOfNextNextMonth.Day()
		if request.OpenDay > day || request.OpenDay < 1 {
			return fmt.Errorf("RiQiBiXuManZuXiaGeYueDeTianShuFanWei")
		}
	}

	reg := regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d:[0-5]\d$`)
	if !reg.MatchString(request.OpenTime) {
		return fmt.Errorf("ShiJianGeShiCuoWu")
	}

	return nil
}

// DeleteActiveReliefRule 删除救济活动规则
func (s *ActiveReliefService) DeleteActiveReliefRule(id int) error {
	affected, err := models.CreateActiveReliefRuleModel().QueryTable(&models.ActiveReliefRule{}).Filter("id", id).Delete()
	if err != nil || affected == 0 {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

type AllActiveReliefRewardRequest struct {
	ActiveId int `json:"active_id"`
}

// AllActiveReliefReward 救济活动奖励列表
func (s *ActiveReliefService) AllActiveReliefReward(request AllActiveReliefRewardRequest) ([]models.ActiveReliefRewards, error) {
	reliefRewardsModel := models.CreateActiveReliefRewardsModel()
	var rewards []models.ActiveReliefRewards
	_, err := reliefRewardsModel.QueryTable(models.ActiveReliefRewards{}).
		Filter("active_id", request.ActiveId).
		All(&rewards)
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("WeiZhiDeCuoWu")
		}

	}
	return rewards, nil
}

type AddActiveReliefRewardRequest struct {
	ActiveId      int     `json:"active_id" orm:"column(active_id)"`           // 活动id
	Amount        float64 `json:"amount" orm:"column(amount)"`                 // 亏损金额
	RebatePercent float64 `json:"rebate_percent" orm:"column(rebate_percent)"` // 返利比例（百分数）
	Status        int     `json:"status" orm:"column(status)"`                 // 状态,0-关闭，1-开启
}

// AddActiveReliefReward 添加救济活动奖励
func (s *ActiveReliefService) AddActiveReliefReward(request AddActiveReliefRewardRequest) error {
	rule, err := s.getAndValidateRule(request.ActiveId)
	if err != nil {
		return err
	}
	_, err = models.CreateActiveReliefRewardsModel().Insert(&models.ActiveReliefRewards{
		ActiveId:      request.ActiveId,
		RuleId:        rule.Id,
		Amount:        request.Amount,
		RebatePercent: request.RebatePercent,
		Status:        request.Status,
		Ctime:         time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// EditActiveReliefReward 编辑救济活动奖励
func (s *ActiveReliefService) EditActiveReliefReward(id int, request AddActiveReliefRewardRequest) error {
	fields := []string{
		"active_id",
		"amount",
		"rebate_percent",
		"status",
	}
	_, err := models.CreateActiveReliefRewardsModel().Update(&models.ActiveReliefRewards{
		Id:            id,
		ActiveId:      request.ActiveId,
		Amount:        request.Amount,
		RebatePercent: request.RebatePercent,
		Status:        request.Status,
	}, fields...)
	if err != nil {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// DeleteActiveReliefReward 删除救济活动奖励
func (s *ActiveReliefService) DeleteActiveReliefReward(id int) error {
	affected, err := models.CreateActiveReliefRewardsModel().QueryTable(&models.ActiveReliefRewards{}).Filter("id", id).Delete()
	if err != nil || affected == 0 {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

type ChangeActiveReliefRewardStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangeActiveReliefRewardStatus 修改救济活动奖励状态
func (s *ActiveReliefService) ChangeActiveReliefRewardStatus(request ChangeActiveReliefRewardStatusRequest) error {
	taskModel := models.CreateActiveReliefRewardsModel()
	if request.Id == 0 {
		return fmt.Errorf("JiangLiIDBiTian")
	}
	_, err := taskModel.Update(&models.ActiveReliefRewards{
		Id:     request.Id,
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}
	return nil
}

// 获取规则
func (s *ActiveReliefService) getAndValidateRule(activityId int) (*models.ActiveReliefRule, error) {
	activityExists := models.CreateActivesModel().QueryTable(new(models.Actives)).
		Filter("id", activityId).
		Filter("is_deleted", 0).
		Exist()
	if !activityExists {
		logs.Error("救济活动不存在")
		return nil, fmt.Errorf("JiuJiHuoDongBuCunZai")
	}
	rule := &models.ActiveReliefRule{}
	err := models.CreateActiveReliefRuleModel().QueryTable(new(models.ActiveReliefRule)).
		Filter("active", activityId).
		One(rule)

	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, fmt.Errorf("JiuJiHuoDongBuCunZai")
		} else {
			return nil, fmt.Errorf("WeiZhiDeCuoWu")
		}
	}
	return rule, nil
}
