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

type ActiveFirstRechargeService struct{}

type ActiveFirstRechargeRuleListRequest struct {
	ActivityId int `json:"activity_id"` // 活动ID
}

// ActiveFirstRechargeRule 首充活动规则
func (s ActiveFirstRechargeService) ActiveFirstRechargeRule(request ActiveFirstRechargeRuleListRequest) (*models.FirstRechargeRule, error) {
	firstRechargeRuleModel := models.CreateFirstRechargeRuleModel()
	rule := &models.FirstRechargeRule{}
	err := firstRechargeRuleModel.QueryTable(models.FirstRechargeRule{}).
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

// AllActiveFirstRechargeRuleBillingType 所有首充统计方式
func (s ActiveFirstRechargeService) AllActiveFirstRechargeRuleBillingType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.BillingDaily, "name": i18n.Tr(language, "MeiRi")},
	}
}

// AllActiveFirstRechargeRuleReceiveType 所有首充奖励接收方式
func (s ActiveFirstRechargeService) AllActiveFirstRechargeRuleReceiveType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.ReceiveNextDayMidnight, "name": i18n.Tr(language, "CiRiLingChen")},
	}
}

// AllActiveFirstRechargeRuleTaskType 所有首充任务类型
func (s ActiveFirstRechargeService) AllActiveFirstRechargeRuleTaskType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.TaskBetNumber, "name": i18n.Tr(language, "TouZhuCiShu")},
		{"id": config.TaskBetAmount, "name": i18n.Tr(language, "TouZhuJinE")},
	}
}

// 获取规则
func (s ActiveFirstRechargeService) getAndValidateRule(activityId int) (*models.FirstRechargeRule, error) {
	activityExists := models.CreateActivesModel().QueryTable(new(models.Actives)).
		Filter("id", activityId).
		Filter("is_deleted", 0).
		Exist()
	if !activityExists {
		logs.Error("首充活动不存在")
		return nil, fmt.Errorf("ShouChongHuoDongBuCunZai")
	}
	rule := &models.FirstRechargeRule{}
	err := models.CreateFirstRechargeRuleModel().QueryTable(new(models.FirstRechargeRule)).
		Filter("activity_id", activityId).
		Filter("is_deleted", 0).
		One(rule)

	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, fmt.Errorf("ShouChongHuoDongGuiZeBuCunZai")
		} else {
			return nil, fmt.Errorf("WeiZhiDeCuoWu")
		}
	}
	return rule, nil
}

type AddActiveFirstRechargeRuleRequest struct {
	ActivityId      int     `json:"activity_id"`      // 活动ID
	BillingType     int     `json:"billing_type"`     // 统计方式
	ReceiveType     int     `json:"receive_type"`     // 领取方式
	TaskType        int     `json:"task_type"`        // 领取任务要求 1投注次数 2投注金额
	BetNumber       int     `json:"bet_number"`       // 投注次数
	BetAmount       float64 `json:"bet_amount"`       // 投注金额
	UpdateFrequency int     `json:"update_frequency"` // 更新频率
	RepeatActive    int     `json:"repeat_active"`    // 是否可以重复参与
}

