package models

type PackageGroupModel struct {
	*Base
}

type PackageGroup struct {
	Id         int    `json:"id" orm:"auto;column(id)"`                        // 主键
	Title      string `json:"title" orm:"column(title);size(100)"`             // 名称
	Status     int    `json:"status" orm:"column(status);int"`                 // 状态
	PackageIds string `json:"package_ids" orm:"column(package_ids);size(256)"` // 分包ids
	CreateTime int64  `json:"create_time" orm:"column(create_time);bigint"`    // 创建时间
	UpdateTime int64  `json:"update_time" orm:"column(update_time);bigint"`    // 更新时间
}

func CreatePackageGroupModel() *PackageGroupModel {
	return &PackageGroupModel{CreateBase()}
}
