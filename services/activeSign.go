package services

import (
	"api/config"
	"api/models"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
)

type ActiveSignService struct{}

type ActiveSignRule struct {
	ActivityId int `json:"activity_id"`
}

// ActiveSignRule 签到活动规则列表
func (s ActiveSignService) ActiveSignRule(request ActiveSignRule) (*models.SignRule, error) {
	signRuleModel := models.CreateSignRuleModel()
	rule := &models.SignRule{}
	err := signRuleModel.QueryTable(models.SignRule{}).
		Filter("activity_id", request.ActivityId).
		Filter("is_deleted", 0).
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

// 获取规则
func (s ActiveSignService) getAndValidateRule(activityId int) (*models.SignRule, error) {
	activityExists := models.CreateActivesModel().QueryTable(new(models.Actives)).
		Filter("id", activityId).
		Filter("is_deleted", 0).
		Exist()
	if !activityExists {
		logs.Error("签到活动不存在")
		return nil, fmt.Errorf("QianDaoHuoDongBuCunZai")
	}
	rule := &models.SignRule{}
	err := models.CreateSignRuleModel().QueryTable(new(models.SignRule)).
		Filter("activity_id", activityId).
		Filter("is_deleted", 0).
		One(rule)

	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, fmt.Errorf("QianDaoHuoDongGuiZeBuCunZai")
		} else {
			return nil, fmt.Errorf("WeiZhiDeCuoWu")
		}
	}
	return rule, nil
}

type AddActiveSignRuleRequest struct {
	ActivityId       int `json:"activity_id"`        // 活动id
	Days             int `json:"days"`               // 签到周期
	IsLoop           int `json:"is_loop"`            // 是否循环
	IsInterruptReset int `json:"is_interrupt_reset"` // 是否中断重置
}

