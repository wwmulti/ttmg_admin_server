package models

type UserOperateLogModel struct {
	*Base
}

type UserOperateLog struct {
	Id        int    `json:"id" orm:"auto;column(id)"`                 // 主键，自增
	UserId    int    `json:"user_id" orm:"column(user_id)"`            // 用户ID
	RoleId    int64  `json:"role_id" orm:"column(role_id)"`            // 用户RoleID
	Uri       string `json:"uri" orm:"column(uri);size(512);null"`     // 请求路径
	Body      string `json:"body" orm:"column(body);type(json);null"`  // 请求参数
	CreatedAt int64  `json:"created_at" orm:"column(created_at);null"` // 创建时间
}

func CreateUserOperateLogModel() *UserOperateLogModel {
	return &UserOperateLogModel{CreateBase()}
}
