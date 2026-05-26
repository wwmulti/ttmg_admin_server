package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/redis/go-redis/v9"
)

type ConfigService struct{}

// GetValue 获取配置值
type GetConfigValueParams struct {
	Key          string
	DefaultValue string
	TypeId       int
	PackageId    int
}

// GetValue 获取配置值
func (s *ConfigService) GetValue(params GetConfigValueParams) (string, error) {
	model := models.CreateConfigModel()

	var row models.Config
	err := model.QueryTable(new(models.Config)).Filter("ConfigKey", params.Key).Filter("PackageId", params.PackageId).One(&row, "config_value")

	if errors.Is(err, orm.ErrNoRows) {
		setting := &models.Config{
			ConfigKey:   params.Key,
			ConfigValue: params.DefaultValue,
			TypeId:      params.TypeId,
			PackageId:   params.PackageId,
		}
		if _, err := model.Insert(setting); err != nil {
			return "", err
		}
		return params.DefaultValue, nil
	} else if err != nil {
		return "", err
	}
	return row.ConfigValue, nil
}

// GetLuckyWheelClearTime 获取转盘活动的积分清除周期
func (s *ConfigService) GetLuckyWheelClearTime(packageId int) int {
	params := GetConfigValueParams{
		Key:          config.System.LuckyWheelClearTimeType.Key,
		DefaultValue: config.System.LuckyWheelClearTimeType.Value,
		TypeId:       int(config.SystemConfigTypeLuckWheel),
		PackageId:    packageId,
	}
	timeType, _ := s.GetValue(params)
	num, err := strconv.Atoi(timeType)
	if err != nil {
		logs.Error("幸运转盘转盘活动的积分清除周期==>%v, 转换错误: %v，使用默认值0", num, err)
		return 0.0
	}
	return num
}

// GetLuckyWheelScore 获取转盘消耗
func (s *ConfigService) GetLuckyWheelScore(gameType int64, packageId int) (float64, error) {
	params := GetConfigValueParams{
		Key:          config.System.LuckyWheelCost.Key,
		DefaultValue: config.System.LuckyWheelCost.Value,
		TypeId:       int(config.SystemConfigTypeLuckWheel),
		PackageId:    packageId,
	}
	cost, _ := s.GetValue(params)

	costArray := strings.Split(cost, ",")
	if len(costArray) != int(config.LuckyWheelTypeMax) {
		return 0.0, fmt.Errorf("ZhuanPanXiaoHaoPeiZhiCuoWu")
	}

	r, err := strconv.ParseFloat(costArray[gameType], 64)
	if err != nil {
		return 0.0, err
	}
	return r, nil
}

// GetSystemConfig 获取系统配置
func (s *ConfigService) GetSystemConfig(RedisCache *redis.Client, key string, defaultValue string, packageId int, typeId int, tx ...orm.TxOrmer) (string, error) {
	ctx := context.Background()
	packageKey := "_package_" + strconv.Itoa(packageId)
	systemConfigKey := config.RedisKeyName.SystemConfig + packageKey
	val, err := RedisCache.HGet(ctx, systemConfigKey, key).Result()
	if err == nil {
		return val, nil
	}
	if !errors.Is(err, redis.Nil) {
		return "", err
	}
	configModel := models.CreateConfigModel()
	setting := &models.Config{}
	if len(tx) > 0 {
		err = tx[0].QueryTable(new(models.Config)).
			Filter("package_id", packageId).
			Filter("config_key", key).
			One(setting)
	} else {
		err = configModel.QueryTable(new(models.Config)).
			Filter("package_id", packageId).
			Filter("config_key", key).
			One(setting)
	}

	if err == nil {
		RedisCache.HSet(ctx, systemConfigKey, key, setting.ConfigValue)
		return setting.ConfigValue, nil
	}
	if !errors.Is(err, orm.ErrNoRows) {
		return "", err
	}
	setting.ConfigKey = key
	setting.ConfigValue = defaultValue
	setting.TypeId = typeId
	setting.PackageId = packageId
	if len(tx) > 0 {
		if _, err := tx[0].Insert(setting); err != nil {
			return "", err
		}
	} else {
		if _, err := configModel.Insert(setting); err != nil {
			return "", err
		}
	}
	RedisCache.HSet(ctx, systemConfigKey, key, defaultValue)
	return defaultValue, nil
}

