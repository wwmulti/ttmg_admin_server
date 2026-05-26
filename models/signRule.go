package models

type SignRuleModel struct {
	*Base
}

type SignRule struct {
	Id               int `json:"id" orm:"auto;column(id)"`                                       // 主键，自增
	ActivityId       int `json:"activity_id" orm:"column(activity_id)"`                          // 所属活动Id
	Days             int `json:"days" orm:"column(days);default(7)"`                             // 签到周期
	IsLoop           int `json:"is_loop" orm:"column(is_loop);default(1)"`                       // 是否循环 0否 1是
	IsInterruptReset int `json:"is_interrupt_reset" orm:"column(is_interrupt_reset);default(1)"` // 断签是否中断连续签到 0否 1是
	Status           int `json:"status" orm:"column(status);default(1)"`                         // 活动状态 0关闭 1开启
	IsDeleted        int `json:"is_deleted" orm:"column(is_deleted);default(0)"`                 // 是否删除
}

func CreateSignRuleModel() *SignRuleModel {
	return &SignRuleModel{CreateBase()}
}
