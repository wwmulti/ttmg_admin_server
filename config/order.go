package config

type RechargeOrderStatus int

const (
	RechargeOrderStatusUnpaid   RechargeOrderStatus = iota + 1 // 未支付
	RechargeOrderStatusPaid                                    // 已支付
	RechargeOrderStatusNotified                                // 已回调
)

type WithdrawOrderStatus int

const (
	WithdrawOrderStatusCreated        WithdrawOrderStatus = iota // 提交
	WithdrawOrderStatusPass                                      // 审核通过
	WithdrawOrderStatusRejectedBack                              // 拒绝并退回
	WithdrawOrderStatusRejectedFrozen                            // 拒绝并冻结
	WithdrawOrderStatusProcessing                                // 处理中
	WithdrawOrderStatusFailBack                                  // 处理失败退回
	WithdrawOrderStatusAddWorkReview                             // 加入工单审核
	WithdrawOrderStatusWarning                                   // 风险挂起
)
const WithdrawOrderStatusSuccess WithdrawOrderStatus = 100 // 成功
