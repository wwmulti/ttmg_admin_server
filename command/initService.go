package command

import (
	"api/config"
	"api/models"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

func initActivity() error {
	activeModel := models.CreateActivesModel()
	tx, err := activeModel.Begin()
	if err != nil {
		logs.Error("初始化活动失败: %v", err)
		return err
	}
	for i := 1; i < int(config.ActiveTypeEnd); i++ {
		activity, isCreated, err := getActivity(i, tx)
		if err != nil {
			tx.Rollback()
			logs.Error("初始化活动失败: %v", err)
			return err
		}
		if isCreated == true {
			switch i {
			case int(config.ActiveTypeInvite):
				err = initInviteActivity(activity.Id, tx)
				if err != nil {
					tx.Rollback()
					logs.Error("初始化邀请注册活动失败: %v", err)
					return err
				}
			case int(config.ActiveTypeSign):
				err = initSignActivity(activity.Id, tx)
				if err != nil {
					tx.Rollback()
					logs.Error("初始化签到活动失败: %v", err)
					return err
				}
			case int(config.ActiveTypeFirstRecharge):
				err = initFirstRechargeActivity(activity.Id, tx)
				if err != nil {
					tx.Rollback()
					logs.Error("初始化首充活动失败: %v", err)
					return err
				}
			case int(config.ActiveTypeRelief):
				err = initReliefActivity(activity.Id, tx)
				if err != nil {
					tx.Rollback()
					logs.Error("初始化救济活动失败: %v", err)
					return err
				}
			case int(config.ActiveTypeInterest):
				err = initInterestActivity(activity.Id, tx)
				if err != nil {
					tx.Rollback()
					logs.Error("初始化利息宝活动失败: %v", err)
					return err
				}
			case int(config.ActiveTypeLuckyWheel):
				err = initLuckyWheelActivity(activity.Id, tx)
				if err != nil {
					tx.Rollback()
					logs.Error("初始化转盘活动失败: %v", err)
					return err
				}
			case int(config.ActiveTypeVipLv):
				err = initVipLvActivity(activity.Id, tx)
				if err != nil {
					tx.Rollback()
					logs.Error("初始化Vip等级失败: %v", err)
					return err
				}
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		logs.Error("初始化活动失败: %v", err)
		return err
	}
	return nil
}

// 获取或创建活动
func getActivity(typeId int, tx orm.TxOrmer) (*models.Actives, bool, error) {
	isCreated := false
	activityExists := models.CreateActivesModel().QueryTable(new(models.Actives)).
		Filter("active_type_id", typeId).
		Filter("package_id", 0).
		Exist()
	activityName, err := getActivityName(typeId)
	if err != nil {
		return nil, false, err
	}
	if !activityExists {
		activity := &models.Actives{
			Name:         activityName,
			PackageId:    0,
			ActiveTypeId: typeId,
		}
		_, err := tx.Insert(activity)
		isCreated = true
		if err != nil {
			logs.Error("创建活动类型为%d的初始活动模板失败，", typeId)
		}
	}
	activity := &models.Actives{}
	err = tx.QueryTable(new(models.Actives)).
		Filter("active_type_id", typeId).
		Filter("package_id", 0).
		One(activity)
	if err != nil {
		logs.Error("获取类型为%d的活动失败", typeId)
		return nil, isCreated, err
	}
	return activity, isCreated, nil
}

// 获取活动名称
func getActivityName(typeId int) (string, error) {
	switch typeId {
	case int(config.ActiveTypeInvite):
		return "邀请注册活动", nil
	case int(config.ActiveTypeSign):
		return "签到活动", nil
	case int(config.ActiveTypeFirstRecharge):
		return "首充活动", nil
	case int(config.ActiveTypeRelief):
		return "救济活动", nil
	case int(config.ActiveTypeInterest):
		return "利息宝活动", nil
	case int(config.ActiveTypeLuckyWheel):
		return "幸转动盘活动", nil
	case int(config.ActiveTypeVipLv):
		return "Vip等级活动", nil
	default:
		return "", fmt.Errorf("未知活动类型: %d", typeId)
	}
}

// 初始化邀请注册活动
func initInviteActivity(activityId int, tx orm.TxOrmer) error {
	rule := &models.ActiveShareRule{
		ActiveId:       activityId,
		ActiveTypeId:   int(config.ActiveTypeInvite),
		ShareUrl:       "https://www.google.com",
		ScanInterval:   1,
		Scope:          1,
		MiniTotalPays:  0,
		MiniTotalWater: 0,
		Condition:      0,
		RewardType:     0,
		ExpireType:     0,
	}
	id, err := tx.Insert(rule)
	if err != nil {
		logs.Error("初始化邀请注册活动失败: %v", err)
		return err
	}
	reward := &models.ActiveShareRewards{
		ActiveId: activityId,
		RuleId:   int(id),
		Mens:     1,
		Rewards:  1,
	}
	_, err = tx.Insert(reward)
	if err != nil {
		logs.Error("初始化邀请注册活动奖励失败: %v", err)
		return err
	}
	return nil
}

// 初始化签到活动
func initSignActivity(activityId int, tx orm.TxOrmer) error {
	rule := &models.SignRule{
		ActivityId:       activityId,
		Days:             7,
		IsLoop:           0,
		IsInterruptReset: 0,
	}
	id, err := tx.Insert(rule)
	if err != nil {
		logs.Error("初始化签到活动失败: %v", err)
		return err
	}
	for i := 1; i < 8; i++ {
		amount := float64(i)
		reward := &models.SignRuleReward{
			SignRuleId:      int(id),
			Day:             i,
			RewardTypeId:    int(config.SignRewardMoney),
			RewardAmount:    amount,
			Icon:            "",
			DayRunningLine:  amount,
			DayRechargeLine: amount,
		}
		_, err = tx.Insert(reward)
		if err != nil {
			logs.Error("初始化签到活动奖励失败: %v", err)
			return err
		}
	}

	return nil
}

// 初始化首充活动
func initFirstRechargeActivity(activityId int, tx orm.TxOrmer) error {
	rule := &models.FirstRechargeRule{
		ActivityId:      activityId,
		BillingType:     int(config.BillingDaily),
		ReceiveType:     int(config.ReceiveNextDayMidnight),
		TaskType:        int(config.TaskBetNumber),
		BetNumber:       1,
		BetAmount:       0,
		UpdateFrequency: 10,
		RepeatActive:    0,
	}
	id, err := tx.Insert(rule)
	if err != nil {
		logs.Error("初始化首充活动失败: %v", err)
		return err
	}
	for i := 1; i < 8; i++ {
		amount := float64(i)
		reward := &models.FirstRechargeRuleReward{
			FirstRechargeRuleId: int(id),
			SerialNumber:        i,
			TotalRechargeAmount: amount,
			RewardAmount:        amount,
		}
		_, err = tx.Insert(reward)
		if err != nil {
			logs.Error("初始化首充活动失败: %v", err)
			return err
		}
	}
	return nil
}

// 初始化救济活动
func initReliefActivity(activityId int, tx orm.TxOrmer) error {
	now := time.Now().Unix()
	rule := &models.ActiveReliefRule{
		ActiveId: activityId,
		Cycle:    2,
		OpenDay:  1,
		OpenTime: "00:00:00",
		IsRepeat: 1,
		Ctime:    now,
	}
	id, err := tx.Insert(rule)
	if err != nil {
		logs.Error("初始化救济活动失败: %v", err)
		return err
	}
	for i := 1; i < 7; i++ {
		amount := float64(i)
		reward := &models.ActiveReliefRewards{
			ActiveId:      activityId,
			RuleId:        int(id),
			Amount:        amount,
			RebatePercent: 1,
			Ctime:         now,
		}
		_, err = tx.Insert(reward)
		if err != nil {
			logs.Error("初始化救济活动失败: %v", err)
			return err
		}
	}

	return nil
}

// 初始化利息宝活动
func initInterestActivity(activityId int, tx orm.TxOrmer) error {
	rule := &models.InterestRule{
		ActivityId:          activityId,
		InterestRate:        10,
		DepositAmount:       50,
		Interval:            1,
		ReceiveType:         int(config.InterestRuleReceiveTypeRealTime),
		InterestLimitType:   int(config.InterestRuleInterestLimitTypePercent),
		InterestLimitAmount: 50,
	}
	_, err := tx.Insert(rule)
	if err != nil {
		logs.Error("初始化利息宝失败: %v", err)
		return err
	}
	return nil
}

// 初始化转盘活动
func initLuckyWheelActivity(activityId int, tx orm.TxOrmer) error {
	now := time.Now().Unix()
	for i := 1; i < 11; i++ {
		amount := float64(i)
		reward := &models.ActiveLuckyReward{
			ActivityId: activityId,
			WheelType:  int(config.LuckyWheelSliver),
			Reward:     amount,
			Weight:     10,
			CreateTime: now,
		}
		_, err := tx.Insert(reward)
		if err != nil {
			logs.Error("初始化转盘活动失败: %v", err)
			return err
		}
	}
	for i := 1; i < 11; i++ {
		amount := float64(i)
		reward := &models.ActiveLuckyReward{
			ActivityId: activityId,
			WheelType:  int(config.LuckyWheelGold),
			Reward:     amount,
			Weight:     10,
			CreateTime: now,
		}
		_, err := tx.Insert(reward)
		if err != nil {
			logs.Error("初始化转盘活动失败: %v", err)
			return err
		}
	}
	for i := 1; i < 11; i++ {
		amount := float64(i)
		reward := &models.ActiveLuckyReward{
			ActivityId: activityId,
			WheelType:  int(config.LuckyWheelDiamond),
			Reward:     amount,
			Weight:     10,
			CreateTime: now,
		}
		_, err := tx.Insert(reward)
		if err != nil {
			logs.Error("初始化转盘活动失败: %v", err)
			return err
		}
	}
	return nil
}

// 初始化Vip等级活动
func initVipLvActivity(activityId int, tx orm.TxOrmer) error {
	for i := 1; i < 11; i++ {
		amount := float64(i)
		rule := &models.ActiveVipRule{
			ActiveId:            activityId,
			Lv:                  i,
			TotalBets:           amount,
			TotalPays:           amount,
			ConAnd:              1,
			Rewards:             amount,
			WithdrawNumLimit:    -1,
			WithdrawAmountLimit: -1,
			WithdrawFreeNum:     -1,
			WithdrawFee:         1,
		}
		_, err := tx.Insert(rule)
		if err != nil {
			logs.Error("初始化Vip等级活动失败: %v", err)
		}
		reward := &models.ActiveVipWelfare{
			ActiveId:  activityId,
			Cycle:     2,
			Lv:        i,
			ConAnd:    1,
			TotalBets: amount,
			TotalPays: amount,
			Rewards:   amount,
		}
		_, err = tx.Insert(reward)
		if err != nil {
			logs.Error("初始化Vip等级活动失败: %v", err)
		}
	}
	return nil
}
