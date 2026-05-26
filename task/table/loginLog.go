package table

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

type LoginLogTask struct {
	*BaseTask
}

const LoginLogTableName = "uk_login_log"

func (t *LoginLogTask) GetCreateTableSQL(tableName string) string {
	sql, err := t.GetTableSchema(LoginLogTableName)
	if err != nil {
		logs.Error("创建uk_login_log错误:%v", err)
	}

	sql = strings.Replace(sql, LoginLogTableName, tableName, 1)
	return sql
}

// 执行
func (t *LoginLogTask) Run() {
	tableNames := t.GetAllTablesToCreate(LoginLogTableName)
	for _, tableName := range tableNames {
		sql := t.GetCreateTableSQL(tableName)
		// 检查错误并处理
		if err := t.CreateTableIfNotExists(tableName, sql); err != nil {
			logs.Error("创建uk_login_log %s 失败: %v", tableName, err)
		} else {
			logs.Info("创建uk_login_log: %s", tableName)
		}
	}
}
