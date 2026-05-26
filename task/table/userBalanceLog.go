package table

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

type UserBalanceLogTask struct {
	*BaseTask
}

const UserBalanceLogTableName = "uk_user_balance_log"

func (t *UserBalanceLogTask) GetCreateTableSQL(tableName string) string {
	sql, err := t.GetTableSchema(UserBalanceLogTableName)
	if err != nil {
		logs.Error("创建uk_user_balance_log错误:%v", err)
	}

	sql = strings.Replace(sql, UserBalanceLogTableName, tableName, 1)
	return sql
}

// 执行
func (t *UserBalanceLogTask) Run() {
	tableNames := t.GetAllTablesToCreate(UserBalanceLogTableName)
	for _, tableName := range tableNames {
		sql := t.GetCreateTableSQL(tableName)
		// 检查错误并处理
		if err := t.CreateTableIfNotExists(tableName, sql); err != nil {
			logs.Error("创建用户余额变动日志表 %s 失败: %v", tableName, err)
		} else {
			logs.Info("创建用户余额变动日志表: %s", tableName)
		}
	}
}
