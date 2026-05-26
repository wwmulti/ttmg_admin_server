package models

type UserModel struct {
	*Base
}

type User struct {
	Id                      int64  `json:"id" orm:"auto;column(id)"`                                                         // 主键，自增
	Username                string `json:"username" orm:"column(username);size(14)"`                                         // 用户名
	Password                string `json:"-" orm:"column(password);size(512)"`                                               // 密码
	Cpf                     string `json:"cpf" orm:"column(cpf);size(30);null"`                                              // cpf
	Avatar                  string `json:"avatar" orm:"column(avatar);size(255);null"`                                       // 头像
	Status                  int    `json:"status" orm:"column(status);default(0)"`                                           // 状态 0封禁 1正常
	Salt                    string `json:"-" orm:"column(salt);size(10);null"`                                               // salt
	Whatsapp                string `json:"whatsapp" orm:"column(whatsapp);size(50);null"`                                    // WhatsApp账号
	Facebook                string `json:"facebook" orm:"column(facebook);size(50);null"`                                    // Facebook账号
	Telegram                string `json:"telegram" orm:"column(telegram);size(50);null"`                                    // Telegram账号
	Insgram                 string `json:"insgram" orm:"column(insgram);size(50);null"`                                      // insgram账号
	Birthday                *int64 `json:"birthday" orm:"column(birthday);null"`                                             // 用户生日
	WithdrawPassword        string `json:"-" orm:"column(withdraw_password);size(255);null"`                                 // 提现密码
	SafeQuestionId          *int   `json:"safe_question_id" orm:"column(safe_question_id);null"`                             // 安全问题id
	SafeQuestionAnswer      string `json:"safe_question_answer" orm:"column(safe_question_answer);size(255);null"`           // 安全问题答案
	RoleId                  int64  `json:"role_id" orm:"column(role_id);null"`                                               // 角色id
	LoginAt                 int64  `json:"login_at" orm:"column(login_at);null"`                                             // 登录时间
	RegisterAt              int64  `json:"register_at" orm:"column(register_at);null"`                                       // 注册时间
	Balance                 int64  `json:"balance" orm:"column(balance);default(0)"`                                         // 可用余额
	LockBalance             int64  `json:"lock_balance" orm:"column(lock_balance);default(0)"`                               // 锁定金额
	IP                      string `json:"ip" orm:"column(ip);size(46);null"`                                                // 当前用户的ip地址
	RegisterIp              string `json:"register_ip" orm:"column(register_ip);size(46);null"`                              // 当前用户的ip地址
	Remark                  string `json:"remark" orm:"column(remark)"`                                                      // 备注
	UserType                int    `json:"user_type" orm:"column(user_type);default(0)"`                                     // 用户类型 1博主 2经纪人
	IsBanGame               int    `json:"is_ban_game" orm:"column(is_ban_game);default(0)"`                                 // 是否封禁游戏
	IsBanWithdraw           int    `json:"is_ban_withdraw" orm:"column(is_ban_withdraw);default(0)"`                         // 是否封禁提现
	IsBanInviteReward       int    `json:"is_ban_invite_reward" orm:"column(is_ban_invite_reward);default(0)"`               // 是否封禁邀请奖励
	IsBanChildBetCommission int    `json:"is_ban_child_bet_commission" orm:"column(is_ban_child_bet_commission);default(0)"` // 是否封禁下级投注返佣
	IsOnlyAllowOfficialGame int    `json:"is_only_allow_official_game" orm:"column(is_only_allow_official_game);default(0)"` // 是否只允许官方游戏
	IsMock                  int    `json:"is_mock" orm:"column(is_mock);default(0)"`                                         // 是否模拟
	IsDeleted               int    `json:"is_deleted" orm:"column(is_deleted);default(0)"`                                   // 是否删除
	PackageId               int    `json:"package_id" orm:"column(package_id);default(1)"`                                   // 分包id
}

type UserLoginDTO struct {
	Id       int64  `json:"id" orm:"auto;column(id)"`               // 主键，自增
	Username string `json:"username" orm:"column(username);unique"` // 用户名
	Avatar   string `json:"avatar" orm:"column(avatar)"`            // 头像
	RoleId   int    `json:"role_id" orm:"column(role_id)"`          // 角色Id
	Token    string `json:"token" orm:"column(token)"`              // 登录token
	LoginAt  int64  `json:"login_at" orm:"column(login_at)"`        // 登录时间
	Ip       string `json:"ip" orm:"column(ip)"`                    // 登录IP
}

type GetUserInfoDTO struct {
	Id       int64  `json:"role_id" orm:"auto;column(role_id)"`     // 主键，自增
	Username string `json:"username" orm:"column(username);unique"` // 用户名
	Avatar   string `json:"avatar" orm:"column(avatar)"`            // 头像
	Whatsapp string `json:"whatsapp" orm:"column(whatsapp)"`        // WhatsApp账号
	Facebook string `json:"facebook" orm:"column(facebook)"`        // Facebook账号
	Telegram string `json:"telegram" orm:"column(telegram)"`        // Telegram账号
	Insgram  string `json:"insgram" orm:"column(insgram)"`          // Insgram账号
	Birthday *int64 `json:"birthday" orm:"column(birthday)"`        // 用户生日
	Balance  int64  `json:"balance" orm:"column(balance)"`          // 用户余额
}

func CreateUserModel() *UserModel {
	return &UserModel{CreateBase()}
}
