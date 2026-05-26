package models

type UserRelationshipModel struct {
	*Base
}

type UserRelationship struct {
	UserId    int    `json:"user_id" orm:"pk;column(user_id)"`    // 玩家id (主键)
	Pid       int    `json:"pid" orm:"column(pid)"`               // 上级
	Pid2      int    `json:"pid2" orm:"column(pid2)"`             // 上上级
	Pid3      int    `json:"pid3" orm:"column(pid3)"`             // 上上上级
	Parents   string `json:"parents" orm:"column(parents);null"`  // 用户关系结构 ,1,
	PackageId int    `json:"package_id" orm:"column(package_id)"` // 包id,
	PRoleId   int64  `json:"p_role_id" orm:"column(p_role_id)"`   // 上级用户RoleId
}

func CreateUserRelationshipModel() *UserRelationshipModel {
	return &UserRelationshipModel{CreateBase()}
}
