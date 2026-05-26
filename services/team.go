package services

import (
	"api/models"
	"context"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type TeamService struct{}

type TeamListRequestParams struct {
	Title      string `json:"title"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	PackageId  int    `json:"package_id"`
	Pid        int    `json:"pid"`
	NeedReload int    `json:"need_reload"`
}

type TeamListResponse struct {
	Total       int              `json:"total"`
	Lists       []models.Team    `json:"lists"`
	PackageList []models.Package `json:"package_list"`
	TeamList    []models.Team    `json:"team_list"` // 团队
}

type TeamRequestParams struct {
	Id         int    `json:"id,omitempty"`
	Title      string `json:"title"`
	PackageId  int    `json:"package_id"`
	Pid        int    `json:"pid"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	IsOfficial int    `json:"is_official"`
	IsBloger   int    `json:"is_bloger"`
	IsDropUser int    `json:"is_drop_user"`
	Language   string
	Token      string `json:"token"`
	Secret     string `json:"secret"`
	Rtp        int    `json:"rtp"`
}

type SubTeamWithParent struct {
	Pid          int    `json:"pid"`            // 父级id
	Title        string `json:"title"`          // 业务员
	RoleId       int    `json:"role_id"`        // 业务员ID
	ParentTitle  string `json:"parent_title"`   // 团队
	ParentRoleId int    `json:"parent_role_id"` // 团队ID
}

// 获取业务员团队列表
func (s *TeamService) GetSalespersonList(params TeamListRequestParams) TeamListResponse {
	condition := make(map[string]interface{})
	condition["is_deleted"] = 0

	if len(params.Title) > 0 {
		condition["title__contains"] = params.Title
	}

	if params.PackageId > 0 {
		condition["package_id"] = params.PackageId
	}

	if params.Pid == -1 {
		condition["pid__gt"] = 0
	} else { // 查询指定pid
		condition["pid"] = params.Pid
	}

	result, total, _ := models.CreateTeamModel().GetPageList(new(models.Team), condition, params.Page, params.PageSize, "-id")
	teamList, ok := result.([]models.Team)
	if !ok {
		teamList = []models.Team{}
	}

	packageList := []models.Package{}
	parentTeamList := []models.Team{}
	if params.NeedReload > 0 {
		params.Pid = -1
		params.Page = 1
		params.PageSize = 10000000
		result := s.GetList(params)
		packageList = result.PackageList
		parentTeamList = result.Lists
	}

	return TeamListResponse{
		Total:       int(total),
		Lists:       teamList,
		PackageList: packageList,
		TeamList:    parentTeamList,
	}
}

// 获取团队列表
func (s *TeamService) GetList(params TeamListRequestParams) TeamListResponse {
	condition := make(map[string]interface{})
	condition["is_deleted"] = 0

	if len(params.Title) > 0 {
		condition["title__contains"] = params.Title
	}

	if params.PackageId > 0 {
		condition["package_id"] = params.PackageId
	}
	if params.Pid == -1 { // 团队
		condition["pid"] = 0
	} else { // 查询指定pid
		condition["pid"] = params.Pid
	}

	result, total, _ := models.CreateTeamModel().GetPageList(new(models.Team), condition, params.Page, params.PageSize, "-id")
	teamList, ok := result.([]models.Team)
	if !ok {
		teamList = []models.Team{}
	}

	var packageList []models.Package
	if params.NeedReload > 0 {
		packageList, _ = (&PackageService{}).AllPackage()
	}

	return TeamListResponse{
		Total:       int(total),
		Lists:       teamList,
		PackageList: packageList,
	}
}

// 创建团队
func (s *TeamService) AddTeam(params TeamRequestParams) error {
	packageId := params.PackageId

	isOfficial := params.IsOfficial
	isBloger := params.IsBloger
	isDropUser := params.IsDropUser
	model := models.CreateTeamModel()
	parents := ""
	if params.Pid != 0 {
		var info models.Team
		err := model.QueryTable(new(models.Team)).Filter("id", params.Pid).Filter("is_deleted", 0).One(&info)
		if err != nil {
			return fmt.Errorf("ShuJuYiChang")
		}
		packageId = info.PackageId
		isOfficial = info.IsOfficial
		isBloger = info.IsBloger
		isDropUser = info.IsDropUser
		parents = info.Parents
	}
	userParams := AddUserRequestParams{
		Username:  params.Username,
		Password:  params.Password,
		PackageId: packageId,
		Language:  params.Language,
	}
	userErr, roleId, userId := (&UserService{}).AddUser(userParams)
	if userErr != nil {
		return userErr
	}
	trErr := model.OrmerMaster.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		id, err := txOrm.Insert(&models.Team{
			Title:      params.Title,
			Username:   params.Username,
			PackageId:  packageId,
			IsOfficial: isOfficial,
			IsBloger:   isBloger,
			IsDropUser: isDropUser,
			Pid:        params.Pid,
			RoleId:     roleId,
			UserId:     userId,
			Token:      params.Token,
			Secret:     params.Secret,
			Rtp:        params.Rtp,
		})
		if err != nil {
			return err
		}

		if params.Pid > 0 {
			parents = fmt.Sprintf("%v%v,", parents, id)
		} else {
			parents = fmt.Sprintf(",%v,", id)
		}

		fields := []string{
			"parents",
		}
		_, updateErr := txOrm.Update(&models.Team{
			Id:      int(id),
			Parents: parents,
		}, fields...)
		if updateErr != nil {
			return updateErr
		}
		return nil
	})

	if trErr != nil {
		logs.Error("创建团队失败:%v", trErr)
		return fmt.Errorf("ChuangJianShuJuShiBai")
	}
	return nil
}

