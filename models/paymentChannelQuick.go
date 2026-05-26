package models

type PaymentChannelQuickModel struct {
	*Base
}

type PaymentChannelQuick struct {
	Id               int     `json:"id" orm:"auto;column(id)"`                            // 主键，自增
	PaymentChannelId int     `json:"payment_channel_id" orm:"column(payment_channel_id)"` // 支付渠道id
	Amount           float64 `json:"amount" orm:"column(amount)"`                         // 金额数值
	IsRecommend      int     `json:"is_recommend" orm:"column(is_recommend)"`             // 是否推荐 0否 1是
	Sort             int     `json:"sort" orm:"column(sort)"`
}

func CreatePaymentChannelQuickModel() *PaymentChannelQuickModel {
	return &PaymentChannelQuickModel{CreateBase()}
}
