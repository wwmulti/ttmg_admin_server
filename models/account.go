package models

type AccountModel struct {
	*Base
}

type Account struct {
	Id            int    `json:"id" orm:"auto;column(id)"`                      // 主键，自增
	AccountName   string `json:"username" orm:"column(account_name)"`           // 账号名称
	Password      string `json:"-" orm:"column(password)"`                      // 密码
	AccountType   int    `json:"-" orm:"column(account_type)"`                  // 账号类型 1管理员
	Salt          string `json:"-" orm:"column(salt)"`                          // 盐值
	LastLoginIp   string `json:"last_login_ip" orm:"column(last_login_ip)"`     // 上次登陆ip
	LastLoginTime int    `json:"last_login_time" orm:"column(last_login_time)"` // 上次登陆时间
	RegisterTime  int    `json:"register_time" orm:"column(register_time)"`     // 注册时间
	Status        int    `json:"status" orm:"column(status)"`                   // 状态 0禁止 1正常
	Secret        string `json:"-" orm:"column(secret)"`                        // 秘钥
	BindTime      int    `json:"bind_time" orm:"column(bind_time)"`             // 秘钥绑定时间
	IsDeleted     int    `json:"is_deleted" orm:"column(is_deleted)"`           // 是否删除
	Parents       string `json:"-" orm:"column(parents)"`                       // 父级关系
	Pid           int    `json:"pid" orm:"column(pid)"`                         // 父级id
	Token         string `json:"token" orm:"-"`                                 // 请求token

	Group *AuthGroup `json:"group" orm:"rel(fk);column(group_id)"` // 外键关联
}

func CreateAccountModel() *AccountModel {
	return &AccountModel{CreateBase()}
}
