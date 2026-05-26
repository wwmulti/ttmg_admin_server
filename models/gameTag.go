package models

type GameTagModel struct {
	*Base
}

type GameTag struct {
	Id        int    `json:"id" orm:"auto;column(id)"`            // 主键，自增
	Name      string `json:"name" orm:"column(name)"`             // 标签名称
	PtName    string `json:"pt_name" orm:"column(pt_name)"`       // 葡语名称
	Status    int    `json:"status" orm:"column(status)"`         // 状态 1-开启 0-关闭
	PackageId int    `json:"package_id" orm:"column(package_id)"` // 分包id
}

func CreateGameTagModel() *GameTagModel {
	return &GameTagModel{CreateBase()}
}
