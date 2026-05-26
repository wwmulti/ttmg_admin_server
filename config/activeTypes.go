package config

// ActiveType 活动类型
type ActiveType int

const (
	ActiveTypeInvite        ActiveType = iota + 1 // 邀请注册
	ActiveTypeSign                                // 签到活动
	ActiveTypeFirstRecharge                       // 首充活动
	ActiveTypeRelief                              // 损失援助
	ActiveTypeInterest                            // 利息宝
	ActiveTypeLuckyWheel                          // 幸运转盘活动
	ActiveTypeVipLv                               // VIP等级活动
	ActiveTypeEnd                                 // 结束占位符，用于循环所有类型
)

type ActiveTypeItem struct {
	ID   ActiveType `json:"id"`
	Name string     `json:"name"`
}

var ActiveTypeList = []ActiveTypeItem{
	{ID: ActiveTypeInvite, Name: "邀请注册"},
	{ID: ActiveTypeSign, Name: "签到活动"},
	{ID: ActiveTypeFirstRecharge, Name: "首充活动"},
	{ID: ActiveTypeRelief, Name: "损失援助"},
	{ID: ActiveTypeInterest, Name: "利息宝"},
	{ID: ActiveTypeLuckyWheel, Name: "幸运转盘"},
	{ID: ActiveTypeVipLv, Name: "VIP等级活动"},
}

// SignActiveRewardType 签到活动奖励类型
type SignActiveRewardType int

const (
	SignRewardMoney    SignActiveRewardType = iota + 1 // 现金奖励
	SignRewardIntegral                                 // 积分奖励
	SignRewardDiscount                                 // 优惠券奖励
)

// ActivityReceiveType 领取类型
type ActivityReceiveType int

const (
	ActivityReceiveTypeGet        ActivityReceiveType = iota + 1 // 领取
	ActivityReceiveTypeDistribute                                // 分配
)

// LuckyWheelRewardType 幸运转盘奖励类型
type LuckyWheelRewardType int

const (
	LuckyWheelSliver  LuckyWheelRewardType = iota // 银
	LuckyWheelGold                                // 金
	LuckyWheelDiamond                             // 钻石
)

// VipRuleCondition Vip规则条件
type VipRuleCondition int

const (
	VipRuleConditionOr       VipRuleCondition = iota // 满足任意条件
	VipRuleConditionAnd                              // 满足所有条件
	VipRuleConditionRecharge                         // 满足充值条件
	VipRuleConditionBet                              // 满足有效押注条件
)

// VipTaskCondition Vip任务条件
type VipTaskCondition int

const (
	VipTaskConditionOr       VipTaskCondition = iota // 满足任意条件
	VipTaskConditionAnd                              // 满足所有条件
	VipTaskConditionRecharge                         // 满足充值条件
	VipTaskConditionBet                              // 满足有效押注条件
)

// VipTaskCycleType Vip奖励类型
type VipTaskCycleType int

const (
	VipTaskCycleTypeDay   VipTaskCycleType = iota + 1 // 日
	VipTaskCycleTypeWeek                              // 周
	VipTaskCycleTypeMonth                             // 月
)

// ReliefCycleType 救济活动统计周期
type ReliefCycleType int

const (
	ReliefCycleTypeDay   ReliefCycleType = iota + 1 // 日
	ReliefCycleTypeWeek                             // 周
	ReliefCycleTypeMonth                            // 月
)

// ShareScopeType 邀请活动统计范围
type ShareScopeType int

const (
	ShareScopeTypeChild ShareScopeType = iota // 直属下级
	ShareScopeTypeAll                         // 所有下级
)

// ShareExpireType 邀请活动奖励过期策略
type ShareExpireType int

const (
	ShareExpireTypeAutoReceive ShareExpireType = iota // 自动领取
	ShareExpireTypeInvalid                            // 作废
)

// ShareConditionType 邀请活动条件
type ShareConditionType int

const (
	ShareConditionTypeAnd ShareConditionType = iota // 满足所有条件
	ShareConditionTypeOr                            // 满足任意条件
)

// ShareRewardType 邀请活动奖励领取方式
type ShareRewardType int

const (
	ShareRewardTypeUserReceive ShareRewardType = iota // 手动领取
	ShareRewardTypeAutoReceive                        // 系统自动派发
)
