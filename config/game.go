package config

// 供应商类型
type SupplierType int

const (
	ZySupplierType SupplierType = iota // 自研
	OfficialSupplierType
)

var PgSupplierTypeMap = map[SupplierType]string{
	ZySupplierType:       "自研",
	OfficialSupplierType: "官方",
}

var PpSupplierTypeMap = map[SupplierType]string{
	OfficialSupplierType: "官方",
}

// 游戏类型
type GameType int

const (
	PgGameType GameType = iota // pg
	PpGameType                 // pp
)

// GamePlatforms 游戏平台
var GamePlatforms = []string{
	"pg", "jdb", "kess", "zy", "wg", "cp",
}
