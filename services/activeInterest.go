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

type ActiveInterestService struct{}

type ActiveInterestRule struct {
	ActivityId int `json:"activity_id"`
}

// ActiveInterestRule 利息宝活动规则列表
func (s ActiveInterestService) ActiveInterestRule(request ActiveInterestRule) (*models.InterestRule, error) {
	interestRuleModel := models.CreateInterestRuleModel()
	rule := &models.InterestRule{}
	err := interestRuleModel.QueryTable(models.InterestRule{}).
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

// AllActiveInterestReceiveType 所有利息领取类型
func (s ActiveInterestService) AllActiveInterestReceiveType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.InterestRuleReceiveTypeRealTime, "name": i18n.Tr(language, "ShiShi")},
		{"id": config.InterestRuleReceiveTypeNextDayZero, "name": i18n.Tr(language, "DiErTian")},
		{"id": config.InterestRuleReceiveTypeNextWeekZero, "name": i18n.Tr(language, "XiaZhouYi")},
		{"id": config.InterestRuleReceiveTypeNextMonthZero, "name": i18n.Tr(language, "XiaYueYiHao")},
	}
}

// AllActiveInterestInterestLimitType 所有利息最大限制类型
func (s ActiveInterestService) AllActiveInterestInterestLimitType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.InterestRuleInterestLimitTypePercent, "name": i18n.Tr(language, "WanFenBi")},
		{"id": config.InterestRuleInterestLimitTypeFixedValue, "name": i18n.Tr(language, "GuDingZhi")},
		{"id": config.InterestRuleInterestLimitTypeUnlimited, "name": i18n.Tr(language, "WuXianZhi")},
	}
}

// 获取规则
func (s ActiveInterestService) getAndValidateRule(activityId int) (*models.InterestRule, error) {
	activityExists := models.CreateActivesModel().QueryTable(new(models.Actives)).
		Filter("id", activityId).
		Filter("status", 1).
		Filter("is_deleted", 0).
		Exist()
	if !activityExists {
		logs.Error("利息宝活动不存在")
		return nil, fmt.Errorf("QianDaoHuoDongBuCunZai")
	}
	rule := &models.InterestRule{}
	err := models.CreateInterestRuleModel().QueryTable(new(models.InterestRule)).
		Filter("activity_id", activityId).
		Filter("status", 1).
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

type AddActiveInterestRuleRequest struct {
	ActivityId          int   `json:"activity_id"`           // 活动ID
	InterestRate        int   `json:"interest_rate"`         // 周期利率
	DepositAmount       int64 `json:"deposit_amount"`        // 起存金额(分)
	Interval            int   `json:"interval"`              // 结算周期(小时)
	ReceiveType         int   `json:"receive_type"`          // 领取类型
	InterestLimitType   int   `json:"interest_limit_type"`   // 利息上限类型
	InterestLimitAmount int64 `json:"interest_limit_amount"` // 利息上限值
}

// AddActiveInterestRule 添加利息宝活动
func (s ActiveInterestService) AddActiveInterestRule(request AddActiveInterestRuleRequest) error {
	err := s.validateAddActiveInterestRule(request)
	if err != nil {
		return err
	}
	interestRuleModel := models.CreateInterestRuleModel()
	_, err = interestRuleModel.Insert(&models.InterestRule{
		ActivityId:          request.ActivityId,
		InterestRate:        request.InterestRate,
		DepositAmount:       request.DepositAmount,
		Interval:            request.Interval,
		ReceiveType:         request.ReceiveType,
		InterestLimitType:   request.InterestLimitType,
		InterestLimitAmount: request.InterestLimitAmount,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加利息宝活动参数验证
func (s ActiveInterestService) validateAddActiveInterestRule(request AddActiveInterestRuleRequest) error {
	if request.ActivityId == 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	activeService := ActiveService{}
	active, err := activeService.GetActiveById(request.ActivityId)
	if err != nil {
		return fmt.Errorf("HuoDongBuCunZai")
	}
	if active.ActiveTypeId != int(config.ActiveTypeInterest) {
		return fmt.Errorf("HuoDongLeiXingBuZhengQue")
	}
	if request.InterestRate <= 0 {
		return fmt.Errorf("LiLvBiXuDaYuLing")
	}
	if request.DepositAmount <= 0 {
		return fmt.Errorf("CunKuanJinEBiXuDaYuLing")
	}
	if request.Interval <= 0 {
		return fmt.Errorf("JieSuanZhouQiBiXuDaYuLing")
	}
	if request.ReceiveType <= 0 {
		return fmt.Errorf("LingQuLeiXingBiTian")
	}
	if request.InterestLimitType <= 0 {
		return fmt.Errorf("LiXiXianZhiLeiXingBiTian")
	}
	if request.InterestLimitAmount < 0 {
		return fmt.Errorf("LiXiXianZhiJinEBuNengXiaoYuLing")
	}
	return nil
}

type EditActiveInterestRuleRequest struct {
	Id                  int   `json:"id"`                    // 利息宝活动ID
	ActivityId          int   `json:"activity_id"`           // 活动ID
	InterestRate        int   `json:"interest_rate"`         // 周期利率
	DepositAmount       int64 `json:"deposit_amount"`        // 起存金额(分)
	Interval            int   `json:"interval"`              // 结算周期(小时)
	ReceiveType         int   `json:"receive_type"`          // 领取类型
	InterestLimitType   int   `json:"interest_limit_type"`   // 利息上限类型
	InterestLimitAmount int64 `json:"interest_limit_amount"` // 利息上限值
}

// EditActiveInterestRule 编辑利息宝活动规则
func (s ActiveInterestService) EditActiveInterestRule(request EditActiveInterestRuleRequest) error {
	err := s.validateEditActiveInterestRule(request)
	if err != nil {
		return err
	}
	fields := []string{
		"interest_rate", "deposit_amount", "interval", "receive_type", "interest_limit_type", "interest_limit_amount",
	}
	interestRuleModel := models.CreateInterestRuleModel()
	_, err = interestRuleModel.Update(&models.InterestRule{
		Id:                  request.Id,
		ActivityId:          request.ActivityId,
		InterestRate:        request.InterestRate,
		DepositAmount:       request.DepositAmount,
		Interval:            request.Interval,
		ReceiveType:         request.ReceiveType,
		InterestLimitType:   request.InterestLimitType,
		InterestLimitAmount: request.InterestLimitAmount,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑利息宝活动参数验证
func (s ActiveInterestService) validateEditActiveInterestRule(request EditActiveInterestRuleRequest) error {
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
	if active.ActiveTypeId != int(config.ActiveTypeInterest) {
		return fmt.Errorf("HuoDongLeiXingBuZhengQue")
	}
	if request.InterestRate <= 0 {
		return fmt.Errorf("LiLvBiXuDaYuLing")
	}
	if request.DepositAmount <= 0 {
		return fmt.Errorf("CunKuanJinEBiXuDaYuLing")
	}
	if request.Interval <= 0 {
		return fmt.Errorf("JieSuanZhouQiBiXuDaYuLing")
	}
	if request.ReceiveType <= 0 {
		return fmt.Errorf("LingQuLeiXingBiTian")
	}
	if request.InterestLimitType <= 0 {
		return fmt.Errorf("LiXiXianZhiLeiXingBiTian")
	}
	if request.InterestLimitAmount < 0 {
		return fmt.Errorf("LiXiXianZhiJinEBuNengXiaoYuLing")
	}
	return nil
}
