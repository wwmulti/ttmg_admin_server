package config

import (
	"reflect"
)

type SystemConfigMap struct {
	Key   string
	Value string
}

type SystemConfig struct {
	// 系统配置
	GamePictureDomain             SystemConfigMap `key:"game_picture_domain" value:"https://uploads.wwapi.vip"`    // 图片域名
	UserRegisterReward            SystemConfigMap `key:"user_register_reward" value:"200"`                         // 注册赠送金额
	MaxRegisterNumber             SystemConfigMap `key:"max_register_number" value:"3"`                            // 注册最大次数
	EnterpriseRegisterCode        SystemConfigMap `key:"enterprise_register_code" value:"CNPJ:42.737.630/0001-93"` // 企业注册码
	DefaultPlatformGameShowNumber SystemConfigMap `key:"default_platform_game_show_number" value:"12"`             // 默认平台游戏显示数量
	CustomerLink                  SystemConfigMap `key:"customer_link" value:"https://t.me/KurobaKaze"`            // 客服链接
	// 奖池配置
	InitialValue   SystemConfigMap `key:"initial_value" value:"60000000.00"` // 初始底池金额（例如：60000000）
	GrowthInterval SystemConfigMap `key:"growth_interval" value:"1"`         // 增长时间间隔（固定：1秒）
	GrowthMin      SystemConfigMap `key:"growth_min" value:"6"`              // 单次增长随机最小值（固定：6）
	GrowthMax      SystemConfigMap `key:"growth_max" value:"12"`             // 单次增长随机最大值（固定：12）
	ReduceTimeMin  SystemConfigMap `key:"reduce_time_min" value:"10800"`     // 衰减时间间隔随机最小值（固定：3小时 = 10800秒）
	ReduceTimeMax  SystemConfigMap `key:"reduce_time_max" value:"32400"`     // 衰减时间间隔随机最大值（固定：9小时 = 32400秒）
	ReduceRates    SystemConfigMap `key:"reduce_rates" value:"1,3,5"`        // 衰减比例集合1,3,5（纯随机一个作为这次衰减）
	// 幸运转盘配置
	LuckyWheelRate          SystemConfigMap `key:"lucky_wheel_rate" value:"1"`               // 幸运转盘积分转换比例
	LuckyWheelClearTimeType SystemConfigMap `key:"lucky_wheel_clear_time_type" value:"0"`    // 幸运转盘积分清理时间类型
	LuckyWheelCost          SystemConfigMap `key:"lucky_wheel_cost" value:"1500,5555,25555"` // 幸运转盘消耗白银，黄金，钻石
	// 支付提现配置
	PasswordMaxWrongTime       SystemConfigMap   `key:"password_max_wrong_time" value:"5"`       // 提现密码每日错误最大次数
	MaxBindWithdrawAccountNum  SystemConfigMap   `key:"max_bind_withdraw_account_num" value:"2"` // 最大绑定提现账号数
	AllowUnbindWithdrawAccount SystemConfigMap   `key:"allow_unbind_withdraw_account" value:"0"` // 是否允许解绑提现账号
	PaymentMobileAreaList      map[string]string // 提现账号允许的手机号国家地区
	// 弹窗下载配置
	AlertDownloadUrlFirst  SystemConfigMap `key:"alert_download_url_first" value:"https://www.google.com"` // 弹窗下载连接1
	AlertDownloadUrlSecond SystemConfigMap `key:"alert_download_url_second" value:"https://www.baidu.com"`
	// 代理配置
	AgentScanInterval       SystemConfigMap `key:"agent_scan_interval" value:"5"`          // 代理数据扫描间隔（分钟）
	AgentSettlementCycle    SystemConfigMap `key:"agent_settlement_cycle" value:"1"`       // 代理结算周期（1-天 2-周 3-月）
	AgentSettlementTimeZone SystemConfigMap `key:"agent_settlement_time_zone" value:"UTC"` // 代理结算时区 （UTC / UTC+8 等）
	AgentShareDomain        SystemConfigMap `key:"agent_share_domain" value:""`            // 代理推广域名（可配置多个;号分隔）
	AgentSocialChannels     SystemConfigMap `key:"agent_social_channels" value:""`         // 代理社交分享渠道（可配置多个TikTok,http://www.google.com;号分隔）
	AgentMinAmount          SystemConfigMap `key:"agent_min_amount" value:"1"`             // 代理佣金最低领取金额
	AgentCommissionBetRate  SystemConfigMap `key:"agent_commission_rate" value:"1"`        // 代理佣金领取流水倍数（万分比，0表示无限制）
	// 自动出款配置
	AutoOutMoneyDailySusNums   SystemConfigMap `key:"auto_out_money_daily_sus_nums" value:"1000"` // 当日成功提现次数（小于等于）
	AutoOutMoneyDividePays     SystemConfigMap `key:"auto_out_money_divide_pays" value:"0"`       // 累计提现/累计充值的倍数（小于等于）
	AutoOutMoneyMaxAmount      SystemConfigMap `key:"auto_out_money_max_amount" value:"30000"`    // 提现单笔金额（小于等于）
	AutoOutMoneyMaxIps         SystemConfigMap `key:"auto_out_money_max_ips" value:"1000"`        // 最大同ip数量（小于等于）
	AutoOutMoneySubPays        SystemConfigMap `key:"auto_out_money_sub_pays" value:"0"`          // 充值差金额-提现金额（大于等于）
	AutoOutMoneyChannel        SystemConfigMap `key:"auto_out_money_channel" value:"0"`           // 自动提现代付渠道
	AutoOutMoneySwitch         SystemConfigMap `key:"auto_out_money_switch" value:"0"`            // 自动提现功能（0表示关闭，1表示开启）
	AutoOutMoneyTotalPays      SystemConfigMap `key:"auto_out_money_total_pays" value:"0"`        // 累计充值（大于等于）
	AutoOutMoneyTotalWithdraws SystemConfigMap `key:"auto_out_money_total_withdraws" value:"0"`   // 累计提现（小于等于）
	AutoOutMoneyFirstAuth      SystemConfigMap `key:"auto_out_money_first_auth" value:"1"`        // 首次出款人工审核
	AutoOutMoneyType           SystemConfigMap `key:"auto_out_money_type" value:"1"`              // 自动出款类型(0=系统全自动出款，1=手动加入自动出款)

	// 额外配置
	CurlSecretKey SystemConfigMap `key:"curl_secret_key" value:"jfdDJDAdawLAJWdwa146746DIwdaw"` // curl请求密钥
	RetentionDays []int
}

