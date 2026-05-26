package services

import (
	"api/models"
	"api/utils"
)

type UserBalanceLogService struct{}

type UserBalanceLogRequestParams struct {
	Id             int
	Type           int
	Types          []int
	PackageId      int
	PackageIds     []int
	GameTypeId     int
	GameTypeIds    []int
	GameId         int
	ActivityTypeId int
	ActivityId     int
	UserId         int
	UserIds        []int
	RoleId         int
	PlatformId     int
	BeginTime      int64
	EndTime        int64
	Fields         string // 查询字段
	Page           int
	PageSize       int
	OrderBy        string //"c_time DESC"
}

type UserBalanceLogGroupListRequestParams struct {
	Id             int
	Type           int
	Types          []int
	PackageId      int
	PackageIds     []int
	GameTypeId     int
	GameTypeIds    []int
	GameId         int
	ActivityTypeId int
	ActivityId     int
	UserId         int
	UserIds        []int
	RoleId         int
	PlatformId     int
	BeginTime      int64
	EndTime        int64
	List           interface{}
	Fields         string // 查询字段
	GroupBy        string // Group By字段
	GroupFields    string // 查询字段
	OuterGroupBy   string // 外层 Group By字段
}

type UserBalanceLogResult struct {
	BetAmount float64 `json:"bet_amount"`
}

type UserBalanceLogResponse struct {
	List  []models.UserBalanceLog `json:"list"`
	Total int                     `json:"total"`
}

func (s *UserBalanceLogService) buildCondition(params UserBalanceLogRequestParams) ([]string, []interface{}) {
	var conditions []string
	var tableArgs []interface{}

	if params.BeginTime > 0 {
		conditions = append(conditions, "created_time >= ?")
		tableArgs = append(tableArgs, params.BeginTime)
	}

	if params.EndTime > 0 {
		conditions = append(conditions, "created_time < ?")
		tableArgs = append(tableArgs, params.EndTime)
	}

	if params.Id > 0 {
		conditions = append(conditions, "id = ?")
		tableArgs = append(tableArgs, params.Id)
	}

	if params.Type > 0 {
		conditions = append(conditions, "type = ?")
		tableArgs = append(tableArgs, params.Type)
	}

	if params.PackageId > 0 {
		conditions = append(conditions, "package_id = ?")
		tableArgs = append(tableArgs, params.PackageId)
	}

	if params.GameTypeId > 0 {
		conditions = append(conditions, "game_type_id = ?")
		tableArgs = append(tableArgs, params.GameTypeId)
	}

	if params.GameId > 0 {
		conditions = append(conditions, "game_id = ?")
		tableArgs = append(tableArgs, params.GameId)
	}

	if params.ActivityTypeId > 0 {
		conditions = append(conditions, "activity_type_id = ?")
		tableArgs = append(tableArgs, params.ActivityTypeId)
	}

	if params.ActivityId > 0 {
		conditions = append(conditions, "activity_id = ?")
		tableArgs = append(tableArgs, params.ActivityId)
	}

	if params.UserId > 0 {
		conditions = append(conditions, "user_id = ?")
		tableArgs = append(tableArgs, params.UserId)
	}

	if params.PlatformId > 0 {
		conditions = append(conditions, "platform_id = ?")
		tableArgs = append(tableArgs, params.PlatformId)
	}

	// in查询
	if len(params.UserIds) > 0 {
		conditions, tableArgs = utils.BuildInCondition("user_id", params.UserIds, conditions, tableArgs)
	}

	if len(params.PackageIds) > 0 {
		conditions, tableArgs = utils.BuildInCondition("package_id", params.PackageIds, conditions, tableArgs)
	}

	if len(params.GameTypeIds) > 0 {
		conditions, tableArgs = utils.BuildInCondition("game_type_id", params.GameTypeIds, conditions, tableArgs)
	}

	if len(params.Types) > 0 {
		conditions, tableArgs = utils.BuildInCondition("type", params.Types, conditions, tableArgs)
	}

	return conditions, tableArgs
}

// GetTotalBalanceAmount 获取总的变动金额
func (s *UserBalanceLogService) GetTotalBalanceAmount(params UserBalanceLogRequestParams) (float64, error) {
	conditions, tableArgs := s.buildCondition(params)
	queryBuilder, err := models.NewReadQueryBuilder()
	if err != nil {
		return 0.0, err
	}

	var amount []float64
	queryParams := models.QuerySubTableRecordListParams{
		BeginTime:     params.BeginTime,
		EndTime:       params.EndTime,
		TimeField:     "created_time",
		TableName:     "user_balance_log",
		Conditions:    conditions,
		TableArgs:     tableArgs,
		List:          &amount,
		Fields:        "IFNULL(sum(amount*100), 0) as total",
		IsNotGetTotal: true,
	}

	_, queryErr := queryBuilder.QuerySubTableRecordList(queryParams)
	if queryErr != nil {
		return 0, queryErr
	}

	return s.sum(amount), nil
}

