package models

type ActiveReliefRewardsModel struct {
	*Base
}

type ActiveReliefRewards struct {
	Id            int     `json:"id" orm:"auto;column(id)"`                    // 主键，自增
	ActiveId      int     `json:"active_id" orm:"column(active_id)"`           // 活动id
	RuleId        int     `json:"rule_id" orm:"column(rule_id)"`               // 规则id
	Amount        float64 `json:"amount" orm:"column(amount)"`                 // 亏损金额
	RebatePercent float64 `json:"rebate_percent" orm:"column(rebate_percent)"` // 返利比例（百分数）
	Status        int     `json:"status" orm:"column(status)"`                 // 状态,0-关闭，1-开启
	Ctime         int64   `json:"c_time" orm:"column(c_time)"`                 // 时间
}

func CreateActiveReliefRewardsModel() *ActiveReliefRewardsModel {
	return &ActiveReliefRewardsModel{CreateBase()}
}

type ActiveReliefRewardsDTO struct {
	List  []ActiveReliefRewards `json:"list"`
	Total int64                 `json:"total"`
}
