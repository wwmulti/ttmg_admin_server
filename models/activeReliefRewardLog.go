package models

type ActiveReliefRewardLogModel struct {
	*Base
}

type ActiveReliefRewardLog struct {
	Id       int     `json:"id" orm:"auto;column(id)"`          // 主键，自增
	ActiveId int     `json:"active_id" orm:"column(active_id)"` // 活动id
	UserId   int     `json:"user_id" orm:"column(user_id)"`     // 用户id
	WaterId  int     `json:"water_id" orm:"column(water_id)"`   // 满足条件的流水记录id
	Amount   float64 `json:"amount" orm:"column(amount)"`       // 返利金额
	CTime    int64   `json:"c_time" orm:"column(c_time)"`       // 领取时间
}

func CreateActiveReliefRewardLogModel() *ActiveReliefRewardLogModel {
	return &ActiveReliefRewardLogModel{CreateBase()}
}
