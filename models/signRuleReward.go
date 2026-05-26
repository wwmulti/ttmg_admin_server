package models

type SignRuleRewardModel struct {
	*Base
}

type SignRuleReward struct {
	Id              int     `json:"id" orm:"auto;column(id)"`                               // 主键，自增
	SignRuleId      int     `json:"sign_rule_id" orm:"column(sign_rule_id)"`                // 签到规则id
	Day             int     `json:"day" orm:"column(day)"`                                  // 签到天数
	RewardTypeId    int     `json:"reward_type_id" orm:"column(reward_type_id)"`            // 奖励类型
	RewardAmount    float64 `json:"reward_amount" orm:"column(reward_amount);null"`         // 奖励金额
	Icon            string  `json:"icon" orm:"column(icon);size(255);null"`                 // 奖励图标
	DayRechargeLine float64 `json:"day_recharge_line" orm:"column(day_recharge_line);null"` // 每日充值门槛
	DayRunningLine  float64 `json:"day_running_line" orm:"column(day_running_line);null"`   // 每日流水门槛
	IsDeleted       int     `json:"is_deleted" orm:"column(is_deleted);default(0)"`         // 是否删除
}

type GetSignActivitySignRuleRewardDTO struct {
	Id              int     `json:"id"`                // 主键，自增
	SignRuleId      int     `json:"sign_rule_id"`      // 签到规则id
	Day             int     `json:"day"`               // 签到天数
	RewardTypeId    int     `json:"reward_type_id"`    // 奖励类型
	RewardAmount    float64 `json:"reward_amount"`     // 奖励金额
	Icon            string  `json:"icon"`              // 奖励图标
	DayRechargeLine float64 `json:"day_recharge_line"` // 每日充值门槛
	DayRunningLine  float64 `json:"day_running_line"`  // 每日流水门槛
}

func CreateSignRuleRewardModel() *SignRuleRewardModel {
	return &SignRuleRewardModel{CreateBase()}
}