var System *SystemConfig

// 自动初始化函数
func init() {
	config := &SystemConfig{
		PaymentMobileAreaList: map[string]string{
			"+55": "巴西",
		},
		RetentionDays: []int{
			1, 2, 3, 4, 5, 6, 7, 15, 30,
		},
	}

	// 通过反射读取 tag 自动赋值给 struct
	v := reflect.ValueOf(config).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		if fieldType.Type == reflect.TypeOf(SystemConfigMap{}) && field.CanSet() {
			key := fieldType.Tag.Get("key")
			value := fieldType.Tag.Get("value")
			field.Set(reflect.ValueOf(SystemConfigMap{
				Key:   key,
				Value: value,
			}))
		}
	}

	System = config
}

// SystemConfigType 系统配置类型
type SystemConfigType int

const (
	SystemConfigTypeSystem       SystemConfigType = iota + 1 // 系统基础配置
	SystemConfigTypePrize                                    // 奖池配置
	SystemConfigTypeBet                                      // 投注配置
	SystemConfigTypeRecharge                                 // 充值
	SystemConfigTypeWithdraw                                 // 提现
	SystemConfigTypeLuckWheel                                // 幸运转盘
	SystemConfigTypeAlert                                    // 弹窗下载链接
	SystemConfigTypeAgent                                    // 代理配置
	SystemConfigTypeAutoOutMoney                             // 自动出款配置
)

// LuckyWheelClearTimeType 幸运转盘积分清理时间类型
type LuckyWheelClearTimeType int

const (
	LuckyWheelClearTimeTypeDay   LuckyWheelClearTimeType = iota // 次日零点
	LuckyWheelClearTimeTypeWeek                                 // 周一零点
	LuckyWheelClearTimeTypeMonth                                // 次月1号零点
)

// LuckyWheelType 幸运转盘积分清理时间类型
type LuckyWheelType int

const (
	LuckyWheelTypeSilver  LuckyWheelType = iota // 白银
	LuckyWheelTypeGold                          // 黄金
	LuckyWheelTypeDiamond                       // 砖石
	LuckyWheelTypeMax
)

// GameBetRatePercent 游戏有效投注打码比例(万分比)
const GameBetRatePercent int64 = 10000

