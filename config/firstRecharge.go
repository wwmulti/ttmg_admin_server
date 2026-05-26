package config

// BillingType 统计方式类型
type BillingType int

const (
	BillingDaily BillingType = iota + 1 // 每日
)

type ReceiveType int

const (
	ReceiveNextDayMidnight ReceiveType = iota + 1 // 次日凌晨
)

// TaskType 任务类型
type TaskType int

const (
	TaskBetNumber TaskType = iota + 1 // 投注次数
	TaskBetAmount                     // 投注金额
)
