package models

type UserAttrModel struct {
	*Base
}

type UserAttr struct {
	Id                 int   `json:"id" orm:"auto;column(id)"`                                // 主键
	AwardNeedBets      int64 `json:"award_need_bets" orm:"column(award_need_bets)"`           // 奖励打码
	EffectiveBet       int64 `json:"effective_bet" orm:"column(effective_bet)"`               // 累计有效押注(总打码值)
	TotalRecharge      int64 `json:"total_recharge" orm:"column(total_recharge)"`             // 累计充值
	TotalWithdraw      int64 `json:"total_withdraw" orm:"column(total_withdraw)"`             // 累计提现
	TotalRechargeCount int   `json:"total_recharge_count" orm:"column(total_recharge_count)"` // 累计充值次数
	TotalWithdrawCount int   `json:"total_withdraw_count" orm:"column(total_withdraw_count)"` // 累计提现次数
	TotalProfit        int64 `json:"total_profit" orm:"column(total_profit)"`                 // 累计盈亏
	VipLevel           int   `json:"vip_level" orm:"column(vip_level)"`                       // VIP等级
	User               *User `json:"user" orm:"rel(fk);column(user_id)"`                      // 关联用户
	TotalGmSend        int64 `json:"total_gm_send" orm:"column(total_gm_send)"`               // 累计GM发送
	TotalAgentBonus    int64 `json:"total_agent_bonus" orm:"column(total_agent_bonus)"`       // 累计代理奖励
}

func CreateUserAttrModel() *UserAttrModel {
	return &UserAttrModel{CreateBase()}
}
