package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type AccountService struct{}

type LoginRequestParams struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	Code     string `json:"code"`     // 验证码
}

type UserListRequestParams struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Username string `json:"username,omitempty"`
	Status   int    `json:"status,,omitempty"`
}

// 列表返回结构
type AccountListResponse struct {
	Lists []models.Account `json:"lists"`
	Total int              `json:"total"`
}

// 获取权限组id
func (s *AccountService) GetGroupId(uid int) int {
	model := models.CreateAccountModel()
	var info models.Account
	err := model.QueryTable(new(models.Account)).Filter("id", uid).One(&info)
	if err != nil {
		return 0
	}
	return info.Group.Id
}

// 管理员请求参数
type AccountRequestParams struct {
	Id         int    `json:"id,omitempty"`
	Username   string `json:"username"`
	Pwd        string `json:"pwd"`
	ConfirmPwd string `json:"conf_pwd"`
	GroupId    int    `json:"group_id"`
}

// 解绑google
func (s *AccountService) Unbindgoogle(accountId, uid int) error {
	model := models.CreateAccountModel()
	var info models.Account
	err := model.QueryTable(new(models.Account)).Filter("id", accountId).Filter("parents__contains", fmt.Sprintf(",%d,", uid)).Filter("is_deleted", 0).One(&info)
	if err != nil {
		return fmt.Errorf("ShuJuYiChang")
	}

	fields := []string{
		"secret", "bind_time",
	}
	_, updateErr := model.Update(&models.Account{
		Id:       accountId,
		Secret:   "",
		BindTime: 0,
	}, fields...)
	if updateErr != nil {
		logs.Error("解绑管理员失败:%v", updateErr)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

// 修改账号
func (s *AccountService) EditAccount(params AccountRequestParams, uid int) error {
	model := models.CreateAccountModel()
	var info models.Account
	err := model.QueryTable(new(models.Account)).Filter("id", params.Id).Filter("parents__contains", fmt.Sprintf(",%d,", uid)).Filter("is_deleted", 0).One(&info)
	if err != nil {
		return fmt.Errorf("ShuJuYiChang")
	}

	fields := []string{}
	if len(params.Username) > 0 && params.Username != info.AccountName {
		isExists := model.QueryTable(new(models.Account)).Filter("account_name", params.Username).Filter("is_deleted", 0).Exist()
		if isExists {
			return fmt.Errorf("YongHuMingYiCunZai")
		}
		fields = append(fields, "account_name")
	}

	password := ""
	salt := ""
	if len(params.Pwd) > 0 {
		if params.Pwd != params.ConfirmPwd {
			return fmt.Errorf("liangCiMiMaBuXiangTong")
		} else {
			salt = utils.RandString(10)
			password = s.generatePwd(params.Pwd, salt)
			fields = append(fields, "password")
			fields = append(fields, "salt")
		}
	}
	// 只能修改其他用户的群组
	groupInfo := models.AuthGroup{}
	if params.GroupId != info.Group.Id && params.Id != uid {
		groupInfo = s.GetGroupInfo(uid, params.GroupId)
		if groupInfo.Id == 0 {
			return fmt.Errorf("ShuJuYiChang")
		}
		fields = append(fields, "group_id")
	}

	_, updateErr := model.Update(&models.Account{
		Id:          params.Id,
		Salt:        salt,
		Password:    password,
		Group:       &groupInfo,
		AccountName: params.Username,
	}, fields...)
	if updateErr != nil {
		logs.Error("更新管理员失败:%v", updateErr)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

// 删除账号
func (s *AccountService) DelAccount(accountId, uid int) error {
	if accountId == uid {
		return fmt.Errorf("BuNengShanChuZiJi")
	}
	// 是否是自己下级
	model := models.CreateAccountModel()
	isExists := model.QueryTable(new(models.Account)).Filter("parents__contains", fmt.Sprintf(",%d,", uid)).Filter("is_deleted", 0).Exist()
	if !isExists {
		return fmt.Errorf("ShuJuYiChang")
	}
	fields := []string{
		"is_deleted",
	}
	_, err := model.Update(&models.Account{
		Id:        accountId,
		IsDeleted: 1,
	}, fields...)
	if err != nil {
		logs.Error("删除管理员失败:%v", err)
		return fmt.Errorf("ShuJuShanChuShiBai")
	}
	return nil
}

// 添加账号
func (s *AccountService) AddAccount(params AccountRequestParams, accountId int) error {
	if params.ConfirmPwd != params.Pwd {
		return fmt.Errorf("liangCiMiMaBuXiangTong")
	}

	model := models.CreateAccountModel()
	// 用户名验证
	isExists := model.QueryTable(new(models.Account)).Filter("account_name", params.Username).Filter("is_deleted", 0).Exist()
	if isExists {
		return fmt.Errorf("YongHuMingYiCunZai")
	}

	var accountInfo models.Account
	err := model.QueryTable(new(models.Account)).Filter("id", accountId).One(&accountInfo)
	if err != nil {
		return fmt.Errorf("ShuJuGengXinShiBai")
	}

	if accountInfo.Group.Id != 1 && accountInfo.Group.Id == params.GroupId {
		return fmt.Errorf("JiaoSeBuNengYuDangQianYongHuXiangTong")
	}

	if len(params.Pwd) == 0 {
		return fmt.Errorf("MiMaBuNengWeiKong")
	}

	salt := utils.RandString(10)
	password := s.generatePwd(params.Pwd, salt)
	// 获取自己的下级群组
	groupInfo := s.GetGroupInfo(accountId, params.GroupId)

	trErr := model.OrmerMaster.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		uid, err := txOrm.Insert(&models.Account{
			AccountName:  params.Username,
			Salt:         salt,
			Password:     password,
			RegisterTime: int(time.Now().Unix()),
			Status:       1,
			IsDeleted:    0,
			Pid:          accountId,
			Group:        &groupInfo,
		})
		if err != nil {
			return err
		}

		parents, perr := s.CreateParents(txOrm, accountId, int(uid))
		if perr != nil {
			return perr
		}

		fields := []string{
			"parents",
		}
		_, updateErr := txOrm.Update(&models.Account{
			Id:      int(uid),
			Parents: parents,
		}, fields...)
		return updateErr
	})
	if trErr != nil {
		logs.Error("创建管理员失败:%v", trErr)
		return fmt.Errorf("ChuangJianShuJuShiBai")
	}
	return nil
}

// 获取指定群组
func (s *AccountService) GetGroupInfo(accountId, groupId int) models.AuthGroup {
	groupList := (&AuthGroupService{}).GetGroupList(accountId)
	for _, group := range groupList {
		if group.Id == groupId {
			return group
		}
	}
	return models.AuthGroup{}
}

func (s *AccountService) CreateParents(txOrm orm.TxOrmer, pid, id int) (string, error) {
	var info models.Account
	err := txOrm.QueryTable(new(models.Account)).Filter("id", pid).One(&info)
	if err != nil {
		logs.Error("获取角色父级pid:%v信息失败:%v", pid, err)
		return "", err
	}
	return fmt.Sprintf("%s%d,", info.Parents, id), nil
}

// 修改用户状态
func (s *AccountService) EditAccountStatus(accountId, uid int) error {
	if accountId == uid {
		return fmt.Errorf("WuFaCaoZuoDangQianDengLuYongHu")
	}
	model := models.CreateAccountModel()
	var info models.Account
	err := model.QueryTable(new(models.Account)).Filter("id", accountId).Filter("is_deleted", 0).Filter("parents__contains", fmt.Sprintf(",%d,", uid)).One(&info)
	if err != nil {
		return fmt.Errorf("ShuJuYiChang")
	}

	status := 1
	if info.Status > 0 {
		status = 0
	}
	fields := []string{
		"status",
	}
	_, updateErr := model.Update(&models.Account{
		Id:     accountId,
		Status: status,
	}, fields...)
	if updateErr != nil {
		logs.Error("账号状态更新失败：%v", updateErr)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

// 获取账号列表
func (s *AccountService) GetUserList(params UserListRequestParams, uid int) AccountListResponse {
	condition := make(map[string]interface{})
	condition["parents__contains"] = fmt.Sprintf(",%d,", uid)
	condition["is_deleted"] = 0

	if params.Status > -1 {
		condition["status"] = params.Status
	}
	if len(params.Username) > 0 {
		condition["account_name__contains"] = params.Username
	}

	result, total, _ := models.CreateAccountModel().GetPageList(new(models.Account), condition, params.Page, params.PageSize, "-id")
	accounts, ok := result.([]models.Account)
	if !ok {
		accounts = []models.Account{}
	}
	return AccountListResponse{
		Total: int(total),
		Lists: accounts,
	}
}

// 用户登陆
func (s *AccountService) Login(params LoginRequestParams) (models.Account, error) {
	accountName := strings.TrimSpace(params.Username)
	code := strings.TrimSpace(params.Code)

	model := models.CreateAccountModel()

	var result models.Account
	err := model.QueryTable(new(models.Account)).Filter("account_name", accountName).Filter("is_deleted", 0).One(&result)
	if err != nil {
		return models.Account{}, fmt.Errorf("MiMaCuoWu")
	}

	// 验证密码
	if !s.VerifyPassword(params.Password, result.Salt, result.Password) {
		return models.Account{}, fmt.Errorf("MiMaCuoWu")
	}

	if result.BindTime > 0 && !s.CheckAccountCode(code, result.Secret) {
		return models.Account{}, fmt.Errorf("YanZhengMaCuoWu")
	}

	token, _ := utils.GenerateToken(int64(result.Id), result.AccountName)
	result.Token = token
	return result, nil
}

// 生成密码
func (s *AccountService) generatePwd(password, salt string) string {
	// 拼接字符串
	combined := fmt.Sprintf("%s%s", password, salt)
	hash := md5.Sum([]byte(combined))
	md5Str := hex.EncodeToString(hash[:])
	return md5Str
}

func (s *AccountService) CheckAccountCode(code, secret string) bool {
	status, _ := (&utils.GoogleAuth{}).VerifyCode(secret, code)
	return status
}

type GooleSecretResponse struct {
	Secret string `json:"secret"`
	Url    string `json:"code"`
}

func (s *AccountService) CreateGoogleSecrect(uid int, username string) (error, GooleSecretResponse) {
	secret := s.getSecretFromCache(uid)
	if len(secret) == 0 {
		secret = (&utils.GoogleAuth{}).GetSecret()
	}

	qrcode := (&utils.GoogleAuth{}).GetQrcode(username, secret)
	key := s.getSecretCacheKey(uid)
	RedisCache := utils.GetRedisClient()
	ctx := context.Background()
	_, err := RedisCache.Set(ctx, key, secret, 24*time.Hour).Result()
	if err != nil {
		return fmt.Errorf("MiYueChuangJianShiBai"), GooleSecretResponse{}
	}
	return nil, GooleSecretResponse{
		Secret: secret,
		Url:    qrcode,
	}
}

// 绑定google 秘钥
func (s *AccountService) Bindgoogle(uid int, code string) error {
	secret := s.getSecretFromCache(uid)
	if len(secret) == 0 {
		return fmt.Errorf("MiYueYiGuoQiQingChongXinShengCheng")
	}
	status, err := (&utils.GoogleAuth{}).VerifyCode(secret, code)
	if err != nil {
		logs.Error("绑定验证码失败：%v", err)
		return fmt.Errorf("YanZhengMaBangDingShiBai")
	}

	if !status {
		return fmt.Errorf("YanZhengMaBangDingShiBai")
	}

	model := models.CreateAccountModel()
	fields := []string{
		"secret", "bind_time",
	}
	_, updateErr := model.Update(&models.Account{
		Id:       uid,
		Secret:   secret,
		BindTime: int(time.Now().Unix()),
	}, fields...)
	if updateErr != nil {
		logs.Error("保存秘钥失败:%v", updateErr)
		return fmt.Errorf("YanZhengMaBangDingShiBai")
	}
	s.delSecretFromCache(uid)
	return nil
}

// 获取缓存秘钥
func (s *AccountService) getSecretFromCache(uid int) string {
	RedisCache := utils.GetRedisClient()
	key := s.getSecretCacheKey(uid)
	ctx := context.Background()
	secret, _ := RedisCache.Get(ctx, key).Result()
	return secret
}

// 获取秘钥缓存key
func (s *AccountService) getSecretCacheKey(uid int) string {
	return fmt.Sprintf("%s_%v", config.RedisKeyName.UserGoogleSecret, uid)
}

// 删除秘钥缓存
func (s *AccountService) delSecretFromCache(uid int) {
	RedisCache := utils.GetRedisClient()
	key := s.getSecretCacheKey(uid)
	ctx := context.Background()
	RedisCache.Del(ctx, key).Result()
}

// 验证密码
func (s *AccountService) VerifyPassword(password, salt, expectedPassword string) bool {
	return s.generatePwd(password, salt) == expectedPassword
}

// 验证用户密码
func (s *AccountService) VerifyUserPassword(uid int, password string) bool {
	var result models.Account
	err := models.CreateAccountModel().QueryTable(new(models.Account)).Filter("id", uid).Filter("is_deleted", 0).One(&result)
	if err != nil {
		return false
	}
	return s.VerifyPassword(password, result.Salt, result.Password)
}
