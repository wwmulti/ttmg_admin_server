package models

type SafeQuestionModel struct {
	*Base
}

type SafeQuestion struct {
	Id        int    `json:"id" orm:"auto;column(id)"`            // 主键，自增
	PackageId string `json:"package_id" orm:"column(package_id)"` // 包ID
	Question  string `json:"question" orm:"column(question)"`     // 安全问题
	IsDeleted int    `json:"is_deleted" orm:"column(is_deleted)"` // 是否删除
}

func CreateSafeQuestionModel() *SafeQuestionModel {
	return &SafeQuestionModel{CreateBase()}
}
