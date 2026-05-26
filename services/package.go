package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type PackageService struct{}

type PackageRequestParams struct {
	Page     int
	PageSize int
	Title    string
	Domain   string
}

type AddPlatformRequestParams struct {
	Id               int    `json:"id,omitempty"`
	GroupId          int    `json:"group_id"` // 分组id
	Title            string `json:"title"`
	Domain           string `json:"domain"`
	ApiDomain        string `json:"api_domain"`
	Icon             string `json:"icon"`                              // icon图片地址
	Logo             string `json:"logo"`                              // logo图片地址
	LoadPic          string `json:"load_pic"`                          // 平台加载图片
	Pg               string `json:"pg"`                                // pg商户 结构1,2
	Pp               string `json:"pp"`                                // pp商户 结构1,2
	OpenRegister     int    `json:"open_register"`                     // 打开注册 0关闭 1打开
	IsModifyRegister bool   `json:"is_modify_open_register,omitempty"` // 是否修改注册开关
}

type PackageListResponse struct {
	Total     int                   `json:"total"`
	Lists     []models.Package      `json:"list"`
	Pg        []models.Merchant     `json:"pg"`
	Pp        []models.Merchant     `json:"pp"`
	GroupList []models.PackageGroup `json:"group_list"`
}

type PackageConfigListResponse struct {
	Id          int  `json:"int"`
	ActiveType1 bool `json:"activeType_1"` // 活动类型对应的状态值
	ActiveType2 bool `json:"activeType_2"`
	ActiveType3 bool `json:"activeType_3"`
	ActiveType4 bool `json:"activeType_4"`
	ActiveType5 bool `json:"activeType_5"`
	ActiveType6 bool `json:"activeType_6"`
	ActiveType7 bool `json:"activeType_7"`
}

type EditPackageConfigRequestParams struct {
	Id         int   `json:"id"`
	ActiveType []int `json:"activeType"` // 活动类型
}

func (s *PackageService) BuildCondition(params PackageRequestParams, userId int) map[string]interface{} {
	condition := make(map[string]interface{})
	condition["is_deleted"] = 0

	if len(params.Domain) > 0 {
		condition["domain__icontains"] = params.Domain
	}

	if len(params.Title) > 0 {
		condition["title__icontains"] = params.Title
	}

	if userId > 0 {
		// 只展示授权的分包
		packageIds := (&AuthGroupService{}).GetPackageIds(userId)
		if len(packageIds) == 0 { // 没有分配
			condition["id__lt"] = 0
		} else if packageIds != "-1" {
			condition["id__in"] = utils.StringSliceToInt(packageIds, ",")
		}
	}
	return condition
}

// GetList 获取平台列表
func (s *PackageService) GetList(params PackageRequestParams, userId int) PackageListResponse {
	condition := s.BuildCondition(params, userId)

	result, total, _ := models.CreatePackageModel().GetPageList(new(models.Package), condition, params.Page, params.PageSize, "-id")
	accounts, ok := result.([]models.Package)
	if !ok {
		accounts = []models.Package{}
	}

	packageGroupParams := PackageGroupListRequestParams{
		Page:     1,
		PageSize: 10000,
		Status:   -1,
	}
	groupList := (&PackageGroupService{}).GetList(packageGroupParams, userId)

	return PackageListResponse{
		Total:     int(total),
		Lists:     accounts,
		Pg:        (&MerchantService{}).GetSupportList(int(config.PgGameType)),
		Pp:        (&MerchantService{}).GetSupportList(int(config.PpGameType)),
		GroupList: groupList.List,
	}
}

// 获取自己的所有的平台
func (s *PackageService) GetMyAllPackageList(userId int) []models.Package {
	condition := make(map[string]interface{})
	condition["is_deleted"] = 0

	// 只展示授权的分包
	packageIds := (&AuthGroupService{}).GetPackageIds(userId)
	if len(packageIds) == 0 { // 没有分配
		condition["id__lt"] = 0
	} else if packageIds != "-1" {
		condition["id__in"] = utils.StringSliceToInt(packageIds, ",")
	}

	packageModel := models.CreatePackageModel()
	db := packageModel.QueryTable(&models.Package{})
	list := make([]models.Package, 0)
	packageModel.Where(db, condition).All(&list, "id", "title")

	return list
}

