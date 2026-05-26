package models

type FirstRechargeRuleModel struct {
	*Base
}

type FirstRechargeRule struct {
	Id              int     `json:"id" orm:"auto;column(id)"`                             // 主键，自增
	ActivityId      int     `json:"activity_id" orm:"column(activity_id);default(1)"`     // 活动ID
	BillingType     int     `json:"billing_type" orm:"column(billing_type);default(1)"`   // 统计方式 1按日累计
	ReceiveType     int     `json:"receive_type" orm:"column(receive_type)"`              // 领取方式 1次日零点领取
	TaskType        int     `json:"task_type" orm:"column(task_type)"`                    // 领取任务要求 1投注次数 2投注金额
	BetNumber       int     `json:"bet_number" orm:"column(bet_number);null"`             // 投注次数
	BetAmount       float64 `json:"bet_amount" orm:"column(bet_amount);null"`             // 投注金额
	Status          int     `json:"status" orm:"column(status);default(1)"`               // 状态 0未开启 1已开启
	UpdateFrequency int     `json:"update_frequency" orm:"column(update_frequency)"`      // 更新统计频率（分钟）
	RepeatActive    int     `json:"repeat_active" orm:"column(repeat_active);default(0)"` // 是否可重复参与
	IsDeleted       int     `json:"is_deleted" orm:"column(is_deleted);default(0)"`       // 是否删除
}

func CreateFirstRechargeRuleModel() *FirstRechargeRuleModel {
	return &FirstRechargeRuleModel{CreateBase()}
}
