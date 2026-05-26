package models

type ActiveShareRewardsModel struct {
	*Base
}

type ActiveShareRewards struct {
	Id        int     `json:"id" orm:"auto;column(id)"`            // ID
	ActiveId  int     `json:"active_id" orm:"column(active_id)"`   // 活动ID
	RuleId    int     `json:"rule_id" orm:"column(rule_id)"`       // 活动规则id
	Mens      int     `json:"mens" orm:"column(mens)"`             // 达成人数条件（人）
	Rewards   float64 `json:"rewards" orm:"column(rewards)"`       // 奖励金额（元）
	IconOn    string  `json:"icon_on" orm:"column(icon_on)"`       // 图标-开
	IconOff   string  `json:"icon_off" orm:"column(icon_off)"`     // 图标-关
	Status    int     `json:"status" orm:"column(status)"`         // 状态0-关闭，1-开启
	IsDeleted int     `json:"is_deleted" orm:"column(is_deleted)"` // 是否删除
	Ctime     int64   `json:"c_time" orm:"column(c_time)"`         // 创建时间
}

func CreateActiveShareRewardsModel() *ActiveShareRewardsModel {
	return &ActiveShareRewardsModel{CreateBase()}
}

type ActiveShareRewardsDTO struct {
	List  []ActiveShareRewards `json:"list"`
	Total int64                `json:"total"`
}
