package models

type UserShareLogsModel struct {
	*Base
}

type UserShareLogs struct {
	Id       int   `json:"id" orm:"auto;column(id)"`          // ID
	UserId   int   `json:"user_id" orm:"column(user_id)"`     // 用户id
	ActiveId int   `json:"active_id" orm:"column(active_id)"` // 活动id
	Pid1     int   `json:"pid1" orm:"column(pid1)"`           // 直属邀请父id
	Pid2     int   `json:"pid2" orm:"column(pid2)"`           // 二级父id
	Pid3     int   `json:"pid3" orm:"column(pid3)"`           // 三级父id
	Type     int   `json:"type" orm:"column(type)"`           // 邀请类型：0-邀请链接分享
	Ctime    int64 `json:"c_time" orm:"column(c_time)"`       // 时间
}

func CreateUserShareLogsModel() *UserShareLogsModel {
	return &UserShareLogsModel{CreateBase()}
}
