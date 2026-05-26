package services

import (
	"api/config"
	"api/controllers/gameApi"
	"api/models"
	"api/utils"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/xuri/excelize/v2"
)

type GameService struct {
	BaseService
}

type GameListRequest struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	Recommend    int    `json:"recommend"`
	Maintain     int    `json:"maintain"`
	GameTypeId   int    `json:"game_type_id"`
	GameTypeName string `json:"game_type_name"`
	PlatformId   int    `json:"platform_id"`
	PlatformName string `json:"platform_name"`
	Status       int    `json:"status"`
	Tag          string `json:"tag"`
	CodeRule     int    `json:"code_rule"`
	NeedReload   int    `json:"need_reload"`
	Event        string `json:"event"`
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
	PackageId    int    `json:"package_id"`
	PackageIds   []int
}

func (s GameService) BuildCondition(params GameListRequest) map[string]interface{} {
	condition := map[string]interface{}{
		"is_deleted": 0,
	}

	if len(params.Name) > 0 {
		condition["name__icontains"] = params.Name
	}

	if len(params.Code) > 0 {
		condition["code__icontains"] = params.Code
	}

	if params.Recommend >= 0 {
		condition["recommend"] = params.Recommend
	}

	if params.Maintain >= 0 {
		condition["maintain"] = params.Maintain
	}

	if params.GameTypeId > 0 {
		condition["game_type_id"] = params.GameTypeId
	}
	if len(params.GameTypeName) > 0 {
		condition["game_type_name"] = params.GameTypeName
	}

	if params.PlatformId > 0 {
		condition["platform_id"] = params.PlatformId
	}
	if len(params.PlatformName) > 0 {
		condition["platform_name"] = params.PlatformName
	}

	if params.Status >= 0 {
		condition["status"] = params.Status
	}

	if params.CodeRule > 0 {
		condition["code_rule"] = params.CodeRule
	}

	if params.PackageId > 0 {
		condition["package_id"] = params.PackageId
	}

	condition = s.LimitPackageId(condition, params.PackageIds)

	return condition
}

// GameList 游戏列表
func (s GameService) GameList(params GameListRequest, userId int) (map[string]interface{}, error) {
	gameModel := models.CreateGameModel()

	condition := s.BuildCondition(params)

	if params.Event == "asyncexport" {
		exportConfig := ExportConfig{
			FileName: "game",
			Model:    models.Game{},
			Header: []HeaderItem{
				{Field: "id", Title: "ID"},
				{Field: "name", Title: "游戏名"},
				{Field: "code_rule", Title: "打码规则"},
			},
			DataQuery: func(res interface{}) error {
				_, err := gameModel.Where(gameModel.QueryTable(new(models.Game)), condition).All(res)
				return err
			},
			DataFormat: func(item interface{}) map[string]interface{} {
				data := item.(models.Game)

				codeRule := ""
				switch data.CodeRule {
				case 1:
					codeRule = "流水"
				case 2:
					codeRule = "净赢"
				case 3:
					codeRule = "赢金打码"
				}

				return map[string]interface{}{
					"id":        data.Id,
					"name":      data.Name,
					"code_rule": codeRule,
				}
			},
		}
		taskId, _ := (&TaskService{}).AddTask(TaskTypeExport, map[string]interface{}{"config": exportConfig})
		return map[string]interface{}{"taskId": taskId}, nil
	}

	data, total, err := gameModel.GetPageList(&models.Game{}, condition, params.Page, params.PageSize, "-id")

	if nil != err {
		return nil, err
	}

	// 下拉选项
	var gameTypes, platforms interface{}
	params.PackageIds = []int{1}
	if params.NeedReload == 1 {
		gameTypes = (GameTypeService{}).GetAllGameTypes(params.PackageIds)
		platforms = (PlatformService{}).GetAllPlatforms(params.PackageIds)
	}

	list := data.([]models.Game)
	for i := range list {
		list[i].Cover = s.GetCoverUrl(list[i], list[i].PackageId)
	}

	return map[string]interface{}{
		"list":       list,
		"total":      total,
		"game_types": gameTypes,
		"platforms":  platforms,
	}, nil
}

