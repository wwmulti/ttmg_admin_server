package models

type RegisterLogModel struct {
	*Base
}

type RegisterLog struct {
	Id           int    `json:"id" orm:"auto;column(id)"`                       // 主键
	UserId       int    `json:"user_id" orm:"column(user_id);null"`             // 用户id
	RoleId       int    `json:"role_id" orm:"column(role_id);null"`             // 角色id
	IP           string `json:"ip" orm:"column(ip);size(46);null"`              // 注册IP
	RegisterDate int64  `json:"register_date" orm:"column(register_date);null"` // 注册日期
	CreatedAt    int64  `json:"created_at" orm:"column(created_at);null"`       // 创建时间
}

func CreateRegisterLogModel() *RegisterLogModel {
	return &RegisterLogModel{CreateBase()}
}
