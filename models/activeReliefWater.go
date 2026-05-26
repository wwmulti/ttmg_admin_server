package models

type ActiveReliefWaterModel struct {
	*Base
}

type ActiveReliefWater struct {
	Id             int   `json:"id" orm:"auto;column(id)"`                      // 主键
	ActiveId       int   `json:"active_id" orm:"column(active_id)"`             // 活动id
	UserId         int   `json:"user_id" orm:"column(user_id)"`                 // 用户id
	Amount         int64 `json:"amount" orm:"column(amount)"`                   // 活动时间内累计亏损
	Cycle          int   `json:"cycle" orm:"column(cycle)"`                     // 统计周期：1-按日，2-按周，3-按月
	Period         int   `json:"period" orm:"column(period)"`                   // 活动内期数：日-年月日，周-月周，月-年月
	ActivationTime int64 `json:"activation_time" orm:"column(activation_time)"` // 可领取时间
	Ctime          int64 `json:"c_time" orm:"column(c_time)"`                   // 创建时间
}

func CreateActiveReliefWaterModel() *ActiveReliefWaterModel {
	return &ActiveReliefWaterModel{CreateBase()}
}

type ReliefWaterDTO struct {
	Current    ReliefRewardActiveDTO   `json:"current"`     // 当前数据
	AllRewards []ReliefRewardActiveDTO `json:"all_rewards"` // 可领取返利档位
	Rules      []ReliefRewardsRuleDTO  `json:"rules"`       // 返利规则
}

type ReliefRewardsRuleDTO struct {
	Amount        float64 `json:"amount" orm:"column(amount)"`                 // 亏损金额
	RebatePercent string  `json:"rebate_percent" orm:"column(rebate_percent)"` // 返利比例（百分数）
}

type ReliefRewardActiveDTO struct {
	Id             int     `json:"id"`              // 领取id
	Amount         float64 `json:"amount"`          // 上周损失
	Rewards        float64 `json:"rewards"`         // 本周可领取返利
	Status         int     `json:"status"`          // 0-不可领取，1-可领取 2-已领取
	ActivationTime int64   `json:"activation_time"` // 可领取时间
}
