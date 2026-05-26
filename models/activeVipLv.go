package models

type ActiveVipLvModel struct {
	*Base
}

type ActiveVipLv struct {
	Id        int   `json:"id" orm:"auto;column(id)"`            // 主键
	ActiveId  int   `json:"active_id" orm:"column(active_id)"`   // 活动id
	UserId    int   `json:"user_id" orm:"column(user_id)"`       // 用户id
	Lv        int   `json:"lv" orm:"column(lv)"`                 // 用户实时等级
	TotalPays int64 `json:"total_pays" orm:"column(total_pays)"` // 累计充值
	TotalBets int64 `json:"total_bets" orm:"column(total_bets)"` // 有效下注
	Ctime     int64 `json:"c_time" orm:"column(c_time)"`         // 时间
}

func CreateActiveVipLvModel() *ActiveVipLvModel {
	return &ActiveVipLvModel{CreateBase()}
}
