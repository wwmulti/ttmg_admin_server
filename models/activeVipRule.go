package models

type ActiveVipRuleModel struct {
	*Base
}

type ActiveVipRule struct {
	Id                  int     `json:"id" orm:"auto;column(id)"`                                   // 主键
	ActiveId            int     `json:"active_id" orm:"column(active_id)"`                          // 活动id
	Lv                  int     `json:"lv" orm:"column(lv)"`                                        // 等级
	TotalPays           float64 `json:"total_pays" orm:"column(total_pays);"`                       // 累计充值
	TotalBets           float64 `json:"total_bets" orm:"column(total_bets);"`                       // 累计有效押注
	ConAnd              int     `json:"con_and" orm:"column(con_and)"`                              // 条件：0-或 1-且 2-只需充值 3-只需流水
	Rewards             float64 `json:"rewards" orm:"column(rewards);"`                             // 晋级金
	WithdrawNumLimit    int     `json:"withdraw_num_limit" orm:"column(withdraw_num_limit)"`        // 每日提现次数限制
	WithdrawAmountLimit float64 `json:"withdraw_amount_limit" orm:"column(withdraw_amount_limit);"` // 每日提现金额限制
	WithdrawFreeNum     int     `json:"withdraw_free_num" orm:"column(withdraw_free_num)"`          // 每日免费交易次数
	WithdrawFee         int     `json:"withdraw_fee" orm:"column(withdraw_fee)"`                    // 提现手续，万分比
	Status              int     `json:"status" orm:"column(status)"`                                // 状态：0-关闭 1-启用
	Ctime               int64   `json:"c_time" orm:"column(c_time)"`                                // 时间
}

func CreateActiveVipRuleModel() *ActiveVipRuleModel {
	return &ActiveVipRuleModel{CreateBase()}
}

// 主界面&每日/周/月奖励界面
type VipRuleRewardDTO struct {
	Header  VipRuleRewardHeaderDTO   `json:"header"`
	Details []VipRuleRewardDetailDTO `json:"details"`
}

// Vip特权界面
type VipRulePrivilegeDTO struct {
	Header  VipRuleRewardHeaderDTO      `json:"header"`
	Details []VipRulePrivilegeDetailDTO `json:"details"`
}

type VipRuleRewardHeaderDTO struct {
	Lv         int     `json:"lv"`          //当前等级
	NextLv     int     `json:"next_lv"`     //下一级别
	ConAnd     int     `json:"con_and"`     //条件：0-或 1-且 2-只需充值 3-只需流水
	RemainPays float64 `json:"remain_pays"` //剩余多少充值晋级下一级
	RemainBets float64 `json:"remain_bets"` //剩余多少押注晋级下一级
}

type VipRuleRewardDetailDTO struct {
	Lv          int     `json:"lv"`           //等级
	Cycle       int     `json:"cycle"`        //周期类型 0-晋级金 1-日 2-周 3-月
	Period      int     `json:"period"`       //周期
	ConAnd      int     `json:"con_and"`      //条件：0-或 1-且 2-只需充值 3-只需流水
	TotalPays   float64 `json:"total_pays"`   //累计充值条件
	TotalBets   float64 `json:"total_bets"`   //累计押注条件
	CurrentPays float64 `json:"current_pays"` //充值进度
	CurrentBets float64 `json:"current_bets"` //押注进度
	Rewards     float64 `json:"rewards"`      //奖励
	Status      int     `json:"status"`       //状态: 0-不可领取 1-可领取 2-已领取
}

type VipRulePrivilegeDetailDTO struct {
	Lv                  int     `json:"lv"`                    //等级
	WithdrawAmountLimit float64 `json:"withdraw_amount_limit"` //提现金额限制
	WithdrawNumLimit    int     `json:"withdraw_num_limit"`    //提现次数限制
	WithdrawFreeNum     int     `json:"withdraw_free_num"`     //免费提现次数
	WithdrawFee         int     `json:"withdraw_fee"`          //提现手续费
	Status              int     `json:"status"`                //状态：0-关闭 1-启用
}

// 奖励领取历史记录
type VipRuleRewardLogDTO struct {
	Lv     int     `json:"lv"`     //等级
	Amount float64 `json:"amount"` //金额
	Cycle  int     `json:"cycle"`  //类型 0-晋级金 1-日 2-周 3-月
	Ctime  int64   `json:"c_time"` //时间
}

type ActiveVipRuleDTO struct {
	List  []ActiveVipRule `json:"list"`
	Total int64           `json:"total"`
}