// AddActiveFirstRechargeRule 添加首充活动规则
func (s ActiveFirstRechargeService) AddActiveFirstRechargeRule(request AddActiveFirstRechargeRuleRequest) error {
	err := s.validateAddActiveFirstRechargeRule(request)
	if err != nil {
		return err
	}
	FirstRechargeRuleModel := models.CreateFirstRechargeRuleModel()
	_, err = FirstRechargeRuleModel.Insert(&models.FirstRechargeRule{
		ActivityId:      request.ActivityId,
		BillingType:     request.BillingType,
		ReceiveType:     request.ReceiveType,
		BetNumber:       request.BetNumber,
		BetAmount:       request.BetAmount,
		UpdateFrequency: request.UpdateFrequency,
		TaskType:        request.TaskType,
		RepeatActive:    request.RepeatActive,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加首充活动规则参数验证
func (s ActiveFirstRechargeService) validateAddActiveFirstRechargeRule(request AddActiveFirstRechargeRuleRequest) error {
	if request.ActivityId == 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	activeService := ActiveService{}
	active, err := activeService.GetActiveById(request.ActivityId)
	if err != nil {
		return fmt.Errorf("HuoDongBuCunZai")
	}
	if active.ActiveTypeId != int(config.ActiveTypeFirstRecharge) {
		return fmt.Errorf("HuoDongLeiXingBuZhengQue")
	}
	if request.BillingType <= 0 {
		return fmt.Errorf("TongJiFangShiBiTian")
	}
	if request.ReceiveType <= 0 {
		return fmt.Errorf("LingQuFangShiBiTian")
	}
	if request.TaskType == int(config.TaskBetAmount) && request.BetAmount < 0 {
		return fmt.Errorf("YaZhuJinEBuNengWeiLing")
	}
	if request.TaskType == int(config.TaskBetNumber) && request.BetNumber < 0 {
		return fmt.Errorf("YaZhuCiShuBuNengWeiLing")
	}
	if request.UpdateFrequency <= 0 {
		return fmt.Errorf("GengXinPingLvBuNengXiaoYuDengYuLing")
	}
	return nil
}

type EditActiveFirstRechargeRuleRequest struct {
	Id              int     `json:"id"`               // 首充活动ID
	ActivityId      int     `json:"activity_id"`      // 活动ID
	BillingType     int     `json:"billing_type"`     // 统计方式
	ReceiveType     int     `json:"receive_type"`     // 领取方式
	TaskType        int     `json:"task_type"`        // 领取任务要求 1投注次数 2投注金额
	BetNumber       int     `json:"bet_number"`       // 投注次数
	BetAmount       float64 `json:"bet_amount"`       // 投注金额
	UpdateFrequency int     `json:"update_frequency"` // 更新频率
	RepeatActive    int     `json:"repeat_active"`    // 是否可以重复参与
}

// EditActiveFirstRechargeRule 编辑首充活动规则
func (s ActiveFirstRechargeService) EditActiveFirstRechargeRule(request EditActiveFirstRechargeRuleRequest) error {
	err := s.validateEditActiveFirstRechargeRule(request)
	if err != nil {
		return err
	}
	fields := []string{
		"activity_id", "billing_type", "receive_type", "bet_number", "bet_amount", "update_frequency", "task_type", "repeat_active",
	}
	FirstRechargeRuleModel := models.CreateFirstRechargeRuleModel()
	_, err = FirstRechargeRuleModel.Update(&models.FirstRechargeRule{
		Id:              request.Id,
		ActivityId:      request.ActivityId,
		BillingType:     request.BillingType,
		ReceiveType:     request.ReceiveType,
		BetNumber:       request.BetNumber,
		BetAmount:       request.BetAmount,
		UpdateFrequency: request.UpdateFrequency,
		TaskType:        request.TaskType,
		RepeatActive:    request.RepeatActive,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑首充活动规则参数验证
func (s ActiveFirstRechargeService) validateEditActiveFirstRechargeRule(request EditActiveFirstRechargeRuleRequest) error {
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
	if active.ActiveTypeId != int(config.ActiveTypeFirstRecharge) {
		return fmt.Errorf("HuoDongLeiXingBuZhengQue")
	}
	if request.BillingType <= 0 {
		return fmt.Errorf("TongJiFangShiBiTian")
	}
	if request.ReceiveType <= 0 {
		return fmt.Errorf("LingQuFangShiBiTian")
	}
	if request.TaskType == int(config.TaskBetAmount) && request.BetAmount < 0 {
		return fmt.Errorf("YaZhuJinEBuNengWeiLing")
	}
	if request.TaskType == int(config.TaskBetNumber) && request.BetNumber < 0 {
		return fmt.Errorf("YaZhuCiShuBuNengWeiLing")
	}
	if request.UpdateFrequency <= 0 {
		return fmt.Errorf("GengXinPingLvBuNengXiaoYuDengYuLing")
	}
	return nil
}

// DeleteActiveFirstRechargeRule 删除首充活动奖励
func (s ActiveFirstRechargeService) DeleteActiveFirstRechargeRule(id int) error {
	FirstRechargeRuleModel := models.CreateFirstRechargeRuleModel()
	if id == 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	_, err := FirstRechargeRuleModel.Update(&models.FirstRechargeRule{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	return nil
}

type AllActiveFirstRechargeRuleRewardRequest struct {
	ActivityId int `json:"activity_id"`
}

// AllActiveFirstRechargeReward 首充活动奖励列表
func (s ActiveFirstRechargeService) AllActiveFirstRechargeReward(request AllActiveFirstRechargeRuleRewardRequest) ([]models.FirstRechargeRuleReward, error) {
	rule, err := s.getAndValidateRule(request.ActivityId)
	if err != nil {
		return nil, err
	}
	FirstRechargeRuleRewardModel := models.CreateFirstRechargeRuleRewardModel()
	var list []models.FirstRechargeRuleReward
	_, err = FirstRechargeRuleRewardModel.QueryTable(&models.FirstRechargeRuleReward{}).
		Filter("first_recharge_rule_id", rule.Id).
		Filter("is_deleted", 0).
		OrderBy("serial_number").
		All(&list)
	if nil != err {
		return nil, err
	}

	return list, nil
}

type AddActiveFirstRechargeRuleRewardRequest struct {
	ActivityId          int     `json:"activity_id"`           // 活动ID
	SerialNumber        int     `json:"serial_number"`         // 序号
	TotalRechargeAmount float64 `json:"total_recharge_amount"` // 累计存款金额
	RewardAmount        float64 `json:"reward_amount"`         // 奖励金额
}

// AddActiveFirstRechargeReward 添加首充活动奖励
func (s ActiveFirstRechargeService) AddActiveFirstRechargeReward(request AddActiveFirstRechargeRuleRewardRequest) error {
	rule, err := s.getAndValidateRule(request.ActivityId)
	if err != nil {
		return err
	}
	err = s.validateAddActiveFirstRechargeRuleReward(request, rule.Id)
	if err != nil {
		return err
	}
	FirstRechargeRuleRewardModel := models.CreateFirstRechargeRuleRewardModel()
	_, err = FirstRechargeRuleRewardModel.Insert(&models.FirstRechargeRuleReward{
		FirstRechargeRuleId: rule.Id,
		SerialNumber:        request.SerialNumber,
		TotalRechargeAmount: request.TotalRechargeAmount,
		RewardAmount:        request.RewardAmount,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加首充活动奖励参数验证
func (s ActiveFirstRechargeService) validateAddActiveFirstRechargeRuleReward(request AddActiveFirstRechargeRuleRewardRequest, firstRechargeRuleId int) error {
	if firstRechargeRuleId <= 0 {
		return fmt.Errorf("ShouChongHuoDongGuiZeIDBiTian")
	}
	if request.SerialNumber <= 0 {
		return fmt.Errorf("ShouChongHuoDongJiangLiXuHaoBiTian")
	}
	isExist := models.CreateFirstRechargeRuleRewardModel().QueryTable(new(models.FirstRechargeRuleReward)).
		Filter("first_recharge_rule_id", firstRechargeRuleId).
		Filter("serial_number", request.SerialNumber).
		Filter("is_deleted", 0).
		Exist()
	if isExist {
		return fmt.Errorf("ShouChongHuoDongJiangLiXuHaoBuNengChongFu")
	}
	if request.TotalRechargeAmount < 0 {
		return fmt.Errorf("ShouChongHuoDongJiangLiZongCunKuanJinEBuNengWeiLing")
	}
	if request.RewardAmount < 0 {
		return fmt.Errorf("ShouChongHuoDongJiangLiJinEBuNengWeiLing")
	}
	return nil
}

type EditActiveFirstRechargeRuleRewardRequest struct {
	Id                  int     `json:"id"`                    // 首充活动奖励ID
	ActivityId          int     `json:"activity_id"`           // 活动ID
	SerialNumber        int     `json:"serial_number"`         // 序号
	TotalRechargeAmount float64 `json:"total_recharge_amount"` // 累计存款金额
	RewardAmount        float64 `json:"reward_amount"`         // 奖励金额
}

// EditActiveFirstRechargeReward 编辑首充活动奖励
func (s ActiveFirstRechargeService) EditActiveFirstRechargeReward(request EditActiveFirstRechargeRuleRewardRequest) error {
	rule, err := s.getAndValidateRule(request.ActivityId)
	if err != nil {
		return err
	}
	err = s.validateEditActiveFirstRechargeRuleReward(request, rule.Id)
	if err != nil {
		return err
	}
	fields := []string{
		"first_recharge_rule_id", "serial_number", "total_recharge_amount", "reward_amount",
	}
	FirstRechargeRuleRewardModel := models.CreateFirstRechargeRuleRewardModel()
	_, err = FirstRechargeRuleRewardModel.Update(&models.FirstRechargeRuleReward{
		Id:                  request.Id,
		FirstRechargeRuleId: rule.Id,
		SerialNumber:        request.SerialNumber,
		TotalRechargeAmount: request.TotalRechargeAmount,
		RewardAmount:        request.RewardAmount,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑首充活动奖励参数验证
func (s ActiveFirstRechargeService) validateEditActiveFirstRechargeRuleReward(request EditActiveFirstRechargeRuleRewardRequest, firstRechargeRuleId int) error {
	if request.Id == 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	if firstRechargeRuleId <= 0 {
		return fmt.Errorf("ShouChongHuoDongGuiZeIDBiTian")
	}
	if request.SerialNumber <= 0 {
		return fmt.Errorf("ShouChongHuoDongJiangLiXuHaoBiTian")
	}
	isExist := models.CreateFirstRechargeRuleRewardModel().QueryTable(new(models.FirstRechargeRuleReward)).
		Filter("first_recharge_rule_id", firstRechargeRuleId).
		Filter("id__ne", request.Id).
		Filter("serial_number", request.SerialNumber).
		Filter("is_deleted", 0).
		Exist()
	if isExist {
		return fmt.Errorf("ShouChongHuoDongJiangLiXuHaoBuNengChongFu")
	}
	if request.TotalRechargeAmount < 0 {
		return fmt.Errorf("ShouChongHuoDongJiangLiZongCunKuanJinEBuNengWeiLing")
	}
	if request.RewardAmount < 0 {
		return fmt.Errorf("ShouChongHuoDongJiangLiJinEBuNengWeiLing")
	}
	return nil
}

// DeleteActiveFirstRechargeReward 删除首充活动奖励
func (s ActiveFirstRechargeService) DeleteActiveFirstRechargeReward(id int) error {
	FirstRechargeRuleRewardModel := models.CreateFirstRechargeRuleRewardModel()
	if id == 0 {
		return fmt.Errorf("JiangLiBuCunZai")
	}
	_, err := FirstRechargeRuleRewardModel.Update(&models.FirstRechargeRuleReward{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	return nil
}