type AllPackageResponse struct {
	Code int              `json:"code" example:"200"`    // 状态码
	Msg  string           `json:"msg" example:"success"` // 提示信息
	Data []models.Package `json:"data"`
}

// AllPackage 获取所有平台
func (s *PackageService) AllPackage(ids ...[]int) ([]models.Package, error) {
	packageModel := models.CreatePackageModel()
	var list []models.Package
	condition := map[string]interface{}{
		"is_deleted": 0,
	}
	if len(ids) > 0 {
		condition["id__in"] = ids
	}
	db := packageModel.QueryTable(&models.Package{})
	_, err := packageModel.Where(db, condition).All(&list)
	if nil != err {
		return nil, err
	}

	return list, nil
}

// AddPackage 创建平台
func (s *PackageService) AddPackage(params AddPlatformRequestParams) error {
	model := models.CreatePackageModel()
	// 判断分组是否存在
	var packageGroupInfo models.PackageGroup
	err := model.QueryTable(new(models.PackageGroup)).Filter("id", params.GroupId).One(&packageGroupInfo)
	if err != nil {
		return fmt.Errorf("PingTaiFenZuBuCunZai")
	}

	isExists := model.QueryTable(new(models.Package)).Filter("title", params.Title).Filter("is_deleted", 0).Exist()
	if isExists {
		return fmt.Errorf("PingTaiMingChengYiCunZai")
	}

	isExists = model.QueryTable(new(models.Package)).Filter("domain", params.Domain).Filter("is_deleted", 0).Exist()
	if isExists {
		return fmt.Errorf("YuMingYiCunZai")
	}

	trErr := model.OrmerMaster.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		id, addErr := txOrm.Insert(&models.Package{
			Title:        params.Title,
			Domain:       params.Domain,
			ApiDomain:    params.ApiDomain,
			Icon:         params.Icon,
			Logo:         params.Logo,
			LoadPic:      params.LoadPic,
			Pg:           params.Pg,
			Pp:           params.Pp,
			OpenRegister: 1,
			CreateTime:   time.Now().Unix(),
			GroupId:      params.GroupId,
		})
		if addErr != nil {
			return addErr
		}
		// 更新平台
		fields := []string{"package_ids"}

		packageIds := fmt.Sprintf("%v,%v,", strings.TrimRight(packageGroupInfo.PackageIds, ","), id)
		_, updateErr := txOrm.Update(&models.PackageGroup{
			Id:         params.GroupId,
			PackageIds: packageIds,
		}, fields...)
		if updateErr != nil {
			return updateErr
		}

		// 判断分组是否存在配置，有配置则同步配置
		syncErr := s.SyncPackageConfig(txOrm, params.GroupId, int(id))
		if syncErr != nil {
			return syncErr
		}
		return nil
	})
	if trErr != nil {
		logs.Error("创建平台失败:%v", trErr)
		return fmt.Errorf("ChuangJianShuJuShiBai")
	}
	return nil
}

