package models

type ActiveVipAwardLogModel struct {
	*Base
}

type ActiveVipAwardLog struct {
	Id       int     `json:"id" orm:"auto;column(id)"`          // 主键
	ActiveId int     `json:"active_id" orm:"column(active_id)"` // 活动id
	UserId   int     `json:"user_id" orm:"column(user_id)"`     // 用户id
	Lv       int     `json:"lv" orm:"column(lv)"`               // 档位
	Cycle    int     `json:"cycle" orm:"column(cycle)"`         // 奖励类型：0-晋级金 1-日 2-周 3-月
	Period   int     `json:"period" orm:"column(period)"`       // 期数
	Amount   float64 `json:"amount" orm:"column(amount);"`      // 奖励金额
	Ctime    int64   `json:"c_time" orm:"column(c_time)"`       // 时间
}

func CreateActiveVipAwardLogModel() *ActiveVipAwardLogModel {
	return &ActiveVipAwardLogModel{CreateBase()}
}