// AddGame 添加游戏
func (s GameService) AddGame(request EditGameRequest) error {
	err := s.validateAddGame(request)
	if err != nil {
		return err
	}

	platformId := (PlatformService{}).GetId(request.PackageId, request.PlatformName)
	gameTypeId := (GameTypeService{}).GetId(request.PackageId, request.GameTypeName)

	/* if platformId == 0 {
		platformId = (PlatformService{}).CreateData(request.PackageId, request.PlatformName, request.PackageIds)
	}
	if gameTypeId == 0 {
		gameTypeId = (GameTypeService{}).CreateData(request.PackageId, request.GameTypeName, request.PackageIds)
	} */

	gameModel := models.CreateGameModel()
	_, err = gameModel.Insert(&models.Game{
		Name:         request.Name,
		PlatformId:   platformId,
		PlatformName: request.PlatformName,
		Cover:        request.Cover,
		GameTypeId:   gameTypeId,
		GameTypeName: request.GameTypeName,
		Status:       request.Status,
		Sort:         request.Sort,
		Tag:          request.Tag,
		Code:         request.Code,
		PtName:       request.PtName,
		Swpg:         request.Swpg,
		Ttpg:         request.Ttpg,
		Mlpg:         request.Mlpg,
		Maintain:     request.Maintain,
		CodeRule:     request.CodeRule,
		PackageId:    request.PackageId,
	})
	if err != nil {
		return err
	}
	return nil
}

// 添加游戏参数验证
func (s GameService) validateAddGame(request EditGameRequest) error {
	if request.Name == "" {
		return fmt.Errorf("YouXiMingChengBiTian")
	}
	/* isExist := models.CreateGameModel().QueryTable(new(models.Game)).
		Filter("name", request.Name).
		Filter("package_id", request.PackageId).
		Exist()
	if isExist {
		return fmt.Errorf("YouXiMingChengYiCunZai")
	}
	if request.Cover == "" {
		return fmt.Errorf("YouXiFengMianBiTian")
	} */
	if len(request.GameTypeName) == 0 {
		return fmt.Errorf("YouXiFenLeiBiTian")
	}
	if len(request.PlatformName) == 0 {
		return fmt.Errorf("YouXiPingTaiBiTian")
	}
	if request.Code == "" {
		return fmt.Errorf("YouXiDaiMaBiTian")
	}
	isExist := models.CreateGameModel().QueryTable(new(models.Game)).
		Filter("code", request.Code).
		Filter("package_id", request.PackageId).
		Exist()
	if isExist {
		return fmt.Errorf("YouXiDaiMaYiCunZai")
	}

	return nil
}

type EditGameRequest struct {
	Id           int    `json:"id"`             // 游戏ID
	Name         string `json:"name"`           // 游戏名
	PlatformName string `json:"platform_name"`  // 平台名称
	Cover        string `json:"cover"`          // 游戏封面
	GameTypeName string `json:"game_type_name"` // 游戏类型名称
	Status       int    `json:"status"`         // 状态 1正常 0禁用
	Sort         int    `json:"sort"`           // 排序
	Tag          string `json:"tag"`            // 标签
	Code         string `json:"code"`           // 游戏code
	PtName       string `json:"pt_name"`        // 葡语名称
	Swpg         int    `json:"swpg"`           // 0-禁用 1-启用
	Ttpg         int    `json:"ttpg"`           // 0-禁用 1-启用
	Mlpg         int    `json:"mlpg"`           // 0-禁用 1-启用
	Maintain     int    `json:"maintain"`       // 维护 1-维护中 0-没有维护
	CodeRule     int    `json:"code_rule"`      // 打码规则 1-流水 2-净赢 3-赢金打码
	PackageId    int    `json:"package_id"`
	PackageIds   []int  // 自身授权的分包id
}

