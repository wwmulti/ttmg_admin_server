package services

import (
	"api/models"
)

type UserOperateLogService struct{}

type UserOperateLogRequestParams struct {
	UserId     int
	RoleId     int
	PackageId  int
	PackageIds string
	BeginTime  int64
	EndTime    int64
	Page       int
	PageSize   int
	OrderBy    string //"created_at DESC"
}
type UserOperateLogResponse struct {
	List  []models.UserOperateLog `json:"list"`
	Total int                     `json:"total"`
}

func (s *UserOperateLogService) buildCondition(params UserOperateLogRequestParams) ([]string, []interface{}) {
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

	if params.UserId > 0 {
		conditions = append(conditions, "user_id = ?")
		tableArgs = append(tableArgs, params.UserId)
	}

	if params.RoleId > 0 {
		conditions = append(conditions, "role_id = ?")
		tableArgs = append(tableArgs, params.RoleId)
	}

	return conditions, tableArgs
}

// GetList 获取列表
func (s *UserOperateLogService) GetList(params UserOperateLogRequestParams) (UserOperateLogResponse, error) {
	// 构建 WHERE 条件
	conditions, tableArgs := s.buildCondition(params)

	queryBuilder, err := models.NewReadQueryBuilder()
	if err != nil {
		return UserOperateLogResponse{
			List:  []models.UserOperateLog{},
			Total: 0,
		}, err
	}

	var list []models.UserOperateLog
	queryParams := models.QuerySubTableRecordListParams{
		BeginTime:  params.BeginTime,
		EndTime:    params.EndTime,
		TimeField:  "created_at",
		TableName:  "user_operate_log",
		Conditions: conditions,
		TableArgs:  tableArgs,
		OrderBy:    params.OrderBy,
		Page:       params.Page,
		PageSize:   params.PageSize,
		List:       &list,
	}

	total, err := queryBuilder.QuerySubTableRecordList(queryParams)

	return UserOperateLogResponse{
		List:  list,
		Total: int(total),
	}, err
}
