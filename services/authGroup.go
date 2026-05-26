package services

import (
	"api/models"
	"api/utils"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

// Menu 菜单结构体
type GroupItem struct {
	ID              int         `json:"id"`
	PID             int         `json:"pid"`
	Title           string      `json:"title"`
	Status          int         `json:"status"`
	Level           int         `json:"level"`
	PackageGroupIds string      `json:"package_group_ids"`
	PackageIds      string      `json:"package_ids"`
	Children        []GroupItem `json:"children,omitempty"`
}

// 角色列表
type AuthGroupListResult struct {
	RoleList         []GroupItem              `json:"role_list"`          // 角色列表
	RoleOptions      []GroupItem              `json:"role_options"`       // 角色选择列表
	PackageGroupList []models.PackageGroup    `json:"package_group_list"` // 平台分组选择列表
	PackageList      map[int][]models.Package `json:"package_list"`       // 平台选择列表
}

// 创建角色参数
type AddRoleRequestParams struct {
	Id              int    `json:"id,omitempty"`
	Pid             int    `json:"pid"`
	Status          int    `json:"status"`
	Title           string `json:"title"`
	PackageGroupIds string `json:"group_ids"`   // 分包分组id
	PackageIds      string `json:"package_ids"` // 分包id
}

type RolerulesResponse struct {
	AuthIds []string `json:"auth_ids"`
	List    []Menu   `json:"list"`
}

type RoleruleseditParams struct {
	Id      int    `json:"id"`
	RuleIds string `json:"auth_rule_ids"`
}

type AuthGroupService struct{}

// 授权
func (s *AuthGroupService) Rolerulesedit(params RoleruleseditParams, uid int) error {
	status := s.IsMyGroup(uid, params.Id)
	if !status {
		return fmt.Errorf("ShuJuYiChang")
	}
	// 判断ruleId
	ruleId := strings.Split(params.RuleIds, ",")
	condition := make(map[string]interface{})
	condition["id__in"] = ruleId
	ruleList := (&AuthRuleService{}).GetAuthList(condition)
	if len(ruleList) != len(ruleId) {
		return fmt.Errorf("ShuJuYiChang")
	}
	model := models.CreateAuthGroupModel()
	fields := []string{
		"rules",
	}
	_, err := model.Update(&models.AuthGroup{
		Id:    params.Id,
		Rules: params.RuleIds,
	}, fields...)
	if err != nil {
		logs.Error("授权失败:%v", err)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

// 获取授权的菜单
func (s *AuthGroupService) Rolerules(groupId int) RolerulesResponse {
	model := models.CreateAuthGroupModel()
	var info models.AuthGroup
	err := model.QueryTable(new(models.AuthGroup)).Filter("id", groupId).One(&info)
	if err != nil {
		return RolerulesResponse{}
	}

	authIds := s.getAuthIds(info)
	list := s.getParentRule(info)
	return RolerulesResponse{
		AuthIds: authIds,
		List:    list,
	}
}

// 获取父级的授权菜单
func (s *AuthGroupService) getParentRule(info models.AuthGroup) []Menu {
	condition := make(map[string]interface{})

	var authRuleList []Menu
	if info.Pid == 0 {
		// 超级 管理员获取全部节点
		authRuleList = (&AuthRuleService{}).GetAuthRuleListForGroup(condition)
	} else {
		model := models.CreateAuthGroupModel()
		var groupInfo models.AuthGroup
		err := model.QueryTable(new(models.AuthGroup)).Filter("id", info.Pid).One(&groupInfo)
		if err != nil {
			return authRuleList
		}
		if groupInfo.Pid == 0 {
			// 超级管理员
			authRuleList = (&AuthRuleService{}).GetAuthRuleListForGroup(condition)
		} else {
			ruleIds := strings.Split(groupInfo.Rules, ",")
			condition["id__in"] = ruleIds
			authRuleList = (&AuthRuleService{}).GetAuthRuleListForGroup(condition)
		}
	}

	return authRuleList
}

// 获取授权的菜单id
func (s *AuthGroupService) getAuthIds(info models.AuthGroup) []string {
	condition := make(map[string]interface{})
	var authRuleList []models.AuthRule
	if info.Pid == 0 {
		// 超级管理员
		authRuleList = (&AuthRuleService{}).GetAuthList(condition)
	} else {
		if len(info.Rules) > 0 {
			ruleIds := strings.Split(info.Rules, ",")
			condition["id__in"] = ruleIds
			authRuleList = (&AuthRuleService{}).GetAuthList(condition)
		}
	}

	return s.convertAuthRule(authRuleList)
}

func (s *AuthGroupService) convertAuthRule(authRuleList []models.AuthRule) []string {
	list := make([]string, 0, len(authRuleList))
	for _, rule := range authRuleList {
		list = append(list, strconv.Itoa(rule.Id))
	}
	return list
}

// 删除角色
func (s *AuthGroupService) DelRole(groupId, uid int) error {
	myGroupId := (&AccountService{}).GetGroupId(uid)
	if myGroupId == groupId {
		return fmt.Errorf("BuNengShanChuZiJi")
	}
	childStatus := s.IsExistsChild(groupId)
	if childStatus {
		return fmt.Errorf("CunZaiZiJi")
	}
	groupStatus := s.IsMyGroup(uid, groupId)
	if !groupStatus {
		return fmt.Errorf("ShuJuYiChang")
	}

	model := models.CreateAuthGroupModel()
	_, err := model.Delete(&models.AuthGroup{
		Id: groupId,
	})
	if err != nil {
		logs.Error("删除角色失败:%v", err)
		return fmt.Errorf("ShuJuShanChuShiBai")
	}
	return nil
}

// 修改角色
func (s *AuthGroupService) EditRole(params AddRoleRequestParams, userId int64) error {
	model := models.CreateAuthGroupModel()
	// 验证平台分组id和分包id是否是自己明下子级
	if !s.isMyPackage(int(userId), params.PackageGroupIds, params.PackageIds) {
		return fmt.Errorf("ShuJuYiChang")
	}

	var info models.AuthGroup
	err := model.QueryTable(new(models.AuthGroup)).Filter("id", params.Id).One(&info)
	if err != nil {
		return fmt.Errorf("ShuJuYiChang")
	}

	if info.Id == params.Pid {
		return fmt.Errorf("ShangJiBuNengXuanZeZiJi")
	}

	childStatus := s.IsExistsChild(params.Id)
	if childStatus {
		return fmt.Errorf("CunZaiZiJi")
	}
	rules := info.Rules
	if info.Pid != params.Pid {
		rules = ""
	}

	fields := []string{
		"title",
		"status",
		"rules",
		"pid",
		"package_group_ids", "package_ids",
	}
	_, dataErr := model.Update(&models.AuthGroup{
		Id:              params.Id,
		Title:           params.Title,
		Status:          params.Status,
		Pid:             params.Pid,
		Rules:           rules,
		PackageGroupIds: params.PackageGroupIds,
		PackageIds:      params.PackageIds,
	}, fields...)
	if dataErr != nil {
		logs.Error("更新角色失败:%v", dataErr)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

// 验证分包信息
func (s *AuthGroupService) isMyPackage(userId int, groupIds string, packageIds string) bool {
	// 验证分组
	packageGroupIds, packageIdLists := s.getPackageGroupInfo(userId)
	if packageGroupIds == "-1" { // 无限制
		return true
	}
	packageGroupList := (&PackageGroupService{}).GetAllGroup(packageGroupIds)
	groupIdList := utils.StringSliceToInt(groupIds, ",")

	isGroupIdExits := false
	for _, packageGroup := range packageGroupList {
		for _, groupId := range groupIdList {
			if groupId == packageGroup.Id {
				isGroupIdExits = true
			}
		}

		if !isGroupIdExits {
			return false
		}
	}
	// 验证分包
	packageLists := s.getPackageList(packageGroupList, packageIdLists)
	packageIdMap := s.packageExistsMap(packageLists)
	packageIdList := utils.StringSliceToInt(packageIds, ",")

	for _, packageId := range packageIdList {
		if !packageIdMap[packageId] {
			return false
		}
	}

	return true
}

// 分包id map
func (s *AuthGroupService) packageExistsMap(packageLists map[int][]models.Package) map[int]bool {
	existsMap := make(map[int]bool)
	for _, packages := range packageLists {
		for _, pkg := range packages {
			existsMap[pkg.Id] = true
		}
	}
	return existsMap
}

// 创建角色
func (s *AuthGroupService) AddRole(params AddRoleRequestParams, userId int64) error {
	model := models.CreateAuthGroupModel()
	// 验证平台分组id和分包id是否是自己明下子级
	if !s.isMyPackage(int(userId), params.PackageGroupIds, params.PackageIds) {
		return fmt.Errorf("ShuJuYiChang")
	}
	err := model.OrmerMaster.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		id, err := txOrm.Insert(&models.AuthGroup{
			Title:           params.Title,
			Status:          params.Status,
			Pid:             params.Pid,
			Rules:           "",
			Parents:         "",
			PackageGroupIds: params.PackageGroupIds,
			PackageIds:      params.PackageIds,
		})
		if err != nil {
			return err
		}
		parents, strErr := s.CreateParents(txOrm, params.Pid, int(id))
		if strErr != nil {
			return strErr
		}
		fields := []string{
			"parents",
		}
		_, updateErr := txOrm.Update(&models.AuthGroup{
			Id:      int(id),
			Parents: parents,
		}, fields...)
		if updateErr != nil {
			return updateErr
		}
		return nil
	})
	if err != nil {
		logs.Error("创建角色失败:%v", err)
		return fmt.Errorf("ChuangJianShuJuShiBai")
	}

	return nil
}

func (s *AuthGroupService) CreateParents(txOrm orm.TxOrmer, pid, id int) (string, error) {
	var info models.AuthGroup
	err := txOrm.QueryTable(new(models.AuthGroup)).Filter("id", pid).One(&info)
	if err != nil {
		logs.Error("获取角色父级pid:%v信息失败:%v", pid, err)
		return "", err
	}
	return fmt.Sprintf("%s%d,", info.Parents, id), nil
}

func (s *AuthGroupService) RoleList(uid int) AuthGroupListResult {
	menus := s.GetGroupList(uid)

	roleList := s.convertToMenus(menus)
	roleOptions := s.array2Level(roleList, roleList[0].PID, 1)

	if roleList != nil {
		roleList = s.tidyGroupTier(roleList, roleList[0].PID)
		sort.Slice(roleList, func(i, j int) bool {
			return roleList[i].PID < roleList[j].PID
		})
	}
	// 获取自己的平台分组id
	packageGroupIds, packageIds := s.getPackageGroupInfo(uid)
	packageGroupList := (&PackageGroupService{}).GetAllGroup(packageGroupIds)

	return AuthGroupListResult{
		RoleList:         roleList,
		RoleOptions:      roleOptions,
		PackageGroupList: packageGroupList,
		PackageList:      s.getPackageList(packageGroupList, packageIds),
	}
}

// 获取指定的分组的分包
func (s *AuthGroupService) getPackageList(packageGroupList []models.PackageGroup, packageIds string) map[int][]models.Package {
	lists := make(map[int][]models.Package)
	packageIdList := utils.StringSliceToInt(packageIds, ",")
	for _, item := range packageGroupList {
		if len(item.PackageIds) > 0 {
			ids := utils.StringSliceToInt(item.PackageIds, ",")
			if packageIds != "-1" { //-1表示无限制
				ids = utils.IntersectSlice(ids, packageIdList)
			}
			if len(ids) > 0 {
				packageList, _ := (&PackageService{}).AllPackage(ids)
				lists[item.Id] = packageList
			}
		}
	}
	return lists
}

// 获取平台分组id
func (s *AuthGroupService) getPackageGroupInfo(uid int) (string, string) {
	groupId := (&AccountService{}).GetGroupId(uid)
	var info models.AuthGroup
	model := models.CreateAuthGroupModel()
	model.QueryTable(new(models.AuthGroup)).Filter("Id", groupId).One(&info)
	return info.PackageGroupIds, info.PackageIds
}

func (s *AuthGroupService) tidyGroupTier(menusList []GroupItem, pid int) []GroupItem {
	var navList []GroupItem

	for _, menu := range menusList {
		if menu.PID == pid {
			// 递归获取子菜单
			menu.Children = s.tidyGroupTier(menusList, menu.ID)
			navList = append(navList, menu)
		}
	}
	return navList
}

func (s *AuthGroupService) array2Level(array []GroupItem, pid int, level int) []GroupItem {
	result := make([]GroupItem, 0)

	for _, v := range array {
		if v.PID == pid {
			// 设置层级
			v.Level = level
			result = append(result, v)

			// 递归查找子节点
			children := s.array2Level(array, v.ID, level+1)
			result = append(result, children...)
		}
	}

	return result
}

// ConvertToMenus 将AuthRule切片转换为Menu切片
func (s *AuthGroupService) convertToMenus(groupList []models.AuthGroup) []GroupItem {
	var menus []GroupItem

	if groupList == nil {
		return menus
	}

	for _, group := range groupList {
		menu := GroupItem{
			ID:              group.Id,
			PID:             group.Pid,
			Title:           group.Title,
			Status:          group.Status,
			PackageGroupIds: group.PackageGroupIds,
			PackageIds:      group.PackageIds,
			Level:           1,
		}
		menus = append(menus, menu)
	}

	return menus
}

// 获取自己的权限组
func (s *AuthGroupService) GetGroupList(uid int) []models.AuthGroup {
	var menus []models.AuthGroup

	groupId := (&AccountService{}).GetGroupId(uid)

	model := models.CreateAuthGroupModel()
	_, groupErr := model.QueryTable(new(models.AuthGroup)).Filter("parents__contains", fmt.Sprintf(",%d,", groupId)).All(&menus)
	if groupErr != nil {
		return nil
	}
	return menus
}

// 是否是自己的权限组
func (s *AuthGroupService) IsMyGroup(uid, groupId int) bool {
	myGroupId := (&AccountService{}).GetGroupId(uid)
	if myGroupId == groupId {
		return true
	}
	model := models.CreateAuthGroupModel()
	status := model.QueryTable(new(models.AuthGroup)).Filter("parents__contains", fmt.Sprintf(",%d,", myGroupId)).Exist()
	return status
}

// 该群组是否存在子级
func (s *AuthGroupService) IsExistsChild(groupId int) bool {
	model := models.CreateAuthGroupModel()
	status := model.QueryTable(new(models.AuthGroup)).Filter("pid", groupId).Exist()
	return status
}

// 获取自己的分包id
func (s *AuthGroupService) GetPackageIds(userId int) string {
	myGroupId := (&AccountService{}).GetGroupId(userId)
	var info models.AuthGroup
	model := models.CreateAuthGroupModel()
	model.QueryTable(new(models.AuthGroup)).Filter("Id", myGroupId).One(&info)
	return info.PackageIds
}

// 获取自己的分包的平台分组id
func (s *AuthGroupService) GetPackageGroupIds(userId int) string {
	myGroupId := (&AccountService{}).GetGroupId(userId)
	var info models.AuthGroup
	model := models.CreateAuthGroupModel()
	model.QueryTable(new(models.AuthGroup)).Filter("Id", myGroupId).One(&info)
	return info.PackageGroupIds
}

// GetAdminIsRootAndPackageIdSlice 获取管理员是否为超级管理员以及包分片
func (s *AuthGroupService) GetAdminIsRootAndPackageIdSlice(userId int64) (bool, []int) {
	authGroupService := AuthGroupService{}
	packageIds := authGroupService.GetPackageIds(int(userId))
	var packageSlice []int
	var isRootAdmin bool
	if packageIds != "-1" {
		isRootAdmin = false
		packageSlice = utils.StringToIntSlice(packageIds)
	} else {
		isRootAdmin = true
	}
	return isRootAdmin, packageSlice
}

// HasPackagePermission 是否有指定包的权限
func (s *AuthGroupService) HasPackagePermission(packageId int, adminId int64) bool {
	isRootAdmin, packageSlice := s.GetAdminIsRootAndPackageIdSlice(adminId)
	if !isRootAdmin && !utils.InArray(packageId, packageSlice) {
		return false
	}
	return true
}

// HasPackagesPermission 是否有指定包数组的权限
func (s *AuthGroupService) HasPackagesPermission(packageIds []int, adminId int64) bool {
	isRootAdmin, packageSlice := s.GetAdminIsRootAndPackageIdSlice(adminId)
	for _, packageId := range packageIds {
		if !isRootAdmin && !utils.InArray(packageId, packageSlice) {
			return false
		}
	}
	return true
}