// EditGame 编辑游戏
func (s GameService) EditGame(request EditGameRequest) error {
	err := s.validateEditGame(request)
	if err != nil {
		return err
	}
	gameModel := models.CreateGameModel()
	fields := []string{
		"name", "platform_id", "platform_name", "cover", "game_type_id", "game_type_name", "status", "sort", "tag", "code", "pt_name", "swpg", "ttpg", "mlpg", "maintain", "code_rule",
	}
	platformId := (PlatformService{}).GetId(request.PackageId, request.PlatformName)
	gameTypeId := (GameTypeService{}).GetId(request.PackageId, request.GameTypeName)

	/* if platformId == 0 {
		platformId = (PlatformService{}).CreateData(request.PackageId, request.PlatformName, request.PackageIds)
	}
	if gameTypeId == 0 {
		gameTypeId = (GameTypeService{}).CreateData(request.PackageId, request.GameTypeName, request.PackageIds)
	} */
	_, err = gameModel.Update(&models.Game{
		Id:           request.Id,
		Name:         request.Name,
		PlatformId:   platformId,
		PlatformName: request.PlatformName,
		Cover:        request.Cover,
		GameTypeId:   gameTypeId,
		GameTypeName: request.GameTypeName,
		Status:       request.Status,
		Sort:         request.Sort,
		Tag:          request.Tag,
		Code:         request.Code,
		PtName:       request.PtName,
		Swpg:         request.Swpg,
		Ttpg:         request.Ttpg,
		Mlpg:         request.Mlpg,
		Maintain:     request.Maintain,
		CodeRule:     request.CodeRule,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑游戏参数验证
func (s GameService) validateEditGame(request EditGameRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("YouXiIDBiTian")
	}
	if request.Name == "" {
		return fmt.Errorf("YouXiMingChengBiTian")
	}
	/* isExist := models.CreateGameModel().QueryTable(new(models.Game)).
		Filter("id__ne", request.Id).
		Filter("package_id", request.PackageId).
		Filter("name", request.Name).
		Exist()
	if isExist {
		return fmt.Errorf("YouXiMingChengYiCunZai")
	}
	if request.Cover == "" {
		return fmt.Errorf("YouXiFengMianBiTian")
	} */
	if len(request.GameTypeName) == 0 {
		return fmt.Errorf("YouXiFenLeiBiTian")
	}
	if len(request.PlatformName) == 0 {
		return fmt.Errorf("YouXiPingTaiBiTian")
	}
	if request.Code == "" {
		return fmt.Errorf("YouXiDaiMaBiTian")
	}
	isExist := models.CreateGameModel().QueryTable(new(models.Game)).
		Filter("id__ne", request.Id).
		Filter("package_id", request.PackageId).
		Filter("code", request.Code).
		Exist()
	if isExist {
		return fmt.Errorf("YouXiDaiMaYiCunZai")
	}
	return nil
}

// DeleteGame 删除游戏
func (s GameService) DeleteGame(id int) error {
	gameModel := models.CreateGameModel()
	if id == 0 {
		return fmt.Errorf("YouXiIDBiTian")
	}
	_, err := gameModel.Update(&models.Game{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	return nil
}

// EditGameAttr 编辑游戏属性
func (s GameService) EditGameAttr(id int, field string, value int) error {
	gameModel := models.CreateGameModel()
	var game models.Game
	err := gameModel.QueryTable(&models.Game{}).Filter("id", id).One(&game)
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return fmt.Errorf("YouXiBuCunZai")
		}
		return err
	}

	switch field {
	case "status":
		game.Status = value
	case "recommend":
		game.Recommend = value
	case "maintain":
		game.Maintain = value
	default:
		return fmt.Errorf("GengXinNeiRongBuCunZai")
	}

	_, err = gameModel.Update(&game, field)
	if err != nil {
		return err
	}
	return nil
}

// GetCoverUrl 获取游戏封面
func (s GameService) GetCoverUrl(game models.Game, packageId int) string {
	if len(game.Cover) != 0 {
		return game.Cover
	} else {
		return (&ConfigService{}).GetGamePictureDomain(packageId) + "/uploads_002/images/" + strings.ToLower(game.PlatformName) + "/" + game.Code + ".png"
	}
}

// 同步游戏相关的配置
func (s GameService) SyncPackageConfig(packageId, copyPackageId int) error {
	model := models.CreateGameModel()

	(GameTypeService{}).InitGameTypeList(packageId)
	(PlatformService{}).InitPlatform(packageId)

	if copyPackageId == 0 {
		logs.Info("当前分包没有对应的复制对象，不同步游戏数据")
		return nil
	}

	condition := map[string]interface{}{
		"package_id": copyPackageId,
	}
	page := 1
	pageSize := 100

	datas, _, err := model.GetPageList(&models.Game{}, condition, page, pageSize, "id")
	if err != nil {
		return err
	}
	lists := datas.([]models.Game)

	for len(lists) > 0 {
		page++

		gameLists := []models.Game{}
		for _, list := range lists {
			game := list
			game.Id = 0
			game.PlatformId = (PlatformService{}).GetId(packageId, list.PlatformName)
			game.GameTypeId = (GameTypeService{}).GetId(packageId, list.GameTypeName)
			game.PackageId = packageId

			gameLists = append(gameLists, game)
		}
		model.InsertMulti(len(gameLists), gameLists)
		datas, _, _ := model.GetPageList(&models.Game{}, condition, page, pageSize, "id")
		lists = datas.([]models.Game)
	}

	return nil
}

// 初始化游戏列表
func (s GameService) InitZyGameList(packageId int) error {
	for _, platform := range config.GamePlatforms {
		data := (&gameApi.ZyGame{}).GetGameList(platform)
		gameTypeId := 0
		gameTypeName := ""

		switch platform {
		case "pg", "cp": // pg 和 cp 都映射到 1
			gameTypeName = "Slot"
			gameTypeId = (GameTypeService{}).GetId(packageId, gameTypeName)
		case "jdb":
			gameTypeName = "Pescaria"
			gameTypeId = (GameTypeService{}).GetId(packageId, gameTypeName)
		case "kess":
			gameTypeName = "Cartas"
			gameTypeId = (GameTypeService{}).GetId(packageId, gameTypeName)
		case "wg":
			gameTypeName = "Ao Vivo"
			gameTypeId = (GameTypeService{}).GetId(packageId, gameTypeName)
		case "zy":
			gameTypeName = "Autoral"
			gameTypeId = (GameTypeService{}).GetId(packageId, gameTypeName)
		default:
			gameTypeId = 0 // 或者处理未知平台的情况
		}
		//"pg", "jdb", "kess", "zy", "wg", "cp",
		params := SyncGameListRequest{
			GameType:     platform,
			GameTypeId:   gameTypeId,
			GameTypeName: gameTypeName,
			SupplierId:   int(config.ZySupplierType),
			GameData:     data,
			PackageId:    packageId,
		}
		err := s.syncGameList(params)
		if err != nil {
			logs.Error("同步游戏列表失败：%v", err)
			return err
		}
	}

	return nil
}

type SyncGameListRequest struct {
	GameType     string // 游戏类型 pg,jdb,kess,zy,wg,cp等
	GameTypeId   int    // 游戏类型归属分类id
	GameTypeName string
	SupplierId   int                  // 游戏供应商ID
	GameData     []gameApi.ZyGameInfo // 游戏详情列表
	PackageId    int                  // 分包id
}

// SyncGameList 同步游戏列表
func (s GameService) syncGameList(data SyncGameListRequest) error {
	gameList := data.GameData
	if len(gameList) == 0 {
		return fmt.Errorf("游戏列表不能为空")
	}

	// 取平台id
	var platforms []models.Platform
	_, err := models.CreatePlatformModel().QueryTable(&models.Platform{}).Filter("package_id", data.PackageId).All(&platforms, "id", "name")
	if err != nil {
		return fmt.Errorf("查询平台信息失败: %v", err.Error())
	}
	platformIds := make(map[string]int)
	for _, platform := range platforms {
		platformIds[platform.Name] = platform.Id
	}

	for _, game := range gameList {
		var gameInfo models.Game
		err := models.CreateGameModel().QueryTable(&models.Game{}).Filter("package_id", data.PackageId).Filter("platform_name", data.GameType).Filter("supplier_id", data.SupplierId).Filter("name", game.Name).One(&gameInfo)
		if err == nil {
			continue
		} else if !errors.Is(err, orm.ErrNoRows) {
			return fmt.Errorf("查询游戏信息失败: %v 错误: %v", game.Name, err.Error())
		}

		platformId, ok := platformIds[data.GameType]
		if !ok {
			return fmt.Errorf("游戏 %v 平台不存在: %v", game.Name, data.GameType)
		}

		gameData := models.Game{
			PlatformId:   platformId,
			PlatformName: data.GameType,
			Name:         game.Name,
			Code:         game.Code,
			GameTypeId:   data.GameTypeId,
			GameTypeName: data.GameTypeName,
			Status:       1, // 正常状态
			Sort:         0,
			IsDeleted:    0,
			SupplierId:   data.SupplierId,
			CodeRule:     1,
			PackageId:    data.PackageId,
		}
		_, err = models.CreateGameModel().Insert(&gameData)
		if err != nil {
			logs.Error("插入游戏信息失败, 记录: %v", utils.ToJson(gameData))
			return fmt.Errorf("插入游戏信息失败, err: %v", err)
		}
	}
	return nil
}

type ExcelImportStatistics struct {
	TotalRecords     int            `json:"total_records"`
	SuccessCount     int            `json:"success_count"`
	SkippedCount     int            `json:"skipped_count"`
	NoCompletedCount int            `json:"no_completed_count"`
	InvalidCount     int            `json:"invalid_count"`
	Categories       map[string]int `json:"categories"`
	Platforms        map[string]int `json:"platforms"`
}

// ImportGamesFromExcel 从Excel导入游戏数据
func (s GameService) ImportGamesFromExcel(file io.Reader, packageId int, packageIds []int) (*ExcelImportStatistics, error) {
	// 读取Excel文件
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("Excel文件解析失败: %v", err)
	}
	defer f.Close()

	// 获取第一个工作表
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("Excel文件中没有工作表")
	}

	// 读取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取Excel数据失败: %v", err)
	}

	if len(rows) <= 1 {
		return nil, fmt.Errorf("Excel文件中没有数据")
	}

	// 第一行是表头，解析列索引
	header := rows[0]
	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[strings.ToLower(strings.TrimSpace(col))] = i
	}

	// 检查必要列
	requiredCols := []string{"name", "game_code", "type", "category"}
	for _, col := range requiredCols {
		if _, ok := colIndex[col]; !ok {
			return nil, fmt.Errorf("Excel文件缺少必要列: %s", col)
		}
	}

	// 统计信息
	stats := &ExcelImportStatistics{
		TotalRecords: len(rows) - 1,
		Categories:   make(map[string]int),
		Platforms:    make(map[string]int),
	}

	// 预先查询所有平台
	var platforms []models.Platform
	_, err = models.CreatePlatformModel().QueryTable(&models.Platform{}).Filter("package_id", packageId).All(&platforms, "id", "name")
	if err != nil {
		return nil, fmt.Errorf("查询平台信息失败: %v", err)
	}
	platformMap := make(map[string]int)
	for _, p := range platforms {
		platformMap[strings.ToLower(p.Name)] = p.Id
	}

	// 处理数据行
	for i := 1; i < len(rows); i++ {
		row := rows[i]

		// 获取单元格值
		getCellValue := func(colName string) string {
			idx := colIndex[colName]
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		// 获取"开发"列的值（如果存在）
		devStatus := ""
		if devIdx, ok := colIndex["开发"]; ok {
			if devIdx < len(row) {
				devStatus = strings.TrimSpace(row[devIdx])
			}
		}

		// 跳过开发状态为"验收中"的记录
		if devStatus == "验收中" {
			stats.NoCompletedCount++
			continue
		}

		name := getCellValue("name")
		gameCode := getCellValue("game_code")
		gameType := strings.ToLower(getCellValue("type"))
		category := strings.ToLower(getCellValue("category"))

		// 验证必要字段
		if name == "" || gameCode == "" || gameType == "" || category == "" {
			stats.InvalidCount++
			continue
		}

		// 检查游戏是否已存在
		var gameInfo models.Game
		err = models.CreateGameModel().QueryTable(&models.Game{}).
			Filter("package_id", packageId).
			Filter("code", gameCode).
			One(&gameInfo)
		if err == nil {
			// 游戏已存在，跳过
			stats.InvalidCount++
			continue
		} else if !errors.Is(err, orm.ErrNoRows) {
			return nil, fmt.Errorf("查询游戏信息失败: %v", err)
		}

		// 获取平台ID
		platformId, ok := platformMap[gameType]
		if !ok {
			logs.Warn("游戏 %s 的平台 %s 不存在，跳过", name, gameType)
			stats.InvalidCount++
			continue
		}

		// 获取游戏类型ID
		gameTypeId := (GameTypeService{}).GetId(packageId, category)
		if gameTypeId == 0 {
			logs.Warn("游戏 %s 的分类 %s 不存在，跳过", name, category)
			stats.InvalidCount++
			continue
		}

		// 插入游戏数据
		gameData := models.Game{
			Name:         name,
			Code:         gameCode,
			PlatformId:   platformId,
			PlatformName: gameType,
			GameTypeId:   gameTypeId,
			GameTypeName: category,
			Status:       1,
			Sort:         0,
			IsDeleted:    0,
			CodeRule:     1,
			PackageId:    packageId,
		}

		_, err = models.CreateGameModel().Insert(&gameData)
		if err != nil {
			logs.Error("插入游戏信息失败, 记录: %v", utils.ToJson(gameData))
			stats.InvalidCount++
			continue
		}

		stats.SuccessCount++
		stats.Categories[category]++
		stats.Platforms[gameType]++
	}

	stats.SkippedCount = stats.NoCompletedCount + stats.InvalidCount

	logs.Info("Excel导入完成: 总记录 %d, 成功 %d, 跳过 %d (验收中: %d, 无效: %d)",
		stats.TotalRecords, stats.SuccessCount, stats.SkippedCount,
		stats.NoCompletedCount, stats.InvalidCount)

	return stats, nil
}

