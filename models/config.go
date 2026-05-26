package models

type ConfigModel struct {
	*Base
}

type Config struct {
	Id          int    `json:"id" orm:"auto;column(id)"`                           // 主键，自增
	TypeId      int    `json:"type_id" orm:"column(type_id)"`                      // 类型ID
	PackageId   int    `json:"package_id" orm:"column(package_id)"`                // 包ID
	ConfigKey   string `json:"config_key" orm:"column(config_key);unique"`         // 配置 key，唯一
	ConfigValue string `json:"config_value" orm:"column(config_value);type(text)"` // 配置值
}

type GetPrizeDTO struct {
	PrizeBalance             int            `json:"prize_balance"`               // 底池金额
	PrizeBalanceChangeAmount map[int]string `json:"prize_balance_change_amount"` // 奖池金额变动
}

func CreateConfigModel() *ConfigModel {
	return &ConfigModel{CreateBase()}
}
