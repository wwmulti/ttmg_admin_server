package services

import (
	"api/config"
	"api/models"
	"fmt"
	"net/url"
	"slices"
)

type BetRateService struct {
	BaseService
}

type BetRateListRequest struct {
	Name       string     `form:"name" op:"like"`
	Type       int        `form:"type"`
	ProRate    float64    `form:"proRate"`
	DevRate    float64    `form:"devRate"`
	Lang       string     `form:"lang" op:"-"`
	Page       int        `form:"page"`
	PageSize   int        `form:"pageSize"`
	PackageId  int        `form:"package_id"`
	PackageIds []int      `form:"package_ids" op:"-"`
	RawQuery   url.Values `form:"-"`
}

// BetRateList 获取打码倍数列表
func (s BetRateService) BetRateList(request BetRateListRequest, userId int) (map[string]interface{}, error) {
	betRateModel := models.CreateBetRateModel()
	condition, sort := betRateModel.BuildCondition(request, "-type") // 默认按类型排序
	condition = s.LimitPackageId(condition, request.PackageIds)

	data, total, err := betRateModel.GetPageList(&models.BetRate{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}

	list := data.([]models.BetRate)
	return map[string]interface{}{
		"packages":  (&PackageService{}).GetMyAllPackageList(userId),
		"type_list": s.GetAllowBetRateTypes(request.Lang),
		"list":      list,
		"total":     total,
	}, nil
}

// GetAllowBetRateTypes 获取允许的打码倍数类型
func (s BetRateService) GetAllowBetRateTypes(lang string) map[int]string {
	allUserBalanceLogType := (LogService{}).AllUserBalanceLogType(lang)
	typeList := make(map[int]string)
	for _, logType := range allUserBalanceLogType {
		if slices.Contains(config.IgnoreBalanceLogTypes, logType["id"].(config.BalanceLogType)) {
			continue
		}
		typeList[int(logType["id"].(config.BalanceLogType))] = logType["name"].(string)
	}
	return typeList
}

type AddBetRateRequest struct {
	Name      string  `json:"name"`       // 名称
	Type      int     `json:"type"`       // 打码类型
	ProRate   float64 `json:"pro_rate"`   // 线上打码倍数
	DevRate   float64 `json:"dev_rate"`   // 测试打码倍数
	PackageId int     `json:"package_id"` // 包id
	Lang      string  `json:"lang"`
}

// AddBetRate 添加打码倍数
func (s BetRateService) AddBetRate(request AddBetRateRequest) error {
	if err := s.validateAddBetRate(request); err != nil {
		return err
	}

	betRateModel := models.CreateBetRateModel()
	_, err := betRateModel.Insert(&models.BetRate{
		Name:      request.Name,
		Type:      request.Type,
		ProRate:   request.ProRate,
		DevRate:   request.DevRate,
		PackageId: request.PackageId,
	})
	return err
}

func (s BetRateService) validateAddBetRate(request AddBetRateRequest) error {
	if request.Name == "" {
		return fmt.Errorf("MingChengBiTian")
	}
	if request.Type == 0 {
		return fmt.Errorf("LeiXingBiTian")
	}
	if request.PackageId == 0 {
		return fmt.Errorf("PingTaiBiTian")
	}

	if _, ok := s.GetAllowBetRateTypes(request.Lang)[request.Type]; !ok {
		return fmt.Errorf("LeiXingBuCunZai")
	}

	// 校验 Type 唯一性
	isExist := models.CreateBetRateModel().QueryTable(new(models.BetRate)).
		Filter("type", request.Type).
		Filter("package_id", request.PackageId).
		Exist()
	if isExist {
		return fmt.Errorf("LeiXingYiCunZai")
	}
	return nil
}

type EditBetRateRequest struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	ProRate float64 `json:"pro_rate"`
	DevRate float64 `json:"dev_rate"`
}

// EditBetRate 编辑打码倍数
func (s BetRateService) EditBetRate(request EditBetRateRequest) error {
	if err := s.validateEditBetRate(request); err != nil {
		return err
	}

	betRateModel := models.CreateBetRateModel()
	fields := []string{"name", "pro_rate", "dev_rate"}
	_, err := betRateModel.Update(&models.BetRate{
		Id:      request.Id,
		Name:    request.Name,
		ProRate: request.ProRate,
		DevRate: request.DevRate,
	}, fields...)
	return err
}

func (s BetRateService) validateEditBetRate(request EditBetRateRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("IDBiTian")
	}
	if request.Name == "" {
		return fmt.Errorf("MingChengBiTian")
	}
	return nil
}

// GetBetRate 获取打码倍数
func (s BetRateService) GetBetRate(packageId int, betRateType int) (models.BetRate, error) {
	betRateModel := models.CreateBetRateModel()
	var betRate models.BetRate
	err := betRateModel.QueryTable(&models.BetRate{}).Filter("type", betRateType).Filter("package_id", packageId).One(&betRate)
	if err != nil {
		return betRate, err
	}
	return betRate, nil
}
