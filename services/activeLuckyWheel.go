package services

import (
	"api/config"
	"api/models"
	"fmt"
	"github.com/beego/i18n"
	"time"
)

type ActiveLuckyWheelService struct{}

type AllActiveLuckyWheelRewardRequest struct {
	ActivityId int `json:"activity_id"` // 活动ID
	WheelType  int `json:"wheel_type"`  // 转盘类型
}

// AllActiveLuckyWheelReward 转盘活动奖励列表
func (s ActiveLuckyWheelService) AllActiveLuckyWheelReward(request AllActiveLuckyWheelRewardRequest) ([]models.ActiveLuckyReward, error) {
	luckyRewardRewardModel := models.CreateActiveLuckyRewardModel()
	var list []models.ActiveLuckyReward
	qs := luckyRewardRewardModel.QueryTable(&models.ActiveLuckyReward{}).
		Filter("activity_id", request.ActivityId)
	if request.WheelType != 0 {
		qs = qs.Filter("wheel_type", request.WheelType)
	}
	_, err := qs.Filter("is_deleted", 0).
		All(&list)
	if nil != err {
		return nil, err
	}

	return list, nil
}

// AllActiveLuckyWheelRewardType 所有转盘奖励类型
func (s ActiveLuckyWheelService) AllActiveLuckyWheelRewardType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.LuckyWheelSliver, "name": i18n.Tr(language, "YinPan")},
		{"id": config.LuckyWheelGold, "name": i18n.Tr(language, "JinPan")},
		{"id": config.LuckyWheelDiamond, "name": i18n.Tr(language, "ZuanShiPan")},
	}
}

type AddActiveLuckyWheelRewardRequest struct {
	ActivityId int     `json:"activity_id"` // 活动ID
	WheelType  int     `json:"wheel_type"`  // 转盘类型
	Reward     float64 `json:"reward"`      // 奖励金额
	Weight     int     `json:"weight"`      // 权重
}

// AddActiveLuckyWheelReward 添加转盘活动奖励
func (s ActiveLuckyWheelService) AddActiveLuckyWheelReward(request AddActiveLuckyWheelRewardRequest) error {
	now := time.Now().Unix()
	err := s.validateAddActiveLuckyWheelReward(request)
	if err != nil {
		return err
	}
	luckyRewardRewardModel := models.CreateActiveLuckyRewardModel()
	_, err = luckyRewardRewardModel.Insert(&models.ActiveLuckyReward{
		ActivityId: request.ActivityId,
		WheelType:  request.WheelType,
		Reward:     request.Reward,
		Weight:     request.Weight,
		CreateTime: now,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加转盘活动奖励参数验证
func (s ActiveLuckyWheelService) validateAddActiveLuckyWheelReward(request AddActiveLuckyWheelRewardRequest) error {
	if request.ActivityId <= 0 {
		return fmt.Errorf("HuoDongBuCunZai")
	}
	activeService := ActiveService{}
	active, err := activeService.GetActiveById(request.ActivityId)
	if err != nil {
		return fmt.Errorf("HuoDongBuCunZai")
	}
	if active.ActiveTypeId != int(config.ActiveTypeLuckyWheel) {
		return fmt.Errorf("HuoDongLeiXingBuZhengQue")
	}
	if request.Reward < 0 {
		return fmt.Errorf("JiangLiJinEBuNengWeiFuShu")
	}
	return nil
}

type EditActiveLuckyWheelRewardRequest struct {
	Id         int     `json:"id"`          // 转盘活动ID
	ActivityId int     `json:"activity_id"` // 活动ID
	WheelType  int     `json:"wheel_type"`  // 转盘类型
	Reward     float64 `json:"reward"`      // 奖励金额
	Weight     int     `json:"weight"`      // 权重
}

// EditActiveLuckyWheelReward 编辑转盘活动奖励
func (s ActiveLuckyWheelService) EditActiveLuckyWheelReward(request EditActiveLuckyWheelRewardRequest) error {
	now := time.Now().Unix()
	err := s.validateEditActiveLuckyWheelReward(request)
	if err != nil {
		return err
	}
	fields := []string{
		"wheel_type", "reward", "weight", "update_time",
	}
	luckyRewardRewardModel := models.CreateActiveLuckyRewardModel()
	_, err = luckyRewardRewardModel.Update(&models.ActiveLuckyReward{
		Id:         request.Id,
		WheelType:  request.WheelType,
		Reward:     request.Reward,
		Weight:     request.Weight,
		UpdateTime: now,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑转盘活动奖励参数验证
func (s ActiveLuckyWheelService) validateEditActiveLuckyWheelReward(request EditActiveLuckyWheelRewardRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("HuoDongIDBiTian")
	}
	if request.ActivityId <= 0 {
		return fmt.Errorf("HuoDongBuCunZai")
	}
	activeService := ActiveService{}
	active, err := activeService.GetActiveById(request.ActivityId)
	if err != nil {
		return fmt.Errorf("HuoDongBuCunZai")
	}
	if active.ActiveTypeId != int(config.ActiveTypeLuckyWheel) {
		return fmt.Errorf("HuoDongLeiXingBuZhengQue")
	}
	if request.Reward < 0 {
		return fmt.Errorf("JiangLiJinEBuNengWeiFuShu")
	}
	return nil
}

// DeleteActiveLuckyWheelReward 删除转盘奖励
func (s ActiveLuckyWheelService) DeleteActiveLuckyWheelReward(id int) error {
	signRuleRewardModel := models.CreateActiveLuckyRewardModel()
	if id == 0 {
		return fmt.Errorf("JiangLiBuCunZai")
	}
	_, err := signRuleRewardModel.Update(&models.ActiveLuckyReward{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	return nil
}
