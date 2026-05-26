package models

type AuthRuleModel struct {
	*Base
}

type AuthRule struct {
	Id     int    `json:"id" orm:"auto;column(id)"`    // 主键
	Title  string `json:"title" orm:"column(title)"`   // 描述
	Name   string `json:"name" orm:"column(name)"`     // 规则名称
	Status int    `json:"status" orm:"column(status)"` // 0禁用 1启动
	Tag    int    `json:"tag" orm:"column(tag)"`       // 1菜单 0数据管控
	Pid    int    `json:"pid" orm:"column(pid)"`       // 父级id
	Icon   string `json:"icon" orm:"column(icon)"`     // 图标
	Sort   int    `json:"sort" orm:"column(sort)"`     // 权重
}

func CreateAuthRuleModel() *AuthRuleModel {
	return &AuthRuleModel{CreateBase()}
}
