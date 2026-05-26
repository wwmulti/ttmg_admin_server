package models

type ActivesModel struct {
	*Base
}

type Actives struct {
	Id              int    `json:"id" orm:"auto;column(id)"`                          // 主键，自增
	Name            string `json:"name" orm:"column(name)"`                           // 活动名称
	PtName          string `json:"pt_name" orm:"column(pt_name)"`                     // 葡语名称
	EnName          string `json:"en_name" orm:"column(en_name)"`                     // 英语名称
	Icon            string `json:"icon" orm:"column(icon)"`                           // 活动图标
	PcIcon          string `json:"pc_icon" orm:"column(pc_icon)"`                     // pc端图片
	RedirectLink    int    `json:"redirect_link" orm:"column(redirect_link)"`         // 跳转地址id
	PtDesc          string `json:"pt_desc" orm:"column(pt_desc)"`                     // 葡语描述
	EnDesc          string `json:"en_desc" orm:"column(en_desc)"`                     // 英语描述
	PackageId       int    `json:"package_id" orm:"column(package_id)"`               // 开放平台id
	StartTime       int64  `json:"start_time" orm:"column(start_time)"`               // 活动开始时间
	EndTime         int64  `json:"end_time" orm:"column(end_time)"`                   // 活动结束时间
	CollecTime      int64  `json:"collec_time" orm:"column(collec_time)"`             // 领取结束时间
	Ctime           int64  `json:"c_time" orm:"column(c_time)"`                       // 创建时间
	Sort            int    `json:"sort" orm:"column(sort)"`                           // 排序
	ActiveTypeId    int    `json:"active_type_id" orm:"column(active_type_id)"`       // 活动类型id
	TimerCheckAward int    `json:"timer_check_award" orm:"column(timer_check_award)"` // 后台定时任务发奖 0-否，1是
	TimerCheckTime  int64  `json:"timer_check_time" orm:"column(timer_check_time)"`   // 定时任务触发时间
	ValidBets       string `json:"valid_bets" orm:"column(valid_bets)"`               // 有效投注条件（存json key有game_types游戏类型,type投注获胜）,
	IsNewcomer      int    `json:"is_newcomer" orm:"column(is_newcomer)"`             // 是否新人 1-是 0-否
	IsPay           int    `json:"is_pay" orm:"column(is_pay)"`                       // 是否充值 1-是 0-否
	IsHideComplete  int    `json:"is_hide_complete" orm:"column(is_hide_complete)"`   // 完成后是否隐藏 1-是 0-否
	Status          int    `json:"status" orm:"column(status)"`                       // 状态0-关闭 1-开启
	IsDeleted       int    `json:"is_deleted" orm:"column(is_deleted)"`               // 删除
}

type PlayLuckyWheelParams struct {
	Uid       int64
	Username  string
	ActiveId  int64
	WheelType int64
}

func CreateActivesModel() *ActivesModel {
	return &ActivesModel{CreateBase()}
}
