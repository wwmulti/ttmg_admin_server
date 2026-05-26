package config

// UserInterestLogType 用户利息日志类型
type UserInterestLogType int

const (
	UserInterestLogTypeDeposit           UserInterestLogType = iota + 1 // 存入
	UserInterestLogTypeWithdraw                                         // 取出
	UserInterestLogTypeCalculateInterest                                // 领取
)

// InterestRuleReceiveType 可领取的时间类型
type InterestRuleReceiveType int

const (
	InterestRuleReceiveTypeRealTime      InterestRuleReceiveType = iota + 1 // 实时
	InterestRuleReceiveTypeNextDayZero                                      // 第二天
	InterestRuleReceiveTypeNextWeekZero                                     // 下周一
	InterestRuleReceiveTypeNextMonthZero                                    // 下月1号
)

// InterestRuleInterestLimitType 利息最大限制类型
type InterestRuleInterestLimitType int

const (
	InterestRuleInterestLimitTypePercent    InterestRuleInterestLimitType = iota + 1 // 万分比
	InterestRuleInterestLimitTypeFixedValue                                          // 固定值
	InterestRuleInterestLimitTypeUnlimited                                           // 不限
)
