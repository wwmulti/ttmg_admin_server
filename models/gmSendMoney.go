package models

type GmSendMoneyModel struct {
	*Base
}

type GmSendMoney struct {
	Id           int     `json:"id" orm:"auto;column(id)"`                    // 主键ID
	UserId       int     `json:"user_id" orm:"column(user_id)"`               // 用户ID
	RoleId       int     `json:"role_id" orm:"column(role_id)"`               // 玩家角色ID
	IfInnerProxy int     `json:"if_inner_proxy" orm:"column(if_inner_proxy)"` // 是否模拟账户: 0-否, 1-是
	PackageId    int     `json:"package_id" orm:"column(package_id)"`         // 游戏包ID
	TeamId       int     `json:"team_id" orm:"column(team_id)"`               // 团队ID
	SaleMenId    int     `json:"sale_men_id" orm:"column(sale_men_id)"`       // 业务员ID
	Money        float64 `json:"money" orm:"column(money)"`                   // 变动金币数量
	Status       int     `json:"status" orm:"column(status)"`                 // 审核状态: 0-未审核, 1-已审核, 2-拒绝
	OperateType  int     `json:"operate_type" orm:"column(operate_type)"`     // 操作类型
	Type         int     `json:"type" orm:"column(type)"`                     // 来源类型: 0-运营后台, 1-渠道后台
	WageMul      float64 `json:"wage_mul" orm:"column(wage_mul)"`             // 打码倍数
	CheckUserId  int     `json:"check_user_id" orm:"column(check_user_id)"`   // 操作人ID
	CheckUser    string  `json:"check_user" orm:"column(check_user)"`         // 操作人名称
	Note         string  `json:"note" orm:"column(note)"`                     // 备注
	CTime        int64   `json:"c_time" orm:"column(c_time)"`                 // 插入时间 (Unix时间戳)
	UTime        int64   `json:"u_time" orm:"column(u_time)"`                 // 更新时间 (Unix时间戳)
}

func CreateGmSendMoneyModel() *GmSendMoneyModel {
	return &GmSendMoneyModel{CreateBase()}
}
