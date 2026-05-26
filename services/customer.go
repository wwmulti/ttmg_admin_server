package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
)

type CustomerService struct {
	BaseService
}

type GmTotalMoney struct {
	TotalMoney  float64 `json:"total_money"`
	OperateType int     `json:"operate_type"`
}

type GmSendMoneyListRequest struct {
	RoleId        int        `form:"role_id"`
	Money         float64    `form:"money"`
	OperateType   int        `form:"operate_type"`
	IfInnerProxy  int        `form:"if_inner_proxy"` // 是否模拟账户
	Status        int        `form:"status"`
	Type          int        `form:"type"`
	TeamId        string     `form:"channel_remark"`         // 所属渠道或业务员备注
	ChannelRoleId int        `form:"channel_role_id" op:"-"` // 所属渠道或业务员角色ID
	CheckUser     string     `form:"check_user" op:"like"`
	CTime         []string   `form:"c_time[]" op:"time_range" orm:"c_time"`
	Event         string     `form:"event" op:"-"`
	Page          int        `form:"page"`
	PageSize      int        `form:"page_size"`
	NeedReload    int        `form:"need_reload" op:"-"`
	RawQuery      url.Values `form:"-"`
	PackageIds    []int      // 游戏包ID列表
}

// GmSendMoneyList 上下分列表
func (s *CustomerService) GmSendMoneyList(req GmSendMoneyListRequest, lang string) (map[string]interface{}, error) {
	gmSendMoneyModel := models.CreateGmSendMoneyModel()
	condition, sort := gmSendMoneyModel.BuildCondition(req, "-id")
	condition = s.LimitPackageId(condition, req.PackageIds)

	if req.Event == "asyncexport" {
		exportConfig := ExportConfig{}
		taskId, _ := (&TaskService{}).AddTask(TaskTypeExport, map[string]interface{}{"config": exportConfig})
		return map[string]interface{}{"taskId": taskId}, nil
	}

	data, total, err := gmSendMoneyModel.GetPageList(&models.GmSendMoney{}, condition, req.Page, req.PageSize, sort)
	if err != nil {
		return nil, err
	}

	list := data.([]models.GmSendMoney)
	newList := make([]interface{}, 0, len(list))
	for _, item := range list {
		temp := map[string]interface{}{
			"id":               item.Id,
			"role_id":          item.RoleId,
			"package_id":       item.PackageId,
			"package_name":     "",
			"channel_role_id":  item.TeamId,
			"channel_remark":   "",
			"salesman_role_id": item.SaleMenId,
			"salesman_remark":  "",
			"operate_type":     item.OperateType,
			"if_inner_proxy":   item.IfInnerProxy,
			"money":            item.Money,
			"wage_mul":         item.WageMul,
			"note":             item.Note,
			"status":           item.Status,
			"insert_time":      item.CTime,
			"update_time":      item.UTime,
			"type":             item.Type,
			// "channel_note":     "", //操作渠道
			"check_user": item.CheckUser,
		}
		newList = append(newList, temp)
	}

	ret := make(map[string]interface{})

	// 统计数据
	var sumList []GmTotalMoney
	var TotalMoney, AdvOut, proxyOut float64
	_, err = gmSendMoneyModel.Where(gmSendMoneyModel.QueryTable(new(models.GmSendMoney)), condition).GroupBy("operate_type").
		Aggregate("sum(money*100) as total_money, operate_type").All(&sumList)
	if err == nil {
		for _, item := range sumList {
			TotalMoney += item.TotalMoney

			switch config.GmOperateType(item.OperateType) {
			case config.GmOperateTypeAddAccount, config.GmOperateTypeSubAccount:
				proxyOut += item.TotalMoney
			case config.GmOperateTypeAddAdReward, config.GmOperateTypeSubAdReward:
				AdvOut += item.TotalMoney
			}
		}
	}

	other := map[string]interface{}{
		"TotalCount": total,
		"TotalMoney": TotalMoney / 100,
		"AdvOut":     AdvOut / 100,
		"proxyOut":   proxyOut / 100,
	}

	var (
		gmTypes  []map[string]interface{}
		packages []models.Package
		teams    []models.Team
	)

	if req.NeedReload == 1 {
		gmTypes = s.GetAllGmTypes(lang)
		packages, _ = (&PackageService{}).AllPackage()
		teams, _ = (&TeamService{}).GetAllTeams()
	}

	// 上下分类型
	ret["list"] = newList
	ret["total"] = total
	ret["gm_send_money_type"] = gmTypes
	ret["packages"] = packages
	ret["teams"] = teams
	ret["other"] = other
	return ret, nil
}

