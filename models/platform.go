package models

type PlatformModel struct {
	*Base
}

type Platform struct {
	Id            int     `json:"id" orm:"auto;column(id)"`                      // 主键，自增
	Name          string  `json:"name" orm:"column(name)"`                       // 平台名称
	Logo          string  `json:"logo" orm:"column(logo)"`                       // 平台logo
	ClickLogo     string  `json:"click_logo" orm:"column(click_logo)"`           // 点击后图标
	Image         string  `json:"image" orm:"column(image)"`                     // 供应商图片
	FrontColor    string  `json:"front_color" orm:"column(front_color)"`         // 前端显示颜色
	MiniMoney     float64 `json:"mini_money" orm:"column(mini_money)"`           // 最小金额
	ApiRate       float64 `json:"api_rate" orm:"column(api_rate)"`               // API费率（单位%）
	Alias         string  `json:"alias" orm:"column(alias)"`                     // 供应商代号
	Description   string  `json:"description" orm:"column(description)"`         // 平台描述
	Status        int     `json:"status" orm:"column(status)"`                   // 平台状态 1正常 0禁用
	Sort          int     `json:"sort" orm:"column(sort)"`                       // 排序权重
	GameShowCount int     `json:"game_show_count" orm:"column(game_show_count)"` // 初始游戏展示数量
	IsDeleted     int     `json:"is_deleted" orm:"column(is_deleted)"`           // 是否删除
	GameShowMore  int     `json:"game_show_more" orm:"column(game_show_more)"`   // 加载更多每页游戏数量
	PackageId     int     `json:"package_id" orm:"column(package_id)"`           // 分包id
}

type PlatformDTO struct {
	Id            int    `json:"id"`                                    // 平台id
	Name          string `json:"name"`                                  // 平台名称
	Logo          string `json:"logo"`                                  // 平台logo
	Description   string `json:"description" orm:"column(description)"` // 平台描述
	Status        int    `json:"status" orm:"column(status)"`           // 平台状态 1正常 0禁用
	Sort          int    `json:"sort"`                                  // 排序权重
	GameShowCount int    `json:"game_show_count"`                       // 初始游戏展示数量
	GameShowMore  int    `json:"game_show_more"`                        // 加载更多每页游戏数量
}

type GameTypesDTO struct {
	Id             int    `json:"id"`               // 游戏类型id
	Name           string `json:"name"`             // 游戏类型名称
	Logo           string `json:"logo"`             // 游戏类型logo
	Sort           int    `json:"sort"`             // 排序权重
	PageGameNumber int    `json:"page_game_number"` // 每页游戏数量
}

type SupporterPlatformGamesDTO struct {
	PlatformList []SupportPlatformDTO `json:"platform_list"` // 平台列表
	GameList     []Game               `json:"game_list"`     // 游戏列表
	Total        int                  `json:"total"`         // 总数量
}

type SupportPlatformDTO struct {
	Id             int    `json:"id"`               // 平台id
	Name           string `json:"name"`             // 平台名称
	Logo           string `json:"logo"`             // 平台logo
	Sort           int    `json:"sort"`             // 排序权重
	PageGameNumber int    `json:"page_game_number"` // 每页游戏数量
}

func CreatePlatformModel() *PlatformModel {
	return &PlatformModel{CreateBase()}
}