// AddActiveSignRule 添加签到活动
func (s ActiveSignService) AddActiveSignRule(request AddActiveSignRuleRequest) error {
	err := s.validateAddActiveSignRule(request)
	if err != nil {
		return err
	}
	signRuleModel := models.CreateSignRuleModel()
	_, err = signRuleModel.Insert(&models.SignRule{
		ActivityId:       request.ActivityId,
		Days:             request.Days,
		IsLoop:           request.IsLoop,
		IsInterruptReset: request.IsInterruptReset,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加签到活动参数验证
func (s ActiveSignService) validateAddActiveSignRule(request AddActiveSignRuleRequest) error {
	if request.ActivityId == 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	activeService := ActiveService{}
	active, err := activeService.GetActiveById(request.ActivityId)
	if err != nil {
		return fmt.Errorf("HuoDongBuCunZai")
	}
	if active.ActiveTypeId != int(config.ActiveTypeSign) {
		return fmt.Errorf("HuoDongLeiXingBuZhengQue")
	}
	if request.Days <= 0 {
		return fmt.Errorf("QianDaoTianShuBiXuDaYuLing")
	}
	return nil
}

type EditActiveSignRuleRequest struct {
	Id               int `json:"id"`                 // 签到活动ID
	ActivityId       int `json:"activity_id"`        // 活动id
	Days             int `json:"days"`               // 签到周期
	IsLoop           int `json:"is_loop"`            // 是否循环
	IsInterruptReset int `json:"is_interrupt_reset"` // 是否中断重置
}

// EditActiveSignRule 编辑签到活动规则
func (s ActiveSignService) EditActiveSignRule(request EditActiveSignRuleRequest) error {
	err := s.validateEditActiveSignRule(request)
	if err != nil {
		return err
	}
	fields := []string{
		"activity_id", "days", "is_loop", "is_interrupt_reset",
	}
	signRuleModel := models.CreateSignRuleModel()
	_, err = signRuleModel.Update(&models.SignRule{
		Id:               request.Id,
		ActivityId:       request.ActivityId,
		Days:             request.Days,
		IsLoop:           request.IsLoop,
		IsInterruptReset: request.IsInterruptReset,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑签到活动参数验证
func (s ActiveSignService) validateEditActiveSignRule(request EditActiveSignRuleRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	if request.ActivityId == 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	activeService := ActiveService{}
	active, err := activeService.GetActiveById(request.ActivityId)
	if err != nil {
		return fmt.Errorf("HuoDongBuCunZai")
	}
	if active.ActiveTypeId != int(config.ActiveTypeSign) {
		return fmt.Errorf("HuoDongLeiXingBuZhengQue")
	}
	if request.Days <= 0 {
		return fmt.Errorf("QianDaoTianShuBiXuDaYuLing")
	}
	return nil
}

type AllActiveSignRewardRequest struct {
	ActivityId int `json:"activity_id"`
}

// AllActiveSignReward 签到活动奖励列表
func (s ActiveSignService) AllActiveSignReward(request AllActiveSignRewardRequest) ([]models.SignRuleReward, error) {
	rule, err := s.getAndValidateRule(request.ActivityId)
	if err != nil {
		return nil, err
	}
	signRuleRewardModel := models.CreateSignRuleRewardModel()
	var list []models.SignRuleReward
	_, err = signRuleRewardModel.QueryTable(&models.SignRuleReward{}).
		Filter("sign_rule_id", rule.Id).
		Filter("is_deleted", 0).
		All(&list)
	if nil != err {
		return nil, err
	}

	return list, nil
}

// AllActiveSignRewardType 所有签到奖励类型
func (s ActiveSignService) AllActiveSignRewardType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.SignRewardMoney, "name": i18n.Tr(language, "XianJin")},
		{"id": config.SignRewardIntegral, "name": i18n.Tr(language, "JiFen")},
		{"id": config.SignRewardDiscount, "name": i18n.Tr(language, "YouHuiQuan")},
	}
}

type AddActiveSignRewardRequest struct {
	ActivityId      int     `json:"activity_id"`       // 活动ID
	Day             int     `json:"day"`               // 天数
	RewardTypeId    int     `json:"reward_type_id"`    // 奖励类型ID
	RewardAmount    float64 `json:"reward_amount"`     // 奖励金额
	Icon            string  `json:"icon"`              // 奖励图标
	DayRechargeLine float64 `json:"day_recharge_line"` // 每日充值门槛
	DayRunningLine  float64 `json:"day_running_line"`  // 每日流水门槛
}

// AddActiveSignReward 添加签到活动奖励
func (s ActiveSignService) AddActiveSignReward(request AddActiveSignRewardRequest) error {
	rule, err := s.getAndValidateRule(request.ActivityId)
	if err != nil {
		return err
	}
	err = s.validateAddActiveSignReward(request, rule.Id)
	if err != nil {
		return err
	}
	signRuleRewardModel := models.CreateSignRuleRewardModel()
	_, err = signRuleRewardModel.Insert(&models.SignRuleReward{
		SignRuleId:      rule.Id,
		Day:             request.Day,
		RewardTypeId:    request.RewardTypeId,
		RewardAmount:    request.RewardAmount,
		Icon:            request.Icon,
		DayRechargeLine: request.DayRechargeLine,
		DayRunningLine:  request.DayRunningLine,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加签到活动奖励参数验证
func (s ActiveSignService) validateAddActiveSignReward(request AddActiveSignRewardRequest, signRuleId int) error {
	if request.Day <= 0 {
		return fmt.Errorf("QianDaoTianShuBiXuDaYuLing")
	}
	isExist := models.CreateSignRuleRewardModel().QueryTable(new(models.SignRuleReward)).
		Filter("sign_rule_id", signRuleId).
		Filter("day", request.Day).
		Filter("is_deleted", 0).
		Exist()
	if isExist {
		return fmt.Errorf("QianDaoTianShuYiCunTian")
	}
	return nil
}

type EditActiveSignRewardRequest struct {
	Id              int     `json:"id"`                // 签到活动ID
	ActivityId      int     `json:"activity_id"`       // 活动ID
	Day             int     `json:"day"`               // 天数
	RewardTypeId    int     `json:"reward_type_id"`    // 奖励类型ID
	RewardAmount    float64 `json:"reward_amount"`     // 奖励金额
	Icon            string  `json:"icon"`              // 奖励图标
	DayRechargeLine float64 `json:"day_recharge_line"` // 每日充值门槛
	DayRunningLine  float64 `json:"day_running_line"`  // 每日流水门槛
}

// EditActiveSignReward 编辑签到活动奖励
func (s ActiveSignService) EditActiveSignReward(request EditActiveSignRewardRequest) error {
	rule, err := s.getAndValidateRule(request.ActivityId)
	if err != nil {
		return err
	}
	err = s.validateEditActiveSignReward(request, rule.Id)
	if err != nil {
		return err
	}
	fields := []string{
		"day", "reward_type_id", "reward_amount", "icon", "day_recharge_line", "day_running_line",
	}
	signRuleRewardModel := models.CreateSignRuleRewardModel()
	_, err = signRuleRewardModel.Update(&models.SignRuleReward{
		Id:              request.Id,
		Day:             request.Day,
		RewardTypeId:    request.RewardTypeId,
		RewardAmount:    request.RewardAmount,
		Icon:            request.Icon,
		DayRechargeLine: request.DayRechargeLine,
		DayRunningLine:  request.DayRunningLine,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑签到活动奖励参数验证
func (s ActiveSignService) validateEditActiveSignReward(request EditActiveSignRewardRequest, signRuleId int) error {
	if request.Id == 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	if request.Day <= 0 {
		return fmt.Errorf("QianDaoTianShuBiXuDaYuLing")
	}
	isExist := models.CreateSignRuleRewardModel().QueryTable(new(models.SignRuleReward)).
		Filter("sign_rule_id", signRuleId).
		Filter("id__ne", request.Id).
		Filter("day", request.Day).
		Filter("is_deleted", 0).
		Exist()
	if isExist {
		return fmt.Errorf("QianDaoTianShuYiCunTian")
	}
	return nil
}

// DeleteActiveSignReward 删除转盘奖励
func (s ActiveSignService) DeleteActiveSignReward(id int) error {
	signRuleRewardModel := models.CreateSignRuleRewardModel()
	if id == 0 {
		return fmt.Errorf("JiangLiBuCunZai")
	}
	_, err := signRuleRewardModel.Update(&models.SignRuleReward{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	return nil
}