// GetAllGmTypes 获取所有上下分类型
func (s *CustomerService) GetAllGmTypes(lang string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.GmOperateTypeAdd, "title": i18n.Tr(lang, "GMShangFen")},
		{"id": config.GmOperateTypeSubtract, "title": i18n.Tr(lang, "GMXiaFen")},
		{"id": config.GmOperateTypeAddInviteReward, "title": i18n.Tr(lang, "GMZengJiaDaiLiJieMianYaoQingJiangLi")},
		{"id": config.GmOperateTypeSubInviteReward, "title": i18n.Tr(lang, "GMJianShaoDaiLiJieMianYaoQingJiangLi")},
		{"id": config.GmOperateTypeAddAdReward, "title": i18n.Tr(lang, "GMGuangGaoFeiXiaFa")},
		{"id": config.GmOperateTypeSubAdReward, "title": i18n.Tr(lang, "GMGuangGaoFeiKouChu")},
		{"id": config.GmOperateTypeAddAccount, "title": i18n.Tr(lang, "GMMoNiZhangHuJiaKuan")},
		{"id": config.GmOperateTypeSubAccount, "title": i18n.Tr(lang, "GMMoNiZhangHuJianKuan")},
	}
}

type GmOperateAddRequest struct {
	Descript    string  `json:"descript"`     // 备注
	Money       string  `json:"money"`        // 变动金币数量 (逗号分隔)
	OperateType int     `json:"operate_type"` // 操作类型
	Password    string  `json:"password"`     // 密码
	RoleId      string  `json:"role_id"`      // 角色ID (逗号分隔)
	WageMul     float64 `json:"wage_mul"`     // 打码倍数
	AdminId     int     `json:"admin_id"`     // 管理员ID
	Ip          string  `json:"ip"`           // IP地址
	PackageIds  []int   `json:"package_ids"`  // 包ID
}

// GmOperateAdd 上下分操作
func (s *CustomerService) GmOperateAdd(req GmOperateAddRequest) error {
	err := s.gmOperateAddValid(req)
	if err != nil {
		return err
	}

	roleIds := strings.Split(req.RoleId, ",")
	moneys := strings.Split(req.Money, ",")

	model := models.CreateGmSendMoneyModel()
	txErr := model.OrmerMaster.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {

		for i := 0; i < len(roleIds); i++ {
			roleId, _ := strconv.Atoi(roleIds[i])
			money, _ := strconv.ParseFloat(moneys[i], 64)

			user, err := (&UserService{}).GetUserByRoleId(roleId)
			if err != nil {
				return err
			}

			if !slices.Contains(req.PackageIds, user.PackageId) {
				logs.Error("GM操作用户不属于你管理的分包下:roleId %v, packageId %v", roleId, user.PackageId)
				return fmt.Errorf("YongHuBuShuYuNiGuanLiDeFenBaoXia")
			}

			// 团队信息
			teamInfo, err := (&TeamService{}).GetUserTeamInfo(roleId)
			if err != nil && !errors.Is(err, orm.ErrNoRows) {
				return err
			}

			_, err = txOrm.Insert(&models.GmSendMoney{
				UserId:       int(user.Id),
				RoleId:       roleId,
				IfInnerProxy: user.IsMock,
				PackageId:    user.PackageId,
				TeamId:       teamInfo.Pid,
				SaleMenId:    teamInfo.Id,
				Money:        money,
				Status:       0, // 未审核
				OperateType:  req.OperateType,
				Type:         0,
				WageMul:      req.WageMul,
				CheckUserId:  0,
				CheckUser:    "",
				Note:         req.Descript,
				CTime:        time.Now().Unix(),
				UTime:        0,
			})
			if err != nil {
				logs.Error("插入上下分记录失败: %v", err)
				return err
			}

			_, err = txOrm.Insert(&models.AdminLog{
				AdminId:    req.AdminId,
				PlayerId:   int(user.Id),
				RoleId:     user.RoleId,
				Path:       "api/customer/gmOperateAdd",
				Controller: "CustomerController",
				Action:     "GmOperateAdd",
				Body:       utils.ToJson(req),
				Status:     1,
				LogType:    int(config.AdminLogTypeGmOperateAdd),
				Ip:         req.Ip,
				CTime:      time.Now().Unix(),
			})
			if err != nil {
				logs.Error("gm上下分插入管理员日志失败: %v", err)
				return err
			}
		}

		return nil
	})

	return txErr
}

// gmOperateAddValid 上下分操作添加验证
func (s *CustomerService) gmOperateAddValid(req GmOperateAddRequest) error {
	if req.RoleId == "" {
		return fmt.Errorf("IDBiTian")
	}
	ids := strings.Split(req.RoleId, ",")
	for _, id := range ids {
		roleId, _ := strconv.Atoi(id)
		_, err := (&UserService{}).GetUserByRoleId(roleId)
		if err != nil {
			return err
		}
	}

	if req.Money == "" {
		return fmt.Errorf("JinEBuNengXiaoYuDengYuLing")
	}

	if req.OperateType == 0 {
		return fmt.Errorf("LeiXingBiTian")
	}

	if req.Password == "" {
		return fmt.Errorf("MiMaBuNengWeiKong")
	}

	if !(&AccountService{}).VerifyUserPassword(req.AdminId, req.Password) {
		return fmt.Errorf("MiMaCuoWu")
	}

	return nil
}

