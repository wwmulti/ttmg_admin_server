package services

import (
	"api/models"
	"api/utils"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type ImageService struct{}

var (
	storageType int = utils.StorageLocal // 默认使用本地存储
	allowTypes      = map[string]bool{   // 指定图片格式
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	maxFileSize int64 = 5 * 1024 * 1024 // 最大上传限制5MB
)

// BatchUpload 批量上传图片
func (s *ImageService) BatchUpload(files []*multipart.FileHeader, categoryId int) ([]models.ImageMaterial, error) {
	model := models.CreateImageMaterialModel()
	factory := &utils.StorageFactory{}

	currentStorage := storageType
	driver, err := factory.GetDriver(currentStorage)
	if err != nil {
		return nil, err
	}

	// 获取目录名
	var category models.ImageCategory
	models.CreateImageCategoryModel().QueryTable(new(models.ImageCategory)).Filter("id", categoryId).One(&category)
	if category.Dir == "" {
		category.Dir = "default"
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("MeiYouWenJianBeiShangChuan")
	}

	var savedImages []models.ImageMaterial
	for _, fh := range files {
		f, _ := fh.Open()
		content, _ := io.ReadAll(f)
		hashVal, _ := utils.GenerateImageHashFromReader(f)
		f.Close()

		mimeType := http.DetectContentType(content)
		if !allowTypes[mimeType] {
			return nil, fmt.Errorf("BaoHanYouBuZhiChiDeTuPianGeShi")
		}

		if fh.Size > maxFileSize {
			return nil, fmt.Errorf("WenJianDaXiaoChaoGuoZuiDaShangChuanXianZhi")
		}

		var existImg models.ImageMaterial
		err = model.QueryTable(&models.ImageMaterial{}).Filter("hash", hashVal).One(&existImg)
		if err == nil {
			// 数据库已存在该图片，直接引用，跳过 Upload 步骤
			savedImages = append(savedImages, existImg)
			continue
		}

		absUrl, relPath, err := driver.Upload(content, category.Dir, fh.Filename)
		if err != nil {
			continue
		}

		img := models.ImageMaterial{
			CategoryId:   categoryId,
			StorageType:  currentStorage,
			OriginName:   fh.Filename,
			RelativePath: relPath,
			AbsoluteUrl:  absUrl,
			Ext:          filepath.Ext(fh.Filename),
			Size:         fh.Size,
			Hash:         hashVal,
			Ctime:        time.Now().Unix(),
		}
		if _, err := model.Insert(&img); err == nil {
			savedImages = append(savedImages, img)
		}
	}
	return savedImages, nil
}

// BatchDelete 批量删除图片
func (s *ImageService) BatchDelete(ids []int) error {
	model := models.CreateImageMaterialModel()
	var list []models.ImageMaterial
	model.QueryTable(new(models.ImageMaterial)).Filter("id__in", ids).All(&list)

	factory := &utils.StorageFactory{}
	for _, img := range list {
		driver, err := factory.GetDriver(img.StorageType)
		if err == nil {
			_ = driver.Delete(img.RelativePath)
			model.Delete(&img)
		}
	}
	return nil
}

type GetImageListRequest struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Keyword    string `json:"keyword"`
	CategoryId int    `json:"category_id"`
}

// 获取图片列表
func (s *ImageService) GetImageList(param GetImageListRequest) (models.ImageMaterialListDTO, error) {
	model := models.CreateImageMaterialModel()
	condition := make(map[string]interface{})
	if param.CategoryId > 0 {
		condition["category_id"] = param.CategoryId
	}

	keyWord := strings.TrimSpace(param.Keyword)
	if keyWord != "" {
		condition["origin_name__icontains"] = keyWord
	}

	data, total, err := model.GetPageList(&models.ImageMaterial{}, condition, param.Page, param.Limit, "-id")
	if err != nil {
		return models.ImageMaterialListDTO{}, fmt.Errorf("WeiZhiDeCuoWu")
	}

	lists := data.([]models.ImageMaterial)

	return models.ImageMaterialListDTO{Lists: lists, Total: total}, nil
}
