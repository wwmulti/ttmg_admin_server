package models

type ActiveVipWelfareAchieveModel struct {
	*Base
}

type ActiveVipWelfareAchieve struct {
	Id        int   `json:"id" orm:"auto;column(id)"`            // 主键
	ActiveId  int   `json:"active_id" orm:"column(active_id)"`   // 活动id
	UserId    int   `json:"user_id" orm:"column(user_id)"`       // 用户id
	Cycle     int   `json:"cycle" orm:"column(cycle)"`           // 周期 1-日 2-周 3-月
	Period    int   `json:"period" orm:"column(period)"`         // 期数
	TotalPays int64 `json:"total_pays" orm:"column(total_pays)"` // 周期内累计充值
	TotalBets int64 `json:"total_bets" orm:"column(total_bets)"` // 周期内累计有效押注
	Ctime     int64 `json:"c_time" orm:"column(c_time)"`         // 时间
}

func CreateActiveVipWelfareAchieveModel() *ActiveVipWelfareAchieveModel {
	return &ActiveVipWelfareAchieveModel{CreateBase()}
}
