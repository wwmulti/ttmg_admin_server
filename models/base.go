package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

var TablePrefix = "" // 表前缀

func init() {
	if sqlconnPrefix, err := beego.AppConfig.String("sqlconn_prefix"); nil == err {
		TablePrefix = sqlconnPrefix + "_"
	}
	sqlConnMaster, err := beego.AppConfig.String("sqlconn_master")
	if sqlConnMaster == "" || nil != err {
		panic("sqlconn_master配置没有设置！" + err.Error())
	}
	orm.RegisterDataBase("default", "mysql", sqlConnMaster)
	orm.SetMaxIdleConns("default", 10)
	orm.SetMaxOpenConns("default", 100)
	sqlConnSlave, err := beego.AppConfig.String("sqlconn_slave")
	if sqlConnSlave == "" || nil != err {
		panic("sqlconn_slave配置没有设置！" + err.Error())
	}
	orm.RegisterDataBase("slave", "mysql", sqlConnSlave)
	orm.SetMaxIdleConns("slave", 10)
	orm.SetMaxOpenConns("slave", 100)
	orm.RegisterModelWithPrefix(TablePrefix, GetRegisterAllModels()...)
}

type QuerySeterEx struct {
	orm.QuerySeter
	base *Base
	md   interface{}
}

func (self QuerySeterEx) Filter(key string, values ...interface{}) QuerySeterEx {
	self.QuerySeter = self.QuerySeter.Filter(key, values...)
	return self
}

func (self QuerySeterEx) Update(values orm.Params) (int64, error) {
	return self.base.OrmerMaster.QueryTable(self.md).SetCond(self.GetCond()).Update(values)
}

func (self QuerySeterEx) Delete() (int64, error) {
	return self.base.OrmerMaster.QueryTable(self.md).SetCond(self.GetCond()).Delete()
}

type Base struct {
	OrmerMaster orm.Ormer
	OrmerSlave  orm.Ormer
}

func (self *Base) GetMasterDb() orm.Ormer {
	// 关键：获取【主库 Ormer】，确保写入主库
	db, _ := orm.GetDB("default")
	o, _ := orm.NewOrmWithDB("mysql", "default", db)
	return o
}

func (self *Base) GetSlaveOrm() orm.Ormer {
	// 关键：获取【从库 Ormer】
	db, _ := orm.GetDB("slave")
	o, _ := orm.NewOrmWithDB("mysql", "slave", db)
	return o
}

func (self *Base) Read(md interface{}, cols ...string) error {
	return self.OrmerSlave.Read(md, cols...)
}

func (self *Base) QueryTable(ptrStructOrTableName interface{}) QuerySeterEx {
	return QuerySeterEx{self.OrmerSlave.QueryTable(ptrStructOrTableName), self, ptrStructOrTableName}
}

func (self *Base) Insert(md interface{}) (int64, error) {
	return self.OrmerMaster.Insert(md)
}

func (self *Base) Update(md interface{}, cols ...string) (int64, error) {
	return self.OrmerMaster.Update(md, cols...)
}

func (self *Base) Delete(md interface{}, cols ...string) (int64, error) {
	return self.OrmerMaster.Delete(md, cols...)
}

func (self *Base) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	return self.OrmerMaster.ReadOrCreate(md, col1, cols...)
}

func (self *Base) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	return self.OrmerMaster.InsertOrUpdate(md, colConflitAndArgs...)
}

func (self *Base) InsertMulti(bulk int, mds interface{}) (int64, error) {
	return self.OrmerMaster.InsertMulti(bulk, mds)
}

func (self *Base) Begin() (orm.TxOrmer, error) {
	return self.OrmerMaster.Begin()
}

func (self *Base) Raw(query string, args ...interface{}) orm.RawSeter {
	return self.OrmerSlave.Raw(query, args...)
}

func (self *Base) Where(db QuerySeterEx, params map[string]interface{}) QuerySeterEx {
	for key, value := range params {
		db = db.Filter(key, value)
	}
	return db
}

