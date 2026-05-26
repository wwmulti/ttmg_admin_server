package models

type LoginLogModel struct {
	*Base
}

type LoginLog struct {
	Id         int    `json:"-" orm:"auto;column(id)"`                               // 主键
	UserId     int    `json:"user_id" orm:"column(user_id)"`                         // 用户id
	RoleId     int    `json:"role_id" orm:"column(role_id)"`                         // 角色id
	IP         string `json:"ip" orm:"column(ip);size(46)"`                          // 登录IP
	RegisterAt int64  `json:"register_at" orm:"column(register_at);null"`            // 注册时间
	DeviceInfo string `json:"device_info" orm:"column(device_info);type(text);null"` // 设备信息
	CreatedAt  int64  `json:"created_at" orm:"column(created_at);null"`
}

func CreateLoginLogModel() *LoginLogModel {
	return &LoginLogModel{CreateBase()}
}
