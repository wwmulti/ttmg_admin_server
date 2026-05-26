package services

import (
	"api/config"
	"api/models"
	"fmt"
	"github.com/beego/i18n"
	"time"
)

type ActiveVipService struct{}

type AllActiveVipRuleRequest struct {
	ActiveId int `json:"active_id"`
}

// AllActiveVipRule 所有VIP规则
func (s *ActiveVipService) AllActiveVipRule(request AllActiveVipRuleRequest) ([]models.ActiveVipRule, error) {
	ActiveVipRuleModel := models.CreateActiveVipRuleModel()
	var list []models.ActiveVipRule
	_, err := ActiveVipRuleModel.QueryTable(&models.ActiveVipRule{}).
		Filter("active_id", request.ActiveId).
		All(&list)
	if nil != err {
		return nil, err
	}

	return list, nil
}

// AllActiveVipRuleConditionType 所有条件类型
func (s *ActiveVipService) AllActiveVipRuleConditionType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.VipRuleConditionOr, "name": i18n.Tr(language, "ManZuRenYiTiaoJian")},
		{"id": config.VipRuleConditionAnd, "name": i18n.Tr(language, "ManZuSuoYouTiaoJian")},
		{"id": config.VipRuleConditionRecharge, "name": i18n.Tr(language, "ManZuChongZhiTiaoJian")},
		{"id": config.VipRuleConditionBet, "name": i18n.Tr(language, "ManZuLiuShuiTiaoJian")},
	}
}

type AddActiveVipRuleRequest struct {
	ActiveId            int     `json:"active_id"`             // 活动id
	Lv                  int     `json:"lv"`                    // 等级
	TotalPays           float64 `json:"total_pays"`            // 累计充值
	TotalBets           float64 `json:"total_bets"`            // 累计有效押注
	ConAnd              int     `json:"con_and"`               // 条件：0-或 1-且 2-只需充值 3-只需流水
	Rewards             float64 `json:"rewards"`               // 晋级金
	WithdrawNumLimit    int     `json:"withdraw_num_limit"`    // 每日提现次数限制
	WithdrawAmountLimit float64 `json:"withdraw_amount_limit"` // 每日提现金额限制
	WithdrawFreeNum     int     `json:"withdraw_free_num"`     // 每日免费交易次数
	WithdrawFee         int     `json:"withdraw_fee"`          // 提现手续，万分比
	Status              int     `json:"status"`                // 状态：0-关闭 1-启用
}

