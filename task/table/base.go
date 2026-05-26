package table

import (
	"api/models"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type BaseTask struct{}

// GetDailyTableName 获取按天分表名
// 格式: uk_login_log_20250325
func (t *BaseTask) GetDailyTableName(tableName string, date string) string {
	return fmt.Sprintf("%s_%s", tableName, date)
}

// 获取最近N天的日期列表
func (t *BaseTask) getRecentDates(days int) []string {
	dates := make([]string, days)
	now := time.Now()

	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, i)
		dates[i] = date.Format("20060102") // 日表格式：YYYYMMDD
	}

	return dates
}

// 动态生成最近3天的日表
func (t *BaseTask) GetAllTablesToCreate(tableName string) []string {
	dates := t.getRecentDates(3) // 包括今天在内的3天

	tables := make([]string, 0)
	for _, date := range dates {
		tableName := t.GetDailyTableName(tableName, date)
		tables = append(tables, tableName)
	}
	return tables
}

// createTableIfNotExists 如果表不存在则创建
func (t *BaseTask) CreateTableIfNotExists(tableName string, createSQL string) error {
	// 检查表是否存在
	exists, err := t.tableExists(tableName)
	if err != nil {
		return err
	}

	if exists {
		logs.Info("表 %s 已存在，跳过创建", tableName)
		return nil
	}

	// 创建表
	qb, _ := models.NewWriteQueryBuilder()
	_, createErr := qb.GetMasterDb().Raw(createSQL).Exec()
	if createErr != nil {
		return fmt.Errorf("创建表 %s 失败: %v", tableName, createErr)
	}

	logs.Info("✓ 成功创建表: %s", tableName)
	return nil
}

// tableExists 检查表是否存在
func (t *BaseTask) tableExists(tableName string) (bool, error) {
	var count int64
	sql := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"

	qb, _ := models.NewWriteQueryBuilder()
	err := qb.GetMasterDb().Raw(sql, tableName).QueryRow(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 获取数据表的schema
func (b *BaseTask) GetTableSchema(tableName string) (string, error) {
	// 查询表结构
	sql := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	var result []orm.Params

	qb, _ := models.NewWriteQueryBuilder()
	_, err := qb.GetMasterDb().Raw(sql).Values(&result)
	if err != nil {
		logs.Error("获取表结构失败: %v", err)
		return "", err
	}

	if len(result) == 0 {
		return "", fmt.Errorf("未找到表 %s 的信息", tableName)
	}

	// 方法1：直接获取 Create Table 字段
	if createTable, ok := result[0]["Create Table"]; ok {
		if createTableStr, ok := createTable.(string); ok {
			return createTableStr, nil
		}
	}

	return "", fmt.Errorf("无法解析表结构")
}
