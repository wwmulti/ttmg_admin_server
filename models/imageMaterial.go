package models

type ImageMaterialModel struct {
	*Base
}

type ImageMaterial struct {
	Id           int    `json:"id" orm:"auto;column(id)"`                  // 主键
	CategoryId   int    `json:"category_id" orm:"column(category_id)"`     // 目录ID
	StorageType  int    `json:"-" orm:"column(storage_type)"`              // 1-本地, 2-阿里云OSS, 3-aws
	OriginName   string `json:"-" orm:"column(origin_name)"`               // 原始文件名
	RelativePath string `json:"relative_path" orm:"column(relative_path)"` // 相对路径
	AbsoluteUrl  string `json:"-" orm:"column(absolute_url)"`              // 完整访问URL
	Ext          string `json:"-" orm:"column(ext)"`                       // 扩展名
	Size         int64  `json:"-" orm:"column(size)"`                      // 字节大小
	Hash         string `json:"-" orm:"column(hash)"`                      // 文件哈希
	Ctime        int64  `json:"-" orm:"column(ctime);auto_now_add"`        // 创建时间
}

func CreateImageMaterialModel() *ImageMaterialModel {
	return &ImageMaterialModel{CreateBase()}
}

type ImageMaterialListDTO struct {
	Lists []ImageMaterial `json:"lists"`
	Total int64           `json:"total"`
}
