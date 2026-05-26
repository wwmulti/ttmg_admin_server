package models

type DailyStatModel struct {
	*Base
}

type DailyStat struct {
	Id                   int   `json:"id" orm:"auto;column(id)"`                                      // 主键
	Time                 int64 `json:"time" orm:"column(time)"`                                       // 0点时间戳
	UpdateTime           int64 `json:"update_time" orm:"column(update_time)"`                         // 更新时间
	Date                 int   `json:"date" orm:"column(date)"`                                       // 统计日期 (如 20231027)
	PackageId            int   `json:"package_id" orm:"column(package_id)"`                           // 包ID
	PlatformGet          int64 `json:"platform_get" orm:"column(platform_get)"`                       // 平台获得
	PlatformUse          int64 `json:"platform_use" orm:"column(platform_use)"`                       // 平台消耗
	PlatformProfit       int64 `json:"platform_profit" orm:"column(platform_profit)"`                 // 平台盈亏
	PlatformRtp          int64 `json:"platform_rtp" orm:"column(platform_rtp)"`                       // 平台rtp*10000
	GameGet              int64 `json:"game_get" orm:"column(game_get)"`                               // 游戏获得（总产出）
	GameUse              int64 `json:"game_use" orm:"column(game_use)"`                               // 游戏消耗（总投注）
	GameProfit           int64 `json:"game_profit" orm:"column(game_profit)"`                         // 平台游戏盈亏
	GameRtp              int64 `json:"game_rtp" orm:"column(game_rtp)"`                               // 游戏rtp*10000
	GameGetZy            int64 `json:"game_get_zy" orm:"column(game_get_zy)"`                         // 游戏获得（总产出）-自营
	GameUseZy            int64 `json:"game_use_zy" orm:"column(game_use_zy)"`                         // 游戏消耗（总投注）-自营
	GameProfitZy         int64 `json:"game_profit_zy" orm:"column(game_profit_zy)"`                   // 平台游戏盈亏-自营
	TotalChargeAmount    int64 `json:"total_charge_amount" orm:"column(total_charge_amount)"`         // 总充值金额
	TotalChargeMen       int   `json:"total_charge_men" orm:"column(total_charge_men)"`               // 总充值人数
	NewChargeAmount      int64 `json:"new_charge_amount" orm:"column(new_charge_amount)"`             // 新增充值金额
	NewChargeMen         int   `json:"new_charge_men" orm:"column(new_charge_men)"`                   // 新增充值人数
	FirstChargeAmount    int64 `json:"first_charge_amount" orm:"column(first_charge_amount)"`         // 首充金额
	FirstChargeMen       int   `json:"first_charge_men" orm:"column(first_charge_men)"`               // 首充人数
	TotalWithdraw        int64 `json:"total_withdraw" orm:"column(total_withdraw)"`                   // 总提现金额
	TotalWithdrawMen     int   `json:"total_withdraw_men" orm:"column(total_withdraw_men)"`           // 总提现人数
	RealWithdraw         int64 `json:"real_withdraw" orm:"column(real_withdraw)"`                     // 实际提现金额
	ApplyWithdraw        int64 `json:"apply_withdraw" orm:"column(apply_withdraw)"`                   // 申请提现金额
	WithdrawFee          int64 `json:"withdraw_fee" orm:"column(withdraw_fee)"`                       // 提现手续费
	RealDiff             int64 `json:"real_diff" orm:"column(real_diff)"`                             // 实际存提差
	TotalFakeWithdraw    int64 `json:"total_fake_withdraw" orm:"column(total_fake_withdraw)"`         // 总假提现金额
	TotalProfit          int64 `json:"total_profit" orm:"column(total_profit)"`                       // 总盈利(充值-提现)
	WithdrawPercent      int64 `json:"withdraw_percent" orm:"column(withdraw_percent)"`               // 提现占比*10000
	NewUserMen           int   `json:"new_user_men" orm:"column(new_user_men)"`                       // 新增用户人数
	ActiveUserMen        int   `json:"active_user_men" orm:"column(active_user_men)"`                 // 活跃用户人数
	TotalNormalUserMen   int   `json:"total_normal_user_men" orm:"column(total_normal_user_men)"`     // 总正常用户人数
	TotalDisabledUserMen int   `json:"total_disabled_user_men" orm:"column(total_disabled_user_men)"` // 总禁用用户人数
	TotalBalance         int64 `json:"total_balance" orm:"column(total_balance)"`                     // 用户总余额
	AgentBonus           int64 `json:"agent_bonus" orm:"column(agent_bonus)"`                         // 代理返佣金额
	InviteBonus          int64 `json:"invite_bonus" orm:"column(invite_bonus)"`                       // 邀请奖励
	PlatformSend         int64 `json:"platform_send" orm:"column(platform_send)"`                     // 平台赠送金额
	SimulationSend       int64 `json:"simulation_send" orm:"column(simulation_send)"`                 // 模拟赠送金额
	GmSend               int64 `json:"gm_send" orm:"column(gm_send)"`                                 // GM上分
	TotalRechargeSend    int64 `json:"total_recharge_send" orm:"column(total_recharge_send)"`         // 总充值赠送金额
	TotalSignSend        int64 `json:"total_sign_send" orm:"column(total_sign_send)"`                 // 总签到赠送金额
	TotalVipSend         int64 `json:"total_vip_send" orm:"column(total_vip_send)"`                   // 总vip奖励金额
	TotalSpinSend        int64 `json:"total_spin_send" orm:"column(total_spin_send)"`                 // 总转盘赠送金额
	TotalHelpSend        int64 `json:"total_help_send" orm:"column(total_help_send)"`                 // 总救助金额
	TotalInterestSend    int64 `json:"total_interest_send" orm:"column(total_interest_send)"`         // 总利息赠送金额
	PlatformSendMen      int   `json:"platform_send_men" orm:"column(platform_send_men)"`             // 平台赠送人数
	TotalRechargeSendMen int   `json:"total_recharge_send_men" orm:"column(total_recharge_send_men)"` // 总充值赠送人数
	TotalSignSendMen     int   `json:"total_sign_send_men" orm:"column(total_sign_send_men)"`         // 总签到赠送人数
	TotalVipSendMen      int   `json:"total_vip_send_men" orm:"column(total_vip_send_men)"`           // 总vip人数
	TotalSpinSendMen     int   `json:"total_spin_send_men" orm:"column(total_spin_send_men)"`         // 总转盘赠送人数
	TotalHelpSendMen     int   `json:"total_help_send_men" orm:"column(total_help_send_men)"`         // 总救助人数
	TotalInterestSendMen int   `json:"total_interest_send_men" orm:"column(total_interest_send_men)"` // 总利息赠送人数
}

func CreateDailyStatModel() *DailyStatModel {
	return &DailyStatModel{CreateBase()}
}
