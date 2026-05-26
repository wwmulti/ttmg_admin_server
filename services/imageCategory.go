package services

import (
	"api/models"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
)

type ImageCategoryService struct{}

// 1. 创建目录
func (s *ImageCategoryService) Create(name, dir string) error {
	categoryModel := models.CreateImageCategoryModel()

	// 检查 dir 是否已存在（物理路径唯一性）
	exists := categoryModel.QueryTable(new(models.ImageCategory)).Filter("dir", dir).Exist()
	if exists {
		return fmt.Errorf("MuLuLuJingYiCunZai")
	}
	_, err := categoryModel.Insert(&models.ImageCategory{
		Name:   name,
		Dir:    dir,
		Status: 1,
	})

	return err
}

// 2. 修改目录名称（通常只改显示名称 Name，不改物理路径 Dir）
func (s *ImageCategoryService) UpdateName(id int, newName string) error {
	categoryModel := models.CreateImageCategoryModel()
	_, err := categoryModel.QueryTable(new(models.ImageCategory)).Filter("id", id).Update(orm.Params{
		"name": newName,
	})
	return err
}

// 3. 删除目录
func (s *ImageCategoryService) Delete(id int) error {
	categoryModel := models.CreateImageCategoryModel()
	materialModel := models.CreateImageMaterialModel()

	// 安全检查：如果该目录下还有图片，禁止删除
	hasImages := materialModel.QueryTable(new(models.ImageMaterial)).Filter("category_id", id).Exist()
	if hasImages {
		return fmt.Errorf("GaiMuLuXiaRengYouTuPianQingXianShanChuHuoQianYiTuPianHouZaiShanChuMuLu")
	}

	_, err := categoryModel.QueryTable(new(models.ImageCategory)).Filter("id", id).Delete()
	return err
}

// 4. 获取所有可用目录列表
func (s *ImageCategoryService) GetAllCategories() ([]models.ImageCategory, error) {
	var list []models.ImageCategory
	_, err := models.CreateImageCategoryModel().QueryTable(new(models.ImageCategory)).
		Filter("status", 1).
		All(&list)
	return list, err
}
