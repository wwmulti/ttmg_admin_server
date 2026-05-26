package models

type UserActivityRewardLogModel struct {
	*Base
}

type UserActivityRewardLog struct {
	Id                    int     `json:"id" orm:"auto;column(id)"`                        // 主键
	UserId                int     `json:"user_id" orm:"column(user_id)"`                   // 用户id
	ActivityId            int     `json:"activity_id" orm:"column(activity_id)"`           // 活动id
	ActivityTypeId        int     `json:"activity_type_id" orm:"column(activity_type_id)"` // 活动类型id
	ActivityReceiveTypeId int     `json:"receive_type_id" orm:"column(receive_type_id)"`   // 领取类型id
	Amount                float64 `json:"amount" orm:"column(amount);null"`                // 奖励金额
	CreatedAt             int64   `json:"created_at" orm:"column(created_at);null"`        // 获取时间（时间戳）
}

type UserActivityRewardLogDTO struct {
	Id                    int     `json:"id" orm:"auto;column(id)"`                        // 主键
	UserId                int     `json:"user_id" orm:"column(user_id)"`                   // 用户id
	ActivityId            int     `json:"activity_id" orm:"column(activity_id)"`           // 活动id
	ActivityTypeId        int     `json:"activity_type_id" orm:"column(activity_type_id)"` // 活动类型id
	ActivityReceiveTypeId int     `json:"receive_type_id" orm:"column(receive_type_id)"`   // 领取类型id
	Amount                float64 `json:"amount" orm:"column(amount);null"`                // 奖励金额
	CreatedAt             string  `json:"created_at" orm:"column(created_at);null"`        // 获取时间（时间戳）
}

func CreateUserActivityRewardLogModel() *UserActivityRewardLogModel {
	return &UserActivityRewardLogModel{CreateBase()}
}
