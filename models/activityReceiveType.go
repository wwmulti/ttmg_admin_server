package models

type ActivityReceiveTypeModel struct {
	*Base
}

type ActivityReceiveType struct {
	Id   int    `json:"id" orm:"auto;column(id)"`               // 主键，自增
	Name string `json:"name" orm:"column(name);size(255);null"` // 名称
}

func CreateActivityReceiveTypeModel() *ActivityReceiveTypeModel {
	return &ActivityReceiveTypeModel{CreateBase()}
}
