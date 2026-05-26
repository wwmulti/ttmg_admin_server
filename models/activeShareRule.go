package models

type ActiveShareRuleModel struct {
	*Base
}

type ActiveShareRule struct {
	Id             int     `json:"id" orm:"auto;column(id)"`                        // ID
	ActiveId       int     `json:"active_id" orm:"column(active_id)"`               // 活动ID
	ActiveTypeId   int     `json:"active_type_id" orm:"column(active_type_id)"`     // 活动类型ID
	ShareUrl       string  `json:"share_url" orm:"column(share_url)"`               // 推广域名
	ScanInterval   int     `json:"scan_interval" orm:"column(scan_interval)"`       // 扫描用户下级间隔（分钟）
	Scope          int     `json:"scope" orm:"column(scope)"`                       // 统计范围：0-仅直属下级 1-所有层级下级
	MiniTotalPays  float64 `json:"mini_total_pays" orm:"column(mini_total_pays)"`   // 下级最低累计充值,0为无限制
	MiniTotalWater float64 `json:"mini_total_water" orm:"column(mini_total_water)"` // 押注消费流水，0为无限制
	Condition      int     `json:"condition" orm:"column(condition)"`               // 判断关系，0-同时满足，1-满足任一
	RewardType     int     `json:"reward_type" orm:"column(reward_type)"`           // 奖励领取方式 ，0-手动点击，1-系统自动派发
	ExpireType     int     `json:"expire_type" orm:"column(expire_type)"`           // 过期处理策略，0-结束未领取自动派发，1-过期作废
	Status         int     `json:"status" orm:"column(status)"`                     // 状态：0-关闭，1-开启
	IsDeleted      int     `json:"is_deleted" orm:"column(is_deleted)"`             // 是否删除
	Ctime          int64   `json:"c_time" orm:"column(c_time)"`                     // 创建时间
}

func CreateActiveShareRuleModel() *ActiveShareRuleModel {
	return &ActiveShareRuleModel{CreateBase()}
}

type ActiveShareRuleDTO struct {
	List  []ActiveShareRule `json:"list"`
	Total int64             `json:"total"`
}
