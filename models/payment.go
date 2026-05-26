package models

type PaymentModel struct {
	*Base
}

type Payment struct {
	Id             int    `json:"id" orm:"auto;column(id)"`                                      // 主键，自增
	PayCode        string `json:"pay_code" orm:"column(pay_code)"`                               // 通道编码
	MerchantConfig string `json:"merchant_config" orm:"column(merchant_config);type(json);null"` // 商户配置
	Logo           string `json:"logo" orm:"column(logo)"`                                       // 支付图标
	Remark         string `json:"remark" orm:"column(remark);null"`                              // 备注
	Status         int    `json:"status" orm:"column(status)"`                                   // 通道状态 1开启 0关闭
	IsDeleted      int    `json:"is_deleted" orm:"column(is_deleted)"`                           // 是否删除
}

func CreatePaymentModel() *PaymentModel {
	return &PaymentModel{CreateBase()}
}
