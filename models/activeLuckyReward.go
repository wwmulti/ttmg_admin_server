package models

type ActiveLuckyRewardModel struct {
	*Base
}

type ActiveLuckyReward struct {
	Id         int     `json:"id" orm:"auto;column(id)"`              // 主键，自增
	ActivityId int     `json:"activity_id" orm:"column(activity_id)"` // 活动ID
	WheelType  int     `json:"wheel_type" orm:"column(wheel_type)"`   // 转盘类型
	Reward     float64 `json:"reward" orm:"column(reward)"`           // 奖励
	Weight     int     `json:"weight" orm:"column(weight)"`           // 权重
	CreateTime int64   `json:"-" orm:"column(create_time)"`           // 创建时间
	UpdateTime int64   `json:"-" orm:"column(update_time)"`           // 更新时间
	IsDeleted  int64   `json:"-" orm:"column(is_deleted)"`            // 是否删除
}

type LuckyWheelUserInfo struct {
	BetAmount   float64 `json:"bet_amount"`   // 下注金额
	ExpireScore float64 `json:"expire_score"` // 过期积分
	TotalScore  float64 `json:"total_score"`  // 总积分
	CostScore   float64 `json:"cost_score"`   // 花费积分
	RemainScore float64 `json:"remain_score"` // 剩余积分
}

type LuckyWheelDataResponse struct {
	List     map[int][]ActiveLuckyReward `json:"list"`
	UserInfo LuckyWheelUserInfo          `json:"user_info"`
}

func CreateActiveLuckyRewardModel() *ActiveLuckyRewardModel {
	return &ActiveLuckyRewardModel{CreateBase()}
}
