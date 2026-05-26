package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type MerchantService struct{}

type MerchantRequestParams struct {
	Page     int
	PageSize int
	Title    string
	Status   int
}

type MerchantListResponse struct {
	Total       int               `json:"total"`
	Lists       []models.Merchant `json:"list"`
	PackageList []models.Package  `json:"package_list"`
}

type EditMerchantRequestParams struct {
	Id           int    `json:"id,omitempty"`
	Title        string `json:"title"`
	Status       int    `json:"status"`
	Domain       string `json:"domain"`
	Secret       string `json:"secret"`
	Token        string `json:"token"`
	Currency     string `json:"currency"`
	Type         int    `json:"type"`
	SupplierType int    `json:"supplier_type"`
}

type DeleteMerchantRequestParams struct {
	Ids string `json:"ids"`
}

type PullMerchantGameRequestParams struct {
	PackageId  int `json:"package_id"`
	MerchantId int `json:"merchant_id"`
}

// 获取列表
func (s *MerchantService) GetList(params MerchantRequestParams, userId int) MerchantListResponse {
	condition := make(map[string]interface{})
	condition["is_deleted"] = 0

	if len(params.Title) > 0 {
		condition["title__icontains"] = params.Title
	}

	if params.Status > -1 {
		condition["status"] = params.Status
	}

	result, total, _ := models.CreateMerchantModel().GetPageList(new(models.Merchant), condition, params.Page, params.PageSize, "-id")
	merchantList, ok := result.([]models.Merchant)
	if !ok {
		merchantList = []models.Merchant{}
	}

	return MerchantListResponse{
		Total:       int(total),
		Lists:       merchantList,
		PackageList: (&PackageService{}).GetMyAllPackageList(userId),
	}
}

// 修改商户
func (s *MerchantService) EditMerchant(params EditMerchantRequestParams) error {
	if s.isSupplierTypeExists(params.Id, params.SupplierType, params.Type) {
		return fmt.Errorf("GaiLeiXingDeGongYingShangYiPeiZhi")
	}

	model := models.CreateMerchantModel()

	domain, isMatch := utils.IsValidUrlFormat(params.Domain)
	token := strings.TrimSpace(params.Token)
	secret := strings.TrimSpace(params.Secret)
	currency := strings.TrimSpace(params.Currency)
	if !isMatch {
		return fmt.Errorf("YuMingGeShiBuZhengQue")
	}

	status := 0
	if params.Status > 0 {
		status = 1
	}
	fields := []string{
		"title", "domain", "token", "secret", "currency", "status", "type", "supplier_type",
	}
	_, err := model.Update(&models.Merchant{
		Id:           params.Id,
		Title:        params.Title,
		Domain:       domain,
		Token:        token,
		Secret:       secret,
		Currency:     currency,
		Status:       status,
		Type:         params.Type,
		SupplierType: params.SupplierType,
	}, fields...)
	if err != nil {
		logs.Error("更新商户失败：%v", err)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

func (s *MerchantService) isSupplierTypeExists(id int, supplierType int, gameType int) bool {
	model := models.CreateMerchantModel()
	isExists := model.QueryTable(new(models.Merchant)).Filter("id__ne", id).Filter("supplier_type", supplierType).Filter("type", gameType).Filter("is_deleted", 0).Exist()
	return isExists
}

// 创建商户
func (s *MerchantService) AddMerchant(params EditMerchantRequestParams) error {
	if s.isSupplierTypeExists(0, params.SupplierType, params.Type) {
		return fmt.Errorf("GaiLeiXingDeGongYingShangYiPeiZhi")
	}

	model := models.CreateMerchantModel()
	domain, isMatch := utils.IsValidUrlFormat(params.Domain)
	token := strings.TrimSpace(params.Token)
	secret := strings.TrimSpace(params.Secret)
	currency := strings.TrimSpace(params.Currency)
	if !isMatch {
		return fmt.Errorf("YuMingGeShiBuZhengQue")
	}

	status := 0
	if params.Status > 0 {
		status = 1
	}
	_, err := model.Insert(&models.Merchant{
		Title:        params.Title,
		Domain:       domain,
		Token:        token,
		Secret:       secret,
		Currency:     currency,
		Status:       status,
		Type:         params.Type,
		SupplierType: params.SupplierType,
	})
	if err != nil {
		logs.Error("创建商户失败：%v", err)
		return fmt.Errorf("ChuangJianShuJuShiBai")
	}
	return nil
}

// 删除商户
func (s *MerchantService) DeleteMerchant(params DeleteMerchantRequestParams) error {
	idArray := strings.Split(params.Ids, ",")
	idInts := make([]int, len(idArray))
	for i, id := range idArray {
		idInts[i], _ = strconv.Atoi(id)
	}
	model := models.CreateMerchantModel()
	_, err := model.QueryTable(new(models.Merchant)).
		Filter("id__in", idInts).
		Update(orm.Params{
			"is_deleted": 1,
		})
	if err != nil {
		return err
	}
	return nil
}

// 获取支持的商户
func (s *MerchantService) GetSupportList(gameType int) []models.Merchant {
	model := models.CreateMerchantModel()
	var lists []models.Merchant

	fields := []string{
		"id", "title",
	}
	_, err := model.QueryTable(new(models.Merchant)).Filter("status", 1).Filter("is_deleted", 0).Filter("type", gameType).All(&lists, fields...)
	if err != nil {
		return []models.Merchant{}
	}
	return lists
}

// 拉取商户的游戏
func (s *MerchantService) PullMerchantGame(params PullMerchantGameRequestParams) error {
	model := models.CreateMerchantModel()
	var info models.Merchant
	model.QueryTable(new(models.Merchant)).Filter("id", params.MerchantId).One(&info)

	if info.SupplierType == int(config.ZySupplierType) {
		err := (GameService{}).InitZyGameList(params.PackageId)
		if err != nil {
			return fmt.Errorf("同步游戏列表失败")
		}
	}
	return nil
}
