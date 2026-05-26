package table

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

type UserOperateLogTask struct {
	*BaseTask
}

const UserOperateLogTableName = "uk_user_operate_log"

func (t *UserOperateLogTask) GetCreateTableSQL(tableName string) string {
	sql, err := t.GetTableSchema(UserOperateLogTableName)
	if err != nil {
		logs.Error("创建uk_user_operate_log错误:%v", err)
	}

	sql = strings.Replace(sql, UserOperateLogTableName, tableName, 1)
	return sql
}

// 执行
func (t *UserOperateLogTask) Run() {
	tableNames := t.GetAllTablesToCreate(UserOperateLogTableName)
	for _, tableName := range tableNames {
		sql := t.GetCreateTableSQL(tableName)
		// 检查错误并处理
		if err := t.CreateTableIfNotExists(tableName, sql); err != nil {
			logs.Error("创建用户操作日志表表 %s 失败: %v", tableName, err)
		} else {
			logs.Info("成功创建用户操作日志表表: %s", tableName)
		}
	}
}
