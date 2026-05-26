package models

type MerchantModel struct {
	*Base
}

type Merchant struct {
	Id           int    `json:"id" orm:"auto;column(id)"`                              // 主键
	Title        string `json:"title" orm:"column(title);size(100);not null"`          // 商户名称
	Domain       string `json:"domain" orm:"column(domain);size(256);not null"`        // api地址
	Token        string `json:"token" orm:"column(token);size(256);not null"`          // token
	Secret       string `json:"secret" orm:"column(secret);size(256);not null"`        // 秘钥
	Currency     string `json:"currency" orm:"column(currency);size(50);default(BRL)"` // 币种
	Type         int    `json:"type" orm:"column(type);not null"`                      // 游戏类型 pp pg
	Status       int    `json:"status" orm:"column(status);default(1)"`                // 是否禁用 1启用 0禁用
	SupplierType int    `json:"supplier_type" orm:"column(supplier_type);default(1)"`  // 供应商类型 1自研 2官方
	IsDeleted    int    `json:"-" orm:"column(is_deleted);default(0)"`                 // 是否删除 0未删 1删除
}

func CreateMerchantModel() *MerchantModel {
	return &MerchantModel{CreateBase()}
}