// AddActiveVipRule 添加VIP规则
func (s *ActiveVipService) AddActiveVipRule(request AddActiveVipRuleRequest) error {
	if request.ActiveId <= 0 {
		return fmt.Errorf("HuoDongBuCunZai")
	}
	if request.Lv < 0 {
		return fmt.Errorf("VipDengJiBuHeFa")
	}

	_, err := models.CreateActiveVipRuleModel().Insert(&models.ActiveVipRule{
		ActiveId:            request.ActiveId,
		Lv:                  request.Lv,
		TotalPays:           request.TotalPays,
		TotalBets:           request.TotalBets,
		ConAnd:              request.ConAnd,
		Rewards:             request.Rewards,
		WithdrawNumLimit:    request.WithdrawNumLimit,
		WithdrawAmountLimit: request.WithdrawAmountLimit,
		WithdrawFreeNum:     request.WithdrawFreeNum,
		WithdrawFee:         request.WithdrawFee,
		Status:              request.Status,
		Ctime:               time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// EditActiveVipRule 编辑VIP规则
func (s *ActiveVipService) EditActiveVipRule(id int, request AddActiveVipRuleRequest) error {
	fields := []string{
		"active_id",
		"lv",
		"total_pays",
		"total_bets",
		"con_and",
		"rewards",
		"withdraw_num_limit",
		"withdraw_amount_limit",
		"withdraw_free_num",
		"withdraw_fee",
		"status",
	}
	affected, err := models.CreateActiveVipRuleModel().Update(&models.ActiveVipRule{
		Id:                  id,
		ActiveId:            request.ActiveId,
		Lv:                  request.Lv,
		TotalPays:           request.TotalPays,
		TotalBets:           request.TotalBets,
		ConAnd:              request.ConAnd,
		Rewards:             request.Rewards,
		WithdrawNumLimit:    request.WithdrawNumLimit,
		WithdrawAmountLimit: request.WithdrawAmountLimit,
		WithdrawFreeNum:     request.WithdrawFreeNum,
		WithdrawFee:         request.WithdrawFee,
		Status:              request.Status,
	}, fields...)
	if err != nil || affected == 0 {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// DeleteActiveVipRule 删除VIP规则
func (s *ActiveVipService) DeleteActiveVipRule(id int) error {
	affected, err := models.CreateActiveVipRuleModel().QueryTable(&models.ActiveVipRule{}).Filter("id", id).Delete()
	if err != nil || affected == 0 {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

type ChangeActiveVipRuleStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangeActiveVipRuleStatus 修改Vip等级规则状态
func (s *ActiveVipService) ChangeActiveVipRuleStatus(request ChangeActiveVipRuleStatusRequest) error {
	ruleModel := models.CreateActiveVipRuleModel()
	if request.Id == 0 {
		return fmt.Errorf("VipDengJiGuiZeIDBiTian")
	}
	_, err := ruleModel.Update(&models.ActiveVipRule{
		Id:     request.Id,
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}
	return nil
}

type AllActiveVipTaskRequest struct {
	ActiveId int `json:"active_id"`
}

// AllActiveVipTask VIP福利列表
func (s *ActiveVipService) AllActiveVipTask(request AllActiveVipTaskRequest) ([]models.ActiveVipWelfare, error) {
	ActiveVipWelfareModel := models.CreateActiveVipWelfareModel()
	var list []models.ActiveVipWelfare
	_, err := ActiveVipWelfareModel.QueryTable(&models.ActiveVipWelfare{}).
		Filter("active_id", request.ActiveId).
		All(&list)
	if nil != err {
		return nil, err
	}

	return list, nil
}

// AllActiveVipTaskConditionType 所有条件类型
func (s *ActiveVipService) AllActiveVipTaskConditionType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.VipTaskConditionOr, "name": i18n.Tr(language, "ManZuRenYiTiaoJian")},
		{"id": config.VipTaskConditionAnd, "name": i18n.Tr(language, "ManZuSuoYouTiaoJian")},
		{"id": config.VipTaskConditionRecharge, "name": i18n.Tr(language, "ManZuChongZhiTiaoJian")},
		{"id": config.VipTaskConditionBet, "name": i18n.Tr(language, "ManZuLiuShuiTiaoJian")},
	}
}

// AllActiveVipTaskCycleType 所有奖励类型
func (s *ActiveVipService) AllActiveVipTaskCycleType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.VipTaskCycleTypeDay, "name": i18n.Tr(language, "Ri")},
		{"id": config.VipTaskCycleTypeWeek, "name": i18n.Tr(language, "Zhou")},
		{"id": config.VipTaskCycleTypeMonth, "name": i18n.Tr(language, "Yue")},
	}
}

type AddActiveVipTaskRequest struct {
	ActiveId  int     `json:"active_id"`  // 活动id
	Cycle     int     `json:"cycle"`      // 奖励类型：1-日 2-周 3-月
	Lv        int     `json:"lv"`         // 匹配的等级
	TotalPays float64 `json:"total_pays"` // 累计充值
	TotalBets float64 `json:"total_bets"` // 累计有效押注
	ConAnd    int     `json:"con_and"`    // 条件：0-或 1-且 2-只需充值 3-只需流水
	Rewards   float64 `json:"rewards"`    // 奖励金额
	Status    int     `json:"status"`     // 状态：0-关闭 1-启用
}

// AddActiveVipTask 添加VIP福利
func (s *ActiveVipService) AddActiveVipTask(request AddActiveVipTaskRequest) error {
	if request.Cycle < 1 || request.Cycle > 3 {
		return fmt.Errorf("JiangLiLeiXingCuoWu")
	}

	_, err := models.CreateActiveVipWelfareModel().Insert(&models.ActiveVipWelfare{
		ActiveId:  request.ActiveId,
		Cycle:     request.Cycle,
		Lv:        request.Lv,
		TotalPays: request.TotalPays,
		TotalBets: request.TotalBets,
		ConAnd:    request.ConAnd,
		Rewards:   request.Rewards,
		Status:    request.Status,
		Ctime:     time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// EditActiveVipTask 编辑VIP福利
func (s *ActiveVipService) EditActiveVipTask(id int, request AddActiveVipTaskRequest) error {
	fields := []string{"active_id", "cycle", "lv", "total_pays", "total_bets", "con_and", "rewards", "status"}
	affected, err := models.CreateActiveVipWelfareModel().Update(&models.ActiveVipWelfare{
		Id:        id,
		ActiveId:  request.ActiveId,
		Cycle:     request.Cycle,
		Lv:        request.Lv,
		TotalPays: request.TotalPays,
		TotalBets: request.TotalBets,
		ConAnd:    request.ConAnd,
		Rewards:   request.Rewards,
		Status:    request.Status,
	}, fields...)
	if err != nil || affected == 0 {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// DeleteActiveVipTask 删除VIP福利
func (s *ActiveVipService) DeleteActiveVipTask(id int) error {
	affected, err := models.CreateActiveVipWelfareModel().QueryTable(&models.ActiveVipWelfare{}).Filter("id", id).Delete()
	if err != nil || affected == 0 {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

type ChangeActiveVipTaskStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangeActiveVipTaskStatus 修改Vip等级任务状态
func (s *ActiveVipService) ChangeActiveVipTaskStatus(request ChangeActiveVipTaskStatusRequest) error {
	taskModel := models.CreateActiveVipWelfareModel()
	if request.Id == 0 {
		return fmt.Errorf("VipDengJiRenWuIDBiTian")
	}
	_, err := taskModel.Update(&models.ActiveVipWelfare{
		Id:     request.Id,
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}
	return nil
}

// GetValidVipLvActiveId 获取有效的VIP等级活动id
func (s *ActiveVipService) GetValidVipLvActiveId(userId int) (int, error) {
	userModel := models.CreateUserModel()
	user := models.User{}
	err := userModel.QueryTable(&models.User{}).
		Filter("id", userId).
		One(&user)
	if err != nil {
		return 0, err
	}
	pkgId := user.PackageId
	nowTime := time.Now().Unix()
	var active models.Actives
	err = models.CreateActivesModel().QueryTable(&models.Actives{}).
		Filter("package_id", pkgId).
		Filter("status", 1).
		Filter("is_deleted", 0).
		Filter("active_type_id", config.ActiveTypeVipLv).
		Filter("collec_time__gte", nowTime).
		Filter("start_time__lt", nowTime).
		OrderBy("-id").
		One(&active, "id")
	if err != nil {
		return 0, err
	}
	return active.Id, nil
}
