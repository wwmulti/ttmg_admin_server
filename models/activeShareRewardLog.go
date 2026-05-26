package models

type ActiveShareRewardLogModel struct {
	*Base
}

type ActiveShareRewardLog struct {
	Id          int     `json:"id" orm:"auto;column(id)"`                // ID
	ActiveId    int     `json:"active_id" orm:"column(active_id)"`       // 活动期id
	UserId      int     `json:"user_id" orm:"column(user_id)"`           // 邀请人用户id
	RewardId    int     `json:"reward_id" orm:"column(reward_id)"`       // 奖励配置id
	RewardLevel float64 `json:"reward_level" orm:"column(reward_level)"` // 奖励档位(10/20等)
	ReceiveTime int64   `json:"receive_time" orm:"column(receive_time)"` // 领取时间
	ReceiveType int     `json:"receive_type" orm:"column(receive_type)"` // 领取类型0-无 1-手动领取 2-系统发奖
	Status      int     `json:"status" orm:"column(status)"`             // 状态 0达成未发放 1已发放
	CTime       int64   `json:"c_time" orm:"column(c_time)"`             // 时间
}

func CreateActiveShareRewardLogModel() *ActiveShareRewardLogModel {
	return &ActiveShareRewardLogModel{CreateBase()}
}
