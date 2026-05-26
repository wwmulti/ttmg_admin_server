package models

type GameModel struct {
	*Base
}

type Game struct {
	Id           int    `json:"id" orm:"auto;column(id)"`                    // 主键，自增
	PlatformId   int    `json:"platform_id" orm:"column(platform_id)"`       // 所属平台id
	PlatformName string `json:"platform_name" orm:"column(platform_name)"`   // 所属平台名称
	Name         string `json:"name" orm:"column(name)"`                     // 游戏名称
	PtName       string `json:"pt_name" orm:"column(pt_name)"`               // 葡语名称
	Cover        string `json:"cover" orm:"column(cover)"`                   // 游戏封面
	GameTypeId   int    `json:"game_type_id" orm:"column(game_type_id)"`     // 游戏分类
	GameTypeName string `json:"game_type_name" orm:"column(game_type_name)"` // 游戏分类名称
	Sort         int    `json:"sort" orm:"column(sort)"`                     // 排序权重
	Code         string `json:"code" orm:"column(code)"`                     // 游戏code
	SupplierId   int    `json:"-" orm:"column(supplier_id)"`                 // 供应商id
	Status       int    `json:"status" orm:"column(status)"`                 // 状态 1-正常 0-关闭
	Maintain     int    `json:"maintain" orm:"column(maintain)"`             // 维护 1-维护中 0-没有维护
	Recommend    int    `json:"recommend" orm:"column(recommend)"`           // 头部推荐 1是 0-否
	Swpg         int    `json:"swpg" orm:"column(swpg)"`                     // 0-禁用 1-启用
	Ttpg         int    `json:"ttpg" orm:"column(ttpg)"`                     // 0-禁用 1-启用
	Mlpg         int    `json:"mlpg" orm:"column(mlpg)"`                     // 0-禁用 1-启用
	Tag          string `json:"tag" orm:"column(tag)"`                       // 游戏标签
	CodeRule     int    `json:"code_rule" orm:"column(code_rule)"`           // 打码规则 1-流水 2-净赢 3-赢金打码
	IsDeleted    int    `json:"-" orm:"column(is_deleted)"`                  // 是否删除
	PackageId    int    `json:"package_id" orm:"column(package_id)"`         // 分包id
}

func CreateGameModel() *GameModel {
	return &GameModel{CreateBase()}
}
