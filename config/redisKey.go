package config

type RedisKey struct {
	UserMoney                   string // 用户金额
	SystemAutoDistributeRewards string // 系统自动分配奖励
	GameRecent                  string // 近期游戏
	GameFavorite                string // 收藏游戏
	SystemConfig                string // 系统配置
	PrizeBalance                string // 奖池金额
	PrizeBalanceQueue           string // 奖池变动队列
	PrizeBalanceAttenuationTime string // 奖池金额衰减
	ActiveLuckyWheel            string // 抽奖
	GameBetRate                 string // 游戏投注打码比例
	RateLimit                   string // 接口限流
	GameInfo                    string // 游戏信息
	ActivesInfo                 string // 活动信息
	UserGoogleSecret            string // 用戶google秘钥
	UserAttr                    string // 用户属性信息
	TaskProcess                 string // 异步任务处理状态
}

var RedisKeyName = RedisKey{
	UserMoney:                   "user_money",
	SystemAutoDistributeRewards: "system_auto_distribute_rewards",
	GameRecent:                  "game_recent_",
	GameFavorite:                "game_favorite_",
	SystemConfig:                "system_config",
	PrizeBalance:                "prize_balance",
	PrizeBalanceQueue:           "prize_balance_queue",
	PrizeBalanceAttenuationTime: "prize_balance_attenuation_time",
	ActiveLuckyWheel:            "active_lucky_wheel_score_hash",
	GameBetRate:                 "game_bet_rate:",
	RateLimit:                   "rate_limit:",
	GameInfo:                    "cache:game_info",
	ActivesInfo:                 "cache:actives_info",
	UserGoogleSecret:            "user_google_secrect",
	UserAttr:                    "user_attr:",
	TaskProcess:                 "task_process:",
}

const (
	UserAttr_PackageId       = "package_id"        // 包id
	UserAttr_TotalBet        = "total_bet"         // 累计投注金额
	UserAttr_TotalWin        = "total_win"         // 累计赢金
	UserAttr_TotalRecharge   = "total_recharge"    // 累计充值金额
	UserAttr_TotalWithdraw   = "total_withdraw"    // 累计提现金额
	UserAttr_AwardNeedBets   = "award_need_bets"   // 奖励打码
	UserAttr_EffectiveBet    = "effective_bet"     // 累计有效押注
	UserAttr_VipLevel        = "vip_level"         // Vip等级
	UserAttr_TotalGmSend     = "total_gm_send"     // GM上分数量
	UserAttr_TotalAgentBonus = "total_agent_bonus" // 代理佣金领取
	UserAttr_Parents         = "agent_parents"     // 代理关系
)
