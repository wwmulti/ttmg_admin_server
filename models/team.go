package models

type TeamModel struct {
	*Base
}

type Team struct {
	Id         int    `orm:"auto;column(id);pk" json:"id"`                        // 主键，自增
	PackageId  int    `orm:"column(package_id)" json:"package_id"`                // 分包id
	RoleId     int    `orm:"column(role_id)" json:"role_id"`                      // role id
	Pid        int    `orm:"column(pid);default(0)" json:"pid"`                   // 父级id
	Parents    string `orm:"column(parents);type(text)" json:"-"`                 // 父级关系,1,
	IsOfficial int    `orm:"column(is_official);default(0)" json:"is_official"`   // 是否官网 0-否 1-是
	IsBloger   int    `orm:"column(is_bloger);default(0)" json:"is_bloger"`       // 是否博主 0-否 1-是
	IsDropUser int    `orm:"column(is_drop_user);default(0)" json:"is_drop_user"` // 是否掉绑 0-否 1-是
	Title      string `orm:"column(title);size(100)" json:"title"`                // 团队名称
	Username   string `orm:"column(username);size(30)" json:"username"`           // 用户名
	UserId     int    `orm:"column(user_id)" json:"user_id"`                      // user id
	Token      string `json:"token" orm:"column(token);size(256);not null"`       // token
	Secret     string `json:"secret" orm:"column(secret);size(256);not null"`     // 秘钥
	Rtp        int    `orm:"column(rtp)" json:"rtp"`                              // 概率
	IsDeleted  int    `orm:"column(is_deleted);default(0)" json:"-"`              // 是否删除 0-否 1-是
}

func CreateTeamModel() *TeamModel {
	return &TeamModel{CreateBase()}
}