// WithdrawFeePercent 提现手续费比例(万分比)
const WithdrawFeePercent int64 = 10000

// BalanceLogType 余额变动日志类型
type BalanceLogType int

const (
	BalanceLogTypeRecharge              BalanceLogType = iota + 1 // 充值
	BalanceLogTypeRechargeExtra                                   // 充值额外赠送
	BalanceLogTypeWithdraw                                        // 提现
	BalanceLogTypeActivityInvite                                  // 邀请活动奖励
	BalanceLogTypeActivitySign                                    // 签到活动奖励
	BalanceLogTypeActivityFirstRecharge                           // 首充活动奖励
	BalanceLogTypeActivityReliefDaily                             // 救济金活动-日奖励
	BalanceLogTypeActivityReliefWeekly                            // 救济金活动-周奖励
	BalanceLogTypeActivityReliefMonthly                           // 救济金活动-月奖励
	BalanceLogTypeActivityInterest                                // 利息宝活动奖励
	BalanceLogTypeActivityLuckyWheel                              // 幸运转盘活动奖励
	BalanceLogTypeActivityVipLv                                   // vip活动晋级奖励
	BalanceLogTypeActivityVipLvDaily                              // vip活动任务-日奖励
	BalanceLogTypeActivityVipLvWeekly                             // vip活动任务-周奖励
	BalanceLogTypeActivityVipLvMonthly                            // vip活动任务-月奖励
	BalanceLogTypeRegister                                        // 注册奖励
	BalanceLogTypeSettlement                                      // 结算类奖励
	BalanceLogTypeAgent                                           // 代理佣金奖励
	BalanceLogTypeSetUserBalance                                  // 修改用户余额
	BalanceLogTypeGmAdd                                           // GM上分
	BalanceLogTypeGmSubtract                                      // GM下分
	BalanceLogTypeAddInviteReward                                 // GM增加代理界面邀请奖励
	BalanceLogTypeSubInviteReward                                 // GM减少代理界面邀请奖励
	BalanceLogTypeAddAdReward                                     // GM广告费下发
	BalanceLogTypeSubAdReward                                     // GM广告费扣除
	BalanceLogTypeAddAccount                                      // GM模拟账户加款
	BalanceLogTypeSubAccount                                      // GM模拟账户减款
)

// 打码忽略的余额变动日志类型
var IgnoreBalanceLogTypes = []BalanceLogType{
	BalanceLogTypeWithdraw,
	BalanceLogTypeRegister,
	BalanceLogTypeSettlement,
	BalanceLogTypeSetUserBalance,
	BalanceLogTypeGmSubtract,
	BalanceLogTypeSubInviteReward,
	BalanceLogTypeSubAdReward,
	BalanceLogTypeSubAccount,
}

// AdminLogType 后台操作日志类型
type AdminLogType int

const (
	AdminLogTypeCommon       AdminLogType = iota + 1 // 常规操作
	AdminLogTypeRefreshTable                         // 刷新表数据
	AdminLogTypeGmOperateAdd                         // 上下分操作
)

// UserType 用户类型
type UserType int

const (
	UserTypeNormal  UserType = iota // 普通用户
	UserTypeBlogger                 // 博主
	UserTypeBroker                  // 经纪人
)

// GmOperateType 上下分操作类型
type GmOperateType int

const (
	GmOperateTypeAdd                  GmOperateType = iota + 1 // 上分
	GmOperateTypeSubtract                                      // 下分
	GmOperateTypeAddInviteReward                               // 增加代理界面邀请奖励
	GmOperateTypeAddWithdrawAbleMoney                          // 增加可提金额
	GmOperateTypeSubWithdrawAbleMoney                          // 减少可提金额
	GmOperateTypeSubInviteReward                               // 减少代理界面邀请奖励
	GmOperateTypeAddAdReward                                   // 广告费下发
	GmOperateTypeSubAdReward                                   // 广告费扣除
	GmOperateTypeAddAccount                                    // 模拟账户加款
	GmOperateTypeSubAccount                                    // 模拟账户减款
)

type CurlRequestType int

const (
	CurlRequestTypeGet  CurlRequestType = iota + 1 // get
	CurlRequestTypePost                            // post
)

type RetentionType int

const (
	RetentionTypeLogin    RetentionType = iota + 1 // 登录
	RetentionTypeRecharge                          // 充值
)
