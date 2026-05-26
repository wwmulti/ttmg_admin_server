package models

type PaymentChannelModel struct {
	*Base
}

type PaymentChannel struct {
	Id            int     `json:"id" orm:"auto;column(id)"`                                    // 主键，自增
	PaymentId     int     `json:"payment_id" orm:"column(payment_id)"`                         // 支付id
	PackageId     int     `json:"package_id" orm:"column(package_id)"`                         // 包id
	PaymentTypeId int     `json:"payment_type_id" orm:"column(payment_type_id)"`               // 支付类型id
	Name          string  `json:"name" orm:"column(name)"`                                     // 渠道名称
	PrizePercent  int     `json:"prize_percent" orm:"column(prize_percent)"`                   // 赠送比例
	LimitSmall    float64 `json:"limit_small" orm:"column(limit_small);null"`                  // 支付最小金额
	LimitBig      float64 `json:"limit_big" orm:"column(limit_big);null"`                      // 支付最大金额
	ChannelConfig string  `json:"channel_config" orm:"column(channel_config);type(json);null"` // 支付配置信息
	VipLevel      int     `json:"vip_level" orm:"column(vip_level);null"`                      // vip等级
	ChannelType   int     `json:"channel_type" orm:"column(channel_type)"`                     // 通道类型 1支付 2提现
	Tag           string  `json:"tag" orm:"column(tag);null"`                                  // 标签文案
	Status        int     `json:"status" orm:"column(status)"`                                 // 状态 0关闭 1开启
	IsDeleted     int     `json:"is_deleted" orm:"column(is_deleted)"`                         // 是否删除
}

func CreatePaymentChannelModel() *PaymentChannelModel {
	return &PaymentChannelModel{CreateBase()}
}
