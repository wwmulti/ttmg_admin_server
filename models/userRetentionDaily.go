package models

type UserRetentionDailyModel struct {
	*Base
}

type UserRetentionDaily struct {
	Id            int     `json:"id" orm:"auto;column(id)"`                               // 主键，自增
	PackageId     int     `json:"package_id" orm:"column(package_id)"`                    // 包ID
	Type          int     `json:"type" orm:"column(type)"`                                // 类型 1登录 2充值
	StatisticDate string  `json:"statistic_date" orm:"column(statistic_date);type(date)"` // 统计日期
	RegisterCount int     `json:"register_count" orm:"column(register_count)"`            // 注册用户数
	Day1          float64 `json:"day_1" orm:"column(day_1);digits(5);decimals(2)"`        // 次日留存率
	Day2          float64 `json:"day_2" orm:"column(day_2);digits(5);decimals(2)"`        // 第2日留存率
	Day3          float64 `json:"day_3" orm:"column(day_3);digits(5);decimals(2)"`        // 第3日留存率
	Day4          float64 `json:"day_4" orm:"column(day_4);digits(5);decimals(2)"`        // 第4日留存率
	Day5          float64 `json:"day_5" orm:"column(day_5);digits(5);decimals(2)"`        // 第5日留存率
	Day6          float64 `json:"day_6" orm:"column(day_6);digits(5);decimals(2)"`        // 第6日留存率
	Day7          float64 `json:"day_7" orm:"column(day_7);digits(5);decimals(2)"`        // 第7日留存率
	Day15         float64 `json:"day_15" orm:"column(day_15);digits(5);decimals(2)"`      // 第15日留存率
	Day30         float64 `json:"day_30" orm:"column(day_30);digits(5);decimals(2)"`      // 第30日留存率
	CreatedAt     int64   `json:"created_at" orm:"column(created_at);null"`               // 创建时间
}

func CreateUserRetentionDailyModel() *UserRetentionDailyModel {
	return &UserRetentionDailyModel{CreateBase()}
}
