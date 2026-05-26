package models

type GameTypeModel struct {
	*Base
}

type GameType struct {
	Id             int    `json:"id" orm:"auto;column(id)"`                        // 主键，自增
	Name           string `json:"name" orm:"column(name)"`                         // 分类名称
	Alias          string `json:"alias" orm:"column(alias)"`                       // 简称
	Logo           string `json:"logo" orm:"column(logo)"`                         // 分类图标
	Sort           int    `json:"sort" orm:"column(sort)"`                         // 排序权重
	AgentConfig    string `json:"agent_config" orm:"column(agent_config)"`         // 代理配置
	Status         int    `json:"status" orm:"column(status)"`                     // 分类状态 0关闭 1开启
	IsDeleted      int    `json:"is_deleted" orm:"column(is_deleted)"`             // 是否删除
	PageGameNumber int    `json:"page_game_number" orm:"column(page_game_number)"` // 显示单页游戏数量
	BetRate        int    `json:"bet_rate" orm:"column(bet_rate)"`                 // 游戏类型默认打码比例（万分比）
	GamesBetRate   string `json:"games_bet_rate" orm:"column(games_bet_rate)"`     // 配置的游戏打码比例 格式 id,rate; 拼接
	PackageId      int    `json:"package_id" orm:"column(package_id)"`             // 分包id
	PlatformIds    string `json:"platform_ids" orm:"column(platform_ids)"`         // 平台id字符串
}

func CreateGameTypeModel() *GameTypeModel {
	return &GameTypeModel{CreateBase()}
}
