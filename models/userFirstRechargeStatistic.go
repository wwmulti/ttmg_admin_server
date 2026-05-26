package models

type UserFirstRechargeStatisticModel struct {
	*Base
}

type UserFirstRechargeStatistic struct {
	Id                   int                `json:"id" orm:"auto;column(id)"`                                             // 主键，自增
	ActivityId           int                `json:"activity_id" orm:"column(activity_id);null"`                           // 活动id
	UserId               int64              `json:"user_id" orm:"column(user_id);null"`                                   // 用户id
	TotalBetAmount       int64              `json:"total_bet_amount" orm:"column(total_bet_amount);default(0)"`           // 总投注金额
	TotalBetNumber       int                `json:"total_bet_number" orm:"column(total_bet_number);default(0)"`           // 总投注次数
	TotalRechargeAmount  int64              `json:"total_recharge_amount" orm:"column(total_recharge_amount);default(0)"` // 总充值金额（单位按你业务定义）
	StartTime            int64              `json:"start_time" orm:"column(start_time);null"`                             // 周期开始时间
	EndTime              int64              `json:"end_time" orm:"column(end_time);null"`                                 // 周期结束时间
	AvailableReceiveTime int64              `json:"available_receive_time" orm:"column(available_receive_time);null"`     // 可领取时间
	ReceiveTime          int64              `json:"receive_time" orm:"column(receive_time);null"`                         // 领取时间
	IsReceive            int                `json:"is_receive" orm:"column(is_receive);default(0)"`                       // 是否领取 0否 1是
	Rule                 *FirstRechargeRule `orm:"rel(fk);column(rule_id)"`
}

func CreateUserFirstRechargeStatisticModel() *UserFirstRechargeStatisticModel {
	return &UserFirstRechargeStatisticModel{CreateBase()}
}