// ========== 第三方游戏API代理方法 ==========

func (s GameService) getConfig(key string) string {
	value, _ := beego.AppConfig.String(key)
	return value
}

// GetOperatorToken 获取运营商token
func (s GameService) GetOperatorToken() string {
	return s.getConfig("OPERATOR_TOKEN")
}

// GetSecret 获取密钥
func (s GameService) GetSecret() string {
	return s.getConfig("SECRET")
}

// GetApiBaseUrl 获取API基础URL
func (s GameService) GetApiBaseUrl() string {
	return s.getConfig("BASE_API_URL")
}

// GetProxyUrl 获取代理URL
func (s GameService) GetProxyUrl() string {
	enabled := s.getConfig("PROXY_ENABLED")
	if enabled == "true" {
		return s.getConfig("PROXY_URL")
	}
	return ""
}

// GenerateSign 生成MD5签名
func (s GameService) GenerateSign(params map[string]interface{}, secret string) string {
	tempMap := make(map[string]string)
	for k, v := range params {
		if k != "sign" && k != "key" {
			tempMap[k] = fmt.Sprintf("%v", v)
		}
	}

	keys := make([]string, 0, len(tempMap))
	for k := range tempMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, tempMap[k]))
	}

	signStr := fmt.Sprintf("%s&key=%s", strings.Join(parts, "&"), secret)
	hash := md5.Sum([]byte(signStr))
	return hex.EncodeToString(hash[:])
}

