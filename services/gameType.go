package services

import (
	"api/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type GameTypeService struct {
	BaseService
}

// GetAllGameTypes 获取所有游戏类型
func (s GameTypeService) GetAllGameTypes(packageIds []int) []models.GameType {
	var list []models.GameType
	models.CreateGameTypeModel().QueryTable(&models.GameType{}).
		Filter("is_deleted", 0).
		Filter("package_id__in", packageIds).
		GroupBy("name").
		OrderBy("id").
		All(&list)
	return list
}

type GameTypeListRequest struct {
	Id       int        `form:"id"`
	Name     string     `form:"name" op:"like"`
	Alias    string     `form:"alias" op:"like"`
	Status   int        `form:"status"`
	Page     int        `form:"page"`
	PageSize int        `form:"page_size"`
	RawQuery url.Values `form:"-"`
}

// GameTypeList 游戏类型列表
func (s GameTypeService) GameTypeList(request GameTypeListRequest, userId int) (map[string]interface{}, error) {
	gameTypeModel := models.CreateGameTypeModel()
	condition, sort := gameTypeModel.BuildCondition(request, "-id") // 默认按sort排序

	condition["is_deleted"] = 0
	data, total, err := gameTypeModel.GetPageList(&models.GameType{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}

	platforms, err := (&PlatformService{}).AllPlatform()
	if nil != err {
		return nil, err
	}

	list := data.([]models.GameType)
	for i, gameType := range list {
		list[i].Logo = s.GetCoverUrl(gameType, 1)
	}

	return map[string]interface{}{
		"list":      list,
		"total":     total,
		"platforms": platforms,
	}, nil
}

// GetCoverUrl 获取封面
func (s GameTypeService) GetCoverUrl(game models.GameType, packageId int) string {
	if len(game.Logo) != 0 {
		return game.Logo
	} else {
		return (&ConfigService{}).GetGamePictureDomain(packageId) + "/uploads_002/images/icons/" + game.Name + ".png"
	}
}

// AllGameType 所有游戏类型
func (s GameTypeService) AllGameType(packageIds []int) ([]models.GameType, error) {
	gameTypes := make([]models.GameType, 0)
	gameTypeModel := models.CreateGameTypeModel()
	_, err := gameTypeModel.QueryTable(&models.GameType{}).
		Filter("is_deleted", 0).
		Filter("package_id__in", packageIds).
		GroupBy("name").
		OrderBy("id").
		All(&gameTypes)
	if err != nil {
		return nil, err
	}
	return gameTypes, nil
}

type AddGameTypeRequest struct {
	Name           string `json:"name"`             // 分类名称
	Alias          string `json:"alias"`            // 简称
	Logo           string `json:"logo"`             // 分类图标
	Sort           int    `json:"sort"`             // 排序权重
	Status         int    `json:"status"`           // 分类状态 0关闭 1开启
	PageGameNumber int    `json:"page_game_number"` // 显示单页游戏数量
	BetRate        int    `json:"bet_rate"`         // 游戏打码比例
	PackageId      int    `json:"package_id"`       // 分包id
	PlatformIds    string `json:"platform_ids"`     // 平台id列表
}

// AddGameType 添加游戏类型
func (s GameTypeService) AddGameType(request AddGameTypeRequest) error {
	err := s.validateAddGameType(request)
	if err != nil {
		return err
	}
	gameTypeModel := models.CreateGameTypeModel()
	_, err = gameTypeModel.Insert(&models.GameType{
		Name:           request.Name,
		Alias:          request.Alias,
		Logo:           request.Logo,
		Sort:           request.Sort,
		Status:         request.Status,
		PageGameNumber: request.PageGameNumber,
		BetRate:        request.BetRate,
		PackageId:      request.PackageId,
		PlatformIds:    request.PlatformIds,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加游戏类型参数验证
func (s GameTypeService) validateAddGameType(request AddGameTypeRequest) error {
	if request.Name == "" {
		return fmt.Errorf("YouXiLeiXingMingChengBiTian")
	}
	isExist := models.CreateGameTypeModel().QueryTable(new(models.GameType)).
		Filter("name", request.Name).
		Filter("package_id", request.PackageId).
		Filter("is_deleted", 0).
		Exist()
	if isExist {
		return fmt.Errorf("YouXiLeiXingMingChengYiCunZai")
	}
	return nil
}

type EditGameTypeRequest struct {
	Id             int    `json:"id"`               // 主键ID
	Name           string `json:"name"`             // 分类名称
	Alias          string `json:"alias"`            // 简称
	Logo           string `json:"logo"`             // 分类图标
	Sort           int    `json:"sort"`             // 排序权重
	Status         int    `json:"status"`           // 分类状态 0关闭 1开启
	PageGameNumber int    `json:"page_game_number"` // 显示单页游戏数量
	BetRate        int    `json:"bet_rate"`         // 游戏打码比例
	PlatformIds    string `json:"platform_ids"`     // 平台id列表
}

// EditGameType 编辑游戏类型
func (s GameTypeService) EditGameType(request EditGameTypeRequest) error {
	err := s.validateEditGameType(request)
	if err != nil {
		return err
	}
	gameTypeModel := models.CreateGameTypeModel()
	fields := []string{
		"name", "alias", "logo", "sort", "status", "page_game_number", "bet_rate", "platform_ids",
	}
	_, err = gameTypeModel.Update(&models.GameType{
		Id:             request.Id,
		Name:           request.Name,
		Alias:          request.Alias,
		Logo:           request.Logo,
		Sort:           request.Sort,
		Status:         request.Status,
		PageGameNumber: request.PageGameNumber,
		BetRate:        request.BetRate,
		PlatformIds:    request.PlatformIds,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑游戏类型参数验证
func (s GameTypeService) validateEditGameType(request EditGameTypeRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("YouXiLeiXingIDBiTian")
	}
	if request.Name == "" {
		return fmt.Errorf("YouXiLeiXingMingChengBiTian")
	}
	isExist := models.CreateGameTypeModel().QueryTable(new(models.GameType)).
		Filter("id__ne", request.Id).
		Filter("name", request.Name).
		Filter("is_deleted", 0).
		Exist()
	if isExist {
		return fmt.Errorf("YouXiLeiXingMingChengYiCunZai")
	}
	return nil
}

// DeleteGameType 删除游戏类型
func (s GameTypeService) DeleteGameType(id int) error {
	gameTypeModel := models.CreateGameTypeModel()
	if id == 0 {
		return fmt.Errorf("YouXiLeiXingIDBiTian")
	}
	_, err := gameTypeModel.Update(&models.GameType{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	return nil
}

// EditGameTypeAttr 编辑分类属性（开关状态等）
func (s GameTypeService) EditGameTypeAttr(id int, field string) error {
	gameTypeModel := models.CreateGameTypeModel()
	var info models.GameType
	err := gameTypeModel.QueryTable(&models.GameType{}).Filter("id", id).One(&info)
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return fmt.Errorf("YouXiFenLeiBuCunZai")
		}
		return err
	}

	switch field {
	case "status":
		if info.Status == 1 {
			info.Status = 0
		} else {
			info.Status = 1
		}
	default:
		return fmt.Errorf("GengXinNeiRongBuCunZai")
	}

	_, err = gameTypeModel.Update(&info, field)
	if err != nil {
		return err
	}
	return nil
}

type AgentConfigRequest struct {
	Level      int     `json:"level"`
	Members    int     `json:"members"`
	Commission float64 `json:"commission"`
	Rate       float64 `json:"rate"`
}

// AgentConfig 代理配置
func (s GameTypeService) AgentConfig(id int, request []AgentConfigRequest) error {
	if id == 0 {
		return fmt.Errorf("YouXiLeiXingIDBiTian")
	}

	temp := -1
	for _, v := range request {
		if v.Level <= 0 {
			return fmt.Errorf("DaiLiDengJiBuNengWeiLing")
		}

		if v.Members < 0 || v.Commission < 0 || v.Rate < 0 {
			return fmt.Errorf("SuoTianNeiRongBuNengXiaoYuLing")
		}

		if temp == v.Level {
			return fmt.Errorf("DengJiChongFu")
		}
		temp = v.Level
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("ShuJuGeShiCuoWu")
	}

	gameTypeModel := models.CreateGameTypeModel()
	_, err = gameTypeModel.Update(&models.GameType{
		Id:          id,
		AgentConfig: string(data),
	}, "agent_config")
	if err != nil {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}

	return nil
}

// BetRateConfig 游戏打码比例配置
func (s GameTypeService) BetRateConfig(id int, betRate string) error {
	if id == 0 {
		return fmt.Errorf("YouXiLeiXingIDBiTian")
	}
	gameTypeModel := models.CreateGameTypeModel()
	_, err := gameTypeModel.Update(&models.GameType{
		Id:           id,
		GamesBetRate: betRate,
	}, "games_bet_rate")
	if err != nil {
		return fmt.Errorf("WeiZhiDeCuoWu")
	}

	return nil
}

// 信息Id
func (s GameTypeService) GetId(packageId int, name string) int {
	var info models.GameType
	model := models.CreateGameTypeModel()

	model.QueryTable(new(models.GameType)).Filter("name", name).Filter("package_id", packageId).One(&info)
	return info.Id
}

// 创建一个新的数据
func (s GameTypeService) CreateData(packageId int, name string, packageIds []int) int {
	var info models.GameType
	model := models.CreateGameTypeModel()
	model.QueryTable(new(models.GameType)).Filter("package_id__in", packageIds).Filter("name", name).Filter("is_deleted", 0).One(&info)

	data := info
	data.Id = 0
	data.PackageId = packageId

	id, _ := model.Insert(&data)
	return int(id)
}

// 初始化游戏类型
func (s GameTypeService) InitGameTypeList(packageId int) {
	gameTypeSeed := map[int]string{
		1: "Slot",
		2: "Pescaria", // 捕鱼
		3: "Cartas",   // 卡牌
		4: "Ao Vivo",  // 真人
		5: "BlockChain",
		6: "Autoral", // 自研
	}
	gameTypeModel := models.CreateGameTypeModel()
	logo := ""
	for i := 1; i <= len(gameTypeSeed); i++ {
		name, ok := gameTypeSeed[i]
		if !ok {
			continue
		}
		exist := gameTypeModel.QueryTable(new(models.GameType)).
			Filter("name", name).Filter("package_id", packageId).
			Exist()
		if exist {
			continue
		}
		switch i {
		case 1:
			logo = "https://uploads.wwapi.vip/uploads/category/10026.png"
		case 2:
			logo = "https://uploads.wwapi.vip/uploads/category/10027.png"
		case 3:
			logo = "https://uploads.wwapi.vip/uploads/category/10028.png"
		case 4:
			logo = "https://uploads.wwapi.vip/uploads/category/10029.png"
		case 5:
			logo = "https://uploads.wwapi.vip/uploads/category/10030.png"
		case 6:
			logo = "https://uploads.wwapi.vip/uploads/category/10039.png"
		}
		gameType := &models.GameType{
			Name:           name,
			Logo:           logo,
			Sort:           0,
			Status:         1,
			IsDeleted:      0,
			PageGameNumber: 20,
			PackageId:      packageId,
		}
		if _, err := gameTypeModel.Insert(gameType); err != nil {
			logs.Error("初始化分包%v游戏类型失败：%v", packageId, err)
		}
	}
}

// 获取游戏类型id
func (s GameTypeService) GetIds(packageIds []int, name string) []int {
	if len(name) == 0 {
		return nil
	}
	model := models.CreateGameTypeModel()
	var list []models.GameType
	model.QueryTable(new(models.GameType)).Filter("package_id__in", packageIds).Filter("name", name).All(&list, "id")

	if len(list) == 0 {
		return nil
	}
	ids := make([]int, 0)
	for _, item := range list {
		ids = append(ids, item.Id)
	}
	return ids
}