// RefreshSystemConfig 刷新系统配置
func (s *ConfigService) RefreshSystemConfig(RedisCache *redis.Client, key string, defaultValue string, packageId int, typeId int) error {
	ctx := context.Background()
	packageKey := "_package_" + strconv.Itoa(packageId)
	systemConfigKey := config.RedisKeyName.SystemConfig + packageKey
	configModel := models.CreateConfigModel()
	setting := &models.Config{}
	err := configModel.QueryTable(new(models.Config)).
		Filter("config_key", key).
		Filter("package_id", packageId).
		One(setting)
	if err == nil {
		RedisCache.HSet(ctx, systemConfigKey, key, setting.ConfigValue)
		return nil
	} else {
		if errors.Is(err, orm.ErrNoRows) {
			setting.ConfigKey = key
			setting.ConfigValue = defaultValue
			setting.TypeId = typeId
			if _, err := configModel.Insert(setting); err != nil {
				return err
			}
			RedisCache.HSet(ctx, systemConfigKey, key, defaultValue)
			return nil
		} else {
			return err
		}
	}
}

type AllSystemConfigRequest struct {
	TypeId    int `json:"type_id" example:"1"`
	PackageId int `json:"package_id" example:"1"`
}

// AllSystemConfig 获取所有系统配置
func (s *ConfigService) AllSystemConfig(request AllSystemConfigRequest, needReload int, adminId int64) (map[string]interface{}, error) {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	configModel := models.CreateConfigModel()
	var data []models.Config
	qs := configModel.QueryTable(new(models.Config))
	if !isRootAdmin {
		if len(packageSlice) == 0 {
			return nil, fmt.Errorf("QuanXianBuZu")
		}
		qs = qs.Filter("package_id__in", packageSlice)
	}
	_, err := qs.Filter("type_id", request.TypeId).
		OrderBy("package_id").
		All(&data)
	if err != nil {
		return nil, err
	} else {
		var packages []models.Package
		result := make(map[string][]models.Config)
		for _, cfg := range data {
			result[cfg.ConfigKey] = append(result[cfg.ConfigKey], cfg)
		}
		if needReload == 1 {
			packageService := PackageService{}
			packages = packageService.GetMyAllPackageList(int(adminId))
		}
		return map[string]interface{}{
			"configs":  result,
			"packages": packages,
		}, nil
	}
}

