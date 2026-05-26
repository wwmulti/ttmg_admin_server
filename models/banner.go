package models

type BannerModel struct {
	*Base
}

type Banner struct {
	Id        int    `json:"id" orm:"auto;column(id)"`            // 主键，自增
	Name      string `json:"name" orm:"column(name);null"`        // banner名称
	Type      int    `json:"type" orm:"column(type);null"`        // banner类型
	Url       string `json:"url" orm:"column(url);null"`          // 图片url
	PackageId int    `json:"package_id" orm:"column(package_id)"` // 包ID
	Sort      int    `json:"sort" orm:"column(sort);null"`        // 排序
	Status    int    `json:"status" orm:"column(status)"`         // 0关闭 1开启
	IsDeleted int    `json:"-" orm:"column(is_deleted)"`          // 是否删除
}

func CreateBannerModel() *BannerModel {
	return &BannerModel{CreateBase()}
}
