package services

import (
	"api/models"
	"fmt"
)

type LoginLogService struct{}

type LoginLogRequestParams struct {
	UserId            int
	IP                string
	PackageId         int
	PackageIds        string
	BeginTime         int64
	EndTime           int64
	RegisterBeginTime int64
	RegisterEndTime   int64
	Page              int
	PageSize          int
	GroupBy           string
}

type LoginLogResponse struct {
	List  []models.LoginLog `json:"list"`
	Total int               `json:"total"`
}

func (s *LoginLogService) buildCondition(params LoginLogRequestParams) ([]string, []interface{}) {
	// 构建 WHERE 条件
	var conditions []string
	var tableArgs []interface{}

	if params.PackageId > 0 {
		conditions = append(conditions, "package_id = ?")
		tableArgs = append(tableArgs, params.PackageId)
	}

	if params.PackageIds != "" {
		conditions = append(conditions, "package_id in (?)")
		tableArgs = append(tableArgs, params.PackageIds)
	}

	if params.BeginTime > 0 {
		conditions = append(conditions, "created_at >= ?")
		tableArgs = append(tableArgs, params.BeginTime)
	}
	if params.EndTime > 0 {
		conditions = append(conditions, "created_at < ?")
		tableArgs = append(tableArgs, params.EndTime)
	}

	if params.RegisterBeginTime > 0 {
		conditions = append(conditions, "register_at >= ?")
		tableArgs = append(tableArgs, params.RegisterBeginTime)
	}

	if params.RegisterEndTime > 0 {
		conditions = append(conditions, "register_at < ?")
		tableArgs = append(tableArgs, params.RegisterEndTime)
	}

	if params.UserId > 0 {
		conditions = append(conditions, "user_id = ?")
		tableArgs = append(tableArgs, params.UserId)
	}

	if params.IP != "" {
		conditions = append(conditions, "ip = ?")
		tableArgs = append(tableArgs, params.IP)
	}

	return conditions, tableArgs
}

// GetList 获取列表
func (s *LoginLogService) GetList(params LoginLogRequestParams) (LoginLogResponse, error) {
	fmt.Printf("参数:%v\n", params)
	conditions, tableArgs := s.buildCondition(params)
	queryBuilder, err := models.NewReadQueryBuilder()
	if err != nil {
		return LoginLogResponse{
			List:  []models.LoginLog{},
			Total: 0,
		}, err
	}

	var list []models.LoginLog
	queryParams := models.QuerySubTableRecordListParams{
		BeginTime:  params.BeginTime,
		EndTime:    params.EndTime,
		TimeField:  "created_at",
		TableName:  "login_log",
		Conditions: conditions,
		TableArgs:  tableArgs,
		OrderBy:    "created_at DESC",
		Page:       params.Page,
		PageSize:   params.PageSize,
		List:       &list,
	}

	total, err := queryBuilder.QuerySubTableRecordList(queryParams)

	return LoginLogResponse{
		List:  list,
		Total: int(total),
	}, err
}

type LoginLogGroupListRequestParams struct {
	UserId            int
	RegisterBeginTime int64
	RegisterEndTime   int64
	BeginTime         int64
	EndTime           int64
	PackageId         int
	List              interface{}
	Fields            string // 查询字段
	GroupBy           string // Group By字段
	GroupFields       string // 查询字段
}

func (s *LoginLogService) GroupList(params LoginLogGroupListRequestParams) error {
	// 构建 WHERE 条件
	conditionParams := LoginLogRequestParams{
		UserId:            params.UserId,
		BeginTime:         params.BeginTime,
		EndTime:           params.EndTime,
		PackageId:         params.PackageId,
		RegisterBeginTime: params.RegisterBeginTime,
		RegisterEndTime:   params.RegisterEndTime,
	}
	conditions, tableArgs := s.buildCondition(conditionParams)

	queryBuilder, err := models.NewReadQueryBuilder()
	if err != nil {
		return err
	}

	queryParams := models.GroupSubTableRecordListParams{
		BeginTime:   params.BeginTime,
		EndTime:     params.EndTime,
		TimeField:   "created_at",
		TableName:   "login_log",
		Conditions:  conditions,
		TableArgs:   tableArgs,
		List:        params.List,
		Fields:      params.Fields,
		GroupBy:     params.GroupBy,
		GroupFields: params.GroupFields,
	}

	logErr := queryBuilder.GroupSubTableRecordList(queryParams)

	return logErr
}
