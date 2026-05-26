package models

type ActiveLuckyLogModel struct {
	*Base
}

type ActiveLuckyLog struct {
	Id         int     `json:"-" orm:"auto;column(id)"`               // 主键，自增
	Uid        int     `json:"-" orm:"column(uid)"`                   // 玩家id
	Username   string  `json:"username" orm:"column(username)"`       // 用户账号
	WheelType  int     `json:"game_type" orm:"column(wheel_type)"`    // 游戏类型
	Reward     float64 `json:"reward" orm:"column(reward)"`           // 奖励
	CreateTime int64   `json:"create_time" orm:"column(create_time)"` // 创建时间
	Cost       float64 `json:"cost" orm:"column(cost)"`               // 消耗
	ActiveId   int64   `json:"-" orm:"column(active_id)"`             // 活动id
	Last       float64 `json:"last" orm:"column(last)"`               // 剩余积分
}

type ActiveLuckyLogList struct {
	Lists []ActiveLuckyLog `json:"lists"`
	Rows  int64            `json:"rows"`
}

func CreateActiveLuckyLogModel() *ActiveLuckyLogModel {
	return &ActiveLuckyLogModel{CreateBase()}
}
