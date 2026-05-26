package models

type ActiveVipWelfareModel struct {
	*Base
}

type ActiveVipWelfare struct {
	Id        int     `json:"id" orm:"auto;column(id)"`            // 主键
	ActiveId  int     `json:"active_id" orm:"column(active_id)"`   // 活动id
	Cycle     int     `json:"cycle" orm:"column(cycle)"`           // 奖励类型：1-日 2-周 3-月
	Lv        int     `json:"lv" orm:"column(lv)"`                 // 匹配的等级
	TotalPays float64 `json:"total_pays" orm:"column(total_pays)"` // 累计充值
	TotalBets float64 `json:"total_bets" orm:"column(total_bets)"` // 累计有效押注
	ConAnd    int     `json:"con_and" orm:"column(con_and)"`       // 条件：0-或 1-且 2-只需充值 3-只需流水
	Rewards   float64 `json:"rewards" orm:"column(rewards)"`       // 奖励金额
	Status    int     `json:"status" orm:"column(status)"`         // 状态：0-关闭 1-启用
	Ctime     int64   `json:"c_time" orm:"column(c_time)"`         // 时间
}

func CreateActiveVipWelfareModel() *ActiveVipWelfareModel {
	return &ActiveVipWelfareModel{CreateBase()}
}

type ActiveVipWelfareDTO struct {
	List  []ActiveVipWelfare `json:"list"`
	Total int64              `json:"total"`
}
