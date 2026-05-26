package models

type UserInterestLogModel struct {
	*Base
}

type UserInterestLog struct {
	Id          int     `json:"id" orm:"auto;column(id)"`                     // 主键，自增
	InterestId  int     `json:"interest_id" orm:"column(interest_id)"`        // 用户利息宝表id
	ActivityId  int     `json:"activity_id" orm:"column(activity_id)"`        // 活动id
	UserId      int     `json:"user_id" orm:"column(user_id)"`                // 用户id
	LogType     int     `json:"log_type" orm:"column(log_type)"`              // 日志类型 1存入金额 2取出金额 3生成利息
	Amount      float64 `json:"amount" orm:"column(amount)"`                  // 操作对应的金额
	CalculateAt int64   `json:"calculate_at" orm:"column(calculate_at);null"` // 利息生成时间
	CreatedAt   int64   `json:"created_at" orm:"column(created_at);null"`     // 创建时间
}

type UserInterestLogDTO struct {
	Id          int     `json:"id" orm:"auto;column(id)"`                     // 主键，自增
	UserId      int     `json:"user_id" orm:"column(user_id)"`                // 用户id
	LogType     int     `json:"log_type" orm:"column(log_type)"`              // 日志类型 1存入金额 2取出金额 3生成利息
	Amount      float64 `json:"amount" orm:"column(amount)"`                  // 金额
	CalculateAt int64   `json:"calculate_at" orm:"column(calculate_at);null"` // 操作时间
}

func CreateUserInterestLogModel() *UserInterestLogModel {
	return &UserInterestLogModel{CreateBase()}
}

type UserInterestStatistic struct {
	Balance        float64 `json:"balance" orm:"column(balance)"`                 // 账户余额
	Interest       float64 `json:"interest" orm:"column(interest)"`               // 可提现的利息
	HistoryHistory float64 `json:"history_balance" orm:"column(history_balance)"` // 全部的利息
	LimitType      int     `json:"limit_type" orm:"column(limit_type)"`           // 利息计算类型 1百分比 2定额 3无限制
	LimitAmount    int64   `json:"limit_amount" orm:"column(limit_amount)"`       // 利息上限值
}
