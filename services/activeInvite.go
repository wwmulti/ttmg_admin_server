package services

import (
	"api/config"
	"api/models"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
	"time"
)

type ActiveInviteService struct{}

type ActiveInviteRuleRequest struct {
	ActiveId int `json:"active_id"`
}

// ActiveInviteRule 获取分享活动规则
func (s *ActiveInviteService) ActiveInviteRule(request ActiveInviteRuleRequest) (*models.ActiveShareRule, error) {
	shareRuleModel := models.CreateActiveShareRuleModel()
	rule := &models.ActiveShareRule{}
	err := shareRuleModel.QueryTable(models.ActiveShareRule{}).
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

// ActiveInviteRuleScopeType 所有扫描范围
func (s *ActiveInviteService) ActiveInviteRuleScopeType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.ShareScopeTypeChild, "name": i18n.Tr(language, "ZhiShuXiaJi")},
		{"id": config.ShareScopeTypeAll, "name": i18n.Tr(language, "SuoYouXiaJi")},
	}
}

// ActiveInviteRuleConditionType 所有条件类型
func (s *ActiveInviteService) ActiveInviteRuleConditionType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.ShareConditionTypeAnd, "name": i18n.Tr(language, "ManZuSuoYouTiaoJian")},
		{"id": config.ShareConditionTypeOr, "name": i18n.Tr(language, "ManZuRenYiTiaoJian")},
	}
}

// ActiveInviteRuleRewardType 所有领取方式
func (s *ActiveInviteService) ActiveInviteRuleRewardType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.ShareRewardTypeUserReceive, "name": i18n.Tr(language, "ShouDongLingQu")},
		{"id": config.ShareRewardTypeAutoReceive, "name": i18n.Tr(language, "ZiDongLingQu")},
	}
}

// ActiveInviteRuleExpireType 所有过期策略
func (s *ActiveInviteService) ActiveInviteRuleExpireType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.ShareExpireTypeAutoReceive, "name": i18n.Tr(language, "ZiDongLingQu")},
		{"id": config.ShareExpireTypeInvalid, "name": i18n.Tr(language, "ZuoFei")},
	}
}

type AddActiveInviteRuleRequest struct {
	ActiveId       int     `json:"active_id"`        // 活动ID
	ActiveTypeId   int     `json:"active_type_id"`   // 活动类型ID
	ShareUrl       string  `json:"share_url"`        // 推广域名
	ScanInterval   int     `json:"scan_interval"`    // 扫描用户下级间隔（分钟）
	Scope          int     `json:"scope"`            // 统计范围：0-仅直属下级 1-所有层级下级
	MiniTotalPays  float64 `json:"mini_total_pays"`  // 下级最低累计充值,0为无限制
	MiniTotalWater float64 `json:"mini_total_water"` // 押注消费流水，0为无限制
	Condition      int     `json:"condition"`        // 判断关系，0-同时满足，1-满足任一
	RewardType     int     `json:"reward_type"`      // 奖励领取方式 ，0-手动点击，1-系统自动派发
	ExpireType     int     `json:"expire_type"`      // 过期处理策略，0-结束未领取自动派发，1-过期作废
}

