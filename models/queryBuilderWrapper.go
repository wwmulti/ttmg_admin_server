package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

// QueryBuilder 读写分离封装
type QueryBuilderWrapper struct {
	Base
	qb     orm.QueryBuilder
	isRead bool
}

// 创建读操作的 QueryBuilder
func NewReadQueryBuilder() (*QueryBuilderWrapper, error) {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, err
	}

	return &QueryBuilderWrapper{
		qb:     qb,
		isRead: true,
	}, nil
}

// 创建写操作的 QueryBuilder
func NewWriteQueryBuilder() (*QueryBuilderWrapper, error) {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, err
	}

	return &QueryBuilderWrapper{
		qb:     qb,
		isRead: false,
	}, nil
}

func (w *QueryBuilderWrapper) QueryRawSql(dest interface{}, sql string, args ...interface{}) (int64, error) {
	if !w.isRead {
		return 0, orm.ErrNotImplement
	}
	o := w.GetSlaveOrm()
	return o.Raw(sql, args...).QueryRows(dest)
}

// 执行查询（读操作）
func (w *QueryBuilderWrapper) QueryRows(dest interface{}, args ...interface{}) (int64, error) {
	if !w.isRead {
		return 0, orm.ErrNotImplement
	}

	sql := w.qb.String()
	o := w.GetSlaveOrm() // 使用从库
	return o.Raw(sql, args...).QueryRows(dest)
}

// 执行查询单条记录（读操作）
func (w *QueryBuilderWrapper) QueryRow(dest interface{}, args ...interface{}) error {
	if !w.isRead {
		return orm.ErrNotImplement
	}

	sql := w.qb.String()
	o := w.GetMasterDb() // 使用从库
	return o.Raw(sql, args...).QueryRow(dest)
}

// 执行写操作
func (w *QueryBuilderWrapper) Exec(args ...interface{}) (sql.Result, error) {
	if w.isRead {
		return nil, orm.ErrNotImplement
	}

	sql := w.qb.String()
	o := w.GetMasterDb() // 使用主库
	return o.Raw(sql, args...).Exec()
}

// 获取原生 QueryBuilder
func (w *QueryBuilderWrapper) QB() orm.QueryBuilder {
	return w.qb
}

// 获取 SQL 语句
func (w *QueryBuilderWrapper) String() string {
	return w.qb.String()
}

// 通用 UNION 查询函数
func (w *QueryBuilderWrapper) UnionQuery(queries []string, orderBy string, limit, offset int, unionAll bool) string {
	if len(queries) == 0 {
		return ""
	}

	// 确定 UNION 类型
	unionType := " UNION "
	if unionAll {
		unionType = " UNION ALL "
	}

	// 合并查询
	sql := queries[0]
	for i := 1; i < len(queries); i++ {
		sql += unionType + queries[i]
	}

	// 添加外层排序和分页
	if orderBy != "" {
		sql += " ORDER BY " + orderBy
	}
	if limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", limit)
		if offset > 0 {
			sql += fmt.Sprintf(" OFFSET %d", offset)
		}
	}

	return sql
}

/*
*

	查询分表的列表数据

*
*/
type QuerySubTableRecordListParams struct {
	BeginTime     int64         // 开始时间
	EndTime       int64         // 结束时间
	TimeField     string        // 时间字段名称
	TableName     string        // 表名称(没有表前缀)
	Conditions    []string      // 查询条件
	TableArgs     []interface{} // 查询条件参数
	OrderBy       string        // 排序 原生写法例如 createTime desc
	Page          int           // 页数
	PageSize      int           // 条数
	List          interface{}
	Fields        string // 查询字段
	IsNotGetTotal bool   // 是否获取总条数
}

type GroupSubTableRecordListParams struct {
	BeginTime    int64         // 开始时间
	EndTime      int64         // 结束时间
	TimeField    string        // 时间字段名称
	TableName    string        // 表名称(没有表前缀)
	Conditions   []string      // 查询条件
	TableArgs    []interface{} // 查询条件参数
	OrderBy      string        // 排序 原生写法例如 createTime desc
	Page         int           // 页数
	PageSize     int           // 条数
	List         interface{}
	Fields       string // 查询字段
	GroupBy      string // Group By字段
	GroupFields  string // 排序查询的字段
	OuterGroupBy string // 外层 Group By字段
}

