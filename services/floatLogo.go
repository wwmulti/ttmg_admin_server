package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"net/url"
	"strconv"
	"time"
)

type FloatLogoService struct{}

type FloatLogoListRequest struct {
	Id        int        `form:"id"`         // ID
	PackageId int        `form:"package_id"` // 包ID
	Status    int        `form:"status"`     // 状态
	Page      int        `form:"page"`       // 页码
	PageSize  int        `form:"page_size"`  // 每页数量
	RawQuery  url.Values `form:"-"`
}

// FloatLogoList 浮动图标列表
func (s FloatLogoService) FloatLogoList(request FloatLogoListRequest, needReload int, adminId int64) (map[string]interface{}, error) {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	floatLogoModel := models.CreateFloatLogoModel()
	condition, sort := floatLogoModel.BuildCondition(request, "-id")
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
	data, total, err := floatLogoModel.GetPageList(&models.FloatLogo{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}
	var packages []models.Package
	if needReload == 1 {
		packageService := PackageService{}
		packages = packageService.GetMyAllPackageList(int(adminId))
	}
	return map[string]interface{}{
		"list":         data,
		"total":        total,
		"current_page": request.Page,
		"packages":     packages,
	}, nil
}

// GetFloatLogoById 获取浮动图标信息
func (s FloatLogoService) GetFloatLogoById(id int) (*models.FloatLogo, error) {
	floatLogo := &models.FloatLogo{}
	floatLogoModel := models.CreateFloatLogoModel()
	err := floatLogoModel.QueryTable(new(models.FloatLogo)).
		Filter("id", id).
		One(floatLogo)
	if err != nil {
		return nil, err
	}
	return floatLogo, nil
}

type AddFloatLogoRequest struct {
	PackageIds []int  `json:"package_ids"` // 包ID
	Logo       string `json:"logo"`        // 图标
	Link       string `json:"link"`        // 跳转地址
	Status     int    `json:"status"`      // 状态
}

// AddFloatLogo 添加浮动图标
func (s FloatLogoService) AddFloatLogo(request AddFloatLogoRequest, adminId int64) error {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	err := s.validateAddFloatLogo(request)
	if err != nil {
		return err
	}
	tx, err := models.CreateFloatLogoModel().Begin()
	if err != nil {
		logs.Error("添加浮动图标事务开启失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	for _, packageId := range request.PackageIds {
		if !isRootAdmin && !utils.InArray(packageId, packageSlice) {
			return fmt.Errorf("QuanXianBuZu")
		}
		_, err = tx.Insert(&models.FloatLogo{
			PackageId: packageId,
			Logo:      request.Logo,
			Link:      request.Link,
			Status:    request.Status,
		})
		if err != nil {
			tx.Rollback()
			logs.Error("添加浮动图标失败，失败原因为：%v", err)
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("添加浮动图标事务提交失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	for _, packageId := range request.PackageIds {
		err := s.refreshFloatLogo(packageId)
		if err != nil {
			logs.Error("刷新浮动图标失败，失败原因为：%v", err)
			return fmt.Errorf("ShuaXinHuanCunShiBai")
		}
	}
	return nil
}

// 添加浮动图标参数验证
func (s FloatLogoService) validateAddFloatLogo(request AddFloatLogoRequest) error {
	if request.Logo == "" {
		return fmt.Errorf("TuBiaoBiTian")
	}
	return nil
}

type EditFloatLogoRequest struct {
	Id     int    `json:"id"`
	Logo   string `json:"logo"`   // 图标
	Link   string `json:"link"`   // 跳转地址
	Status int    `json:"status"` // 状态
}

// EditFloatLogo 编辑浮动图标
func (s FloatLogoService) EditFloatLogo(request EditFloatLogoRequest, adminId int64) error {
	floatLogo, err := s.GetFloatLogoById(request.Id)
	if err != nil {
		logs.Error("获取浮动图标失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(floatLogo.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	fmt.Println(request)
	err = s.validateEditFloatLogo(request)
	if err != nil {
		return err
	}
	tx, err := models.CreateFloatLogoModel().Begin()
	if err != nil {
		logs.Error("编辑浮动图标事务开启失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	fields := []string{
		"logo",
		"link",
		"status",
	}
	_, err = tx.Update(&models.FloatLogo{
		Id:     request.Id,
		Logo:   request.Logo,
		Link:   request.Link,
		Status: request.Status,
	}, fields...)
	if err != nil {
		tx.Rollback()
		logs.Error("修改浮动图标失败，失败原因为：%v", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("编辑浮动图标事务提交失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	err = s.refreshFloatLogo(floatLogo.PackageId)
	if err != nil {
		logs.Error("刷新浮动图标失败，失败原因为：%v", err)
		return fmt.Errorf("ShuaXinHuanCunShiBai")
	}
	return nil
}

// 编辑浮动图标参数验证
func (s FloatLogoService) validateEditFloatLogo(request EditFloatLogoRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("TuBiaoIDBiTian")
	}
	if request.Logo == "" {
		return fmt.Errorf("TuBiaoBiTian")
	}
	return nil
}

// DeleteFloatLogo 删除浮动图标
func (s FloatLogoService) DeleteFloatLogo(id int, adminId int64) error {
	if id == 0 {
		return fmt.Errorf("TuBiaoIDBiTian")
	}
	floatLogo, err := s.GetFloatLogoById(id)
	if err != nil {
		logs.Error("获取浮动图标失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(floatLogo.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	tx, err := models.CreateFloatLogoModel().Begin()
	if err != nil {
		logs.Error("删除浮动图标事务开启失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	_, err = tx.Update(&models.FloatLogo{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("删除浮动图标事务提交失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	err = s.refreshFloatLogo(floatLogo.PackageId)
	if err != nil {
		logs.Error("刷新浮动图标失败，失败原因为：%v", err)
		return fmt.Errorf("ShuaXinHuanCunShiBai")
	}
	return nil
}

type ChangeFloatLogoStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangeFloatLogoStatus 修改浮动图标状态
func (s FloatLogoService) ChangeFloatLogoStatus(request ChangeFloatLogoStatusRequest, adminId int64) error {
	if request.Id == 0 {
		return fmt.Errorf("TuBiaoIDBiTian")
	}
	floatLogo, err := s.GetFloatLogoById(request.Id)
	if err != nil {
		logs.Error("获取浮动图标失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(floatLogo.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	tx, err := models.CreateFloatLogoModel().Begin()
	if err != nil {
		logs.Error("删除浮动图标事务开启失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	_, err = tx.Update(&models.FloatLogo{
		Id:     request.Id,
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logs.Error("修改浮动图标状态事务提交失败，失败原因为：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}
	err = s.refreshFloatLogo(floatLogo.PackageId)
	if err != nil {
		logs.Error("刷新浮动图标失败，失败原因为：%v", err)
		return fmt.Errorf("ShuaXinHuanCunShiBai")
	}
	return nil
}

// 刷新浮动图标缓存
func (s FloatLogoService) refreshFloatLogo(packageId int) error {
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
		Url:       "/data/refresh/refreshFloatLogo",
		Type:      int(config.CurlRequestTypePost),
	}
	response, err := packageService.RequestPackageUrl(resp)
	if err != nil {
		logs.Error("刷新浮动图标失败：%v", err)
		return fmt.Errorf("ShuXinHuanCunShiBai")
	}
	if !response.OK() {
		logs.Error("刷新浮动图标失败：%v", response)
		return fmt.Errorf("ShuXinHuanCunShiBai")
	}
	return nil
}

// RefreshPackageFloatLogo 刷新包浮动图标缓存
func (s FloatLogoService) RefreshPackageFloatLogo(packageId int) error {
	err := s.refreshFloatLogo(packageId)
	if err != nil {
		logs.Error("刷新浮动图标失败：%v", err)
		return fmt.Errorf(err.Error())
	}
	return nil
}