// GmOperateCheck 上下分操作检查
func (s *CustomerService) GmOperateCheck(idStr string, checkUserId int, checkUser string) error {
	if idStr == "" {
		return fmt.Errorf("IDBiTian")
	}

	ids := strings.Split(idStr, ",")
	for _, id := range ids {
		idInt, _ := strconv.Atoi(id)
		var gmSendMoney models.GmSendMoney
		model := models.CreateGmSendMoneyModel()
		err := model.QueryTable(new(models.GmSendMoney)).Filter("id", idInt).One(&gmSendMoney)
		if err != nil {
			return err
		}

		if gmSendMoney.Status != 0 {
			return fmt.Errorf("YiCaoZuoGuo")
		}

		err = s.GmOperateRun(gmSendMoney)
		if err != nil {
			return err
		}

		_, err = model.Update(&models.GmSendMoney{
			Id:          idInt,
			Status:      1,
			CheckUserId: checkUserId,
			CheckUser:   checkUser,
			UTime:       time.Now().Unix(),
		}, "status", "check_user_id", "check_user", "u_time")
		if err != nil {
			logs.Error("审核更新上下分记录状态失败: id: %v, err: %v", idInt, err)
			return err
		}
	}

	return nil
}

// GmOperateRun 上下分操作执行
func (s *CustomerService) GmOperateRun(gmSendMoney models.GmSendMoney) error {

	var balanceLogType int // 余额变动类型

	switch config.GmOperateType(gmSendMoney.OperateType) {
	case config.GmOperateTypeAdd:
		err := (&UserService{}).AddUserBalance(int64(gmSendMoney.UserId), gmSendMoney.Money)
		if err != nil {
			logs.Error("gm上分失败: userId: %v, money: %v, err: %v", gmSendMoney.UserId, gmSendMoney.Money, err)
			return err
		}
		balanceLogType = int(config.BalanceLogTypeGmAdd)
	case config.GmOperateTypeSubtract:
		err := (&UserService{}).AddUserBalance(int64(gmSendMoney.UserId), -gmSendMoney.Money)
		if err != nil {
			logs.Error("gm下分失败: userId: %v, money: %v, err: %v", gmSendMoney.UserId, gmSendMoney.Money, err)
			return err
		}
		balanceLogType = int(config.BalanceLogTypeGmSubtract)
	case config.GmOperateTypeAddInviteReward:
		// TODO: 执行上下分操作逻辑
	case config.GmOperateTypeSubInviteReward:
		// TODO: 执行上下分操作逻辑
	case config.GmOperateTypeAddAdReward:
		// TODO: 执行上下分操作逻辑
	case config.GmOperateTypeSubAdReward:
		// TODO: 执行上下分操作逻辑
	case config.GmOperateTypeAddAccount:
		// TODO: 执行上下分操作逻辑
	case config.GmOperateTypeSubAccount:
		// TODO: 执行上下分操作逻辑
	default:
		logs.Error("gm上下分操作类型不存在: %v", gmSendMoney.OperateType)
		return fmt.Errorf("LeiXingBuCunZai")
	}

	// 记录用户用户余额变动日志
	err := (&UserService{}).AddUserBalanceLog(AddUserBalanceLogRequest{
		UserId:    gmSendMoney.UserId,
		RoleId:    int64(gmSendMoney.RoleId),
		PackageId: gmSendMoney.PackageId,
		Amount:    gmSendMoney.Money,
		Type:      balanceLogType,
		BetRate:   gmSendMoney.WageMul,
	})
	if err != nil {
		logs.Error("gm上下分操作余额变动日志失败: userId: %v, money: %v, type: %v, err: %v", gmSendMoney.UserId, gmSendMoney.Money, gmSendMoney.OperateType, err)
		return err
	}
	return nil
}

// GmOperateRefuse 上下分操作拒绝
func (s *CustomerService) GmOperateRefuse(id int, checkUserId int, checkUser string) error {
	var gmSendMoney models.GmSendMoney
	model := models.CreateGmSendMoneyModel()
	err := model.QueryTable(new(models.GmSendMoney)).Filter("id", id).One(&gmSendMoney)
	if err != nil {
		return err
	}

	if gmSendMoney.Status != 0 {
		return fmt.Errorf("YiCaoZuoGuo")
	}

	_, err = model.Update(&models.GmSendMoney{
		Id:          id,
		Status:      2,
		CheckUserId: checkUserId,
		CheckUser:   checkUser,
		UTime:       time.Now().Unix(),
	}, "status", "check_user_id", "check_user", "u_time")
	if err != nil {
		logs.Error("拒绝更新上下分记录状态失败: id: %v, err: %v", id, err)
		return err
	}
	return nil
}
