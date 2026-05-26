package services

import (
	"api/models"
	"api/utils"
)

type GameWaterLogService struct{}

type GameWaterLogRequestParams struct {
	UserId         int
	RoleId         int
	GameCateId     []int
	GamePlatformId []int
	GameId         int
	OrderId        string
	ParentOrderId  string
	PackageId      int
	PackageIds     []int
	BeginTime      int64
	EndTime        int64
	IsEnd          int // 数据库值加1（1：未结束，2：已结束）
	Page           int
	PageSize       int
	OrderBy        string //"c_time DESC"
}

type GameWaterLogGroupListRequestParams struct {
	UserId         int
	GameCateId     interface{}
	GamePlatformId interface{}
	GameId         int
	OrderId        string
	ParentOrderId  string
	BeginTime      int64
	PackageId      []int
	EndTime        int64
	IsEnd          int // 数据库值加1（1：未结束，2：已结束）
	List           interface{}
	Fields         string // 查询字段
	GroupBy        string // Group By字段
	GroupFields    string // 查询字段
}

type GameWaterLogResponse struct {
	List  []models.GameWaterLog `json:"list"`
	Total int                   `json:"total"`
}

func (s *GameWaterLogService) buildCondition(params GameWaterLogRequestParams) ([]string, []interface{}) {
	var conditions []string
	var tableArgs []interface{}

	if len(params.PackageIds) > 0 {
		conditions, tableArgs = utils.BuildInCondition("package_id", params.PackageIds, conditions, tableArgs)
	}

	if params.PackageId > 0 {
		conditions = append(conditions, "package_id = ?")
		tableArgs = append(tableArgs, params.PackageId)
	}

	if params.BeginTime > 0 {
		conditions = append(conditions, "c_time >= ?")
		tableArgs = append(tableArgs, params.BeginTime)
	}
	if params.EndTime > 0 {
		conditions = append(conditions, "c_time < ?")
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

	if len(params.OrderId) > 0 {
		conditions = append(conditions, "order_id = ?")
		tableArgs = append(tableArgs, params.OrderId)
	}

	if len(params.ParentOrderId) > 0 {
		conditions = append(conditions, "parent_order_id = ?")
		tableArgs = append(tableArgs, params.ParentOrderId)
	}

	if params.GameId > 0 {
		conditions = append(conditions, "game_id = ?")
		tableArgs = append(tableArgs, params.GameId)
	}

	if len(params.GameCateId) > 0 {
		conditions, tableArgs = utils.BuildInCondition("game_cat_id", params.GameCateId, conditions, tableArgs)
	}

	if len(params.GamePlatformId) > 0 {
		conditions, tableArgs = utils.BuildInCondition("game_platform_id", params.GamePlatformId, conditions, tableArgs)
	}

	if params.IsEnd >= 1 {
		isEnd := params.IsEnd - 1
		conditions = append(conditions, "is_end = ?")
		tableArgs = append(tableArgs, isEnd)
	}

	return conditions, tableArgs
}

// GetList 获取列表
func (s *GameWaterLogService) GetList(params GameWaterLogRequestParams, uid int) (GameWaterLogResponse, error) {
	// 构建 WHERE 条件
	conditions, tableArgs := s.buildCondition(params)
	queryBuilder, err := models.NewReadQueryBuilder()
	if err != nil {
		return GameWaterLogResponse{
			List:  []models.GameWaterLog{},
			Total: 0,
		}, err
	}

	var list []models.GameWaterLog
	queryParams := models.QuerySubTableRecordListParams{
		BeginTime:  params.BeginTime,
		EndTime:    params.EndTime,
		TimeField:  "c_time",
		TableName:  "game_water_log",
		Conditions: conditions,
		TableArgs:  tableArgs,
		OrderBy:    params.OrderBy,
		Page:       params.Page,
		PageSize:   params.PageSize,
		List:       &list,
	}

	total, err := queryBuilder.QuerySubTableRecordList(queryParams)

	if len(list) == 0 {
		list = make([]models.GameWaterLog, 0)
	}

	return GameWaterLogResponse{
		List:  list,
		Total: int(total),
	}, err
}

func (s *GameWaterLogService) GroupList(params GameWaterLogGroupListRequestParams) error {
	var gameCateId []int
	if v, ok := params.GameCateId.([]int); ok {
		gameCateId = v
	} else if v, ok := params.GameCateId.(int); ok {
		gameCateId = []int{v}
	} else {
		gameCateId = []int{}
	}

	var gamePlatformId []int
	if v, ok := params.GamePlatformId.([]int); ok {
		gamePlatformId = v
	} else if v, ok := params.GamePlatformId.(int); ok {
		gamePlatformId = []int{v}
	} else {
		gamePlatformId = []int{}
	}

	// 构建 WHERE 条件
	conditionParams := GameWaterLogRequestParams{
		UserId:         params.UserId,
		GameCateId:     gameCateId,
		GamePlatformId: gamePlatformId,
		GameId:         params.GameId,
		OrderId:        params.OrderId,
		ParentOrderId:  params.ParentOrderId,
		BeginTime:      params.BeginTime,
		EndTime:        params.EndTime,
	}
	conditions, tableArgs := s.buildCondition(conditionParams)

	queryBuilder, err := models.NewReadQueryBuilder()
	if err != nil {
		return err
	}

	queryParams := models.GroupSubTableRecordListParams{
		BeginTime:   params.BeginTime,
		EndTime:     params.EndTime,
		TimeField:   "c_time",
		TableName:   "game_water_log",
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
