package models

type ActiveReliefRuleModel struct {
	*Base
}

type ActiveReliefRule struct {
	Id       int    `json:"id" orm:"auto;column(id)"`          // 主键，自增
	ActiveId int    `json:"active_id" orm:"column(active_id)"` // 活动id
	Cycle    int    `json:"cycle" orm:"column(cycle)"`         // 统计周期：1-按日，2-按周，3-按月
	OpenDay  int    `json:"open_day" orm:"column(open_day)"`   // 开放领取日期（周：1-7；月：1-31；日：固定 1）
	OpenTime string `json:"open_time" orm:"column(open_time)"` // 开放领取时间点（如 "00:00:00"）
	IsRepeat int    `json:"is_repeat" orm:"column(is_repeat)"` // 是否重复：0-否，1-是
	Status   int    `json:"status" orm:"column(status)"`       // 状态：0-关闭，1-开启
	Ctime    int64  `json:"c_time" orm:"column(c_time)"`       // 时间
}

func CreateActiveReliefRuleModel() *ActiveReliefRuleModel {
	return &ActiveReliefRuleModel{CreateBase()}
}

type ActiveReliefRuleModelDTO struct {
	List  []ActiveReliefRule `json:"list"`
	Total int64              `json:"total"`
}
