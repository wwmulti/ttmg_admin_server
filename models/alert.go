package models

type AlertModel struct {
	*Base
}

type Alert struct {
	Id          int    `json:"id" orm:"auto;column(id)"`                     // 主键，自增
	PackageId   int    `json:"package_id" orm:"column(package_id)"`          // 包ID
	Name        string `json:"name" orm:"column(name)"`                      // 弹窗名称
	Type        int    `json:"type" orm:"column(type);null"`                 // 弹窗类型
	ContentType int    `json:"content_type" orm:"column(content_type);null"` // 内容类型
	EnTitle     string `json:"en_title" orm:"column(en_title)"`              // 英文标题
	PtTitle     string `json:"pt_title" orm:"column(pt_title)"`              // 葡萄牙语标题
	EnContent   string `json:"en_content" orm:"column(en_content)"`          // 英文内容
	PtContent   string `json:"pt_content" orm:"column(pt_content)"`          // 葡萄牙语内容
	Image       string `json:"image" orm:"column(image)"`                    // 图片
	Sort        int    `json:"sort" orm:"column(sort);null"`                 // 排序
	AlertRule   int    `json:"alert_rule" orm:"column(alert_rule);null"`     // 1无限制弹窗 2自然日首次弹窗 3只弹一次 4冷却时间
	AlertHours  int    `json:"alert_hours" orm:"column(alert_hours);null"`   // 冷却小时数
	Status      int    `json:"status" orm:"column(status)"`                  // 0关闭 1开启
	IsDeleted   int    `json:"-" orm:"column(is_deleted)"`                   // 是否删除
}

func CreateAlertModel() *AlertModel {
	return &AlertModel{CreateBase()}
}
