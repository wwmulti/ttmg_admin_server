package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
	"net/url"
	"strconv"
	"time"
)

type BannerService struct{}

type BannerListRequest struct {
	Id        int        `form:"id"`             // ID
	PackageId int        `form:"package_id"`     // 包ID
	Name      string     `form:"name" op:"like"` // 轮播图名称
	Status    int        `form:"status"`         // 状态
	Page      int        `form:"page"`           // 页码
	PageSize  int        `form:"page_size"`      // 每页数量
	RawQuery  url.Values `form:"-"`
}

// BannerList 轮播图列表
func (s BannerService) BannerList(request BannerListRequest, needReload int, lang string, adminId int64) (map[string]interface{}, error) {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	bannerModel := models.CreateBannerModel()
	condition, sort := bannerModel.BuildCondition(request, "-id")
	if !isRootAdmin {
		if request.PackageId != 0 {
			inArray := utils.InArray(request.PackageId, packageSlice)
			if !inArray {
				return nil, fmt.Errorf("QuanXianBuZu")
			}
		} else {
			if len(packageSlice) == 0 {
				return nil, fmt.Errorf("QuanXianBuZu")
			}
			delete(condition, "package_id")
			condition["package_id__in"] = packageSlice
		}
	}
	condition["is_deleted"] = 0
	data, total, err := bannerModel.GetPageList(&models.Banner{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}
	var bannerType []map[string]interface{}
	var packages []models.Package
	if needReload == 1 {
		bannerType = s.AllBannerType(lang)
		packageService := PackageService{}
		packages = packageService.GetMyAllPackageList(int(adminId))
	}
	return map[string]interface{}{
		"list":         data,
		"total":        total,
		"current_page": request.Page,
		"banner_type":  bannerType,
		"packages":     packages,
	}, nil
}

// GetBannerById 获取轮播图信息
func (s BannerService) GetBannerById(id int) (*models.Banner, error) {
	banner := &models.Banner{}
	bannerModel := models.CreateBannerModel()
	err := bannerModel.QueryTable(new(models.Banner)).
		Filter("id", id).
		One(banner)
	if err != nil {
		return nil, err
	}
	return banner, nil
}

// AllBannerType 所有轮播图类型
func (s BannerService) AllBannerType(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{"id": config.BannerTypeIndex, "name": i18n.Tr(language, "ShouYeLunBoTu")},
	}
}

type AddBannerRequest struct {
	PackageIds []int  `json:"package_ids"` // 包ID
	Name       string `json:"name"`        // 轮播图名称
	Type       int    `json:"type"`        // 轮播图类型
	Url        string `json:"url"`         // 跳转地址
	Sort       int    `json:"sort"`        // 排序
	Status     int    `json:"status"`      // 状态
}

