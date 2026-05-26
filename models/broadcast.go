package models

type BroadcastModel struct {
	*Base
}

type Broadcast struct {
	Id        int    `json:"id" orm:"auto;column(id)"`            // 主键，自增
	PackageId int    `json:"package_id" orm:"column(package_id)"` // 包ID
	Name      string `json:"name" orm:"column(name)"`             // 名称
	EnContent string `json:"en_content" orm:"column(en_content)"` // 英语广播内容
	PtContent string `json:"pt_content" orm:"column(pt_content)"` // 葡语广播内容
	Status    int    `json:"status" orm:"column(status)"`         // 状态
	Sort      int    `json:"sort" orm:"column(sort)"`             // 排序
	IsDeleted int    `json:"-" orm:"column(is_deleted)"`          // 是否删除
}

func CreateBroadcastModel() *BroadcastModel {
	return &BroadcastModel{CreateBase()}
}
