package config

// AccountType 账号类型
type AccountType int

const (
	ActiveTypeManager AccountType = iota + 1 // 管理员
	ActiveTypeGame                           // 游戏运营

)
