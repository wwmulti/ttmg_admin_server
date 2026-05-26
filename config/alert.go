package config

// AlertRuleType 弹窗规则类型
type AlertRuleType int

const (
	Always        AlertRuleType = iota + 1 // 无限制弹窗
	DayOneTime                             // 每日弹窗一次
	DeviceOnTime                           // 设备弹窗一次
	PeriodOneTime                          // 周期弹窗一次
)

// AlertType 弹窗类型
type AlertType int

const (
	DownloadAlert      AlertType = iota + 1 // 下载弹窗
	MultiTagAlert                           // 多标签弹窗
	FirstRechargeAlert                      // 首充活动弹窗
)

// AlertContentType 内容类型
type AlertContentType int

const (
	TextAreaContent AlertContentType = iota + 1 // 富文本内容
	ImageContent                                // 图片内容
)

type BannerType int

const (
	BannerTypeIndex BannerType = iota + 1 // 首页轮播图
)
