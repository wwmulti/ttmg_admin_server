package table

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

type GameWaterLogTask struct {
	*BaseTask
}

const GameWaterLogTableName = "uk_game_water_log"

func (t *GameWaterLogTask) GetCreateTableSQL(tableName string) string {
	sql, err := t.GetTableSchema(GameWaterLogTableName)
	if err != nil {
		logs.Error("创建uk_game_water_log错误:%v", err)
	}

	sql = strings.Replace(sql, GameWaterLogTableName, tableName, 1)
	return sql
}

// 执行
func (t *GameWaterLogTask) Run() {
	tableNames := t.GetAllTablesToCreate(GameWaterLogTableName)
	for _, tableName := range tableNames {
		sql := t.GetCreateTableSQL(tableName)
		// 检查错误并处理
		if err := t.CreateTableIfNotExists(tableName, sql); err != nil {
			logs.Error("创建流水表 %s 失败: %v", tableName, err)
		} else {
			logs.Info("创建流水表: %s", tableName)
		}
	}
}
