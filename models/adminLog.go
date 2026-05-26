package models

type AdminLogModel struct {
	*Base
}

type AdminLog struct {
	Id         int    `json:"id" orm:"auto;column(id)"`                 // 主键，自增
	AdminId    int    `json:"admin_id" orm:"column(admin_id)"`          // 操作人id
	PlayerId   int    `json:"player_id" orm:"column(player_id),null"`   // 玩家id
	RoleId     int64  `json:"role_id" orm:"column(role_id),null"`       // 玩家RoleID
	Path       string `json:"path" orm:"column(path);null"`             // 请求路径
	Controller string `json:"controller" orm:"column(controller);null"` // 控制器
	Action     string `json:"action" orm:"column(action);null"`         // 方法
	Body       string `json:"body" orm:"column(body);null"`             // 请求内容 (通常为JSON)
	Result     string `json:"result" orm:"column(result);null"`         // 操作结果 (通常为JSON)
	Status     int    `json:"status" orm:"column(status)"`              // 状态 0-失败 1-成功
	LogType    int    `json:"log_type" orm:"column(log_type)"`          // 日志类型
	Ip         string `json:"ip" orm:"column(ip)"`                      // 操作ip
	CTime      int64  `json:"c_time" orm:"column(c_time)"`              // 操作时间 (Unix时间戳)
	PackageId  int    `json:"package_id" orm:"column(package_id)"`      // 分包id
}

func CreateAdminLogModel() *AdminLogModel {
	return &AdminLogModel{CreateBase()}
}
