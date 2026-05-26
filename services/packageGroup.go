package services

import (
	"api/models"
	"api/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type PackageGroupService struct{}

type PackageGroupListRequestParams struct {
	Title     string `json:"title"`
	Status    int    `json:"status"`
	PackageId int    `json:"package_id"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

type AddPackageGroupListRequestParams struct {
	Id             int    `json:"id,omitempty"`
	Title          string `json:"title"`
	Status         int    `json:"status"`
	IsModifyStatus bool   `json:"is_modify_status,omitempty"` // 是否修改状态
}

type GroupListResponse struct {
	Total       int                   `json:"total"`
	List        []models.PackageGroup `json:"list"`
	PackageList []models.Package      `json:"package_list"`
}

func (s *PackageGroupService) BuildCondition(params PackageGroupListRequestParams, userId int) map[string]interface{} {
	condition := make(map[string]interface{})

	if len(params.Title) > 0 {
		condition["title__icontains"] = params.Title
	}

	if params.PackageId > 0 {
		condition["package_ids__icontains"] = params.PackageId
	}

	if params.Status >= 0 {
		condition["status"] = params.Status
	}

	if userId > 0 {
		// 只展示授权的分包
		groupIds := (&AuthGroupService{}).GetPackageGroupIds(userId)
		if len(groupIds) == 0 { // 没有分配
			condition["id__lt"] = 0
		} else if groupIds != "-1" {
			condition["id__in"] = utils.StringSliceToInt(groupIds, ",")
		}
	}
	return condition
}

// 分组列表
func (s *PackageGroupService) GetList(params PackageGroupListRequestParams, userId int) GroupListResponse {
	condition := s.BuildCondition(params, userId)
	result, total, _ := models.CreatePackageGroupModel().GetPageList(new(models.PackageGroup), condition, params.Page, params.PageSize, "-id")
	groupList, ok := result.([]models.PackageGroup)
	if !ok {
		groupList = []models.PackageGroup{}
	}

	packageList, _ := (&PackageService{}).AllPackage()
	return GroupListResponse{
		Total:       int(total),
		List:        groupList,
		PackageList: packageList,
	}
}

// 创建分组
func (s *PackageGroupService) AddGroup(params AddPackageGroupListRequestParams, language string) error {
	model := models.CreatePackageGroupModel()

	title := strings.TrimSpace(params.Title)
	isExists := model.QueryTable(new(models.PackageGroup)).Filter("title", params.Title).Exist()
	if isExists {
		return fmt.Errorf("PingTaiFenZuCunZai")
	}

	currentTime := time.Now().Unix()
	_, err := model.Insert(&models.PackageGroup{
		Title:      title,
		Status:     1,
		PackageIds: "",
		CreateTime: currentTime,
		UpdateTime: currentTime,
	})
	if err != nil {
		logs.Error("创建平台分组失败：%v", err)
		return fmt.Errorf("ChuangJianShuJuShiBai")
	}
	return nil
}

// 移除分组的分包id
func (s *PackageGroupService) RemovePackageId(txOrm orm.TxOrmer, groupId int, packageId int) error {
	var group models.PackageGroup
	err := txOrm.QueryTable(new(models.PackageGroup)).Filter("id", groupId).One(&group)
	if err != nil {
		return err
	}
	fields := []string{
		"package_ids",
	}
	packageIds := strings.Split(strings.Trim(group.PackageIds, ","), ",")
	var ids []int
	for _, item := range packageIds {
		id, _ := strconv.Atoi(item)
		if id != packageId && id != 0 {
			ids = append(ids, id)
		}
	}

	packageIdsStr := utils.IntSliceToString(ids, ",")
	if len(packageIdsStr) == 0 {
		packageIdsStr = ""
	} else {
		packageIdsStr = fmt.Sprintf(",%v,", utils.IntSliceToString(ids, ","))
	}

	_, updateGroup := txOrm.Update(&models.PackageGroup{
		Id:         groupId,
		PackageIds: packageIdsStr,
	}, fields...)
	if updateGroup != nil {
		return updateGroup
	}

	return nil
}

// 修改分包的分组id
func (s *PackageGroupService) ChangeGroupId(txOrm orm.TxOrmer, packageId int, groupId int) error {
	// 获取之前分包的分组
	var packageInfo models.Package
	err := txOrm.QueryTable(new(models.Package)).Filter("id", packageId).One(&packageInfo)
	if err != nil {
		return err
	}

	var packageGroupInfo models.PackageGroup
	groupErr := txOrm.QueryTable(new(models.PackageGroup)).Filter("id", groupId).One(&packageGroupInfo)
	if groupErr != nil {
		return groupErr
	}

	fields := []string{
		"package_ids",
	}
	if packageInfo.GroupId == 0 {
		// 未分配 直接更新分组的package_ids
		packageIds := packageGroupInfo.PackageIds
		if len(packageIds) == 0 {
			packageIds = fmt.Sprintf(",%v,", packageId)
		} else {
			packageIds = fmt.Sprintf(",%v,%v,", packageIds, packageId)
		}
		_, updateErr := txOrm.Update(&models.PackageGroup{
			Id:         groupId,
			PackageIds: packageIds,
		}, fields...)
		if updateErr != nil {
			return updateErr
		}
	} else {
		// 已分配 先去除旧的分组信息 在更新新的分组信息
		var oldPackageGroupInfo models.PackageGroup
		groupErr := txOrm.QueryTable(new(models.PackageGroup)).Filter("id", packageInfo.GroupId).One(&oldPackageGroupInfo)
		if groupErr != nil {
			return groupErr
		}

		packageIds := strings.Split(strings.Trim(oldPackageGroupInfo.PackageIds, ","), ",")
		var ids []int
		for _, item := range packageIds {
			if len(item) > 0 {
				id, _ := strconv.Atoi(item)
				if id != packageId {
					ids = append(ids, id)
				}
			}
		}

		packageIdsStr := utils.IntSliceToString(ids, ",")
		if len(packageIdsStr) == 0 {
			packageIdsStr = ""
		} else {
			packageIdsStr = fmt.Sprintf(",%v,", packageIdsStr)
		}
		_, updateGroup := txOrm.Update(&models.PackageGroup{
			Id:         packageInfo.GroupId,
			PackageIds: packageIdsStr,
		}, fields...)
		if updateGroup != nil {
			return updateGroup
		}
		packageIdsStr = packageGroupInfo.PackageIds
		if len(packageIdsStr) == 0 {
			packageIdsStr = fmt.Sprintf(",%v,", packageId)
		} else {
			packageIdsStr = fmt.Sprintf("%v%v,", packageGroupInfo.PackageIds, packageId)
		}

		_, updateGroup1 := txOrm.Update(&models.PackageGroup{
			Id:         groupId,
			PackageIds: packageIdsStr,
		}, fields...)
		if updateGroup1 != nil {
			return updateGroup1
		}
	}
	return nil
}

// 修改分组
func (s *PackageGroupService) EditGroup(params AddPackageGroupListRequestParams, language string) error {
	model := models.CreatePackageGroupModel()

	title := strings.TrimSpace(params.Title)
	isExists := model.QueryTable(new(models.PackageGroup)).Filter("id__ne", params.Id).Filter("title", params.Title).Exist()
	if isExists {
		return fmt.Errorf("PingTaiFenZuCunZai")
	}

	currentTime := time.Now().Unix()

	fields := []string{
		"title", "package_ids", "update_time",
	}
	if params.IsModifyStatus {
		fields = []string{
			"status", "update_time",
		}
	}

	_, err := model.Update(&models.PackageGroup{
		Id:         params.Id,
		Title:      title,
		Status:     params.Status,
		PackageIds: "",
		UpdateTime: currentTime,
	}, fields...)
	if err != nil {
		logs.Error("修改平台分组失败：%v", err)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

// 删除分组
func (s *PackageGroupService) DelGroup(ids string) error {
	model := models.CreatePackageGroupModel()

	idArray := strings.Split(ids, ",")
	idInts := make([]int, len(idArray))
	for i, id := range idArray {
		idInts[i], _ = strconv.Atoi(id)
	}

	isExists := model.QueryTable(new(models.PackageGroup)).Filter("id__in", idInts).Filter("package_ids__ne", "").Exist()
	if isExists {
		return fmt.Errorf("QingXianJieBangPingTai")
	}
	_, err := model.QueryTable(new(models.PackageGroup)).
		Filter("id__in", idInts).
		Delete()
	if err != nil {
		logs.Error("删除平台失败:%v", err)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

// GetUserList 用户列表
func (s *PackageGroupService) GetUserList(params UserRequestParams) UserListResponse {
	condition := make(map[string]interface{})
	condition["is_deleted"] = 0

	if params.PackageId > 0 {
		condition["package_id"] = params.PackageId
	}

	if len(params.Username) > 0 {
		condition["username__icontains"] = params.Username
	}
	if params.RoleId > 0 {
		condition["role_id"] = params.RoleId
	}

	result, total, _ := models.CreateUserModel().GetPageList(new(models.User), condition, params.Page, params.PageSize, "-id")
	userList, ok := result.([]models.User)
	if !ok {
		userList = []models.User{}
	}

	return UserListResponse{
		Total: int(total),
		Lists: userList,
	}
}

// 获取所有的分组
func (s *PackageGroupService) GetAllGroup(groupIds string) []models.PackageGroup {
	conditon := map[string]interface{}{
		"Status": 1,
	}
	// 0获取所有的分组
	if len(groupIds) > 0 && groupIds != "-1" {
		conditon["id__in"] = groupIds
	}
	var lists []models.PackageGroup
	model := models.CreatePackageGroupModel()
	db := model.QueryTable(new(models.PackageGroup))
	_, err := model.Where(db, conditon).All(&lists)
	if err == nil {
		return lists
	}
	return []models.PackageGroup{}
}
