package models

type RechargeOrderModel struct {
	*Base
}

type RechargeOrder struct {
	Id               int     `json:"id" orm:"auto;column(id)"`                              // 主键，自增
	UserId           int     `json:"user_id" orm:"column(user_id)"`                         // 用户id
	RoleId           int     `json:"role_id" orm:"column(role_id)"`                         // 用户RoleId
	Username         string  `json:"username" orm:"column(username)"`                       // 用户名
	PaymentId        int     `json:"payment_id" orm:"column(payment_id)"`                   // 支付通道id
	PaymentChannelId int     `json:"payment_channel_id" orm:"column(payment_channel_id)"`   // 支付渠道id
	OrderNo          string  `json:"order_no" orm:"column(order_no);size(50)"`              // 订单号
	OutOrderNo       string  `json:"out_order_no" orm:"column(out_order_no);size(50);null"` // 外部订单号
	Amount           float64 `json:"amount" orm:"column(amount)"`                           // 订单金额
	Status           int     `json:"status" orm:"column(status)"`                           // 订单状态 0未支付 1已支付 2已回调
	IsFirst          int     `json:"is_first" orm:"column(is_first)"`                       // 是否首充 0-否 1-是
	CreatedAt        int64   `json:"created_at" orm:"column(created_at)"`                   // 创建时间
	PayAt            int64   `json:"pay_at" orm:"column(pay_at);null"`                      // 支付完成时间
	PackageId        int     `json:"package_id" orm:"column(package_id)"`                   // 包ID
	RegisterAt       int64   `json:"registerAt" orm:"column(register_at);"`                 // 用户注册时间
}

func CreateRechargeOrderModel() *RechargeOrderModel {
	return &RechargeOrderModel{CreateBase()}
}
