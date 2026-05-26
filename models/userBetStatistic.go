package models

type UserBetStatisticModel struct {
	*Base
}

type UserBetStatistic struct {
	Id         int     `json:"id" orm:"auto;column(id)"`              // 主键，自增
	UserId     int     `json:"user_id" orm:"column(user_id)"`         // 用户id
	RoleId     int64   `json:"role_id" orm:"column(role_id)"`         // 用户RoleId
	Amount     int64   `json:"amount" orm:"column(amount)"`           // 变动金币
	NeedBets   int64   `json:"need_bets" orm:"column(need_bets)"`     // 需要打码的金额
	AlreadyBet int64   `json:"already_bet" orm:"column(already_bet)"` // 已打码金额
	Remaining  int64   `json:"remaining" orm:"column(remaining)"`     // 剩余打码金额
	BetRate    float64 `json:"bet_rate" orm:"column(bet_rate)"`       // 打码倍率
	BetType    int     `json:"bet_type" orm:"column(bet_type)"`       // 打码类型
	CTime      int64   `json:"c_time" orm:"column(c_time)"`           // 领取时间 (Unix时间戳)
}

func CreateUserBetStatisticModel() *UserBetStatisticModel {
	return &UserBetStatisticModel{CreateBase()}
}
