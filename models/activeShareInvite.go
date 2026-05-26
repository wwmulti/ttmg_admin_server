package models

type ActiveShareInviteModel struct {
	*Base
}

type ActiveShareInvite struct {
	Id                int   `json:"id" orm:"auto;column(id)"`                              // ID
	ActiveId          int   `json:"active_id" orm:"column(active_id)"`                     // 活动期id
	Pid1              int   `json:"pid1" orm:"column(pid1)"`                               // 直属邀请父id
	Pid2              int   `json:"pid2" orm:"column(pid2)"`                               // 二级父id
	Pid3              int   `json:"pid3" orm:"column(pid3)"`                               // 三级父id
	UserId            int   `json:"user_id" orm:"column(user_id)"`                         // 被邀请人用户id
	InviteTime        int64 `json:"invite_time" orm:"column(invite_time)"`                 // 邀请时间
	InviteeTotalPay   int64 `json:"invitee_total_pay" orm:"column(invitee_total_pay)"`     // 被邀请人活动期内累计充值
	InviteeTotalWater int64 `json:"invitee_total_water" orm:"column(invitee_total_water)"` // 被邀请人活动期内累计流水
	IsValid           int   `json:"is_valid" orm:"column(is_valid)"`                       // 是否有效用户 0否 1是
	ValidTime         int64 `json:"valid_time" orm:"column(valid_time)"`                   // 达标时间
	CTime             int64 `json:"c_time" orm:"column(c_time)"`                           // 时间
}

func CreateActiveShareInviteModel() *ActiveShareInviteModel {
	return &ActiveShareInviteModel{CreateBase()}
}

type ShareActiveDTO struct {
	Id            int     `json:"id"`           //活动ID
	RuleID        int     `json:"rule_id"`      //使用的规则ID
	Name          string  `json:"name"`         //活动名称
	Icon          string  `json:"icon"`         //活动图标
	ShareUrl      string  `json:"share_url"`    //分享链接
	ValidNum      int     `json:"valid_num"`    //有效数量
	ConTotalPays  float64 `json:"condition1"`   //条件1 累计充值金额
	ConTotalWater float64 `json:"condition2"`   //条件2 累计流水值
	ConRelation   int     `json:"con_relation"` //条件关系 0-同时满足，1-满足任一
	Box           []ShareActiveBoxDTO
}

type ShareActiveBoxDTO struct {
	Id      int     `json:"id"`      //宝箱ID
	Men     int     `json:"men"`     //有效人数
	Rewards float64 `json:"rewards"` //奖励金额
	Icon    string  `json:"icon"`    //宝箱图标
	Status  int     `json:"status"`  //领取状态 -1-不可领取 0-可领取 1-已领取
}

type ActiveShareUsers struct {
	Id                int     `json:"id"`                  //用户ID
	Name              string  `json:"name"`                //用户名称
	InviteeTotalPay   float64 `json:"invitee_total_pay"`   //邀请人累计充值金额
	InviteeTotalWater float64 `json:"invitee_total_water"` //邀请人累计流水值
	IsValid           int     `json:"is_valid"`            //是否有效 0-无效 1-有效
	Relation          int     `json:"relation"`            //邀请关系 0-直属 1-非直属
	InviteTime        int64   `json:"invite_time"`         //邀请时间
}
