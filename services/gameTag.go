package services

import (
	"api/models"
	"fmt"
	"net/url"
)

type GameTagService struct {
	BaseService
}

// GetAllGameTags 获取所有游戏标签
func (s GameTagService) GetAllGameTags(packageIds []int) []models.GameTag {
	var list []models.GameTag
	models.CreateGameTagModel().QueryTable(&models.GameTag{}).Filter("package_id__in", packageIds).
		GroupBy("name").
		OrderBy("id").
		All(&list)
	return list
}

type GameTagListRequest struct {
	Id         int        `form:"id"`
	Name       string     `form:"name" op:"like"`
	PtName     string     `form:"pt_name" op:"like"`
	Page       int        `form:"page"`
	PageSize   int        `form:"page_size"`
	RawQuery   url.Values `form:"-"`
	PackageId  int        `form:"package_id"`
	PackageIds []int
}

// GameTagList 游戏标签列表（带分页）
func (s GameTagService) GameTagList(request GameTagListRequest, userId int) (map[string]interface{}, error) {
	gameTagModel := models.CreateGameTagModel()
	condition, sort := gameTagModel.BuildCondition(request, "-id")

	condition = s.LimitPackageId(condition, request.PackageIds)

	data, total, err := gameTagModel.GetPageList(&models.GameTag{}, condition, request.Page, request.PageSize, sort)
	if nil != err {
		return nil, err
	}

	list := data.([]models.GameTag)
	return map[string]interface{}{
		"list":         list,
		"total":        total,
		"package_list": (&PackageService{}).GetMyAllPackageList(userId),
	}, nil
}

type AddGameTagRequest struct {
	Name      string `json:"name"`    // 标签名称
	PtName    string `json:"pt_name"` // 葡语名称
	Status    int    `json:"status"`  // 状态 1-开启 0-关闭
	PackageId int    `json:"package_id"`
}

// AddGameTag 添加游戏标签
func (s GameTagService) AddGameTag(request AddGameTagRequest) error {
	if err := s.validateAddGameTag(request); err != nil {
		return err
	}

	gameTagModel := models.CreateGameTagModel()
	_, err := gameTagModel.Insert(&models.GameTag{
		Name:      request.Name,
		PtName:    request.PtName,
		Status:    request.Status,
		PackageId: request.PackageId,
	})
	return err
}

// validateAddGameTag 验证添加参数
func (s GameTagService) validateAddGameTag(request AddGameTagRequest) error {
	if request.Name == "" {
		return fmt.Errorf("YouXiBiaoQianMingChengBiTian")
	}
	// 检查名称是否重复
	isExist := models.CreateGameTagModel().QueryTable(new(models.GameTag)).
		Filter("name", request.Name).
		Exist()
	if isExist {
		return fmt.Errorf("YouXiBiaoQianMingChengYiCunZai")
	}
	return nil
}

type EditGameTagRequest struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	PtName string `json:"pt_name"`
	Status int    `json:"status"`
}

// EditGameTag 编辑游戏标签
func (s GameTagService) EditGameTag(request EditGameTagRequest) error {
	if err := s.validateEditGameTag(request); err != nil {
		return err
	}
	gameTagModel := models.CreateGameTagModel()
	fields := []string{"name", "pt_name", "status"}
	_, err := gameTagModel.Update(&models.GameTag{
		Id:     request.Id,
		Name:   request.Name,
		PtName: request.PtName,
		Status: request.Status,
	}, fields...)
	return err
}

// validateEditGameTag 验证编辑参数
func (s GameTagService) validateEditGameTag(request EditGameTagRequest) error {
	if request.Id == 0 {
		return fmt.Errorf("YouXiBiaoQianIDBiTian")
	}
	if request.Name == "" {
		return fmt.Errorf("YouXiBiaoQianMingChengBiTian")
	}
	isExist := models.CreateGameTagModel().QueryTable(new(models.GameTag)).
		Filter("id__ne", request.Id).
		Filter("name", request.Name).
		Exist()
	if isExist {
		return fmt.Errorf("YouXiBiaoQianMingChengYiCunZai")
	}
	return nil
}

// DeleteGameTag 删除游戏标签
func (s GameTagService) DeleteGameTag(id int) error {
	if id == 0 {
		return fmt.Errorf("YouXiBiaoQianIDBiTian")
	}
	gameTagModel := models.CreateGameTagModel()
	_, err := gameTagModel.Delete(&models.GameTag{Id: id})
	return err
}