/*
*
获取分页列表
relatedFields 传表的模型名称例如 Account
*
*/
func (self *Base) GetPageList(model interface{}, condition map[string]interface{}, page, limit int, orderBy string, relatedFields ...string) (interface{}, int64, error) {
	db := self.QueryTable(model)
	db = self.Where(db, condition)
	// 获取总数
	// 获取总数 - 使用正确的类型
	var total int64
	var err error

	if len(relatedFields) > 0 {
		relatedInterfaces := make([]interface{}, len(relatedFields))
		for i, v := range relatedFields {
			relatedInterfaces[i] = v
		}
		// 直接使用 db，不需要重新赋值给 countDB
		// RelatedSel 返回 QuerySeter，可以直接调用 Count()
		total, err = db.RelatedSel(relatedInterfaces...).Count()
		if err != nil {
			return nil, 0, err
		}
	} else {
		total, err = db.Count()
		if err != nil {
			return nil, 0, err
		}
	}

	if limit < 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// 使用反射创建切片
	modelType := reflect.TypeOf(model).Elem()
	sliceType := reflect.SliceOf(modelType)
	resultSlice := reflect.New(sliceType)

	query := db.Limit(limit, offset).OrderBy(orderBy)
	// 指定要加载的关联字段
	if len(relatedFields) > 0 {
		relatedInterfaces := make([]interface{}, len(relatedFields))
		for i, v := range relatedFields {
			relatedInterfaces[i] = v
		}
		query = query.RelatedSel(relatedInterfaces...)
	} else {
		query = query.RelatedSel()
	}

	_, err = query.All(resultSlice.Interface())
	if err != nil {
		return nil, 0, err
	}

	return resultSlice.Elem().Interface(), total, nil
}