// AddActiveInviteRule 添加邀请规则
func (s *ActiveInviteService) AddActiveInviteRule(data AddActiveInviteRuleRequest) error {
	if data.ActiveId <= 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}

	if data.ActiveTypeId <= 0 {
		return fmt.Errorf("HuoDongLeiXingBiTian")
	}

	_, err := models.CreateActiveShareRuleModel().Insert(&models.ActiveShareRule{
		ActiveId:       data.ActiveId,
		ActiveTypeId:   data.ActiveTypeId,
		ShareUrl:       data.ShareUrl,
		ScanInterval:   data.ScanInterval,
		Scope:          data.Scope,
		MiniTotalPays:  data.MiniTotalPays,
		MiniTotalWater: data.MiniTotalWater,
		Condition:      data.Condition,
		RewardType:     data.RewardType,
		ExpireType:     data.ExpireType,
		IsDeleted:      0,
		Ctime:          time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// EditActiveInviteRule 编辑邀请规则
func (s *ActiveInviteService) EditActiveInviteRule(id int, data AddActiveInviteRuleRequest) error {
	if data.ActiveId <= 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}

	if data.ActiveTypeId <= 0 {
		return fmt.Errorf("HuoDongLeiXingBiTian")
	}
	fields := []string{
		"active_id",
		"active_type_id",
		"share_url",
		"scan_interval",
		"scope",
		"mini_total_pays",
		"mini_total_water",
		"condition",
		"reward_type",
		"expire_type",
	}
	affected, err := models.CreateActiveShareRuleModel().Update(&models.ActiveShareRule{
		Id:             id,
		ActiveId:       data.ActiveId,
		ActiveTypeId:   data.ActiveTypeId,
		ShareUrl:       data.ShareUrl,
		ScanInterval:   data.ScanInterval,
		Scope:          data.Scope,
		MiniTotalPays:  data.MiniTotalPays,
		MiniTotalWater: data.MiniTotalWater,
		Condition:      data.Condition,
		RewardType:     data.RewardType,
		ExpireType:     data.ExpireType,
	}, fields...)
	if err != nil || affected == 0 {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

type AllActiveInviteRewardRequest struct {
	ActiveId int `json:"active_id"`
}

// AllActiveInviteReward 获取邀请奖励列表
func (s *ActiveInviteService) AllActiveInviteReward(request AllActiveInviteRewardRequest) ([]models.ActiveShareRewards, error) {
	ActiveShareRewardsModel := models.CreateActiveShareRewardsModel()
	var list []models.ActiveShareRewards
	_, err := ActiveShareRewardsModel.QueryTable(&models.ActiveShareRewards{}).
		Filter("active_id", request.ActiveId).
		All(&list)
	if nil != err {
		return nil, err
	}

	return list, nil
}

type AddActiveInviteRewardRequest struct {
	ActiveId int     `json:"active_id"` // 活动ID
	RuleId   int     `json:"rule_id"`   // 活动规则id
	Mens     int     `json:"mens"`      // 达成人数条件（人）
	Rewards  float64 `json:"rewards"`   // 奖励金额（元）
	IconOn   string  `json:"icon_on"`   // 图标-开
	IconOff  string  `json:"icon_off"`  // 图标-关
	Status   int     `json:"status"`    // 状态0-关闭，1-开启
}

// AddActiveInviteReward 添加邀请奖励
func (s *ActiveInviteService) AddActiveInviteReward(data AddActiveInviteRewardRequest) error {
	if data.ActiveId <= 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	rule, err := s.getAndValidateRule(data.ActiveId)
	if err != nil {
		return err
	}
	_, err = models.CreateActiveShareRewardsModel().Insert(&models.ActiveShareRewards{
		ActiveId:  data.ActiveId,
		RuleId:    rule.Id,
		Mens:      data.Mens,
		Rewards:   data.Rewards,
		IconOn:    data.IconOn,
		IconOff:   data.IconOff,
		Status:    data.Status,
		IsDeleted: 0,
		Ctime:     time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}

	return nil
}

// EditActiveInviteReward 编辑邀请奖励
func (s *ActiveInviteService) EditActiveInviteReward(id int, data AddActiveInviteRewardRequest) error {
	if data.ActiveId <= 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	if data.RuleId <= 0 {
		return fmt.Errorf("GuiZeIdBiTian")
	}
	fields := []string{
		"active_id",
		"rule_id",
		"mens",
		"rewards",
		"icon_on",
		"icon_off",
		"status",
	}
	affected, err := models.CreateActiveShareRewardsModel().Update(&models.ActiveShareRewards{
		Id:       id,
		ActiveId: data.ActiveId,
		RuleId:   data.RuleId,
		Mens:     data.Mens,
		Rewards:  data.Rewards,
		IconOn:   data.IconOn,
		IconOff:  data.IconOff,
		Status:   data.Status,
	}, fields...)
	if err != nil || affected == 0 {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// DeleteActiveInviteReward 删除邀请奖励
func (s *ActiveInviteService) DeleteActiveInviteReward(id int) error {
	affected, err := models.CreateActiveShareRewardsModel().Update(&models.ActiveShareRewards{
		Id:        id,
		IsDeleted: 1,
	})
	if err != nil || affected == 0 {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

type ChangeActiveInviteRewardsStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangeActiveInviteRewardsStatus 修改Vip等级任务状态
func (s *ActiveInviteService) ChangeActiveInviteRewardsStatus(request ChangeActiveInviteRewardsStatusRequest) error {
	taskModel := models.CreateActiveShareRewardsModel()
	if request.Id == 0 {
		return fmt.Errorf("VipDengJiRenWuIDBiTian")
	}
	_, err := taskModel.Update(&models.ActiveShareRewards{
		Id:     request.Id,
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}
	return nil
}

// 获取规则
func (s *ActiveInviteService) getAndValidateRule(activityId int) (*models.ActiveShareRule, error) {
	activityExists := models.CreateActivesModel().QueryTable(new(models.Actives)).
		Filter("id", activityId).
		Filter("is_deleted", 0).
		Exist()
	if !activityExists {
		logs.Error("邀请活动不存在")
		return nil, fmt.Errorf("YaoQingHuoDongBuCunZai")
	}
	rule := &models.ActiveShareRule{}
	err := models.CreateActiveShareRuleModel().QueryTable(new(models.ActiveShareRule)).
		Filter("active_id", activityId).
		Filter("is_deleted", 0).
		One(rule)

	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, fmt.Errorf("YaoQingHuoDongBuCunZai")
		} else {
			return nil, fmt.Errorf("WeiZhiDeCuoWu")
		}
	}
	return rule, nil
}