// EditPackage 编辑平台
func (s *PackageService) EditPackage(params AddPlatformRequestParams) error {
	model := models.CreatePackageModel()
	params.Domain = strings.TrimSpace(params.Domain)
	params.ApiDomain = strings.TrimSpace(params.ApiDomain)

	var packageInfo models.Package
	err := model.QueryTable(new(models.Package)).Filter("id", params.Id).One(&packageInfo)
	if err != nil {
		return fmt.Errorf("ShuJuYiChang")
	}

	isExists := model.QueryTable(new(models.Package)).Filter("title", params.Title).Filter("id__ne", params.Id).Filter("is_deleted", 0).Exist()
	if isExists {
		return fmt.Errorf("PingTaiMingChengYiCunZai")
	}

	isExists = model.QueryTable(new(models.Package)).Filter("domain", params.Domain).Filter("id__ne", params.Id).Filter("is_deleted", 0).Exist()
	if isExists {
		return fmt.Errorf("YuMingYiCunZai")
	}

	if params.IsModifyRegister {
		field := []string{"open_register", "update_time"}

		_, updateErr := model.Update(&models.Package{
			Id:           params.Id,
			OpenRegister: params.OpenRegister,
			UpdateTime:   time.Now().Unix(),
		}, field...)
		if updateErr != nil {
			return updateErr
		}
		return nil
	}

	var packageGroupInfo models.PackageGroup
	// 现在的分组
	groupErr := model.QueryTable(new(models.PackageGroup)).Filter("id", params.GroupId).One(&packageGroupInfo)
	if groupErr != nil {
		return fmt.Errorf("PingTaiFenZuBuCunZai")
	}
	// 之前的分组

	trErr := model.OrmerMaster.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		if !params.IsModifyRegister && params.GroupId != packageInfo.GroupId {
			updateGroupErr := (&PackageGroupService{}).ChangeGroupId(txOrm, params.Id, params.GroupId)
			if updateGroupErr != nil {
				return updateGroupErr
			}
		}

		fields := []string{
			"title", "domain", "api_domain", "icon", "logo", "load_pic", "pg", "pp", "update_time", "group_id",
		}

		if params.IsModifyRegister {
			fields = []string{"open_register"}
		}
		_, updateErr := txOrm.Update(&models.Package{
			Id:           params.Id,
			GroupId:      params.GroupId,
			Title:        params.Title,
			Domain:       params.Domain,
			ApiDomain:    params.ApiDomain,
			Icon:         params.Icon,
			Logo:         params.Logo,
			LoadPic:      params.LoadPic,
			Pg:           params.Pg,
			Pp:           params.Pp,
			OpenRegister: params.OpenRegister,
			UpdateTime:   time.Now().Unix(),
		}, fields...)
		if updateErr != nil {
			return updateErr
		}
		return nil
	})

	if trErr != nil {
		logs.Error("更新平台失败:%v", trErr)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

// DelPackage 删除平台
func (s *PackageService) DelPackage(ids string) error {
	model := models.CreatePackageModel()
	// 移除平台分组
	trErr := model.OrmerMaster.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		idArray := strings.Split(ids, ",")
		idInts := make([]int, len(idArray))
		for i, id := range idArray {
			idInts[i], _ = strconv.Atoi(id)

			var packageInfo models.Package
			pErr := txOrm.QueryTable(new(models.Package)).Filter("id", idInts[i]).One(&packageInfo)
			if pErr != nil {
				return pErr
			}
			removeErr := (&PackageGroupService{}).RemovePackageId(txOrm, packageInfo.GroupId, packageInfo.Id)
			if removeErr != nil {
				return removeErr
			}
			err := (&ConfigService{}).DeletePackageConfig(packageInfo.Id, txOrm)
			if err != nil {
				return err
			}
		}
		_, err := model.QueryTable(new(models.Package)).
			Filter("id__in", idInts).
			Update(orm.Params{
				"is_deleted":  1,
				"group_id":    0,
				"update_time": time.Now().Unix(),
			})
		if err != nil {
			return err
		}
		return nil
	})

	if trErr != nil {
		logs.Error("更新平台失败:%v", trErr)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

// GetInfo 获取分包的信息
func (s *PackageService) GetInfo(id int) models.Package {
	model := models.CreatePackageModel()
	var info models.Package
	err := model.QueryTable(new(models.Package)).Filter("id", id).One(&info)
	if err != nil {
		return models.Package{}
	}
	return info
}

// ConfigList 平台配置列表
func (s *PackageService) ConfigList(id int) map[string]interface{} {
	model := models.CreateActivesModel()
	activityMap := make(map[string]interface{})
	for _, activeType := range config.ActiveTypeList {
		logs.Info("item=>%v", activeType)
		isExists := model.QueryTable(new(models.Actives)).Filter("package_id", id).Filter("active_type_id", int(activeType.ID)).Filter("status", 1).Exist()

		activityMap[fmt.Sprintf("activeType_%d", activeType.ID)] = isExists
	}
	activityMap["id"] = id
	return activityMap
}

// EditConfig 平台配置修改
func (s *PackageService) EditConfig(params EditPackageConfigRequestParams) error {
	model := models.CreatePackageModel()
	trErr := model.OrmerMaster.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		typeIds := make([]int, len(config.ActiveTypeList))
		for index, item := range config.ActiveTypeList {
			typeIds[index] = int(item.ID)
		}

		activeParams := CreateActivityByTypeRequest{
			TypeIds:    typeIds,
			PackageId:  params.Id,
			StatusList: params.ActiveType,
		}
		activeErr := (&ActiveService{}).CreateActivityByType(activeParams, txOrm)
		if activeErr != nil {
			return activeErr
		}
		return nil
	})
	if trErr != nil {
		logs.Error("配置平台失败：%v", trErr)
		return fmt.Errorf("ChuangJianShuJuShiBai")
	}
	return nil
}

