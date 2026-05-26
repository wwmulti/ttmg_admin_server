package models

import "github.com/beego/beego/v2/client/orm"

type ActiveVipLvLogModel struct {
	*Base
}

type ActiveVipLvLog struct {
	Id       int   `json:"id" orm:"auto;column(id)"`          // 主键
	ActiveId int   `json:"active_id" orm:"column(active_id)"` // 活动id
	UserId   int   `json:"user_id" orm:"column(user_id)"`     // 用户id
	Lv       int   `json:"lv" orm:"column(lv)"`               // 达成等级
	Ctime    int64 `json:"c_time" orm:"column(c_time)"`       // 达成时间
}

func CreateActiveVipLvLogModel() *ActiveVipLvLogModel {
	return &ActiveVipLvLogModel{CreateBase()}
}

func (m *ActiveVipLvLogModel) RecordLog(request ActiveVipLvLog, tx ...orm.TxOrmer) error {
	if len(tx) > 0 {
		_, err := tx[0].Insert(&request)
		if err != nil {
			return err
		}
	} else {
		logModel := CreateActiveVipLvLogModel()
		_, err := logModel.Insert(&request)
		if err != nil {
			return err
		}
	}
	return nil
}
