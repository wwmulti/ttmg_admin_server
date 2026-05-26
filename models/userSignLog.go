package models

type UserSignLogModel struct {
	*Base
}

type UserSignLog struct {
	Id         int   `json:"id" orm:"auto;column(id)"`                // 主键，自增
	SignRuleId int   `json:"sign_rule_id" orm:"column(sign_rule_id)"` // 用户签到的活动
	UserId     int   `json:"user_id" orm:"column(user_id)"`           // 用户id
	SignAt     int64 `json:"sign_at" orm:"column(sign_at)"`           // 用户签到时间（时间戳）
	Day        int   `json:"day" orm:"column(day)"`                   // 当前连续签到的天数
}

func CreateUserSignLogModel() *UserSignLogModel {
	return &UserSignLogModel{CreateBase()}
}
