package models

type ImageCategoryModel struct {
	*Base
}

type ImageCategory struct {
	Id     int    `json:"id" orm:"auto;column(id)"`    // 主键
	Name   string `json:"name" orm:"column(name)"`     // 目录名称
	Dir    string `json:"dir" orm:"column(dir)"`       // 实际存储目录名
	Status int    `json:"status" orm:"column(status)"` // 状态 1正常 0禁用
}

func CreateImageCategoryModel() *ImageCategoryModel {
	return &ImageCategoryModel{CreateBase()}
}
