package models

type AuthGroupModel struct {
	*Base
}

type AuthGroup struct {
	Id              int        `json:"id" orm:"auto;column(id)"`                          // 主键
	Title           string     `json:"title" orm:"column(title)"`                         // 描述
	Status          int        `json:"status" orm:"column(status)"`                       // 0禁用 1启动
	Rules           string     `json:"rules" orm:"column(rules)"`                         // 规则id
	Pid             int        `json:"pid" orm:"column(pid)"`                             // 父id
	Parents         string     `json:"-" orm:"column(parents)"`                           // 父级id 例如,1,2,
	PackageGroupIds string     `json:"package_group_ids" orm:"column(package_group_ids)"` // 分包分组ID
	PackageIds      string     `json:"package_ids" orm:"column(package_ids)"`             // 分包ID
	Account         []*Account `orm:"reverse(many)"`
}

func CreateAuthGroupModel() *AuthGroupModel {
	return &AuthGroupModel{CreateBase()}
}
