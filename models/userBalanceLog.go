package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type UserBalanceLogModel struct {
	*Base
}

type UserBalanceLog struct {
	Id             int     `json:"id" orm:"auto;column(id)"`                             // 主键，自增
	Balance        float64 `json:"balance" orm:"column(balance);null"`                   // 当前余额
	Amount         float64 `json:"amount" orm:"column(amount);null"`                     // 变动金额
	PackageId      int     `json:"package_id" orm:"column(package_id);null"`             // 包ID
	UserId         int     `json:"user_id" orm:"column(user_id);null"`                   // 用户ID
	RoleId         int64   `json:"role_id" orm:"column(role_id);null"`                   // 用户RoleId
	Type           int     `json:"type" orm:"column(type);null"`                         // 类型
	ActivityTypeId int     `json:"activity_type_id" orm:"column(activity_type_id);null"` // 活动类型
	ActivityId     int     `json:"activity_id" orm:"column(activity_id);null"`           // 活动ID
	ActivityTitle  string  `json:"activity_title" orm:"column(activity_title);null"`     // 活动名称
	PlatformId     int     `json:"platform_id" orm:"column(platform_id);null"`           // 游戏平台ID
	PlatformTitle  string  `json:"platform_title" orm:"column(platform_title);null"`     // 平台名称
	GameTypeId     int     `json:"game_type_id" orm:"column(game_type_id);null"`         // 游戏类型
	GameId         int     `json:"game_id" orm:"column(game_id);null"`                   // 游戏ID
	GameTitle      string  `json:"game_title" orm:"column(game_title);null"`             // 游戏名称
	Mark           string  `json:"mark" orm:"column(mark);null"`                         // 备注
	CreatedTime    int64   `json:"created_time" orm:"column(created_time);null"`         // 生成时间
}

func (m *UserBalanceLog) TableName() string {
	return (&UserBalanceLogModel{}).GetCurrentTableName()
}

// GetCurrentTableName 获取当前日期的表名
func (m *UserBalanceLogModel) GetCurrentTableName() string {
	tableName := fmt.Sprintf("uk_user_balance_log_%s", time.Now().Format("20060102"))
	return tableName
}

// RecordLog 记录登录日志，支持传入事务对象
// RecordLog 记录登录日志，支持传入事务对象
func (m *UserBalanceLogModel) RecordLog(log UserBalanceLog, tx ...orm.TxOrmer) error {
	tableName := m.GetCurrentTableName()

	// 定义字段和值的映射
	dataMap := map[string]interface{}{
		"balance":          log.Balance,
		"amount":           log.Amount,
		"package_id":       log.PackageId,
		"user_id":          log.UserId,
		"role_id":          log.RoleId,
		"type":             log.Type,
		"activity_type_id": log.ActivityTypeId,
		"activity_id":      log.ActivityId,
		"activity_title":   log.ActivityTitle,
		"platform_id":      log.PlatformId,
		"platform_title":   log.PlatformTitle,
		"game_type_id":     log.GameTypeId,
		"game_id":          log.GameId,
		"game_title":       log.GameTitle,
		"mark":             log.Mark,
		"created_time":     log.CreatedTime,
	}

	sql, args := m.GetCreateRecordSql(dataMap, tableName)

	// 执行插入
	var err error
	if len(tx) > 0 && tx[0] != nil {
		_, err = tx[0].Raw(sql, args...).Exec()
	} else {
		o := m.GetMasterDb()
		_, err = o.Raw(sql, args...).Exec()
	}

	if err != nil {
		logs.Error("插入用户余额日志失败: %v, table: %s, user_id: %d",
			err, tableName, log.UserId)
		return err
	}

	return nil
}

func CreateUserBalanceLogModel() *UserBalanceLogModel {
	return &UserBalanceLogModel{CreateBase()}
}