// SyncPackageConfig 同步分包的配置
func (s *PackageService) SyncPackageConfig(txOrm orm.TxOrmer, groupId int, packageId int) error {
	var info models.Package
	err := txOrm.QueryTable(new(models.Package)).Filter("group_id", groupId).Filter("id__ne", packageId).Filter("is_deleted", 0).One(&info)

	var activeParams CreateActivityByTypeRequest
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) { // 当前分组没有分包不需要处理
			activeParams = CreateActivityByTypeRequest{
				PackageId:     packageId,
				CopyPackageId: 0,
			}
		} else {
			return err
		}
	} else {
		activeParams = CreateActivityByTypeRequest{
			PackageId:     packageId,
			CopyPackageId: info.Id,
		}
	}

	activeErr := (&ActiveService{}).CopyActivity(activeParams, txOrm)
	if activeErr != nil {
		return activeErr
	}

	// 复制分包的配置
	configErr := (&ConfigService{}).SyncPackageConfig(packageId, activeParams.CopyPackageId)
	if configErr != nil {
		return configErr
	}
	// 复制游戏的配置
	gameErr := (&GameService{}).SyncPackageConfig(packageId, activeParams.CopyPackageId)
	return gameErr
}

type RequestPackageUrlRequest struct {
	PackageId int                    `json:"package_id"`
	Url       string                 `json:"url"`
	Params    map[string]interface{} `json:"params"`
	Type      int                    `json:"type"`
	Language  string                 `json:"language"`
}

// RequestPackageUrl 请求分包的接口
func (s *PackageService) RequestPackageUrl(request RequestPackageUrlRequest) (*utils.Response, error) {
	if request.Language == "" {
		request.Language = "en-US"
	}
	if request.PackageId <= 0 {
		return nil, fmt.Errorf("IDBiTian")
	}
	packageInfo := s.GetInfo(request.PackageId)
	if len(packageInfo.ApiDomain) == 0 {
		return nil, fmt.Errorf("YuMingBuCunZai")
	}
	apiUrl := packageInfo.ApiDomain + request.Url
	domain := packageInfo.Domain
	if request.Type == int(config.CurlRequestTypeGet) {
		resp := utils.Get(apiUrl, utils.RequestOption{
			Headers: map[string]string{
				"Accept-Language": request.Language,
				"Do":              domain,
			},
		})
		return resp, nil
	} else {
		resp := utils.PostJSON(apiUrl, request.Params, utils.RequestOption{
			Headers: map[string]string{
				"Accept-Language": request.Language,
				"Do":              domain,
			},
		})
		return resp, nil
	}
}

// GetPackagesByIds 根据包id获取分包信息
func (s *PackageService) GetPackagesByIds(idInts []int) map[int]models.Package {
	ret := make(map[int]models.Package)

	allPackageList, err := s.AllPackage()
	if err != nil || len(allPackageList) == 0 {
		return ret
	}

	for _, item := range allPackageList {
		if slices.Contains(idInts, item.Id) {
			ret[item.Id] = item
		}
	}
	return ret
}
