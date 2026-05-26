package models

// --- 提现订单表 ---

type WithdrawOrderModel struct {
	*Base
}

type WithdrawOrder struct {
	Id                    int     `json:"id" orm:"auto;column(id)"`                                         // 自增主键
	OrderNo               string  `json:"order_no" orm:"column(order_no);"`                                 // 业务订单号
	UserId                int     `json:"user_id" orm:"column(user_id)"`                                    // 用户ID
	RoleId                int64   `json:"role_id" orm:"column(role_id);null"`                               // 用户RoleId
	IMoney                float64 `json:"i_money" orm:"column(i_money)"`                                    // 金币(原始金额)
	Tax                   float64 `json:"tax" orm:"column(tax)"`                                            // 税
	PaymentId             int     `json:"payment_id" orm:"column(payment_id)"`                              // 渠道id
	PaymentChannelId      int     `json:"payment_channel_id" orm:"column(payment_channel_id)"`              // 通道id
	PaymentTypeId         int     `json:"payment_type_id" orm:"column(payment_type_id)"`                    // 提现类型id
	PaymentTypeName       string  `json:"payment_type_name" orm:"column(payment_type_name)"`                // 提现类型
	RealName              string  `json:"real_name" orm:"column(real_name)"`                                // 真实姓名
	CardNo                string  `json:"card_no" orm:"column(card_no)"`                                    // 卡号
	BankName              string  `json:"bank_name" orm:"column(bank_name)"`                                // 银行名字
	Status                int     `json:"status" orm:"column(status)"`                                      // 提现订单状态
	IsFake                int     `json:"is_fake" orm:"column(is_fake)"`                                    // 是否假提现
	CheckTime             int64   `json:"check_time" orm:"column(check_time);"`                             // 审查时间
	CheckUser             string  `json:"check_user" orm:"column(check_user);"`                             // 审查人
	Descript              string  `json:"descript" orm:"column(descript);"`                                 // 描述
	Phone                 string  `json:"phone" orm:"column(phone)"`                                        // 手机号
	Email                 string  `json:"email" orm:"column(email)"`                                        // 邮箱
	TransactionNo         string  `json:"transaction_no" orm:"column(transaction_no);"`                     // 第三方订单号
	ChannelId             int     `json:"channel_id" orm:"column(channel_id);"`                             // 通道ID
	AddTime               int64   `json:"add_time" orm:"column(add_time)"`                                  // 记录时间/下单时间
	UpdateTime            int64   `json:"update_time" orm:"column(update_time);"`                           // 更新时间
	RecordType            int     `json:"record_type" orm:"column(record_type)"`                            // 提现类型
	OperatorId            int     `json:"operator_id" orm:"column(operator_id)"`                            // 操作人ID
	SourceType            int     `json:"source_type" orm:"column(source_type)"`                            // 来源类型
	AutoWithdrawFlag      int     `json:"auto_withdraw_flag" orm:"column(auto_withdraw_flag)"`              // 是否自动出款标记
	AutoWithdrawCheckMsg  string  `json:"auto_withdraw_check_msg" orm:"column(auto_withdraw_check_msg)"`    // 自动出款失败原因
	AutoWithdrawCheckTime int64   `json:"auto_withdraw_check_time" orm:"column(auto_withdraw_check_time);"` // 检查时间
	IdCard                string  `json:"id_card" orm:"column(id_card)"`                                    // CPF/身份证号
	WarningFlag           int     `json:"warning_flag" orm:"column(warning_flag)"`                          // 风险预警标记
	PackageId             int     `json:"package_id" orm:"column(package_id)"`                              // 包ID
}

func CreateWithdrawOrderModel() *WithdrawOrderModel {
	return &WithdrawOrderModel{CreateBase()}
}