func (s *UserBalanceLogService) sum(amount []float64) float64 {
	sum := 0.0
	for _, v := range amount {
		sum += v
	}
	return sum / 100
}

// 获取列表
func (s *UserBalanceLogService) GetList(params UserBalanceLogRequestParams) (UserBalanceLogResponse, error) {
	// 构建 WHERE 条件
	conditions, tableArgs := s.buildCondition(params)

	queryBuilder, err := models.NewReadQueryBuilder()
	if err != nil {
		return UserBalanceLogResponse{
			List:  []models.UserBalanceLog{},
			Total: 0,
		}, err
	}

	var list []models.UserBalanceLog
	queryParams := models.QuerySubTableRecordListParams{
		BeginTime:  params.BeginTime,
		EndTime:    params.EndTime,
		TimeField:  "created_time",
		TableName:  "user_balance_log",
		Conditions: conditions,
		TableArgs:  tableArgs,
		OrderBy:    params.OrderBy,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Fields:     params.Fields,
		List:       &list,
	}

	total, err := queryBuilder.QuerySubTableRecordList(queryParams)

	if list == nil {
		list = []models.UserBalanceLog{}
	}

	return UserBalanceLogResponse{
		List:  list,
		Total: int(total),
	}, err
}

func (s *UserBalanceLogService) GroupList(params UserBalanceLogGroupListRequestParams) error {
	// 构建 WHERE 条件
	conditionParams := UserBalanceLogRequestParams{
		Id:             params.Id,
		Type:           params.Type,
		Types:          params.Types,
		PackageId:      params.PackageId,
		GameTypeId:     params.GameTypeId,
		GameId:         params.GameId,
		ActivityTypeId: params.ActivityTypeId,
		ActivityId:     params.ActivityId,
		UserId:         params.UserId,
		UserIds:        params.UserIds,
		BeginTime:      params.BeginTime,
		EndTime:        params.EndTime,
	}
	conditions, tableArgs := s.buildCondition(conditionParams)

	queryBuilder, err := models.NewReadQueryBuilder()
	if err != nil {
		return err
	}

	queryParams := models.GroupSubTableRecordListParams{
		BeginTime:    params.BeginTime,
		EndTime:      params.EndTime,
		TimeField:    "created_time",
		TableName:    "user_balance_log",
		Conditions:   conditions,
		TableArgs:    tableArgs,
		List:         params.List,
		Fields:       params.Fields,
		GroupBy:      params.GroupBy,
		GroupFields:  params.GroupFields,
		OuterGroupBy: params.OuterGroupBy,
	}

	logErr := queryBuilder.GroupSubTableRecordList(queryParams)

	return logErr
}

type GetUserBalanceLogDailyStatResponse struct {
	TotalAmount float64
	TotalMen    int
	PackageId   int
	Type        int
}

// GetUserBalanceLogDailyStat 分包获取用户的金币统计
func (s *UserBalanceLogService) GetUserBalanceLogDailyStat(params UserBalanceLogGroupListRequestParams) ([]GetUserBalanceLogDailyStatResponse, error) {
	list := make([]GetUserBalanceLogDailyStatResponse, 0)
	params.List = &list
	params.Fields = "sum(amount * 100) as total_amount, user_id, package_id, type"
	params.GroupBy = "package_id, type, user_id"
	params.OuterGroupBy = "package_id, type"
	params.GroupFields = "sum(total_amount)/100 as total_amount, count(distinct user_id) as total_men, package_id, type"
	err := s.GroupList(params)
	if err != nil {
		return list, err
	}
	return list, nil
}

type GetUserBalanceLogDailyMenStatResponse struct {
	TotalMen int
	Type     int
}

// GetUserBalanceLogDailyMenStat 日、月报表统计实时人数
func (s *UserBalanceLogService) GetUserBalanceLogDailyMenStat(params UserBalanceLogGroupListRequestParams) ([]GetUserBalanceLogDailyMenStatResponse, error) {
	list := make([]GetUserBalanceLogDailyMenStatResponse, 0)
	params.List = &list
	params.Fields = "type, user_id"
	params.GroupBy = "type, user_id"
	params.OuterGroupBy = "type"
	params.GroupFields = "count(distinct user_id) as total_men, type"
	err := s.GroupList(params)
	if err != nil {
		return list, err
	}
	return list, nil
}
