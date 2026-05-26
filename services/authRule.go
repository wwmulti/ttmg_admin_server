package services

import (
	"api/models"
	"fmt"
	"sort"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

type AuthAddRequestParams struct {
	Id     int    `json:"id,omitempty"`
	Icon   string `json:"icon,omitempty"`
	IsMenu int    `json:"ismenu"`
	Name   string `json:"name"`
	Pid    int    `json:"pid"`
	Title  string `json:"title"`
	Weigh  int    `json:"weigh"`
	Status int    `json:"status,omitempty"`
}

// Menu 菜单结构体
type Menu struct {
	ID       int    `json:"id"`
	PID      int    `json:"pid"`
	Name     string `json:"name"`
	Sort     int    `json:"sort"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Tag      int    `json:"tag"`
	Level    int    `json:"level"`
	Icon     string `json:"icon"`
	Label    string `json:"label"` // 授权页面标题
	Children []Menu `json:"children,omitempty"`
}

type AuthDelRequestParams struct {
	Id int `json:"id"`
}

type AuthRuleService struct{}

// 删除权限菜单
func (s *AuthRuleService) Authdel(params AuthDelRequestParams) error {
	model := models.CreateAuthRuleModel()
	_, err := model.Delete(&models.AuthRule{
		Id: params.Id,
	})
	if err != nil {
		return fmt.Errorf("ShuJuYiChang")
	}
	return nil
}

// 编辑权限菜单
func (s *AuthRuleService) Authedit(params AuthAddRequestParams) error {
	model := models.CreateAuthRuleModel()
	params.Title = strings.TrimSpace(params.Title)
	params.Name = strings.TrimSpace(params.Name)

	isExists := model.QueryTable(new(models.AuthRule)).Filter("id", params.Id).Exist()
	if !isExists {
		return fmt.Errorf("ShuJuYiChang")
	}

	status := 0
	if params.Status == 1 {
		status = 1
	}
	fields := []string{
		"title",
		"name",
		"status",
		"sort",
		"icon",
		"pid",
		"tag",
	}
	data := &models.AuthRule{
		Id:     params.Id,
		Title:  params.Title,
		Name:   params.Name,
		Status: status,
		Tag:    params.IsMenu,
		Pid:    params.Pid,
		Icon:   params.Icon,
		Sort:   params.Weigh,
	}
	_, err := models.CreateAuthRuleModel().Update(data, fields...)
	if err != nil {
		logs.Error("创建菜单数据失败:%v", err)
		return fmt.Errorf("ChuangJianShuJuShiBai")
	}
	return nil
}

func (s *AuthRuleService) Authadd(params AuthAddRequestParams) error {
	// 验证title
	model := models.CreateAuthRuleModel()
	params.Title = strings.TrimSpace(params.Title)
	params.Name = strings.TrimSpace(params.Name)

	isExists := model.QueryTable(new(models.AuthRule)).Filter("title", params.Title).Exist()
	if isExists {
		return fmt.Errorf("LuYouDiZhiYiPeiZhi")
	}

	data := models.AuthRule{
		Title:  params.Title,
		Name:   params.Name,
		Status: params.Status,
		Tag:    params.IsMenu,
		Pid:    params.Pid,
		Icon:   params.Icon,
		Sort:   params.Weigh,
	}
	_, err := models.CreateAuthRuleModel().Insert(&data)
	if err != nil {
		logs.Error("创建菜单数据失败:%v", err)
		return fmt.Errorf("ChuangJianShuJuShiBai")
	}
	return nil
}

// 获取后台菜单规则
type AuthRuleListResult struct {
	MenuList    []Menu `json:"menuList"`
	MenuOptions []Menu `json:"menuOptions"`
}

func (s *AuthRuleService) GetAuthRuleListForGroup(condition map[string]interface{}) []Menu {
	menus := s.GetAuthList(condition)
	menuList := s.convertToMenus(menus)

	if menuList != nil {
		menuList = s.TidyMenuTier(menuList, 0)
		sort.Slice(menuList, func(i, j int) bool {
			return menuList[i].Sort < menuList[j].Sort
		})
	}
	return menuList
}

// 权限管理菜单
func (s *AuthRuleService) GetAuthRuleList() AuthRuleListResult {
	condition := make(map[string]interface{})
	menus := s.GetAuthList(condition)

	menuList := s.convertToMenus(menus)
	menuOptions := s.array2Level(menuList, 0, 1)

	if menuList != nil {
		menuList = s.TidyMenuTier(menuList, 0)
		sort.Slice(menuList, func(i, j int) bool {
			return menuList[i].Sort < menuList[j].Sort
		})
	}
	return AuthRuleListResult{
		MenuList:    menuList,
		MenuOptions: menuOptions,
	}
}

func (s *AuthRuleService) array2Level(array []Menu, pid int, level int) []Menu {
	result := make([]Menu, 0)

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

// 获取前端菜单
func (s *AuthRuleService) GetMenusTreeByCache(uid int) []Menu {
	menuList := s.GetMenusTreee(1)
	menuList = s.FilterMenus(menuList, uid)
	return menuList
}

// 后台菜单树
func (s *AuthRuleService) GetMenusTreee(status int) []Menu {
	condition := make(map[string]interface{})
	condition["status"] = status
	condition["tag"] = 1
	menus := s.GetAuthList(condition)
	menuList := s.convertToMenus(menus)

	if menus != nil {
		menuList = s.TidyMenuTier(menuList, 0)
		sort.Slice(menuList, func(i, j int) bool {
			return menuList[i].Sort < menuList[j].Sort
		})
	}
	return menuList
}

// 过滤掉没有权限的菜单
func (s *AuthRuleService) FilterMenus(menuList []Menu, uid int) []Menu {
	groupId := (&AccountService{}).GetGroupId(uid)
	var menus []Menu
	for _, menu := range menuList {
		if (&Authservice{}).Check(menu.ID, groupId) {
			menus = append(menus, menu)
		}
	}
	return menus
}

// ConvertToMenus 将AuthRule切片转换为Menu切片
func (s *AuthRuleService) convertToMenus(rules []models.AuthRule) []Menu {
	var menus []Menu

	if rules == nil {
		return menus
	}

	for _, rule := range rules {
		menu := Menu{
			ID:     rule.Id,
			PID:    rule.Pid,
			Name:   rule.Name,
			Label:  rule.Name,
			Sort:   rule.Sort,
			Title:  rule.Title,
			Status: rule.Status,
			Icon:   rule.Icon,
			Tag:    rule.Tag,
			Level:  1,
		}
		menus = append(menus, menu)
	}

	return menus
}

// 后台菜单
func (s *AuthRuleService) GetAuthList(condition map[string]interface{}) []models.AuthRule {
	model := models.CreateAuthRuleModel()
	var menus []models.AuthRule

	dbObj := model.QueryTable(new(models.AuthRule))
	dbObj = model.Where(dbObj, condition)

	_, err := dbObj.All(&menus)
	if err != nil {
		return nil
	}
	return menus
}

// TidyMenuTier 将菜单列表转换为层级结构
func (s *AuthRuleService) TidyMenuTier(menusList []Menu, pid int) []Menu {
	var navList []Menu

	for _, menu := range menusList {
		if menu.PID == pid {
			// 递归获取子菜单
			menu.Children = s.TidyMenuTier(menusList, menu.ID)
			navList = append(navList, menu)
		}
	}
	return navList
}

// 获取菜单id
func (s *AuthRuleService) GetMenuId(title string) int {
	var info models.AuthRule
	model := models.CreateAuthRuleModel()
	model.QueryTable(new(models.AuthRule)).Filter("title", title).One(&info)
	return info.Id
}
