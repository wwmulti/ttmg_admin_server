package services

import (
	"api/models"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

type PlatformService struct {
	BaseService
}

// GetAllPlatforms 获取所有游戏平台
func (s PlatformService) GetAllPlatforms(packageIds []int) []models.Platform {
	var list []models.Platform
	models.CreatePlatformModel().QueryTable(&models.Platform{}).
		Filter("is_deleted", 0).
		Filter("package_id__in", packageIds).
		GroupBy("name").
		OrderBy("id").
		All(&list)
	return list
}

type PlatformListRequest struct {
	Id         int        `form:"id"`
	Name       string     `form:"name" op:"like"`
	Alias      string     `form:"alias" op:"like"`
	Status     int        `form:"status"`
	Page       int        `form:"page"`
	PageSize   int        `form:"page_size"`
	NeedReload int        `form:"need_reload" op:"-"`
	RawQuery   url.Values `form:"-"`
}

// PlatformList 平台列表
func (s PlatformService) PlatformList(request PlatformListRequest, userId int) (map[string]interface{}, error) {
	platformModel := models.CreatePlatformModel()
	condition, sort := platformModel.BuildCondition(request, "-id")

	condition["is_deleted"] = 0
	data, total, err := platformModel.GetPageList(&models.Platform{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}

	list := data.([]models.Platform)
	for i := range list {
		list[i].Logo = s.GetCoverUrl(list[i], 1)
	}

	return map[string]interface{}{
		"list":  list,
		"total": total,
	}, nil
}

// AllPlatformByPackageIds 根据分包id获取游戏平台
func (s PlatformService) AllPlatformByPackageIds(packageIds []int) ([]models.Platform, error) {
	list := make([]models.Platform, 0)
	_, err := models.CreatePlatformModel().QueryTable(&models.Platform{}).
		Filter("is_deleted", 0).
		Filter("package_id__in", packageIds).
		GroupBy("name").
		All(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// AllPlatform 所有游戏平台
func (s PlatformService) AllPlatform() ([]models.Platform, error) {
	var list []models.Platform
	_, err := models.CreatePlatformModel().QueryTable(&models.Platform{}).
		Filter("is_deleted", 0).
		All(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

type AddPlatformRequest struct {
	Name          string  `json:"name"`            // 供应商名称
	Alias         string  `json:"alias"`           // 供应商代号
	Logo          string  `json:"logo"`            // 默认图标
	ClickLogo     string  `json:"click_logo"`      // 点击后图标
	Image         string  `json:"image"`           // 供应商图片
	FrontColor    string  `json:"front_color"`     // 前端显示颜色
	MiniMoney     float64 `json:"mini_money"`      // 最小金额
	Status        int     `json:"status"`          // 状态 1正常 0禁用
	Sort          int     `json:"sort"`            // 排序
	GameShowCount int     `json:"game_show_count"` // 显示游戏数量
	GameShowMore  int     `json:"game_show_more"`  // 显示更多游戏数量
	ApiRate       float64 `json:"api_rate"`        // API费率（单位%）
	PackageId     int     `json:"package_id"`      // 分包id
}

// AddPlatform 添加平台
func (s PlatformService) AddPlatform(request AddPlatformRequest) error {
	err := s.validateAddPlatform(request)
	if err != nil {
		return err
	}

	platformModel := models.CreatePlatformModel()
	_, err = platformModel.Insert(&models.Platform{
		Name:          request.Name,
		Alias:         request.Alias,
		Logo:          request.Logo,
		ClickLogo:     request.ClickLogo,
		Image:         request.Image,
		FrontColor:    request.FrontColor,
		MiniMoney:     request.MiniMoney,
		Status:        request.Status,
		Sort:          request.Sort,
		GameShowCount: request.GameShowCount,
		GameShowMore:  request.GameShowMore,
		ApiRate:       request.ApiRate,
		PackageId:     request.PackageId,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加平台参数验证
func (s PlatformService) validateAddPlatform(request AddPlatformRequest) error {
	if request.Name == "" {
		return fmt.Errorf("PingTaiMingChengBiTian")
	}
	isExist := models.CreatePlatformModel().QueryTable(new(models.Platform)).
		Filter("name", request.Name).
		Filter("package_id", request.PackageId).
		Filter("is_deleted", 0).
		Exist()
	if isExist {
		return fmt.Errorf("PingTaiMingChengYiCunZai")
	}
	/* if request.Alias == "" {
		return fmt.Errorf("PingTaiDaiHaoBuNengWeiKong")
	}
	isExist = models.CreatePlatformModel().QueryTable(new(models.Platform)).
		Filter("alias", request.Alias).
		Filter("is_deleted", 0).
		Exist()
	if isExist {
		return fmt.Errorf("PingTaiDaiHaoYiCunZai")
	}
	if request.Logo == "" {
		return fmt.Errorf("PingTaiLogoBiTian")
	}
	if request.GameShowCount == 0 {
		return fmt.Errorf("ZhanShiYouXiShuLiangBiTian")
	}
	if request.GameShowMore == 0 {
		return fmt.Errorf("ZhanShiGengDuoYouXiShuLiangBiTian")
	}*/
	return nil
}

type EditPlatformRequest struct {
	Id            int     `json:"id"`              // 供应商ID
	Name          string  `json:"name"`            // 供应商名称
	Alias         string  `json:"alias"`           // 供应商代号
	Logo          string  `json:"logo"`            // 默认图标
	ClickLogo     string  `json:"click_logo"`      // 点击后图标
	Image         string  `json:"image"`           // 供应商图片
	FrontColor    string  `json:"front_color"`     // 前端显示颜色
	MiniMoney     float64 `json:"mini_money"`      // 最小金额
	Status        int     `json:"status"`          // 状态 1正常 0禁用
	Sort          int     `json:"sort"`            // 排序
	GameShowCount int     `json:"game_show_count"` // 显示游戏数量
	GameShowMore  int     `json:"game_show_more"`  // 显示更多游戏数量
	ApiRate       float64 `json:"api_rate"`        // API费率（单位%）
}

// EditPlatform 编辑平台
func (s PlatformService) EditPlatform(request EditPlatformRequest) error {
	err := s.validateEditPlatform(request)
	if err != nil {
		return err
	}

	platformModel := models.CreatePlatformModel()
	fields := []string{
		"name", "alias", "logo",
		"click_logo", "image", "front_color",
		"mini_money", "status", "sort", "game_show_count",
		"game_show_more", "api_rate",
	}
	_, err = platformModel.Update(&models.Platform{
		Id:            request.Id,
		Name:          request.Name,
		Alias:         request.Alias,
		Logo:          request.Logo,
		ClickLogo:     request.ClickLogo,
		Image:         request.Image,
		FrontColor:    request.FrontColor,
		MiniMoney:     request.MiniMoney,
		Status:        request.Status,
		Sort:          request.Sort,
		GameShowCount: request.GameShowCount,
		GameShowMore:  request.GameShowMore,
		ApiRate:       request.ApiRate,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑平台参数验证
func (s PlatformService) validateEditPlatform(request EditPlatformRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("PingTaiIDBiTian")
	}
	if request.Name == "" {
		return fmt.Errorf("PingTaiMingChengBiTian")
	}
	isExist := models.CreatePlatformModel().QueryTable(new(models.Platform)).
		Filter("id__ne", request.Id).
		Filter("name", request.Name).
		Filter("is_deleted", 0).
		Exist()
	if isExist {
		return fmt.Errorf("PingTaiMingChengYiCunZai")
	}
	/* if request.Alias == "" {
		return fmt.Errorf("PingTaiDaiHaoBuNengWeiKong")
	}
	isExist = models.CreatePlatformModel().QueryTable(new(models.Platform)).
		Filter("id__ne", request.Id).
		Filter("alias", request.Alias).
		Filter("is_deleted", 0).
		Exist()
	if isExist {
		return fmt.Errorf("PingTaiDaiHaoYiCunZai")
	}
	if request.Logo == "" {
		return fmt.Errorf("PingTaiLogoBiTian")
	}
	if request.GameShowCount == 0 {
		return fmt.Errorf("ZhanShiYouXiShuLiangBiTian")
	}
	if request.GameShowMore == 0 {
		return fmt.Errorf("ZhanShiGengDuoYouXiShuLiangBiTian")
	}*/
	return nil
}

// DeletePlatform 删除平台
func (s PlatformService) DeletePlatform(id int) error {
	platformModel := models.CreatePlatformModel()
	if id == 0 {
		return fmt.Errorf("PingTaiIDBiTian")
	}
	_, err := platformModel.Update(&models.Platform{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	return nil
}

// EditGameAttr 编辑游戏属性
func (s PlatformService) EditPlatformAttr(id int, field string) error {
	platformModel := models.CreatePlatformModel()
	var platform models.Platform
	err := platformModel.QueryTable(&models.Platform{}).Filter("id", id).One(&platform)
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return fmt.Errorf("PingTaiBuCunZai")
		}
		return err
	}

	switch field {
	case "status":
		if platform.Status == 1 {
			platform.Status = 0
		} else {
			platform.Status = 1
		}
	default:
		return fmt.Errorf("GengXinNeiRongBuCunZai")
	}

	_, err = platformModel.Update(&platform, field)
	if err != nil {
		return err
	}
	return nil
}

// GetCoverUrl 获取平台封面
func (s PlatformService) GetCoverUrl(platform models.Platform, packageId int) string {
	if len(platform.Logo) != 0 {
		return platform.Logo
	} else {
		return (&ConfigService{}).GetGamePictureDomain(packageId) + "/uploads_002/images/icons/" + strings.ToUpper(platform.Name) + ".png"
	}
}

// 初始化游戏平台表
func (s PlatformService) InitPlatform(packageId int) {
	/* gameTypeIds := fmt.Sprintf(",%v,", (GameTypeService{}).GetId(packageId, "Slot"))
	for _, platform := range config.GamePlatforms {
		var info models.Platform
		err := models.CreatePlatformModel().QueryTable(&models.Platform{}).Filter("name", platform).Filter("package_id", packageId).One(&info)
		if err == nil {
			continue
		} else if !errors.Is(err, orm.ErrNoRows) {
			logs.Error("初始化分包%v游戏平台表查询失败:%v", packageId, err.Error())
		}
		switch platform {
		case "pg":
			gameTypeIds = fmt.Sprintf(",%v,", (GameTypeService{}).GetId(packageId, "Slot"))
		case "jdb":
			gameTypeIds = fmt.Sprintf(",%v,%v,%v,%v,", (GameTypeService{}).GetId(packageId, "Slot"), (GameTypeService{}).GetId(packageId, "Pescaria"), (GameTypeService{}).GetId(packageId, "Cartas"), (GameTypeService{}).GetId(packageId, "BlockChain"))
		case "kess":
			gameTypeIds = fmt.Sprintf(",%v,", (GameTypeService{}).GetId(packageId, "Slot"))
		case "zy":
			gameTypeIds = fmt.Sprintf(",%v,", (GameTypeService{}).GetId(packageId, "Slot"))
		case "wg":
			gameTypeIds = fmt.Sprintf(",%v,%v,%v,%v,", (GameTypeService{}).GetId(packageId, "Slot"), (GameTypeService{}).GetId(packageId, "Pescaria"), (GameTypeService{}).GetId(packageId, "Cartas"), (GameTypeService{}).GetId(packageId, "BlockChain"))
		case "cp":
			gameTypeIds = fmt.Sprintf(",%v,%v,", (GameTypeService{}).GetId(packageId, "Slot"), (GameTypeService{}).GetId(packageId, "Cartas"))
		}
		_, err = models.CreatePlatformModel().Insert(&models.Platform{
			Name:          platform,
			Status:        1,
			Sort:          0,
			GameShowCount: 8,
			GameShowMore:  12,
			PackageId:     packageId,
		})
		if err != nil {
			logs.Error("初始化分包%v游戏平台表写入失败:%v", packageId, err.Error())
		}
	} */
}

// 信息Id
func (s PlatformService) GetId(packageId int, name string) int {
	var info models.Platform
	model := models.CreatePlatformModel()

	model.QueryTable(new(models.Platform)).Filter("name", name).Filter("package_id", packageId).One(&info)
	return info.Id
}

// 创建一个新的数据
func (s PlatformService) CreateData(packageId int, name string, packageIds []int) int {
	var info models.Platform
	model := models.CreatePlatformModel()
	model.QueryTable(new(models.Platform)).Filter("package_id__in", packageIds).Filter("name", name).Filter("is_deleted", 0).One(&info)

	data := info
	data.Id = 0
	data.PackageId = packageId

	id, _ := model.Insert(&data)
	return int(id)
}

// 获取游戏平台ids
func (s PlatformService) GetIds(packageIds []int, name string) []int {
	if len(name) == 0 {
		return nil
	}
	model := models.CreatePlatformModel()
	var list []models.Platform
	model.QueryTable(new(models.Platform)).Filter("package_id__in", packageIds).Filter("name", name).All(&list, "id")

	if len(list) == 0 {
		return nil
	}
	ids := make([]int, 0)
	for _, item := range list {
		ids = append(ids, item.Id)
	}
	return ids
}
