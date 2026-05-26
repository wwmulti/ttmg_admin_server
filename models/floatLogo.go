package models

type FloatLogoModel struct {
	*Base
}

type FloatLogo struct {
	Id        int    `json:"id" orm:"auto;column(id)"`            // 主键，自增
	Logo      string `json:"logo" orm:"column(logo);null"`        // 图标
	Link      string `json:"link" orm:"column(link);null"`        // 跳转连接
	PackageId int    `json:"package_id" orm:"column(package_id)"` // 包ID
	Status    int    `json:"status" orm:"column(status)"`         // 0关闭 1开启
	IsDeleted int    `json:"-" orm:"column(is_deleted)"`          // 是否删除
}

func CreateFloatLogoModel() *FloatLogoModel {
	return &FloatLogoModel{CreateBase()}
}
