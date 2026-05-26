package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
	"mime/multipart"
	"strconv"
	"strings"
)

type ImageController struct {
	BaseController
}

type UploadImageRequest struct {
	Images     []*multipart.FileHeader `form:"images"`
	CategoryId int                     `form:"category_id"`
}

type UploadImageResponse struct {
	JSONResponseTpl
	Data []models.ImageMaterial `json:"data"`
}

// @Summary	批量上传图片
// @Param request body controllers.UploadImageRequest true "请求参数"
// @Success 200 {object} controllers.UploadImageResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /uploadAction [post]
func (c *ImageController) UploadAction() {
	files, _ := c.GetFiles("images") // 获取批量文件
	catId, _ := c.GetInt("category_id", 0)

	service := services.ImageService{}
	data, err := service.BatchUpload(files, catId)
	if err != nil {
		c.JSONError(500, err.Error())
		return
	}
	c.JSONSuccess(data)
}

type DeleteImagesRequest struct {
	Ids string `json:"ids"` // 图片ID列表，英文逗号分隔
}

// @Summary 批量删除图片
// @Param request body controllers.DeleteImagesRequest true "请求参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /batchDelete [post]
func (c *ImageController) BatchDelete() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "参数解析失败")
	}

	idsStr := getString(params, "ids")
	if idsStr == "" {
		c.JSONError(500, "请选择要删除的图片")
	}

	// 2. 将字符串解析为 int 切片
	idStrs := strings.Split(idsStr, ",")
	var ids []int
	for _, s := range idStrs {
		if id, err := strconv.Atoi(s); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		c.JSONError(500, "无效的ID参数")
	}

	// 3. 调用 Service 执行批量删除
	service := services.ImageService{}
	err = service.BatchDelete(ids)
	if err != nil {
		c.JSONError(500, err.Error())
		return
	}
	c.JSONSuccess(nil)
}

type GetAllCategoriesResponse struct {
	JSONResponseTpl
	Data []models.ImageCategory `json:"data"`
}

// @Summary 获取目录分类列表
// @Success 200 {object} controllers.GetAllCategoriesResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /getCategoryList [get]
func (c *ImageController) GetCategoryList() {
	service := services.ImageCategoryService{}
	data, err := service.GetAllCategories()
	if err != nil {
		c.JSONError(500, err.Error())
		return
	}
	c.JSONSuccess(data)
}

type GetImageListResponse struct {
	JSONResponseTpl
	Data models.ImageMaterialListDTO `json:"data"`
}

// @Summary 获取图片列表
// @Param page query int true "页数"
// @Param page_size query int true "行数"
// @Param keyword query string false "图片名"
// @Param category_id query int false "分类目录id"
// @Success 200 {object} controllers.GetImageListResponse
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /getImageList [get]
func (c *ImageController) GetImageList() {
	page, limit := c.GetPagination()
	keyword := c.GetString("keyword", "")
	categoryId, _ := c.GetInt("category_id", 0)

	service := services.ImageService{}
	data, err := service.GetImageList(services.GetImageListRequest{
		Page:       page,
		Limit:      limit,
		Keyword:    keyword,
		CategoryId: categoryId,
	})
	if err != nil {
		c.JSONError(500, err.Error())
		return
	}
	c.JSONSuccess(data)
}

type CreateCategoryRequest struct {
	Name string `json:"name"` // 目录名称
	Dir  string `json:"dir"`  // 目录路径
}

// @Summary 创建分类新目录
// @Param request body controllers.CreateCategoryRequest true "请求json参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /createCategory [post]
func (c *ImageController) CreateCategory() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "参数解析失败")
	}

	name := getString(params, "name")
	dir := getString(params, "dir")

	if name == "" || dir == "" {
		c.JSONError(500, "名称和目录不能为空")
	}

	service := services.ImageCategoryService{}
	err = service.Create(name, dir)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(nil)
}

type UpdateCategoryRequest struct {
	Id   int    `json:"id"`   // 目录ID
	Name string `json:"name"` // 目录名称
}

// @Summary 修改分类目录名称
// @Param request body controllers.UpdateCategoryRequest true "请求json参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /updateCategory [post]
func (c *ImageController) UpdateCategory() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "参数解析失败")
	}

	id := getInt(params, "id")
	newName := getString(params, "name")

	if id <= 0 || newName == "" {
		c.JSONError(500, "参数错误")
	}

	service := services.ImageCategoryService{}
	err = service.UpdateName(id, newName)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(nil)
}

type DeleteCategoryRequest struct {
	Id int `json:"id"` // 目录ID
}

// @Summary 删除分类目录
// @Param request body controllers.DeleteCategoryRequest true "请求json参数"
// @Success 200 {object} controllers.JSONResponseTpl
// @Failure 500 {"code": 0,"msg": "error message"}
// @router /deleteCategory [post]
func (c *ImageController) DeleteCategory() {
	var params map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if err != nil {
		c.JSONError(500, "参数解析失败")
	}

	id := getInt(params, "id")
	if id <= 0 {
		c.JSONError(500, "无效ID")
	}

	service := services.ImageCategoryService{}
	err = service.Delete(id)
	if err != nil {
		c.JSONError(500, err.Error())
	}
	c.JSONSuccess(nil)
}
