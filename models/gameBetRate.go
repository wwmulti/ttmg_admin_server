package models

// --- 游戏有效投注打码配置表 ---

type GameBetRateModel struct {
	*Base
}

type GameBetRate struct {
	Id           int    `json:"id" orm:"auto;column(id)"`                    // 主键
	GameTypeId   int    `json:"game_type_id" orm:"column(game_type_id)"`     // 游戏类型
	BetRate      int    `json:"bet_rate" orm:"column(bet_rate)"`             // 默认打码比例（万分比）
	GamesBetRate string `json:"games_bet_rate" orm:"column(games_bet_rate)"` // 配置的游戏打码比例
}

func CreateGameBetRateModel() *GameBetRateModel {
	return &GameBetRateModel{CreateBase()}
}
