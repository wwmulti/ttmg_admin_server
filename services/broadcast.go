package services

import (
	"api/models"
	"api/utils"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"net/url"
)

type BroadcastService struct{}

type BroadcastListRequest struct {
	ID        int    `form:"id"`
	PackageId int    `form:"package_id"`
	Name      string `form:"name"`
	Status    int    `form:"status"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	RawQuery  url.Values
}

// BroadcastList 广播列表
func (s BroadcastService) BroadcastList(request BroadcastListRequest, needReload int, adminId int64) (map[string]interface{}, error) {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	broadcastModel := models.CreateBroadcastModel()
	condition, sort := broadcastModel.BuildCondition(request, "-id")
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
	data, total, err := broadcastModel.GetPageList(&models.Broadcast{}, condition, request.Page, request.PageSize, sort)
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

// GetBroadcastById 获取广播信息
func (s BroadcastService) GetBroadcastById(id int) (*models.Broadcast, error) {
	floatLogo := &models.Broadcast{}
	floatLogoModel := models.CreateBroadcastModel()
	err := floatLogoModel.QueryTable(new(models.Broadcast)).
		Filter("id", id).
		One(floatLogo)
	if err != nil {
		return nil, err
	}
	return floatLogo, nil
}

type AddBroadcastRequest struct {
	PackageIds []int  `json:"package_ids"` // 包ID
	Name       string `json:"name"`
	EnContent  string `json:"en_content"`
	PtContent  string `json:"pt_content"`
	Sort       int    `json:"sort"`
	Status     int    `json:"status"`
}

// AddBroadcast 添加广播
func (s BroadcastService) AddBroadcast(request AddBroadcastRequest, adminId int64) error {
	authGroupService := AuthGroupService{}
	isRootAdmin, packageSlice := authGroupService.GetAdminIsRootAndPackageIdSlice(adminId)
	err := s.validateAddBroadcast(request)
	if err != nil {
		return err
	}
	broadcastModel := models.CreateBroadcastModel()
	for _, packageId := range request.PackageIds {
		if !isRootAdmin && !utils.InArray(packageId, packageSlice) {
			return fmt.Errorf("QuanXianBuZu")
		}
		_, err = broadcastModel.Insert(&models.Broadcast{
			PackageId: packageId,
			Name:      request.Name,
			EnContent: request.EnContent,
			PtContent: request.PtContent,
			Sort:      request.Sort,
			Status:    request.Status,
		})
	}
	if err != nil {
		return err
	}
	return nil
}

// 添加广播参数验证
func (s BroadcastService) validateAddBroadcast(request AddBroadcastRequest) error {
	return nil
}

type EditBroadcastRequest struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	EnContent string `json:"en_content"`
	PtContent string `json:"pt_content"`
	Sort      int    `json:"sort"`
	Status    int    `json:"status"`
}

// EditBroadcast 编辑广播
func (s BroadcastService) EditBroadcast(request EditBroadcastRequest, adminId int64) error {
	err := s.validateEditBroadcast(request)
	if err != nil {
		return err
	}
	broadcast, err := s.GetBroadcastById(request.Id)
	if err != nil {
		logs.Error("获取广播标失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(broadcast.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	fields := []string{"name", "en_content", "pt_content", "sort", "status"}
	broadcastModel := models.CreateBroadcastModel()
	_, err = broadcastModel.Update(&models.Broadcast{
		Id:        request.Id,
		Name:      request.Name,
		EnContent: request.EnContent,
		PtContent: request.PtContent,
		Sort:      request.Sort,
		Status:    request.Status,
	}, fields...)
	if err != nil {
		return err
	}
	return nil
}

// 编辑广播参数验证
func (s BroadcastService) validateEditBroadcast(request EditBroadcastRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("GuangBoIDBiTian")
	}
	return nil
}

// DeleteBroadcast 删除广播
func (s BroadcastService) DeleteBroadcast(id int, adminId int64) error {
	broadcast, err := s.GetBroadcastById(id)
	if err != nil {
		logs.Error("获取广播标失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(broadcast.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	broadcastModel := models.CreateBroadcastModel()
	if id == 0 {
		return fmt.Errorf("GuangBoIDBiTian")
	}
	_, err = broadcastModel.Update(&models.Broadcast{
		Id:        id,
		IsDeleted: 1,
	}, "is_deleted")
	if err != nil {
		return err
	}
	return nil
}

type ChangeBroadcastStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"` // 状态
}

// ChangeBroadcastStatus 修改广播状态
func (s BroadcastService) ChangeBroadcastStatus(request ChangeBroadcastStatusRequest, adminId int64) error {
	broadcast, err := s.GetBroadcastById(request.Id)
	if err != nil {
		logs.Error("获取广播标失败，失败原因为：%v", err)
		return err
	}
	authGroupService := AuthGroupService{}
	hasAuthority := authGroupService.HasPackagePermission(broadcast.PackageId, adminId)
	if !hasAuthority {
		return fmt.Errorf("QuanXianBuZu")
	}
	broadcastModel := models.CreateBroadcastModel()
	if request.Id == 0 {
		return fmt.Errorf("GuangBoIDBiTian")
	}
	_, err = broadcastModel.Update(&models.Broadcast{
		Id:     request.Id,
		Status: request.Status,
	}, "status")
	if err != nil {
		return err
	}
	return nil
}
