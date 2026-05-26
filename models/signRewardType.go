package models

type SignRewardTypeModel struct {
	*Base
}

type SignRewardType struct {
	Id        int    `json:"id" orm:"auto;column(id)"`                       // 主键，自增
	Name      string `json:"name" orm:"column(name);size(255)"`              // 类型名称
	IsDeleted int    `json:"is_deleted" orm:"column(is_deleted);default(0)"` // 是否删除
}

func CreateSignRewardTypeModel() *SignRewardTypeModel {
	return &SignRewardTypeModel{CreateBase()}
}
