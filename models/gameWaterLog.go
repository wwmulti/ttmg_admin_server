package models

type GameWaterLogModel struct {
	*Base
}

type GameWaterLog struct {
	Id                int     `json:"id" orm:"auto;column(id)"`                              // 主键，自增
	UserId            int     `json:"user_id" orm:"column(user_id)"`                         // 用户id
	RoleId            int     `json:"role_id" orm:"column(role_id);null"`                    // 用户RoleId
	GameId            int     `json:"game_id" orm:"column(game_id)"`                         // 游戏id
	GameTitle         string  `json:"game_title" orm:"column(game_title)"`                   // 游戏名称
	GetAmount         float64 `json:"get_amount" orm:"column(get_amount)"`                   // 单局获得
	BetAmount         float64 `json:"bet_amount" orm:"column(bet_amount)"`                   // 单局押注
	CTime             int64   `json:"c_time" orm:"column(c_time)"`                           // 时间
	GamePlatformId    int     `json:"game_platform_id" orm:"column(game_platform_id)"`       // 游戏平台id
	GamePlatformTitle string  `json:"game_platform_title" orm:"column(game_platform_title)"` // 平台名称
	GameCatId         int     `json:"game_cat_id" orm:"column(game_cat_id)"`                 // 游戏分类id
	GameTypeTitle     string  `json:"game_type_title" orm:"column(game_type_title)"`         // 游戏类型名称
	OrderId           string  `json:"order_id" orm:"column(order_id)"`                       // 订单号
	ParentOrderId     string  `json:"parent_order_id" orm:"column(parent_order_id)"`         // 主订单号
	IsEnd             int     `json:"is_end" orm:"column(is_end)"`                           // 是否最后一局
	PackageId         int     `json:"package_id" orm:"column(package_id)"`                   // 包id
}

func CreateGameWaterLogModel() *GameWaterLogModel {
	return &GameWaterLogModel{CreateBase()}
}