// AddBanner 添加轮播图
func (s BannerService) AddBanner(request AddBannerRequest, adminId int64) error {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	err := s.validateAddBanner(request)
	if err != nil {
		return err
	}
	tx, err := models.CreateBannerModel().Begin()
	if err != nil {
		logs.Error("添加轮播图事务开启失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	for _, packageId := range request.PackageIds {
		if !isRootAdmin && !utils.InArray(packageId, packageSlice) {
			return fmt.Errorf("QuanXianBuZu")
		}
		_, err = tx.Insert(&models.Banner{
			PackageId: packageId,
			Name:      request.Name,
			Type:      request.Type,
			Url:       request.Url,
			Sort:      request.Sort,
			Status:    request.Status,
		})
		if err != nil {
			tx.Rollback()
			logs.Error("添加轮播图失败，失败原因为：%v", err)
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("添加轮播图事务提交失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	for _, packageId := range request.PackageIds {
		err := s.refreshBanner(packageId)
		if err != nil {
			logs.Error("刷新轮播图失败，失败原因为：%v", err)
			return fmt.Errorf("ShuaXinHuanCunShiBai")
		}
	}
	return nil
}

// 添加轮播图参数验证
func (s BannerService) validateAddBanner(request AddBannerRequest) error {
	if request.Type == 0 {
		return fmt.Errorf("LunBoTuLeiXingBiTian")
	}
	if request.Url == "" {
		return fmt.Errorf("LunBoTuDiZhiBiTian")
	}
	return nil
}

type EditBannerRequest struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`   // 轮播图名称
	Type   int    `json:"type"`   // 轮播图类型
	Url    string `json:"url"`    // 跳转地址
	Sort   int    `json:"sort"`   // 排序
	Status int    `json:"status"` // 状态
}

// EditBanner 编辑轮播图
func (s BannerService) EditBanner(request EditBannerRequest, adminId int64) error {
	banner, err := s.GetBannerById(request.Id)
	if err != nil {
		logs.Error("获取轮播图失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(banner.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	err = s.validateEditBanner(request)
	if err != nil {
		return err
	}
	tx, err := models.CreateBannerModel().Begin()
	if err != nil {
		logs.Error("编辑轮播图事务开启失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	fields := []string{
		"name",
		"type",
		"url",
		"sort",
		"status",
	}
	_, err = tx.Update(&models.Banner{
		Id:     request.Id,
		Name:   request.Name,
		Type:   request.Type,
		Url:    request.Url,
		Sort:   request.Sort,
		Status: request.Status,
	}, fields...)
	if err != nil {
		tx.Rollback()
		logs.Error("修改轮播图失败，失败原因为：%v", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("编辑轮播图事务提交失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	err = s.refreshBanner(banner.PackageId)
	if err != nil {
		logs.Error("刷新轮播图失败，失败原因为：%v", err)
		return fmt.Errorf("ShuaXinHuanCunShiBai")
	}
	return nil
}

// 编辑轮播图参数验证
func (s BannerService) validateEditBanner(request EditBannerRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("LunBoIDBiTian")
	}
	if request.Type == 0 {
		return fmt.Errorf("LunBoTuLeiXingBiTian")
	}
	if request.Url == "" {
		return fmt.Errorf("LunBoTuDiZhiBiTian")
	}
	return nil
}

// DeleteBanner 删除轮播图
func (s BannerService) DeleteBanner(id int, adminId int64) error {
	if id == 0 {
		return fmt.Errorf("LunBoIDBiTian")
	}
	banner, err := s.GetBannerById(id)
	if err != nil {
		logs.Error("获取轮播图失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(banner.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	tx, err := models.CreateBannerModel().Begin()
	if err != nil {
		logs.Error("删除轮播图事务开启失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	_, err = tx.Update(&models.Banner{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("删除轮播图事务提交失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}

	err = s.refreshBanner(banner.PackageId)
	if err != nil {
		logs.Error("刷新轮播图失败，失败原因为：%v", err)
		return fmt.Errorf("ShuaXinHuanCunShiBai")
	}
	return nil
}

type ChangeBannerStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangeBannerStatus 修改轮播图状态
func (s BannerService) ChangeBannerStatus(request ChangeBannerStatusRequest, adminId int64) error {
	if request.Id == 0 {
		return fmt.Errorf("LunBoIDBiTian")
	}
	banner, err := s.GetBannerById(request.Id)
	if err != nil {
		logs.Error("获取轮播图失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(banner.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	tx, err := models.CreateBannerModel().Begin()
	if err != nil {
		logs.Error("删除轮播图事务开启失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	_, err = tx.Update(&models.Banner{
		Id:     request.Id,
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("修改轮播图状态事务提交失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	err = s.refreshBanner(banner.PackageId)
	if err != nil {
		logs.Error("刷新轮播图失败，失败原因为：%v", err)
		return fmt.Errorf("ShuaXinHuanCunShiBai")
	}
	return nil
}

// 刷新轮播缓存
func (s BannerService) refreshBanner(packageId int) error {
	now := time.Now().Unix()
	secret := config.System.CurlSecretKey.Value
	params := map[string]interface{}{
		"package_id": strconv.Itoa(packageId),
		"time":       now,
	}
	request, _ := utils.SignMap(params, secret)
	packageService := PackageService{}
	resp := RequestPackageUrlRequest{
		PackageId: packageId,
		Params:    request,
		Url:       "/data/refresh/refreshBanner",
		Type:      int(config.CurlRequestTypePost),
	}
	response, err := packageService.RequestPackageUrl(resp)
	if err != nil {
		logs.Error("刷新轮播失败：%v", err)
		return fmt.Errorf("ShuXinHuanCunShiBai")
	}
	if !response.OK() {
		logs.Error("刷新轮播失败：%v", response)
		return fmt.Errorf("ShuXinHuanCunShiBai")
	}
	return nil
}

// RefreshPackageBanner 刷新包打脸缓存
func (s BannerService) RefreshPackageBanner(packageId int) error {
	err := s.refreshBanner(packageId)
	if err != nil {
		logs.Error("刷新浮动图标失败：%v", err)
		return fmt.Errorf(err.Error())
	}
	return nil
}