type EditConfigRequest struct {
	ID          *int   `json:"id"`
	TypeId      int    `json:"type_id"`
	PackageId   int    `json:"package_id"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
}

// EditSystemConfig 编辑系统配置
func (s *ConfigService) EditSystemConfig(request []EditConfigRequest, adminId int64) error {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	configModel := models.CreateConfigModel()
	tx, err := configModel.Begin()
	if err != nil {
		return err
	}
	for _, req := range request {
		if !isRootAdmin && !utils.InArray(req.PackageId, packageSlice) {
			return fmt.Errorf("QuanXianBuZu")
		}
		if req.ID != nil {
			configRow := &models.Config{
				Id:          *req.ID,
				PackageId:   req.PackageId,
				TypeId:      req.TypeId,
				ConfigKey:   req.ConfigKey,
				ConfigValue: req.ConfigValue,
			}
			_, err := tx.Update(configRow)
			if err != nil {
				return err
			}
		} else {
			configRow := &models.Config{
				PackageId:   req.PackageId,
				TypeId:      req.TypeId,
				ConfigKey:   req.ConfigKey,
				ConfigValue: req.ConfigValue,
			}
			_, err := tx.Insert(configRow)
			if err != nil {
				return err
			}
		}
		s.clearRedisConfigCache(req.ConfigKey, req.ConfigValue, req.PackageId)
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		logs.Error("后台提交配置事务出错", err.Error())
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	return nil
}

// 清理 Redis 缓存
func (s *ConfigService) clearRedisConfigCache(key string, value string, packageId int) {
	RedisCache := utils.GetRedisClient()
	ctx := context.Background()
	packageKey := "_package_" + strconv.Itoa(packageId)
	systemConfigKey := config.RedisKeyName.SystemConfig + packageKey
	RedisCache.HSet(ctx, systemConfigKey, key, value)
}

// InitPackage 初始化包所有数据
func (s *ConfigService) InitPackage(packageId int) error {
	err := s.InitPackageConfig(packageId)
	if err != nil {
		return err
	}
	return nil
}

// InitPackageConfig 初始化包的系统配置
func (s *ConfigService) InitPackageConfig(packageId int) error {
	RedisCache := utils.GetRedisClient()
	// 系统配置
	systemConfigs := [][]string{
		{config.System.GamePictureDomain.Key, config.System.GamePictureDomain.Value},
		{config.System.UserRegisterReward.Key, config.System.UserRegisterReward.Value},
		{config.System.MaxRegisterNumber.Key, config.System.MaxRegisterNumber.Value},
		{config.System.EnterpriseRegisterCode.Key, config.System.EnterpriseRegisterCode.Value},
		{config.System.DefaultPlatformGameShowNumber.Key, config.System.DefaultPlatformGameShowNumber.Value},
		{config.System.CustomerLink.Key, config.System.CustomerLink.Value},
	}
	for _, item := range systemConfigs {
		key := item[0]
		value := item[1]
		_, err := s.GetSystemConfig(RedisCache, key, value, packageId, int(config.SystemConfigTypeSystem))
		if err != nil {
			logs.Error("初始化系统配置失败: %v", err)
			return err
		}
	}
	// 奖池配置
	prizeConfigs := [][]string{
		{config.System.InitialValue.Key, config.System.InitialValue.Value},
		{config.System.GrowthInterval.Key, config.System.GrowthInterval.Value},
		{config.System.GrowthMin.Key, config.System.GrowthMin.Value},
		{config.System.GrowthMax.Key, config.System.GrowthMax.Value},
		{config.System.ReduceTimeMin.Key, config.System.ReduceTimeMin.Value},
		{config.System.ReduceTimeMax.Key, config.System.ReduceTimeMax.Value},
		{config.System.ReduceRates.Key, config.System.ReduceRates.Value},
	}
	for _, item := range prizeConfigs {
		key := item[0]
		value := item[1]
		_, err := s.GetSystemConfig(RedisCache, key, value, packageId, int(config.SystemConfigTypePrize))
		if err != nil {
			logs.Error("初始化奖池配置失败: %v", err)
			return err
		}
	}
	// 提现配置
	withdrawConfigs := [][]string{
		{config.System.PasswordMaxWrongTime.Key, config.System.PasswordMaxWrongTime.Value},
		{config.System.MaxBindWithdrawAccountNum.Key, config.System.MaxBindWithdrawAccountNum.Value},
		{config.System.AllowUnbindWithdrawAccount.Key, config.System.AllowUnbindWithdrawAccount.Value},
	}
	for _, item := range withdrawConfigs {
		key := item[0]
		value := item[1]
		_, err := s.GetSystemConfig(RedisCache, key, value, packageId, int(config.SystemConfigTypeWithdraw))
		if err != nil {
			logs.Error("初始化提现配置失败: %v", err)
			return err
		}
	}
	// 转盘配置
	luckyWheelConfigs := [][]string{
		{config.System.LuckyWheelRate.Key, config.System.LuckyWheelRate.Value},
		{config.System.LuckyWheelClearTimeType.Key, config.System.LuckyWheelClearTimeType.Value},
		{config.System.LuckyWheelCost.Key, config.System.LuckyWheelCost.Value},
	}
	for _, item := range luckyWheelConfigs {
		key := item[0]
		value := item[1]
		_, err := s.GetSystemConfig(RedisCache, key, value, packageId, int(config.SystemConfigTypeLuckWheel))
		if err != nil {
			logs.Error("初始化转盘配置失败: %v", err)
			return err
		}
	}
	// 弹窗配置
	alertConfigs := [][]string{
		{config.System.AlertDownloadUrlFirst.Key, config.System.AlertDownloadUrlFirst.Value},
		{config.System.AlertDownloadUrlSecond.Key, config.System.AlertDownloadUrlSecond.Value},
	}
	for _, item := range alertConfigs {
		key := item[0]
		value := item[1]
		_, err := s.GetSystemConfig(RedisCache, key, value, packageId, int(config.SystemConfigTypeAlert))
		if err != nil {
			logs.Error("初始化弹窗配置失败: %v", err)
			return err
		}
	}
	// 代理配置
	agentConfigs := [][]string{
		{config.System.AgentScanInterval.Key, config.System.AgentScanInterval.Value},
		{config.System.AgentSettlementCycle.Key, config.System.AgentSettlementCycle.Value},
		{config.System.AgentSettlementTimeZone.Key, config.System.AgentSettlementTimeZone.Value},
		{config.System.AgentShareDomain.Key, config.System.AgentShareDomain.Value},
		{config.System.AgentSocialChannels.Key, config.System.AgentSocialChannels.Value},
		{config.System.AgentMinAmount.Key, config.System.AgentMinAmount.Value},
		{config.System.AgentCommissionBetRate.Key, config.System.AgentCommissionBetRate.Value},
	}
	for _, item := range agentConfigs {
		key := item[0]
		value := item[1]
		_, err := s.GetSystemConfig(RedisCache, key, value, packageId, int(config.SystemConfigTypeAgent))
		if err != nil {
			logs.Error("初始化代理配置失败: %v", err)
			return err
		}
	}
	// 自动出款配置
	autoOutMoneyConfigs := [][]string{
		{config.System.AutoOutMoneyDailySusNums.Key, config.System.AutoOutMoneyDailySusNums.Value},
		{config.System.AutoOutMoneyDividePays.Key, config.System.AutoOutMoneyDividePays.Value},
		{config.System.AutoOutMoneyMaxAmount.Key, config.System.AutoOutMoneyMaxAmount.Value},
		{config.System.AutoOutMoneyMaxIps.Key, config.System.AutoOutMoneyMaxIps.Value},
		{config.System.AutoOutMoneySubPays.Key, config.System.AutoOutMoneySubPays.Value},
		{config.System.AutoOutMoneyChannel.Key, config.System.AutoOutMoneyChannel.Value},
		{config.System.AutoOutMoneySwitch.Key, config.System.AutoOutMoneySwitch.Value},
		{config.System.AutoOutMoneyTotalPays.Key, config.System.AutoOutMoneyTotalPays.Value},
		{config.System.AutoOutMoneyTotalWithdraws.Key, config.System.AutoOutMoneyTotalWithdraws.Value},
		{config.System.AutoOutMoneyFirstAuth.Key, config.System.AutoOutMoneyFirstAuth.Value},
		{config.System.AutoOutMoneyType.Key, config.System.AutoOutMoneyType.Value},
	}
	for _, item := range autoOutMoneyConfigs {
		key := item[0]
		value := item[1]
		_, err := s.GetSystemConfig(RedisCache, key, value, packageId, int(config.SystemConfigTypeAutoOutMoney))
		if err != nil {
			logs.Error("初始化自动出款配置失败: %v", err)
			return err
		}
	}
	return nil
}

// SyncPackageConfig 同步分包配置
func (s *ConfigService) SyncPackageConfig(packageId, copyPackageId int) error {
	var lists []models.Config
	model := models.CreateConfigModel()
	_, err := model.QueryTable(new(models.Config)).Filter("PackageId", copyPackageId).All(&lists)
	if err != nil {
		return err
	}
	var configList []models.Config
	for _, item := range lists {
		item.PackageId = packageId
		v := models.Config{
			PackageId:   packageId,
			ConfigKey:   item.ConfigKey,
			ConfigValue: item.ConfigValue,
			TypeId:      item.TypeId,
		}
		configList = append(configList, v)
	}
	_, insertErr := model.InsertMulti(len(configList), configList)
	return insertErr
}

// GetGamePictureDomain 获取图片域名地址
func (s *ConfigService) GetGamePictureDomain(packageId int) string {
	params := GetConfigValueParams{
		Key:          config.System.GamePictureDomain.Key,
		DefaultValue: config.System.GamePictureDomain.Value,
		TypeId:       int(config.SystemConfigTypeSystem),
		PackageId:    packageId,
	}
	domain, _ := s.GetValue(params)
	return domain
}

// DeletePackageConfig 删除指定包的配置
func (s *ConfigService) DeletePackageConfig(packageId int, tx ...orm.TxOrmer) error {
	if len(tx) > 0 {
		_, err := tx[0].QueryTable(new(models.Config)).
			Filter("package_id", packageId).
			Delete()
		if err != nil {
			return err
		}
	} else {
		configModel := models.CreateConfigModel()
		_, err := configModel.QueryTable(new(models.Config)).
			Filter("package_id", packageId).
			Delete()
		if err != nil {
			return err
		}
	}
	return nil
}
