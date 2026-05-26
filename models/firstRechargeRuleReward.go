package models

type FirstRechargeRuleRewardModel struct {
	*Base
}

type FirstRechargeRuleReward struct {
	Id                  int     `json:"id" orm:"auto;column(id)"`                                    // 主键，自增
	FirstRechargeRuleId int     `json:"first_recharge_rule_id" orm:"column(first_recharge_rule_id)"` // 首充规则id
	SerialNumber        int     `json:"serial_number" orm:"column(serial_number)"`                   // 序号
	TotalRechargeAmount float64 `json:"total_recharge_amount" orm:"column(total_recharge_amount)"`   // 累计存款金额
	RewardAmount        float64 `json:"reward_amount" orm:"column(reward_amount)"`                   // 奖励金额
	IsDeleted           int     `json:"is_deleted" orm:"column(is_deleted);default(0)"`              // 是否删除
}

type FirstRechargeRuleRewardDTO struct {
	Id                  int     `json:"id" orm:"auto;column(id)"`                                  // 主键，自增
	SerialNumber        int     `json:"serial_number" orm:"column(serial_number)"`                 // 序号
	TotalRechargeAmount float64 `json:"total_recharge_amount" orm:"column(total_recharge_amount)"` // 累计存款金额
	RewardAmount        float64 `json:"reward_amount" orm:"column(reward_amount)"`                 // 奖励金额
}

func CreateFirstRechargeRuleRewardModel() *FirstRechargeRuleRewardModel {
	return &FirstRechargeRuleRewardModel{CreateBase()}
}
