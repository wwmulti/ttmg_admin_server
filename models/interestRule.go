package models

type InterestRuleModel struct {
	*Base
}

type InterestRule struct {
	Id                  int   `json:"id" orm:"auto;column(id)"`                                  // 主键，自增
	ActivityId          int   `json:"activity_id" orm:"column(activity_id)"`                     // 活动id
	InterestRate        int   `json:"interest_rate" orm:"column(interest_rate)"`                 // 周期利率
	DepositAmount       int64 `json:"deposit_amount" orm:"column(deposit_amount)"`               // 起存金额
	Interval            int   `json:"interval" orm:"column(interval)"`                           // 结算周期(小时)
	ReceiveType         int   `json:"receive_type" orm:"column(receive_type)"`                   // 可领取时间 1实时 2次日零点 3次周零点 4次月零点
	InterestLimitType   int   `json:"interest_limit_type" orm:"column(Interest_limit_type)"`     // 利息上限类型 1百分比 2定额 3无限制
	InterestLimitAmount int64 `json:"interest_limit_amount" orm:"column(Interest_limit_amount)"` // 利息限制值 0为无限制
	Status              int   `json:"status" orm:"column(status)"`                               // 状态 0下架 1上架
	IsDeleted           int   `json:"is_deleted" orm:"column(is_deleted)"`                       // 是否删除
}

func CreateInterestRuleModel() *InterestRuleModel {
	return &InterestRuleModel{CreateBase()}
}
