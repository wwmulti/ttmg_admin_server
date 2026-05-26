package models

type ActiveTypesModel struct {
	*Base
}

type ActiveTypes struct {
	Id        int    `json:"id" orm:"auto;column(id)"`            // 主键，自增
	Name      string `json:"name" orm:"column(name)"`             // 类型名称
	Status    int    `json:"status" orm:"column(status)"`         // 状态0-关闭 1-开启
	IsDeleted int    `json:"is_deleted" orm:"column(is_deleted)"` // 删除
}

func CreateActiveTypes() *ActiveTypesModel {
	return &ActiveTypesModel{CreateBase()}
}