// CallWebApi 调用第三方游戏API
func (s GameService) CallWebApi(path string, params map[string]interface{}) (map[string]interface{}, error) {
	secret := s.GetSecret()
	if secret == "" {
		return nil, fmt.Errorf("VITE_SECRET环境变量未配置")
	}

	params["operator_token"] = s.GetOperatorToken()
	params["ts"] = time.Now().Unix()
	params["sign"] = s.GenerateSign(params, secret)

	baseUrl := s.GetApiBaseUrl()
	url := baseUrl + path

	var resp *utils.Response
	if s.GetProxyUrl() != "" {
		resp = utils.PostJSON(url, params, utils.RequestOption{
			Proxy:   s.GetProxyUrl(),
			Timeout: 15 * time.Second,
		})
	} else {
		resp = utils.PostJSON(url, params, utils.RequestOption{
			Timeout: 15 * time.Second,
		})
	}

	if !resp.OK() {
		logs.Error("调用第三方API失败: %s, 错误: %v", url, resp.Error)
		return nil, fmt.Errorf("API请求失败")
	}

	var result map[string]interface{}
	if err := resp.JSON(&result); err != nil {
		logs.Error("解析API响应失败: %v", err)
		return nil, fmt.Errorf("解析响应失败")
	}

	if result["code"].(float64) != 0 {
		return nil, fmt.Errorf("第三方API返回错误: %v", result)
	}

	return result["data"].(map[string]interface{}), nil
}

