package models

type PaymentTypeModel struct {
	*Base
}

type PaymentType struct {
	Id        int    `json:"id" orm:"auto;column(id)"`            // 主键，自增
	Name      string `json:"name" orm:"column(name)"`             // 账号类型
	Status    int    `json:"status" orm:"column(status)"`         // 状态 0关闭 1开启
	Sort      int    `json:"sort" orm:"column(sort)"`             // 状态 0关闭 1开启
	IsDeleted int    `json:"is_deleted" orm:"column(is_deleted)"` // 是否删除
}

func CreatePaymentTypeModel() *PaymentTypeModel {
	return &PaymentTypeModel{CreateBase()}
}
