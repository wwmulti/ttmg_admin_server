package models

type ActiveLuckyClearLogModel struct {
	*Base
}

type ActiveLuckyClearLog struct {
	Id         int     `json:"-" orm:"auto;column(id)"`                            // 主键，自增
	Uid        int64   `json:"-" orm:"column(uid)"`                                // 玩家id
	Score      float64 `json:"score" orm:"column(score)"`                          // 积分
	CreateTime int64   `json:"create_time" orm:"column(create_time);auto_now_add"` // 创建时间
	ActiveId   int64   `json:"-" orm:"column(active_id)"`                          // 活动id
}

func CreateActiveLuckyClearLogModel() *ActiveLuckyClearLogModel {
	return &ActiveLuckyClearLogModel{CreateBase()}
}
