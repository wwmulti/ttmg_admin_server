package models

type BetRateModel struct {
	*Base
}

type BetRate struct {
	Id        int     `json:"id" orm:"auto;column(id)"`            // 主键，自增
	ProRate   float64 `json:"pro_rate" orm:"column(pro_rate)"`     // 线上打码倍数
	DevRate   float64 `json:"dev_rate" orm:"column(dev_rate)"`     // 测试打码倍数
	Type      int     `json:"type" orm:"column(type)"`             // 打码类型(唯一)
	Name      string  `json:"name" orm:"column(name)"`             // 名称
	PackageId int     `json:"package_id" orm:"column(package_id)"` // 包id
}

func CreateBetRateModel() *BetRateModel {
	return &BetRateModel{CreateBase()}
}