// 通用条件查询构造
// 用法：op: like, gte, lte, gt, lt, ne, in, range, time_range, -
// 默认 结构体op为空就是 等于, op 为 "-" 不处理
// like （Name string `json:"name" op:"like"`）
// range (SortRange []int `json:"sort_range" op:"range"`)  [1, 2]
// range (SortRange string `json:"sort_range" op:"range"`) {"vaule1":1, "value2":2}
// time_range (CreateTimeRange []string `json:"create_time_range" op:"time_range"`)
// in (Id string `json:"id" op:"in"`)  1,2,3逗号分隔
// gte, lte, gt, lt, ne (Number int `json:"number" op:"gte"`)
// 注意：
// 1、字段用form 类型
// 2、RawQuery 必传，用于判断前端是否传参
// 3、有orm标签则优先使用orm标签，否则使用form标签
// 4、page, pageSize 分页参数不处理
func (self *Base) BuildCondition(req interface{}, sorts ...string) (map[string]interface{}, string) {
	condition := make(map[string]interface{})
	val := reflect.ValueOf(req)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	// 获取 RawQuery 字段
	rawQueryVal := val.FieldByName("RawQuery")
	var query url.Values
	if rawQueryVal.IsValid() && !rawQueryVal.IsNil() {
		query = rawQueryVal.Interface().(url.Values)
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("orm")
		if tag == "" || tag == "-" {
			tag = fieldType.Tag.Get("form")
		}
		if tag == "" || tag == "page" || tag == "page_size" || tag == "-" {
			continue
		}

		op := fieldType.Tag.Get("op")
		if op == "-" {
			continue
		}

		// 只有当 Query 中确实存在这个 Key 时，才允许处理 0 或 空值
		isPresent := false
		if query != nil {
			_, isPresent = query[tag]
		}

		if op == "range" {
			if (field.Kind() == reflect.Slice || field.Kind() == reflect.Array) && field.Len() >= 2 {
				start := field.Index(0).Interface()
				end := field.Index(1).Interface()
				condition[tag+"__gte"] = start
				condition[tag+"__lte"] = end
				continue
			}
		}

		if op == "time_range" {
			if (field.Kind() == reflect.Slice || field.Kind() == reflect.Array) && field.Len() >= 2 {
				s, _ := field.Index(0).Interface().(string)
				e, _ := field.Index(1).Interface().(string)

				if s != "" || e != "" {
					// 尝试解析并转为时间戳
					condition[tag+"__gte"] = self.parseToUnix(s, false) // 开始时间
					condition[tag+"__lt"] = self.parseToUnix(e, true)   // 结束时间
				}
				continue
			}
		}

		// --- 常规类型处理 ---
		switch field.Kind() {
		case reflect.String:
			strVal := strings.TrimSpace(field.String()) // 去除空格
			if strVal != "" {

				// range {value1:1,value2:2}时处理
				if op == "range" && strings.HasPrefix(strVal, "{") {
					var temp map[string]interface{}
					if err := json.Unmarshal([]byte(strVal), &temp); err == nil {
						if v1, ok := temp["value1"]; ok && v1 != "" {
							condition[tag+"__gte"] = v1
						}
						if v2, ok := temp["value2"]; ok && v2 != "" {
							condition[tag+"__lte"] = v2
						}
						continue
					}
				}

				if op == "in" {
					strList := strings.Split(strVal, ",")
					var finalIn []string
					for _, s := range strList {
						s = strings.TrimSpace(s)
						if s != "" {
							finalIn = append(finalIn, s)
						}
					}
					if len(finalIn) > 0 {
						condition[tag+"__in"] = finalIn
					}
				} else {
					key := tag
					if op == "like" {
						key += "__icontains"
					} else if op != "" {
						key = self.mapOp(tag, op)
					}
					condition[key] = strVal
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal := field.Int()
			if intVal > 0 || (intVal == 0 && isPresent) {
				condition[self.mapOp(tag, op)] = intVal
			}
		case reflect.Float32, reflect.Float64:
			floatVal := field.Float()
			if floatVal > 0 || (floatVal == 0 && isPresent) {
				condition[self.mapOp(tag, op)] = floatVal
			}
			/* case reflect.Slice:
			if field.Len() > 0 && op == "in" {
				condition[tag+"__in"] = field.Interface()
			} */
		}
	}

	// 排序
	sort := query.Get("sort")
	order := query.Get("order")
	orderBy := ""
	if sort != "" && order != "" {
		if order == "desc" {
			orderBy = "-" + sort
		} else {
			orderBy = sort
		}
	} else {
		orderBy = sorts[0]
		if sorts[0] == "" {
			orderBy = "-id"
		}
	}
	return condition, orderBy
}

func (self *Base) mapOp(tag, op string) string {
	switch op {
	case "gte":
		return tag + "__gte"
	case "lte":
		return tag + "__lte"
	case "gt":
		return tag + "__gt"
	case "lt":
		return tag + "__lt"
	case "ne":
		return tag + "__ne"
	default:
		return tag
	}
}

// parseToUnix 辅助方法：解析字符串，isEnd 为 true 时，如果只有日期则补全到当天深夜
func (self *Base) parseToUnix(s string, isEnd bool) int64 {
	if s == "" {
		return 0
	}

	nextDay := false
	// 如果是 "2006-01-02" (长度为10)
	if len(s) == 10 {
		if isEnd {
			// s += " 23:59:59" // 结束时间补全到深夜
			s += " 00:00:00"
			nextDay = true
		} else {
			s += " 00:00:00" // 开始时间补全到凌晨
		}
	}

	layout := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		return 0
	}

	if nextDay {
		return t.Unix() + 24*60*60
	}
	return t.Unix()
}

// GetTableNameByTime 根据时间获取表名
func (self *Base) GetTableNameByTime(tableName string, t time.Time) string {
	return fmt.Sprintf("%v%v_%s", TablePrefix, tableName, t.Format("20060102"))
}

// 获取按日期分表的表名称
func (self *Base) GetTableNamesByDateRange(tableName string, startTime, endTime time.Time) []string {
	var tableNames []string

	// 获取本地时区的零点时刻
	startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.Local)
	endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, time.Local)

	// 判断结束时间是否正好是零点
	isEndTimeZero := endTime.Hour() == 0 && endTime.Minute() == 0 && endTime.Second() == 0 && endTime.Nanosecond() == 0

	// 确定结束日期边界
	var endBoundary time.Time
	if isEndTimeZero {
		// 如果是零点，不包含今天，使用 Before
		endBoundary = endDate
	} else {
		// 如果大于零点，包含今天，需要加一天
		endBoundary = endDate.AddDate(0, 0, 1)
	}

	for d := startDate; d.Before(endBoundary); d = d.AddDate(0, 0, 1) {
		curTableName := self.GetTableNameByTime(tableName, d)

		if len(tableNames) == 0 || tableNames[len(tableNames)-1] != curTableName {
			if self.TableExistsInDB(curTableName) {
				tableNames = append(tableNames, curTableName)
			}
		}
	}

	return tableNames
}

// 表是否存在
func (self *Base) TableExistsInDB(tableName string) bool {
	sql := `SELECT COUNT(*) FROM information_schema.tables  WHERE table_schema = DATABASE() AND table_name = ?`

	var count int
	err := self.GetSlaveOrm().Raw(sql, tableName).QueryRow(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// 获取创建分表数据的sql
func (self *Base) GetCreateRecordSql(dataMap map[string]interface{}, tableName string) (string, []interface{}) {
	// 构建 SQL
	fields := make([]string, 0, len(dataMap))
	placeholders := make([]string, 0, len(dataMap))
	args := make([]interface{}, 0, len(dataMap))

	for field, value := range dataMap {
		fields = append(fields, field)
		placeholders = append(placeholders, "?")
		args = append(args, value)
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "))

	return sql, args
}

func CreateBase() *Base {
	return &Base{orm.NewOrm(), orm.NewOrmUsingDB("slave")}
}