// CreateUserSession 创建用户会话
func (s GameService) CreateUserSession(userId string, userName string, rtp int) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"user_id":   userId,
		"user_name": userName,
		"rtp":       rtp,
	}
	return s.CallWebApi("/api/web/user_session", params)
}

// GetGameUrl 获取游戏URL
func (s GameService) GetGameUrl(gameCode, language, gameType, userId, userToken string) (map[string]interface{}, error) {
	var game models.Game
	err := models.CreateGameModel().QueryTable(&models.Game{}).
		Filter("code", gameCode).
		Filter("is_deleted", 0).
		Filter("status", 1).
		One(&game)
	if err != nil {
		return nil, fmt.Errorf("查询游戏失败: %v", err)
	}
	if game.Maintain == 1 {
		return nil, fmt.Errorf("游戏维护中")
	}

	params := map[string]interface{}{
		"game_code":  gameCode,
		"language":   language,
		"type":       gameType,
		"user_id":    userId,
		"user_token": userToken,
	}
	return s.CallWebApi("/api/web/game_url", params)
}

// SetUserRtp 设置用户RTP
func (s GameService) SetUserRtp(userId string, rtp int) (map[string]interface{}, error) {
	if userId == "" {
		return nil, fmt.Errorf("缺少必要参数: user_id")
	}
	if rtp > 300 {
		return nil, fmt.Errorf("RTP不能大于300")
	}
	params := map[string]interface{}{
		"user_id": userId,
		"rtp":     rtp,
	}
	return s.CallWebApi("/api/web/set_user_rtp", params)
}