// 修改团队
func (s *TeamService) EditTeam(params TeamRequestParams) error {
	model := models.CreateTeamModel()
	field := []string{
		"package_id",
		"is_official",
		"is_bloger",
		"is_drop_user",
		"title",
		"token",
		"secret",
		"rtp",
	}

	if params.Pid > 0 {
		// 业务员修改
		field = []string{
			"title",
		}
	}
	trErr := model.OrmerMaster.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		_, err := txOrm.Update(&models.Team{
			Id:         params.Id,
			Title:      params.Title,
			PackageId:  params.PackageId,
			IsOfficial: params.IsOfficial,
			IsBloger:   params.IsBloger,
			IsDropUser: params.IsDropUser,
			Token:      params.Token,
			Secret:     params.Secret,
			Rtp:        params.Rtp,
		}, field...)
		if err != nil {
			return err
		}
		// 父级更新，修改子级信息
		_, updateErr := txOrm.QueryTable(new(models.Team)).Filter("pid", params.Id).Update(orm.Params{
			"package_id":   params.PackageId,
			"is_official":  params.IsOfficial,
			"is_bloger":    params.IsBloger,
			"is_drop_user": params.IsDropUser,
			"token":        params.Token,
			"secret":       params.Secret,
			"rtp":          params.Rtp,
		})
		if updateErr != nil {
			return updateErr
		}
		return nil
	})

	if trErr != nil {
		logs.Error("更新团队失败:%v", trErr)
		return fmt.Errorf("ShuJuGengXinShiBai")
	}
	return nil
}

// 删除团队
func (s *TeamService) DelTeam(id int) error {
	model := models.CreateTeamModel()
	field := []string{
		"is_deleted",
	}
	_, err := model.Update(&models.Team{
		Id:        id,
		IsDeleted: 1,
	}, field...)
	if err != nil {
		logs.Error("删除团队失败:%v", err)
		return fmt.Errorf("ShuJuShanChuShiBai")
	}
	return nil
}

// GetUserTeamInfo 获取用户团队信息
func (s *TeamService) GetUserTeamInfo(roleID int) (models.Team, error) {
	model := models.CreateTeamModel()
	var info models.Team
	err := model.QueryTable(new(models.Team)).Filter("role_id", roleID).Filter("is_deleted", 0).One(&info)
	if err != nil {
		return info, err
	}
	return info, nil
}

// 获取代理和团队信息
func (s *TeamService) GetAllSubTeamsWithParent(parents string) (SubTeamWithParent, error) {
	if len(parents) == 0 {
		return SubTeamWithParent{}, nil
	}
	parents = strings.Trim(parents, ",")
	userIds := strings.Split(parents, ",")

	queryBuilder, _ := models.NewReadQueryBuilder()
	queryBuilder.QB().Select(
		"t1.title",
		"t1.role_id",
		"t1.pid",
		"COALESCE(t2.role_id, 0) AS parent_role_id",
		"COALESCE(t2.title, '') AS parent_title").
		From(models.TablePrefix + "team t1").
		LeftJoin(models.TablePrefix + "team t2").On("t1.pid = t2.id").
		Where("t1.is_deleted = 0").
		And(fmt.Sprintf("t1.user_id IN (%s)", strings.Trim(strings.Repeat("?,", len(userIds)), ","))).OrderBy("t1.pid desc")

	// 获取 SQL 字符串
	sql := queryBuilder.String()

	// 转换参数
	args := make([]interface{}, len(userIds))
	for i, v := range userIds {
		args[i] = v
	}

	// 执行查询
	var results []SubTeamWithParent
	_, err := queryBuilder.QueryRawSql(&results, sql, args)

	if len(results) == 0 {
		return SubTeamWithParent{}, nil
	}
	result := results[0]

	if result.Pid == 0 {
		info := SubTeamWithParent{
			Title:        "",
			RoleId:       0,
			ParentTitle:  result.Title,
			ParentRoleId: result.RoleId,
		}
		result = info
	}
	return result, err
}

// GetAllTeams 获取所有团队
func (s *TeamService) GetAllTeams(packageIds ...int) ([]models.Team, error) {
	model := models.CreateTeamModel()
	var teams []models.Team
	condition := map[string]interface{}{
		"is_deleted": 0,
	}
	if len(packageIds) > 0 {
		condition["package_id__in"] = packageIds
	}
	db := model.QueryTable(&models.Team{})
	_, err := model.Where(db, condition).All(&teams)
	if err != nil {
		return teams, err
	}
	return teams, nil
}
