package models

type UserInterestModel struct {
	*Base
}

type UserInterest struct {
	Id                 int   `json:"id" orm:"auto;column(id)"`                                       // 主键，自增
	RuleId             int   `json:"rule_id" orm:"column(rule_id)"`                                  // 规则id
	ActivityId         int   `json:"activity_id" orm:"column(activity_id)"`                          // 活动id
	UserId             int   `json:"user_id" orm:"column(user_id)"`                                  // 用户id
	Amount             int64 `json:"amount" orm:"column(amount)"`                                    // 存款金额
	InterestAmount     int64 `json:"interest_amount" orm:"column(interest_amount)"`                  // 利息金额
	DepositAt          int64 `json:"deposit_at" orm:"column(deposit_at)"`                            // 存入时间
	AvailableTakeOutAt int64 `json:"available_take_out_at" orm:"column(available_take_out_at);null"` // 可以提取利息的时间
	LastCalculateTime  int64 `json:"last_calculate_time" orm:"column(last_calculate_time);null"`     // 上一次计算利息的时间
}

type UserInterestBalance struct {
	Balance  float64
	Interest float64
}

func CreateUserInterestModel() *UserInterestModel {
	return &UserInterestModel{CreateBase()}
}