// AllGamesResponse 所有游戏列表响应结构
type AllGamesResponse struct {
	Games           []GameTypeGroup `json:"games"`
	CustomerService string          `json:"customer_service"`
}

// GameTypeGroup 游戏类型分组
type GameTypeGroup struct {
	Name      string          `json:"name"`
	Icon      string          `json:"icon"`
	Platforms []PlatformGroup `json:"platforms"`
}

// PlatformGroup 平台分组
type PlatformGroup struct {
	Name  string     `json:"name"`
	Icon  string     `json:"icon"`
	Games []GameItem `json:"games"`
}

// GameItem 游戏项
type GameItem struct {
	Name     string `json:"name"`
	GameCode string `json:"game_code"`
	Icon     string `json:"icon"`
}

// GetAllGamesList 获取所有游戏列表（按游戏类型和平台sort倒序）
func (s GameService) GetAllGamesList(packageId int) (*AllGamesResponse, error) {
	// 查询所有游戏类型（按sort倒序）
	var gameTypes []models.GameType
	_, err := models.CreateGameTypeModel().QueryTable(&models.GameType{}).
		Filter("package_id", packageId).
		Filter("status", 1).
		Filter("is_deleted", 0).
		OrderBy("-sort").
		All(&gameTypes, "id", "name")
	if err != nil {
		return nil, fmt.Errorf("查询游戏类型失败: %v", err)
	}
	gameTypeMap := make(map[int]string)
	for _, gt := range gameTypes {
		gameTypeMap[gt.Id] = strings.ToLower(gt.Name)
	}

	// 查询所有平台（按sort倒序）
	var platforms []models.Platform
	_, err = models.CreatePlatformModel().QueryTable(&models.Platform{}).
		Filter("package_id", packageId).
		Filter("status", 1).
		Filter("is_deleted", 0).
		OrderBy("-sort").
		All(&platforms, "id", "name")
	if err != nil {
		return nil, fmt.Errorf("查询平台失败: %v", err)
	}
	platformMap := make(map[int]string)
	for _, p := range platforms {
		platformMap[p.Id] = strings.ToLower(p.Name)
	}

	// 查询所有游戏
	var games []models.Game
	_, err = models.CreateGameModel().QueryTable(&models.Game{}).
		Filter("package_id", packageId).
		Filter("is_deleted", 0).
		Filter("status", 1).
		All(&games, "id", "name", "code", "cover", "game_type_id", "platform_id")
	if err != nil {
		return nil, fmt.Errorf("查询游戏列表失败: %v", err)
	}

	// 按 game_type -> platform_name 分组
	typeGroupMap := make(map[string]map[string][]GameItem)

	for _, game := range games {
		gameTypeName := gameTypeMap[game.GameTypeId]
		platformName := platformMap[game.PlatformId]

		if gameTypeName == "" || platformName == "" {
			continue
		}

		if typeGroupMap[gameTypeName] == nil {
			typeGroupMap[gameTypeName] = make(map[string][]GameItem)
		}

		game.PlatformName = platformName
		typeGroupMap[gameTypeName][platformName] = append(typeGroupMap[gameTypeName][platformName], GameItem{
			Name:     game.Name,
			GameCode: game.Code,
			Icon:     s.GetCoverUrl(game, packageId),
		})
	}

	// 按 gameTypes 顺序构建响应
	var result AllGamesResponse
	for _, gt := range gameTypes {
		gtName := strings.ToLower(gt.Name)
		platformsMap, ok := typeGroupMap[gtName]
		if !ok {
			continue
		}

		var pg GameTypeGroup
		pg.Name = gtName
		pg.Icon = (&GameTypeService{}).GetCoverUrl(gt, packageId)

		// 按 platforms 顺序构建平台分组
		for _, p := range platforms {
			pName := strings.ToLower(p.Name)
			if games, ok := platformsMap[pName]; ok {
				pg.Platforms = append(pg.Platforms, PlatformGroup{
					Name:  pName,
					Games: games,
					Icon:  (&PlatformService{}).GetCoverUrl(p, packageId),
				})
			}
		}

		if len(pg.Platforms) > 0 {
			result.Games = append(result.Games, pg)
		}
	}

	// 查询客服地址
	result.CustomerService, err = s.GetCustomerService(packageId)
	if err != nil {
		return nil, fmt.Errorf("查询客服地址失败: %v", err)
	}

	return &result, nil
}

// CustomerServiceConfigKey 客服配置key
const CustomerServiceConfigKey = "customer_service_url"

// GetCustomerService 获取客服地址
func (s GameService) GetCustomerService(packageId int) (string, error) {
	model := models.CreateConfigModel()
	var row models.Config
	err := model.QueryTable(new(models.Config)).
		Filter("config_key", CustomerServiceConfigKey).
		Filter("package_id", packageId).
		One(&row, "config_value")
	if errors.Is(err, orm.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return row.ConfigValue, nil
}

// SetCustomerService 设置客服地址
func (s GameService) SetCustomerService(packageId int, url string) error {
	model := models.CreateConfigModel()

	var row models.Config
	err := model.QueryTable(new(models.Config)).
		Filter("config_key", CustomerServiceConfigKey).
		Filter("package_id", packageId).
		One(&row, "id", "config_value")

	if errors.Is(err, orm.ErrNoRows) {
		// 不存在则插入
		_, err = model.Insert(&models.Config{
			ConfigKey:   CustomerServiceConfigKey,
			ConfigValue: url,
			PackageId:   packageId,
		})
		return err
	} else if err != nil {
		return err
	}

	// 存在则更新
	row.ConfigValue = url
	_, err = model.Update(&row, "config_value")
	return err
}