// 分表分组查询
func (w *QueryBuilderWrapper) GroupSubTableRecordList(params GroupSubTableRecordListParams) error {
	condition := make(map[string]interface{})

	if params.BeginTime <= 0 || params.EndTime <= 0 {
		return fmt.Errorf("ChaXunShiJianBiChuan")
	}

	condition[fmt.Sprintf("%v__gte", params.TimeField)] = params.BeginTime
	condition[fmt.Sprintf("%v__lt", params.TimeField)] = params.EndTime

	tableList := w.GetTableNamesByDateRange(params.TableName, time.Unix(params.BeginTime, 0), time.Unix(params.EndTime, 0))

	if len(tableList) == 0 {
		return nil
	}

	// ===================== 构建所有分表查询 =====================
	var queries []string
	var allArgs []interface{}

	// 构建基础 SELECT
	fields := params.Fields
	if len(fields) == 0 {
		fields = "*"
	}

	for _, table := range tableList {
		queryBuilder, err := NewReadQueryBuilder()
		if err != nil {
			return err
		}

		qb := queryBuilder.QB()

		qb.Select(fields).From(table)

		// 拼接 WHERE 条件
		if len(params.Conditions) > 0 {
			whereClause := strings.Join(params.Conditions, " AND ")
			qb.Where(whereClause)
		}

		// 收集 SQL 和参数
		queries = append(queries, "("+qb.String()+" group by "+params.GroupBy+")")
		allArgs = append(allArgs, params.TableArgs...)
	}

	// ===================== 使用你自己的 UnionQuery 合并 =====================
	qbw, _ := NewReadQueryBuilder()

	unionSQL := qbw.UnionQuery(queries, "", 0, 0, true)

	if len(params.GroupFields) == 0 {
		params.GroupFields = fields
	}

	// 外层 Group By 字段
	outerGroupBy := params.OuterGroupBy
	if len(outerGroupBy) == 0 {
		outerGroupBy = params.GroupBy
	}

	groupSql := "SELECT " + params.GroupFields + " FROM (" + unionSQL + ") AS t group by " + outerGroupBy
	if params.OrderBy != "" {
		groupSql += " ORDER BY " + params.OrderBy
	}
	if params.PageSize > 0 {
		groupSql += fmt.Sprintf(" LIMIT %d", params.PageSize)
		// 排序分页
		offset := params.PageSize * (params.Page - 1)
		if offset < 0 {
			offset = 0
		}
		groupSql += fmt.Sprintf(" OFFSET %d", offset)
	}

	// ===================== 执行查询 =====================
	// 从库查询 + 自动绑定结构体
	logs.Info("sql:%v,args:%v", groupSql, allArgs)

	_, err := qbw.QueryRawSql(params.List, groupSql, allArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (w *QueryBuilderWrapper) QuerySubTableRecordList(params QuerySubTableRecordListParams) (int64, error) {
	condition := make(map[string]interface{})

	if params.BeginTime <= 0 || params.EndTime <= 0 {
		return 0, fmt.Errorf("ChaXunShiJianBiChuan")
	}

	condition[fmt.Sprintf("%v__gte", params.TimeField)] = params.BeginTime
	condition[fmt.Sprintf("%v__lt", params.TimeField)] = params.EndTime

	tableList := w.GetTableNamesByDateRange(params.TableName, time.Unix(params.BeginTime, 0), time.Unix(params.EndTime, 0))

	if len(tableList) == 0 {
		return 0, nil
	}

	// ===================== 构建所有分表查询 =====================
	var queries []string
	var allArgs []interface{}

	for _, table := range tableList {
		queryBuilder, err := NewReadQueryBuilder()
		if err != nil {
			return 0, err
		}

		qb := queryBuilder.QB()

		// 构建基础 SELECT
		fields := params.Fields
		if len(fields) == 0 {
			fields = "*"
		}
		qb.Select(fields).From(table)

		// 拼接 WHERE 条件
		if len(params.Conditions) > 0 {
			whereClause := strings.Join(params.Conditions, " AND ")
			qb.Where(whereClause)
		}

		// 收集 SQL 和参数
		queries = append(queries, "("+qb.String()+")")
		allArgs = append(allArgs, params.TableArgs...)
	}

	// ===================== 使用你自己的 UnionQuery 合并 =====================
	qbw, _ := NewReadQueryBuilder()

	// 排序分页
	offset := params.PageSize * (params.Page - 1)
	if offset < 0 {
		offset = 0
	}
	unionSQL := qbw.UnionQuery(queries, params.OrderBy, params.PageSize, offset, true)

	// ===================== 执行查询 =====================
	// 从库查询 + 自动绑定结构体
	logs.Info("sql:%v,args:%v", unionSQL, allArgs)

	rowNumber, err := qbw.QueryRawSql(params.List, unionSQL, allArgs...)
	if err != nil {
		return 0, err
	}

	if !params.IsNotGetTotal {
		//========总条数==========
		countUnionSQL := qbw.UnionQuery(queries, "", 0, 0, true)
		countSql := "SELECT COUNT(*) FROM (" + countUnionSQL + ") AS t"
		var total []int
		_, countErr := qbw.QueryRawSql(&total, countSql, allArgs...)
		if countErr != nil {
			return 0, countErr
		}
		return int64(total[0]), nil
	}

	return rowNumber, nil
}
