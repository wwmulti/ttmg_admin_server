package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/i18n"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/redis/go-redis/v9"
)

type UserService struct{}

type AddUserRequestParams struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	PackageId int    `json:"package_id"`
	Language  string `json:"language,omitempty"`
}

type UserRequestParams struct {
	VipLevel  int    `json:"vip_level"`
	PackageId int    `json:"package_id"`
	RoleId    int    `json:"role_id"`
	Status    int    `json:"status"`
	Username  string `json:"username"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

type UserListResponse struct {
	Total    int              `json:"total"`
	Lists    []models.User    `json:"list"`
	Packages []models.Package `json:"packages"`
}

type UserAttrDTO struct {
	Id                 int                     `json:"id" orm:"auto;column(id)"`                                // 主键
	AwardNeedBets      int64                   `json:"award_need_bets" orm:"column(award_need_bets)"`           // 奖励打码
	EffectiveBet       int64                   `json:"effective_bet" orm:"column(effective_bet)"`               // 累计有效押注(总打码值)
	TotalRecharge      int64                   `json:"total_recharge" orm:"column(total_recharge)"`             // 累计充值
	TotalWithdraw      int64                   `json:"total_withdraw" orm:"column(total_withdraw)"`             // 累计提现
	TotalRechargeCount int                     `json:"total_recharge_count" orm:"column(total_recharge_count)"` // 累计充值次数
	TotalWithdrawCount int                     `json:"total_withdraw_count" orm:"column(total_withdraw_count)"` // 累计提现次数
	TotalProfit        int64                   `json:"total_profit" orm:"column(total_profit)"`                 // 累计盈亏
	VipLevel           int                     `json:"vip_level" orm:"column(vip_level)"`                       // VIP等级
	User               *models.User            `json:"user" orm:"column(user)"`                                 // 关联用户
	UserRelationship   models.UserRelationship `json:"relationship" orm:"column(relationship)"`                 // 上级关系
	Team               SubTeamWithParent       `json:"team" orm:"column(team)"`                                 // 团队
}

// UserList 用户列表
func (s *UserService) UserList(params UserRequestParams, needReload int, language string, adminId int64) (map[string]interface{}, error) {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	userAttrModel := models.CreateUserAttrModel()
	var userAttrList []models.UserAttr
	qs := userAttrModel.QueryTable(new(models.UserAttr)).
		RelatedSel("User")
	if params.Status != 0 {
		qs = qs.Filter("User__Status", params.Status)
	}
	if params.VipLevel != 0 {
		qs = qs.Filter("vip_level", params.VipLevel)
	}
	if params.PackageId != 0 {
		if !isRootAdmin {
			if len(packageSlice) == 0 {
				return nil, nil
			}
			inArray := utils.InArray(params.PackageId, packageSlice)
			if inArray {
				qs = qs.Filter("User__PackageId__In", packageSlice)
			} else {
				return nil, fmt.Errorf("QuanXianBuZu")
			}
		} else {
			qs = qs.Filter("User__PackageId", params.PackageId)
		}
	} else {
		if !isRootAdmin {
			qs = qs.Filter("User__PackageId__in", packageSlice)
		}
	}
	if params.RoleId != 0 {
		qs = qs.Filter("User__RoleId", params.RoleId)
	}
	if params.Username != "" {
		qs = qs.Filter("User__Username__contains", params.Username)
	}
	_, err := qs.Limit(params.PageSize, (params.Page-1)*params.PageSize).
		OrderBy("-User__Id").
		All(&userAttrList)
	fmt.Println(userAttrList)
	if err != nil {
		logs.Error("获取用户列表失败:%v", err)
		return nil, fmt.Errorf("WeiZhiDeCuoWu")
	}
	total, err := qs.Count()
	if err != nil {
		logs.Error("计算用户总数失败:%v", err)
		return nil, fmt.Errorf("WeiZhiDeCuoWu")
	}
	var packages []models.Package
	var userType []map[string]interface{}
	if needReload == 1 {
		packageService := PackageService{}
		packages = packageService.GetMyAllPackageList(int(adminId))
		userType = s.AllUserType(language)
	}
	var userIds []int
	for _, user := range userAttrList {
		userIds = append(userIds, int(user.User.Id))
	}
	userIds = utils.UniqueInt(userIds)
	var userRelationships []models.UserRelationship
	relationshipModel := models.CreateUserRelationshipModel()
	if len(userIds) > 0 {
		_, err = relationshipModel.QueryTable(new(models.UserRelationship)).
			Filter("user_id__in", userIds).
			All(&userRelationships)
		if err != nil {
			logs.Error("获取用户关系失败:%v", err)
			return nil, fmt.Errorf("WeiZhiDeCuoWu")
		}
	}
	teamService := TeamService{}
	userAttrDTO, err := utils.MapSliceWithError(userAttrList, func(userAttr models.UserAttr) (UserAttrDTO, error) {
		var userRelationship models.UserRelationship
		userTeam := SubTeamWithParent{}
		for _, relationship := range userRelationships {
			if relationship.UserId == int(userAttr.User.Id) {
				userRelationship = relationship
				userTeam, err = teamService.GetAllSubTeamsWithParent(userRelationship.Parents)
				if err != nil {
					logs.Error("获取用户团队失败:%v", err)
					return UserAttrDTO{}, fmt.Errorf("WeiZhiDeCuoWu")
				}
			}
		}
		return UserAttrDTO{
			User:               userAttr.User,
			AwardNeedBets:      userAttr.AwardNeedBets,
			EffectiveBet:       userAttr.EffectiveBet,
			TotalRecharge:      userAttr.TotalRecharge,
			TotalWithdraw:      userAttr.TotalWithdraw,
			VipLevel:           userAttr.VipLevel,
			UserRelationship:   userRelationship,
			TotalProfit:        userAttr.TotalProfit,
			TotalRechargeCount: userAttr.TotalRechargeCount,
			TotalWithdrawCount: userAttr.TotalWithdrawCount,
			Team:               userTeam,
		}, nil
	})
	return map[string]interface{}{
		"total":     total,
		"list":      userAttrDTO,
		"packages":  packages,
		"user_type": userType,
	}, nil
}

type GetUserResponse struct {
	User     models.User     `json:"user"`
	UserAttr models.UserAttr `json:"user_attr"`
}

// GetUser 获取用户信息
func (s *UserService) GetUser(userId int, adminId int64) (map[string]interface{}, error) {
	user, err := s.GetUserById(userId)
	if err != nil {
		logs.Error("获取用户信息失败:%v", err)
		return nil, fmt.Errorf("YongHuBuCunZai")
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return nil, fmt.Errorf("QuanXianBuZu")
	}
	userAttrModel := models.CreateUserAttrModel()
	userAttr := models.UserAttr{}
	err = userAttrModel.QueryTable(new(models.UserAttr)).
		Filter("user_id", userId).
		One(&userAttr)
	if err != nil {
		logs.Error("获取用户属性信息失败:%v", err)
		return nil, fmt.Errorf("YongHuBuCunZai")
	}
	packageModel := models.CreatePackageModel()
	userPackage := models.Package{}
	err = packageModel.QueryTable(new(models.Package)).
		Filter("id", user.PackageId).
		One(&userPackage)
	if err != nil {
		logs.Error("获取用户平台信息失败:%v", err)
		return nil, fmt.Errorf("WeiZhiDeCuoWu")
	}
	relationshipModel := models.CreateUserRelationshipModel()
	userRelationship := models.UserRelationship{}
	err = relationshipModel.QueryTable(new(models.UserRelationship)).
		Filter("user_id", userId).
		One(&userRelationship)
	if err != nil {
		logs.Error("获取用户关系信息失败:%v", err)
		return nil, fmt.Errorf("WeiZhiDeCuoWu")
	}
	userTeam := SubTeamWithParent{}
	teamService := TeamService{}
	userTeam, err = teamService.GetAllSubTeamsWithParent(userRelationship.Parents)
	if err != nil {
		logs.Error("获取用户团队失败:%v", err)
		return nil, fmt.Errorf("WeiZhiDeCuoWu")
	}

	return map[string]interface{}{
		"user":         user,
		"user_attr":    userAttr,
		"package":      userPackage,
		"relationship": userRelationship,
		"team":         userTeam,
	}, nil
}

// AddUser 创建用户
func (s *UserService) AddUser(params AddUserRequestParams) (error, int, int) {
	resp := RequestPackageUrlRequest{
		Language:  params.Language,
		PackageId: params.PackageId,
		Type:      int(config.CurlRequestTypePost),
		Params: map[string]interface{}{
			"username":        params.Username,
			"password":        params.Password,
			"account_type":    1,
			"inviter_id":      0,
			"active_id":       0,
			"is_need_user_id": 1,
		},
		Url: "/api/user/register",
	}
	packageService := PackageService{}
	response, err := packageService.RequestPackageUrl(resp)
	if err != nil {
		logs.Error("后台创建用户失败：%v", err)
		return fmt.Errorf("ChuangJianShuJuShiBai"), 0, 0
	}

	// 获取并解析JSON响应
	var userData struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data map[string]interface{} `json:"data"`
	}

	err = response.JSON(&userData)
	logs.Info("创建用户：%v", userData.Data)
	if userData.Code == 200 || err != nil {
		var roleId, userId int

		if id, ok := userData.Data["id"].(float64); ok {
			roleId = int(id)
		}
		if uid, ok := userData.Data["user_id"].(float64); ok {
			userId = int(uid)
		}

		return nil, roleId, userId
	}
	return fmt.Errorf("%s", userData.Msg), 0, 0
}

// EditPassword 修改密码
func (s *UserService) EditPassword(id int, password string, adminId int64) error {
	if len(password) == 0 {
		return nil
	}
	model := models.CreateUserModel()
	user, err := s.GetUserById(id)
	if err != nil {
		return fmt.Errorf("ShuJuYiChang")
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	fields := []string{
		"password",
	}
	pwd := s.generatePassword(password, user.Salt)
	_, updateErr := model.Update(&models.User{
		Id:       int64(id),
		Password: pwd,
	}, fields...)
	if updateErr != nil {
		logs.Error("更新用户密码失败:%v", updateErr)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

// AllUserType 所有轮播图类型
func (s *UserService) AllUserType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.UserTypeNormal, "name": i18n.Tr(language, "puTongYongHu")},
		{"id": config.UserTypeBlogger, "name": i18n.Tr(language, "boZhuYongHu")},
		{"id": config.UserTypeBroker, "name": i18n.Tr(language, "jingJiRenYongHu")},
	}
}

// generatePassword 生成密码
func (s *UserService) generatePassword(password string, salt string) string {
	first := md5.Sum([]byte(password + salt))
	firstStr := hex.EncodeToString(first[:])

	second := md5.Sum([]byte(firstStr))
	return hex.EncodeToString(second[:])
}

// GetUserEffectiveBet 获取用户打码值
func (s *UserService) GetUserEffectiveBet(userId int) float64 {
	val := s.GetUserAttr(userId, config.UserAttr_EffectiveBet)
	bets, _ := strconv.ParseInt(val, 10, 64)
	return utils.ItoExFloat64(bets)
}

type SetUserEffectiveBetRequest struct {
	UserId       int     `json:"user_id"`
	EffectiveBet float64 `json:"effective_bet"`
}

// SetUserEffectiveBet 设置用户打码值
func (s *UserService) SetUserEffectiveBet(request SetUserEffectiveBetRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	val := utils.FtoExInt64(request.EffectiveBet)
	err = s.SetUserAttr(request.UserId, config.UserAttr_EffectiveBet, val)
	if err != nil {
		return err
	}

	_, err = models.CreateUserAttrModel().QueryTable(&models.UserAttr{}).
		Filter("user_id", request.UserId).Update(orm.Params{
		"effective_bet": val,
	})
	if err != nil {
		return err
	}
	return nil
}

// GetUserAttr 获取用户单个属性缓存
func (s *UserService) GetUserAttr(userId int, fieldName string) string {
	ctx := context.Background()
	redisCache := utils.GetRedisClient()
	key := fmt.Sprintf("%s%d", config.RedisKeyName.UserAttr, userId)

	val, err := redisCache.HGet(ctx, key, fieldName).Result()
	if err == redis.Nil {
		if fieldName == config.UserAttr_PackageId {
			return "1" // 默认包ID
		}
		return "0"
	} else if err != nil {
		logs.Error("获取用户属性失败, userId: %d, fieldName: %s, error: %v", userId, fieldName, err)
		return "0"
	}
	return val
}

// SetUserAttr 设置用户单个属性缓存
func (s *UserService) SetUserAttr(userId int, attrName string, value interface{}) error {
	ctx := context.Background()
	redisCache := utils.GetRedisClient()
	key := fmt.Sprintf("%s%d", config.RedisKeyName.UserAttr, userId)
	return redisCache.HSet(ctx, key, attrName, value).Err()
}

type SetUserBalanceRequest struct {
	UserId  int     `json:"user_id"` // 用户ID
	Balance float64 `json:"balance"` // 金额
}

// SetUserBalance 设置用户余额
func (s *UserService) SetUserBalance(request SetUserBalanceRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		logs.Error("设置用户余额获取用户失败：用户id %d, error: %v", request.UserId, err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	balance := utils.FtoExInt64(request.Balance)
	userModel := models.CreateUserModel()
	txErr := userModel.OrmerMaster.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		err = s.AddUserBalanceLog(AddUserBalanceLogRequest{
			UserId:    request.UserId,
			RoleId:    user.RoleId,
			PackageId: user.PackageId,
			Amount:    request.Balance,
			Type:      int(config.BalanceLogTypeSetUserBalance),
		}, txOrm)
		if err != nil {
			logs.Error("设置用户余额添加余额日志失败：用户id %d, error: %v", request.UserId, err)
			return err
		}

		_, err = txOrm.QueryTable(new(models.User)).Filter("id", request.UserId).Update(orm.Params{"balance": balance})
		if err != nil {
			logs.Error("设置用户余额更新用户余额失败：用户id %d, 余额 %v, error: %v", request.UserId, balance, err)
			return err
		}

		redisCache := utils.GetRedisClient()
		ctx = context.Background()
		err = redisCache.HSet(ctx, config.RedisKeyName.UserMoney, strconv.Itoa(request.UserId), balance).Err()
		if err != nil {
			logs.Error("设置用户余额更新用户余额缓存失败：用户id %d, 余额 %v, error: %v", request.UserId, balance, err)
			return err
		}
		return nil
	})

	return txErr
}

// AddUserBalance 添加用户余额 操作要放在最后不然缓存不会滚需要手动回滚
func (s *UserService) AddUserBalance(userId int64, amount float64, isCalculateReward ...bool) error {
	money := utils.FtoExInt64(amount)
	// 缓存保存金额
	balance, cacheErr := s.saveMoney(int(userId), money)
	if cacheErr != nil {
		return cacheErr
	}
	if balance < 0 {
		// 金额扣除过多,还原金额
		money = 0 - money
		_, cacheError := s.saveMoney(int(userId), money)
		if cacheError != nil {
			logs.Error("用户ID:%v扣除金额:%v还原失败:%v", userId, money, cacheError)
		}
		return fmt.Errorf("YuEBuZu")
	}
	go s.updateUserBalance(userId, balance)
	return nil
}

// 异步更新用户金额
func (s *UserService) updateUserBalance(uid int64, money int64) {
	model := models.CreateUserModel()
	model.QueryTable(new(models.User)).Filter("id", uid).Update(orm.Params{"balance": money})
}

// GetBalance 获取用户金额
func (s *UserService) GetBalance(uid int) (float64, error) {
	// model := models.CreateUserModel()
	// user := models.User{Id: uid}
	// model.Read(&user)
	// return utils.ItoExFloat64(user.Balance)
	v, err := s.getMoney(uid)
	return v, err
}

// 获取用户金额
func (s *UserService) getMoney(userId int) (float64, error) {
	redisClient := utils.GetRedisClient()
	ctx := context.Background()
	str, redisErr := redisClient.HGet(ctx, config.RedisKeyName.UserMoney, strconv.Itoa(userId)).Result()

	if errors.Is(redisErr, redis.Nil) {
		return 0, nil
	}
	if redisErr != nil {
		return 0, redisErr
	}
	v, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return utils.ItoExFloat64(int64(v)), nil
}

// 保存用户金额
func (s *UserService) saveMoney(userId int, money int64) (int64, error) {
	redisCache := utils.GetRedisClient()
	ctx := context.Background()
	v, err := redisCache.HIncrBy(ctx, config.RedisKeyName.UserMoney, strconv.Itoa(userId), money).Result()
	return v, err
}

type SetUserLvRequest struct {
	UserId int `json:"user_id"` // 用户ID
	Lv     int `json:"lv"`      // 等级
}

// SetUserLv 设置用户等级
func (s *UserService) SetUserLv(request SetUserLvRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userLvModel := models.CreateActiveVipLvModel()
	tx, err := userLvModel.Begin()
	if err != nil {
		logs.Error("事务开启失败：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	activityId, err := (&ActiveVipService{}).GetValidVipLvActiveId(request.UserId)
	if err != nil {
		return err
	}
	activeVipRule := models.ActiveVipRule{}
	activeVipRuleModel := models.CreateActiveVipRuleModel()
	err = activeVipRuleModel.QueryTable(new(models.ActiveVipRule)).
		Filter("active_id", activityId).
		OrderBy("-lv").
		One(&activeVipRule)
	if err != nil {
		logs.Error("获取会员等级失败：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	if request.Lv > activeVipRule.Lv {
		return fmt.Errorf("BuNengChaoGuoZuiDaDengJi")
	}
	_, err = tx.QueryTable(new(models.ActiveVipLv)).
		Filter("user_id", request.UserId).
		Filter("active_id", activityId).
		Update(orm.Params{
			"lv": request.Lv,
		})
	if err != nil {
		tx.Rollback()
		logs.Error("更新会员等级失败：%v", err)
		return err
	}
	_, err = tx.QueryTable(new(models.UserAttr)).
		Filter("user_id", request.UserId).
		Update(orm.Params{
			"vip_level": request.Lv,
		})
	if err != nil {
		tx.Rollback()
		logs.Error("更新会员等级失败：%v", err)
		return err
	}
	log := models.ActiveVipLvLog{
		UserId:   request.UserId,
		ActiveId: activityId,
		Ctime:    time.Now().Unix(),
		Lv:       request.Lv,
	}
	logModel := models.CreateActiveVipLvLogModel()
	err = logModel.RecordLog(log, tx)
	if err != nil {
		tx.Rollback()
		logs.Error("插入日志失败：%v", err)
		return err
	}
	// 更新用户等级缓存
	err = s.SetUserAttr(request.UserId, config.UserAttr_VipLevel, request.Lv)
	if err != nil {
		tx.Rollback()
		logs.Error("更新用户等级缓存失败：%v", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("事务提交失败：%v", err)
		return err
	}
	return nil
}

// AddAdminLog 后台操作日志
func (s *UserService) AddAdminLog(data models.AdminLog) {
	_, err := models.CreateAdminLogModel().Insert(&data)
	if err != nil {
		logs.Error("后台操作日志失败:%v", err)
	}
}

type ChangeUserStatusRequest struct {
	UserId int `json:"user_id"`
	Status int `json:"status"` // 状态
}

// ChangeUserStatus 修改用户封禁状态
func (s *UserService) ChangeUserStatus(request ChangeUserStatusRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	if request.UserId == 0 {
		return fmt.Errorf("IDBiTian")
	}
	_, err = userModel.Update(&models.User{
		Id:     int64(request.UserId),
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}
	return nil
}

type SetUserPidRequest struct {
	UserId  int `json:"user_id"`   // 用户ID
	PRoleId int `json:"p_role_id"` // 上级角色ID
}

// SetUserPid 设置用户上级
func (s *UserService) SetUserPid(request SetUserPidRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userId := request.UserId
	pRoleId := request.PRoleId
	parentUser, err := s.GetUserByRoleId(pRoleId)
	if err != nil {
		logs.Error("获取父级用户失败：%v", err)
		return fmt.Errorf("ShangJiYongHuBuCunZai")
	}
	pid := parentUser.Id
	relationshipModel := models.CreateUserRelationshipModel()
	userModel := models.CreateUserModel()
	// 2. 当前用户关系
	current := models.UserRelationship{UserId: userId}
	if err := relationshipModel.Read(&current); err != nil {
		return fmt.Errorf("DangQianYongHuBuCunZai")
	}
	oldPrefix := current.Parents
	// 3. 父级关系
	parent := models.UserRelationship{UserId: int(pid)}
	if err := relationshipModel.Read(&parent); err != nil {
		return fmt.Errorf("FuJiGuanXiBuCunZai")
	}
	// 4. 防止挂到自己的子孙节点
	if strings.Contains(parent.Parents, fmt.Sprintf(",%d,", userId)) {
		return fmt.Errorf("BuNengGuaDaoZiJiDeZiSun")
	}
	// 5. 构建新的关系链
	newPrefix := fmt.Sprintf("%s%d,", parent.Parents, userId)
	current.Parents = newPrefix
	current.PRoleId = int64(pRoleId)
	current.Pid = parent.UserId
	current.Pid2 = parent.Pid
	current.Pid3 = parent.Pid2
	tx, err := relationshipModel.Begin()
	if err != nil {
		return fmt.Errorf("事务开启失败：%v", err)
	}
	if _, err := tx.Update(
		&current,
		"Parents",
		"Pid",
		"Pid2",
		"Pid3",
		"PRoleId",
	); err != nil {
		_ = tx.Rollback()
		return err
	}
	// 更新用户上级缓存
	if err := s.SetUserAttr(userId, config.UserAttr_Parents, newPrefix); err != nil {
		tx.Rollback()
		logs.Error("更新用户上级缓存失败：%v", err)
		return err
	}

	// 6. 查询所有子孙节点
	var children []models.UserRelationship
	_, err = relationshipModel.QueryTable(new(models.UserRelationship)).
		Filter("Parents__startswith", oldPrefix).
		Exclude("UserId", userId).
		All(&children)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	// 7. 批量更新子孙节点
	for _, item := range children {
		item.Parents = strings.Replace(item.Parents, oldPrefix, newPrefix, 1)
		p1, p2, p3 := s.parseTopParents(item.Parents, item.UserId)
		item.Pid = p1
		item.Pid2 = p2
		item.Pid3 = p3
		if item.Pid > 0 {
			pUser := models.User{Id: int64(item.Pid)}
			if err := userModel.Read(&pUser); err == nil {
				item.PRoleId = pUser.RoleId
			}
		} else {
			item.PRoleId = 0
		}
		if _, err := tx.Update(
			&item,
			"Parents",
			"Pid",
			"Pid2",
			"Pid3",
			"PRoleId",
		); err != nil {
			_ = tx.Rollback()
			return err
		}
		// 更新子孙节点上级缓存
		if err := s.SetUserAttr(item.UserId, config.UserAttr_Parents, item.Parents); err != nil {
			tx.Rollback()
			logs.Error("更新子孙节点上级缓存失败：%v", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}

// 解析顶级父级
func (s *UserService) parseTopParents(parents string, selfId int) (int, int, int) {
	path := strings.Trim(parents, ",")
	if path == "" {
		return 0, 0, 0
	}
	parts := strings.Split(path, ",")
	// 去掉自己，只保留祖先链
	ancestors := make([]int, 0, len(parts))
	for _, part := range parts {
		id, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		if id != selfId {
			ancestors = append(ancestors, id)
		}
	}
	var pid, pid2, pid3 int
	n := len(ancestors)
	if n >= 1 {
		pid = ancestors[n-1]
	}
	if n >= 2 {
		pid2 = ancestors[n-2]
	}
	if n >= 3 {
		pid3 = ancestors[n-3]
	}
	return pid, pid2, pid3
}

type SetUserRemarkRequest struct {
	UserId int    `json:"user_id"` // 用户ID
	Remark string `json:"remark"`
}

// SetUserRemark 设置用户备注
func (s *UserService) SetUserRemark(request SetUserRemarkRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	_, err = userModel.Update(&models.User{
		Id:     int64(request.UserId),
		Remark: request.Remark,
	}, "remark")
	if err != nil {
		logs.Error("修改用户备注失败：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// SetUserBlogger 设置用户博主
func (s *UserService) SetUserBlogger(userId int, adminId int64) error {
	user, err := s.GetUserById(userId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	_, err = userModel.Update(&models.User{
		Id:       int64(userId),
		UserType: int(config.UserTypeBlogger),
	}, "user_type")
	if err != nil {
		logs.Error("设置用户博主失败：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// SetUserBroker 设置用户经纪人
func (s *UserService) SetUserBroker(userId int, adminId int64) error {
	user, err := s.GetUserById(userId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	_, err = userModel.Update(&models.User{
		Id:       int64(userId),
		UserType: int(config.UserTypeBroker),
	}, "user_type")
	if err != nil {
		logs.Error("设置用户经纪人失败：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// CancelUserBlogger 取消用户博主
func (s *UserService) CancelUserBlogger(userId int, adminId int64) error {
	user, err := s.GetUserById(userId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	_, err = userModel.Update(&models.User{
		Id:       int64(userId),
		UserType: 0,
	}, "user_type")
	if err != nil {
		logs.Error("取消用户博主失败：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// CancelUserBroker 取消用户经纪人
func (s *UserService) CancelUserBroker(userId int, adminId int64) error {
	user, err := s.GetUserById(userId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	_, err = userModel.Update(&models.User{
		Id:       int64(userId),
		UserType: 0,
	}, "user_type")
	if err != nil {
		logs.Error("取消用户经纪人失败：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

type ChangeUserGameStatusRequest struct {
	UserId    int `json:"user_id"`
	IsBanGame int `json:"is_ban_game"` // 是否禁止游戏
}

// ChangeUserGameStatus 修改用户游戏封禁状态
func (s *UserService) ChangeUserGameStatus(request ChangeUserGameStatusRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	if request.UserId == 0 {
		return fmt.Errorf("IDBiTian")
	}
	_, err = userModel.Update(&models.User{
		Id:        int64(request.UserId),
		IsBanGame: request.IsBanGame,
	}, "is_ban_game")
	if err != nil {
		return err
	}
	return nil
}

type ChangeUserWithdrawStatusRequest struct {
	UserId        int `json:"user_id"`
	IsBanWithdraw int `json:"is_ban_withdraw"` // 是否禁止提现
}

// ChangeUserWithdrawStatus 修改用户提现封禁状态
func (s *UserService) ChangeUserWithdrawStatus(request ChangeUserWithdrawStatusRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	if request.UserId == 0 {
		return fmt.Errorf("IDBiTian")
	}
	_, err = userModel.Update(&models.User{
		Id:            int64(request.UserId),
		IsBanWithdraw: request.IsBanWithdraw,
	}, "is_ban_withdraw")
	if err != nil {
		return err
	}
	return nil
}

type ChangeUserInviteRewardStatusRequest struct {
	UserId            int `json:"user_id"`
	IsBanInviteReward int `json:"is_ban_invite_reward"` // 是否禁止提现
}

// ChangeUserInviteRewardStatus 修改用户邀请奖励封禁状态
func (s *UserService) ChangeUserInviteRewardStatus(request ChangeUserInviteRewardStatusRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	if request.UserId == 0 {
		return fmt.Errorf("IDBiTian")
	}
	_, err = userModel.Update(&models.User{
		Id:                int64(request.UserId),
		IsBanInviteReward: request.IsBanInviteReward,
	}, "is_ban_invite_reward")
	if err != nil {
		return err
	}
	return nil
}

type ChangeUserChildBetCommissionStatusRequest struct {
	UserId                  int `json:"user_id"`
	IsBanChildBetCommission int `json:"is_ban_child_bet_commission"` // 是否禁止提现
}

// ChangeUserChildBetCommissionStatus 修改用户禁止下级返佣
func (s *UserService) ChangeUserChildBetCommissionStatus(request ChangeUserChildBetCommissionStatusRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	if request.UserId == 0 {
		return fmt.Errorf("IDBiTian")
	}
	_, err = userModel.Update(&models.User{
		Id:                      int64(request.UserId),
		IsBanChildBetCommission: request.IsBanChildBetCommission,
	}, "is_ban_child_bet_commission")
	if err != nil {
		return err
	}
	return nil
}

type ChangeUserOnlyAllowOfficialGameStatusRequest struct {
	UserId                  int `json:"user_id"`
	IsOnlyAllowOfficialGame int `json:"is_only_allow_official_game"` // 是否禁止提现
}

// ChangeUserOnlyAllowOfficialGameStatus 修改用户邀请奖励封禁状态
func (s *UserService) ChangeUserOnlyAllowOfficialGameStatus(request ChangeUserOnlyAllowOfficialGameStatusRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	if request.UserId == 0 {
		return fmt.Errorf("IDBiTian")
	}
	_, err = userModel.Update(&models.User{
		Id:                      int64(request.UserId),
		IsOnlyAllowOfficialGame: request.IsOnlyAllowOfficialGame,
	}, "is_only_allow_official_game")
	if err != nil {
		return err
	}
	return nil
}

type ChangeUserMockStatusRequest struct {
	UserId int `json:"user_id"`
	IsMock int `json:"is_mock"` // 是否模拟账户
}

// ChangeUserMockStatus 修改用户模拟账户状态
func (s *UserService) ChangeUserMockStatus(request ChangeUserMockStatusRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	userModel := models.CreateUserModel()
	if request.UserId == 0 {
		return fmt.Errorf("IDBiTian")
	}
	_, err = userModel.Update(&models.User{
		Id:     int64(request.UserId),
		IsMock: request.IsMock,
	}, "is_mock")
	if err != nil {
		return err
	}
	return nil
}

type ChangeUserPasswordRequest struct {
	UserId           int    `json:"user_id"`
	Password         string `json:"password"`
	WithdrawPassword string `json:"withdraw_password"`
}

// ChangeUserPassword 修改用户密码
func (s *UserService) ChangeUserPassword(request ChangeUserPasswordRequest, adminId int64) error {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	password := ""
	withdrawPassword := ""
	if request.UserId == 0 {
		return fmt.Errorf("IDBiTian")
	}
	if request.Password == "" {
		password = user.Password
	} else {
		password = s.generatePassword(request.Password, user.Salt)
	}
	if request.WithdrawPassword == "" {
		withdrawPassword = user.WithdrawPassword
	} else {
		withdrawPassword = s.generatePassword(request.WithdrawPassword, user.Salt)
	}
	userModel := models.CreateUserModel()
	_, err = userModel.Update(&models.User{
		Id:               int64(request.UserId),
		Password:         password,
		WithdrawPassword: withdrawPassword,
	}, "password", "withdraw_password")
	if err != nil {
		logs.Error("修改用户密码失败：%v", err)
		return err
	}
	return nil
}

// GetUserById 获取用户信息
func (s *UserService) GetUserById(userId int) (*models.User, error) {
	user := &models.User{}
	userModel := models.CreateUserModel()
	err := userModel.QueryTable(new(models.User)).
		Filter("id", userId).Filter("is_deleted", 0).
		One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByRoleId 获取用户信息
func (s *UserService) GetUserByRoleId(roleId int) (*models.User, error) {
	user := &models.User{}
	userModel := models.CreateUserModel()
	err := userModel.QueryTable(new(models.User)).
		Filter("role_id", roleId).Filter("is_deleted", 0).
		One(user)
	if err != nil {
		return nil, fmt.Errorf("YongHuBuCunZai")
	}
	return user, nil
}

type UserBetStatisticListRequest struct {
	UserId   int        `form:"user_id"`   // ID
	Page     int        `form:"page"`      // 页码
	PageSize int        `form:"page_size"` // 每页数量
	RawQuery url.Values `form:"-"`
}

// UserBetStatisticList 打码记录列表
func (s *UserService) UserBetStatisticList(request UserBetStatisticListRequest, adminId int64) (map[string]interface{}, error) {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	userBetStatisticModel := models.CreateUserBetStatisticModel()
	condition, sort := userBetStatisticModel.BuildCondition(request, "-id")
	if !isRootAdmin {
		condition["package_id__in"] = packageSlice
	}
	data, total, err := userBetStatisticModel.GetPageList(&models.UserBetStatistic{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}

	return map[string]interface{}{
		"list":         data,
		"total":        total,
		"current_page": request.Page,
	}, nil
}

type AddUserBalanceLogRequest struct {
	UserId         int     // 用户ID
	RoleId         int64   // 用户RoleId
	PackageId      int     // 包ID
	Amount         float64 // 变动金额
	Type           int     // 类型
	ActivityTypeId int     // 活动类型
	ActivityTitle  string  // 活动名称
	ActivityId     int     // 活动ID
	PlatformId     int     // 平台ID
	PlatformTitle  string  // 平台名称
	GameTypeId     int     // 游戏类型ID
	GameId         int     // 游戏ID
	GameTitle      string  // 游戏名称
	Mark           string  // 备注
	BetRate        float64 // 打码倍数
}

// AddUserBalanceLog 添加用户余额变动日志
func (s *UserService) AddUserBalanceLog(params AddUserBalanceLogRequest, tx ...orm.TxOrmer) error {
	// 计算打码金额
	go func() {
		s.CalculateBetAmount(params)
	}()

	// 累计金额
	go func() {
		s.AccumulateUserAmount(params)
	}()

	now := time.Now().Unix()
	balance, err := s.GetBalance(params.UserId)
	if err != nil {
		logs.Error("获取用户余额失败,错误内容为%v", err)
		return err
	}

	log := models.UserBalanceLog{
		UserId:         params.UserId,
		Balance:        balance,
		RoleId:         params.RoleId,
		PackageId:      params.PackageId,
		GameTypeId:     params.GameTypeId,
		GameId:         params.GameId,
		GameTitle:      params.GameTitle,
		Amount:         params.Amount,
		Type:           params.Type,
		ActivityTypeId: params.ActivityTypeId,
		ActivityTitle:  params.ActivityTitle,
		ActivityId:     params.ActivityId,
		PlatformId:     params.PlatformId,
		PlatformTitle:  params.PlatformTitle,
		Mark:           params.Mark,
		CreatedTime:    now,
	}
	err = (&models.UserBalanceLogModel{}).RecordLog(log, tx...)
	if err != nil {
		return err
	}
	return nil
}

// CalculateBetAmount 计算打码金额
func (s *UserService) CalculateBetAmount(params AddUserBalanceLogRequest) float64 {
	if params.Amount <= 0 {
		return 0.0
	}

	// 忽略的金额变动类型
	exists := slices.Contains(config.IgnoreBalanceLogTypes, config.BalanceLogType(params.Type))
	if exists {
		return 0.0
	}

	var rate float64
	if params.BetRate > 0 {
		rate = params.BetRate
	} else {
		// 获取打码倍数
		betRate, err := (BetRateService{}).GetBetRate(params.PackageId, params.Type)
		if err != nil {
			logs.Error("获取打码倍数失败, packageId: %v, type: %v, error: %v", params.PackageId, params.Type, err)
			return 0.0
		}

		// TODO: 根据用户id判断取哪个打码配置
		rate = betRate.ProRate
	}

	amount := utils.FtoExInt64(params.Amount) * utils.FtoExInt64(rate) / 1000

	// redis累计
	ctx := context.Background()
	redisCache := utils.GetRedisClient()
	key := fmt.Sprintf("%s%d", config.RedisKeyName.UserAttr, params.UserId)
	redisCache.HIncrBy(ctx, key, config.UserAttr_AwardNeedBets, amount)

	// 更新打码金额
	affected, err := models.CreateUserAttrModel().QueryTable(new(models.UserAttr)).
		Filter("user_id", params.UserId).
		Update(orm.Params{
			"award_need_bets": orm.ColValue(orm.ColAdd, amount),
		})
	if err != nil || affected == 0 {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		logs.Error("累计用户打码金额失败 用户id:%v, 金额变动类型:%v, 打码金额:%v, 错误:%v", params.UserId, params.Type, utils.ItoExFloat64(amount), errMsg)
	}

	// 记录打码日志
	log := models.UserBetStatistic{
		UserId:    params.UserId,
		Amount:    utils.FtoExInt64(params.Amount),
		NeedBets:  amount,
		BetRate:   rate,
		Remaining: amount,
		BetType:   params.Type,
		CTime:     time.Now().Unix(),
	}
	_, err = models.CreateUserBetStatisticModel().Insert(&log)
	if err != nil {
		logs.Error("记录打码日志失败 err:%v", err)
		logs.Error("记录打码日志失败 记录:%v", utils.ToJson(log))
	}

	return utils.ItoExFloat64(amount)
}

// AccumulateUserAmount 累计金额
func (s *UserService) AccumulateUserAmount(params AddUserBalanceLogRequest) error {
	var mysqlColumn, redisField string
	amountInt := utils.FtoExInt64(params.Amount)

	switch params.Type {
	case int(config.BalanceLogTypeAgent): // 代理奖励
		mysqlColumn = "total_agent_bonus"
		redisField = config.UserAttr_TotalAgentBonus
	case int(config.BalanceLogTypeGmAdd): // GM上分
		mysqlColumn = "total_gm_send"
		redisField = config.UserAttr_TotalGmSend
	default:
		return nil
	}

	if mysqlColumn != "" {
		_, err := models.CreateUserAttrModel().QueryTable(new(models.UserAttr)).
			Filter("user_id", params.UserId).
			Update(orm.Params{
				mysqlColumn: orm.ColValue(orm.ColAdd, amountInt),
			})
		if err != nil {
			logs.Error("更新数据库累计金额失败: %v", err)
			return err
		}
	}

	if redisField != "" {
		key := fmt.Sprintf("%s%d", config.RedisKeyName.UserAttr, params.UserId)
		err := utils.GetRedisClient().HIncrBy(context.Background(), key, redisField, amountInt).Err()
		if err != nil {
			logs.Error("更新Redis累计金额失败: %v", err)
		}
	}

	return nil
}

type UserSameIpLoginLogRequest struct {
	UserId    int
	BeginTime int64
	EndTime   int64
	Page      int
	PageSize  int
}

// UserSameIpLoginLogList 同一个IP登录日志列表
func (s *UserService) UserSameIpLoginLogList(request UserSameIpLoginLogRequest, adminId int64) (map[string]interface{}, error) {
	user, err := s.GetUserById(request.UserId)
	if err != nil {
		return nil, err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(user.PackageId, adminId)
	if !hasAuthority {
		return nil, fmt.Errorf("QuanXianBuZu")
	}
	loginIp := user.IP
	loginLogService := &LoginLogService{}
	data, err := loginLogService.GetList(LoginLogRequestParams{
		BeginTime: request.BeginTime,
		EndTime:   request.EndTime,
		IP:        loginIp,
		Page:      request.Page,
		PageSize:  request.PageSize,
	})
	if err != nil {
		logs.Error("获取用户登录日志失败 userId:%v, error:%v", request.UserId, err)
		return nil, fmt.Errorf("WeiZhiDeCuoWu")
	}
	return map[string]interface{}{
		"list":         data.List,
		"total":        data.Total,
		"current_page": request.Page,
	}, nil
}

// UserLoginLogList 获取用户登录日志
func (s *UserService) UserLoginLogList(request LoginLogRequestParams, needReload int, adminId int64) (map[string]interface{}, error) {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	if !isRootAdmin {
		if request.PackageId != 0 {
			if !utils.InArray(request.PackageId, packageSlice) {
				return nil, fmt.Errorf("QuanXianBuZu")
			}
		} else {
			if len(packageSlice) == 0 {
				return nil, fmt.Errorf("QuanXianBuZu")
			}
			request.PackageIds = utils.IntSliceToString(packageSlice, ",")
		}
	}
	result, err := (&LoginLogService{}).GetList(request)
	if nil != err {
		return nil, err
	}
	var packages []models.Package
	if needReload == 1 {
		packageService := PackageService{}
		packages = packageService.GetMyAllPackageList(int(adminId))
	}
	return map[string]interface{}{
		"list":     result.List,
		"total":    result.Total,
		"packages": packages,
	}, nil
}

// UserOperateLogList 获取用户登录日志
func (s *UserService) UserOperateLogList(request UserOperateLogRequestParams, needReload int, adminId int64) (map[string]interface{}, error) {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	if !isRootAdmin {
		if request.PackageId != 0 {
			if !utils.InArray(request.PackageId, packageSlice) {
				return nil, fmt.Errorf("QuanXianBuZu")
			}
		} else {
			if len(packageSlice) == 0 {
				return nil, fmt.Errorf("QuanXianBuZu")
			}
			request.PackageIds = utils.IntSliceToString(packageSlice, ",")
		}
	}
	result, err := (&UserOperateLogService{}).GetList(request)
	if nil != err {
		return nil, err
	}
	var packages []models.Package
	if needReload == 1 {
		packageService := PackageService{}
		packages = packageService.GetMyAllPackageList(int(adminId))
	}
	return map[string]interface{}{
		"list":     result.List,
		"total":    result.Total,
		"packages": packages,
	}, nil
}

type GetUserTotalBalanceResponse struct {
	TotalBalance int64
	PackageId    int
}

// GetUserTotalBalance 获取用户总余额
func (s *UserService) GetUserTotalBalance() ([]GetUserTotalBalanceResponse, error) {
	totalAmount := make([]GetUserTotalBalanceResponse, 0)
	_, err := models.CreateUserModel().QueryTable(new(models.User)).
		Filter("is_deleted", 0).
		GroupBy("package_id").
		Aggregate("sum(balance) as total_balance, package_id").
		All(&totalAmount)
	if err != nil {
		return totalAmount, err
	}
	return totalAmount, nil
}

type GetUserDailyCountResponse struct {
	TotalCount int
	PackageId  int
}

// GetUserDailyCount 分包统计用户人数
func (s *UserService) GetUserDailyCount(condition map[string]interface{}) ([]GetUserDailyCountResponse, error) {
	model := models.CreateUserModel()
	condition["is_deleted"] = 0
	totalCount := make([]GetUserDailyCountResponse, 0)
	_, err := model.Where(model.QueryTable(new(models.User)), condition).GroupBy("package_id").
		Aggregate("count(id) as total_count, package_id").All(&totalCount)
	if err != nil {
		return totalCount, err
	}
	return totalCount, nil
}

type UserRegisterRetentionResponse struct {
	TotalCount int `json:"total_count"`
	UserId     int `json:"user_id"`
}

type UserLoginRetentionResponse struct {
	TotalCount int `json:"total_count"`
	UserId     int `json:"user_id"`
}

// CalculateAllPackageUserLoginRetention 计算所有包用户登录留存
func (s *UserService) CalculateAllPackageUserLoginRetention(startTime int64, days []int) error {
	packageService := PackageService{}
	packages, err := packageService.AllPackage()
	if err != nil {
		logs.Error("获取所有包错误:", err)
		return err
	}
	// 插入0，0为总平台
	packages = append(packages, models.Package{Id: 0})
	var wg sync.WaitGroup
	errChan := make(chan error, len(packages))
	for _, packageInfo := range packages {
		packageId := packageInfo.Id
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.CalculateUserLoginRetention(startTime, days, packageId)
			if err != nil {
				logs.Error("计算平台用户登录留存错误，平台ID %d, 错误内容 %v", packageId, err)
				errChan <- err
			}
		}()
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

type FixUserLoginRetentionRequest struct {
	PackageIds []int
	StartDate  string
	EndDate    string
}

// CalculatePackagesUserLoginRetention 计算某些包用户登录留存
func (s *UserService) CalculatePackagesUserLoginRetention(request FixUserLoginRetentionRequest) error {
	days := config.System.RetentionDays
	packageIds := request.PackageIds
	timeList, err := utils.GetNextDaysZeroTimestamps(request.StartDate, request.EndDate)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	errChan := make(chan error, len(packageIds))
	for _, packageId := range packageIds {
		for _, dayTime := range timeList {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := s.CalculateUserLoginRetention(dayTime, days, packageId)
				if err != nil {
					logs.Error("计算平台用户登录留存错误，平台ID %d, 错误内容 %v", packageId, err)
					errChan <- err
				}
			}()
		}
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// CalculateUserLoginRetention 计算用户登录留存
func (s *UserService) CalculateUserLoginRetention(startTime int64, days []int, packageId int) error {
	now := time.Now().Unix()
	loginLogService := &LoginLogService{}
	loginStartTime, loginEndTime, ranges := utils.BuildTimePoints(startTime, days)
	yesterday := time.Unix(loginStartTime, 0).AddDate(0, 0, 0).Format("2006-01-02")
	var yesterdayRegisterCount int64
	userRetentionDailyModel := models.CreateUserRetentionDailyModel()
	qs := userRetentionDailyModel.QueryTable(new(models.UserRetentionDaily)).
		Filter("type", config.RetentionTypeLogin)
	if packageId > 0 {
		qs = qs.Filter("package_id", packageId)
	}
	exist := qs.Filter("statistic_date", yesterday).Exist()
	if exist {
		// 如果已创建就跳过
		return nil
	}
	userModel := models.CreateUserModel()
	yesterdayRegisterQs := userModel.QueryTable(new(models.User))
	if packageId > 0 {
		yesterdayRegisterQs = yesterdayRegisterQs.Filter("package_id", packageId)
	}
	yesterdayRegisterCount, _ = yesterdayRegisterQs.Filter("register_at__gte", loginStartTime).
		Filter("register_at__lt", loginEndTime).
		Count()
	userRetentionData := models.UserRetentionDaily{
		StatisticDate: yesterday,
		RegisterCount: int(yesterdayRegisterCount),
		PackageId:     packageId,
		CreatedAt:     now,
		Type:          int(config.RetentionTypeLogin),
	}
	userRetention := 0.00
	var userLoginRetentionCount int
	var userRegisterRetentionCount int
	for _, r := range ranges {
		// 昨日登录 Day天注册的用户数量
		var userLoginRetentionTotal []UserRegisterRetentionResponse
		loginParams := LoginLogGroupListRequestParams{
			BeginTime:         loginStartTime,
			EndTime:           loginEndTime,
			RegisterBeginTime: r.RegisterStart,
			RegisterEndTime:   r.RegisterEnd,
			GroupBy:           "user_id",
			Fields:            "count(user_id) as total_count,user_id",
			List:              &userLoginRetentionTotal,
		}
		var userRegisterRetentionTotal []UserRegisterRetentionResponse
		registerParams := LoginLogGroupListRequestParams{
			RegisterBeginTime: r.RegisterStart,
			RegisterEndTime:   r.RegisterEnd,
			BeginTime:         r.RegisterStart,
			EndTime:           r.RegisterEnd,
			GroupBy:           "user_id",
			Fields:            "count(user_id) as total_count,user_id",
			List:              &userRegisterRetentionTotal,
		}
		if packageId > 0 {
			loginParams.PackageId = packageId
			registerParams.PackageId = packageId
		}
		err := loginLogService.GroupList(loginParams)
		if err != nil {
			logs.Error("统计用户登录留存时获取登录用户失败: %v", err)
			return err
		}
		err = loginLogService.GroupList(registerParams)
		if err != nil {
			logs.Error("统计用户登录留存时获取注册用户失败: %v", err)
			return err
		}

		if len(userLoginRetentionTotal) == 0 {
			userLoginRetentionCount = 0
		} else {
			userLoginRetentionCount = userLoginRetentionTotal[0].TotalCount
		}
		if len(userRegisterRetentionTotal) == 0 {
			userRegisterRetentionCount = 0
		} else {
			userRegisterRetentionCount = userRegisterRetentionTotal[0].TotalCount
		}
		userRetention = s.getUserRetention(userLoginRetentionCount, userRegisterRetentionCount)
		s.setRetentionValue(&userRetentionData, r.Day, userRetention)
	}
	_, err := userRetentionDailyModel.Insert(&userRetentionData)
	if err != nil {
		return err
	}
	return nil
}

// CalculateAllPackageUserRechargeRetention 计算所有包用户充值留存
func (s *UserService) CalculateAllPackageUserRechargeRetention(startTime int64, days []int) error {
	packageService := PackageService{}
	packages, err := packageService.AllPackage()
	if err != nil {
		logs.Error("获取所有包错误:", err)
		return err
	}
	// 插入0，0为总平台
	packages = append(packages, models.Package{Id: 0})
	var wg sync.WaitGroup
	errChan := make(chan error, len(packages))
	for _, packageInfo := range packages {
		packageId := packageInfo.Id
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.CalculateUserRechargeRetention(startTime, days, packageId)
			if err != nil {
				logs.Error("计算平台用户充值留存错误，平台ID %d, 错误内容 %v", packageId, err)
				errChan <- err
			}
		}()
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

type FixUserRechargeRetentionRequest struct {
	PackageIds []int
	StartDate  string
	EndDate    string
}

// CalculatePackagesUserRechargeRetention 计算某些包用户充值留存
func (s *UserService) CalculatePackagesUserRechargeRetention(request FixUserRechargeRetentionRequest) error {
	days := config.System.RetentionDays
	packageIds := request.PackageIds
	timeList, err := utils.GetNextDaysZeroTimestamps(request.StartDate, request.EndDate)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	errChan := make(chan error, len(packageIds))
	for _, packageId := range packageIds {
		for _, dayTime := range timeList {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := s.CalculateUserRechargeRetention(dayTime, days, packageId)
				if err != nil {
					logs.Error("计算平台用户充值留存错误，平台ID %d, 错误内容 %v", packageId, err)
					errChan <- err
				}
			}()
		}
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// CalculateUserRechargeRetention 计算用户充值留存
func (s *UserService) CalculateUserRechargeRetention(startTime int64, days []int, packageId int) error {
	now := time.Now().Unix()
	loginLogService := &LoginLogService{}
	RechargeStartTime, RechargeEndTime, ranges := utils.BuildTimePoints(startTime, days)
	yesterday := time.Unix(RechargeStartTime, 0).AddDate(0, 0, 0).Format("2006-01-02")
	var yesterdayRegisterCount int64
	userRetentionDailyModel := models.CreateUserRetentionDailyModel()
	qs := userRetentionDailyModel.QueryTable(new(models.UserRetentionDaily)).
		Filter("type", config.RetentionTypeRecharge)
	if packageId > 0 {
		qs = qs.Filter("package_id", packageId)
	}
	exist := qs.Filter("statistic_date", yesterday).Exist()
	if exist {
		// 如果已创建就跳过
		return nil
	}
	userModel := models.CreateUserModel()
	yesterdayRegisterQs := userModel.QueryTable(new(models.User))
	if packageId > 0 {
		yesterdayRegisterQs = yesterdayRegisterQs.Filter("package_id", packageId)
	}
	yesterdayRegisterCount, _ = yesterdayRegisterQs.Filter("register_at__gte", RechargeStartTime).
		Filter("register_at__lt", RechargeEndTime).
		Count()
	userRetentionData := models.UserRetentionDaily{
		StatisticDate: yesterday,
		RegisterCount: int(yesterdayRegisterCount),
		PackageId:     packageId,
		CreatedAt:     now,
		Type:          int(config.RetentionTypeRecharge),
	}
	userRetention := 0.00
	rechargeOrderModel := models.CreateRechargeOrderModel()
	var userRegisterRetentionCount int
	for _, r := range ranges {
		// 昨日充值 Day天注册的用户数量
		rechargeQuery := rechargeOrderModel.QueryTable(new(models.RechargeOrder))
		if packageId > 0 {
			rechargeQuery = rechargeQuery.Filter("package_id", packageId)
		}
		rechargeTotal, err := rechargeQuery.Filter("register_at__gte", r.RegisterStart).
			Filter("register_at__lt", r.RegisterEnd).
			Filter("pay_at__gte", RechargeStartTime).
			Filter("pay_at__lt", RechargeEndTime).
			GroupBy("user_id").
			Count()
		if err != nil {
			logs.Error("统计用户充值留存时获取充值用户失败: %v", err)
			continue
		}
		var userRegisterRetentionTotal []UserRegisterRetentionResponse
		registerParams := LoginLogGroupListRequestParams{
			RegisterBeginTime: r.RegisterStart,
			RegisterEndTime:   r.RegisterEnd,
			BeginTime:         r.RegisterStart,
			EndTime:           r.RegisterEnd,
			GroupBy:           "user_id",
			Fields:            "count(user_id) as total_count,user_id",
			List:              &userRegisterRetentionTotal,
		}

		if packageId > 0 {
			registerParams.PackageId = packageId
		}
		err = loginLogService.GroupList(registerParams)
		if err != nil {
			logs.Error("统计用户充值留存时获取注册用户失败: %v", err)
			continue
		}

		if len(userRegisterRetentionTotal) == 0 {
			userRegisterRetentionCount = 0
		} else {
			userRegisterRetentionCount = userRegisterRetentionTotal[0].TotalCount
		}
		userRetention = s.getUserRetention(int(rechargeTotal), userRegisterRetentionCount)
		s.setRetentionValue(&userRetentionData, r.Day, userRetention)
	}
	_, err := userRetentionDailyModel.Insert(&userRetentionData)
	if err != nil {
		return err
	}
	return nil
}

// 获取用户留存
func (s *UserService) getUserRetention(child int, mother int) float64 {
	if mother == 0 {
		return 0
	} else {
		return math.Round(float64(child)/float64(mother)*10000) / 100
	}
}

func (s *UserService) setRetentionValue(data *models.UserRetentionDaily, day int, value float64) {
	switch day {
	case 1:
		data.Day1 = value
	case 2:
		data.Day2 = value
	case 3:
		data.Day3 = value
	case 4:
		data.Day4 = value
	case 5:
		data.Day5 = value
	case 6:
		data.Day6 = value
	case 7:
		data.Day7 = value
	case 15:
		data.Day15 = value
	case 30:
		data.Day30 = value
	}
}

type UserRetentionRequest struct {
	PackageId int
	Page      int
	PageSize  int
}

// UserRetentionList 获取用户留存统计
func (s *UserService) UserRetentionList(request UserRetentionRequest, retentionType int, needReload int, adminId int64) (map[string]interface{}, error) {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	retentionModel := models.CreateUserRetentionDailyModel()
	var list []models.UserRetentionDaily
	qs := retentionModel.QueryTable(new(models.UserRetentionDaily))
	if request.PackageId == -1 {
		if !isRootAdmin {
			qs = qs.Filter("package_id__in", packageSlice)
		}
	} else {
		if !isRootAdmin {
			inArray := utils.InArray(request.PackageId, packageSlice)
			if !inArray {
				return nil, fmt.Errorf("QuanXianBuZu")
			} else {
				qs = qs.Filter("package_id", request.PackageId)
			}
		} else {
			qs = qs.Filter("package_id", request.PackageId)
		}
	}

	_, err := qs.Filter("type", retentionType).
		Limit(request.PageSize, (request.Page-1)*request.PageSize).
		OrderBy("-id").
		All(&list)
	if err != nil {
		return nil, err
	}
	total, err := qs.Count()
	if nil != err {
		return nil, err
	}
	var packages []models.Package
	if needReload == 1 {
		packageService := PackageService{}
		packages = packageService.GetMyAllPackageList(int(adminId))
	}

	return map[string]interface{}{
		"list":     list,
		"total":    total,
		"packages": packages,
	}, nil
}

type GetUserDailyMenStatResponse struct {
	TotalMen int
}

// GetUserDailyMenStat 日、月报表统计用户总人数
func (s *UserService) GetUserDailyMenStat(condition map[string]interface{}) (int, error) {
	model := models.CreateUserModel()
	condition["is_deleted"] = 0
	totalCount := make([]GetUserDailyMenStatResponse, 0)
	_, err := model.Where(model.QueryTable(new(models.User)), condition).
		Aggregate("count(id) as total_men").All(&totalCount)
	if err != nil {
		return 0, err
	}
	return totalCount[0].TotalMen, nil
}
