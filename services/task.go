package services

import (
	"api/config"
	"api/models"
	"api/utils"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type TaskService struct{}

const (
	TaskTypeExport  = "export"  // 数据导出
	TaskTypeRepair  = "repair"  // 数据修复
	TaskTypeRefresh = "refresh" // 刷新数据 (HTTP)

	TaskStatusProcessing = "processing"
	TaskStatusCompleted  = "completed"
	TaskStatusError      = "error"
)

// TaskStatus 对应前端 getMultiTaskProcess 的数据结构
type TaskStatus struct {
	Status       string      `json:"status"`                  // 状态：processing, completed, error
	TipType      string      `json:"tipType,omitempty"`       // 提醒方式：showMsgBox (存入消息列表), toast (气泡), showDownloadMsg (下载)
	Tips         string      `json:"tips"`                    // 显示给用户看的文字内容
	Process      int         `json:"process,omitempty"`       // 进度百分比
	DownloadLink string      `json:"download_link,omitempty"` // 下载链接
	Result       interface{} `json:"result,omitempty"`        // 只有 tipType 为 showMsgBox 时才需要详细结构(result.alertBoxString 若值为 settlement 则触发前端结算详情弹窗result.autoAlertInt 1: 自动弹出弹窗；0: 仅存入)
}

// AddTask 通用添加任务入口
func (s *TaskService) AddTask(taskType string, params map[string]interface{}) (string, error) {
	taskId := fmt.Sprintf("task_%d", time.Now().UnixNano())

	// 初始化任务状态为处理中
	s.updateStatus(taskId, TaskStatus{
		Status:  TaskStatusProcessing,
		Tips:    "starting...",
		Process: 5,
	})

	// 异步分发处理
	go func() {
		var err error
		switch taskType {
		case TaskTypeExport:
			err = s.handleExportTask(taskId, params)
		case TaskTypeRefresh:
			err = s.handleRefreshTask(taskId, params)
		case TaskTypeRepair:
			// 数据修复先不处理
			err = nil
		default:
			logs.Error("异步任务未知任务类型: %s", taskType)
		}

		if err != nil {
			s.updateStatus(taskId, TaskStatus{
				Status: TaskStatusError,
				Tips:   err.Error(),
			})
		}
	}()

	return taskId, nil
}

type HeaderItem struct {
	Field string
	Title string
}

// ExportConfig 导出配置
type ExportConfig struct {
	FileName   string // 文件名
	Model      interface{}
	Header     []HeaderItem                                  // Excel 表头
	DataQuery  func(results interface{}) error               // 数据获取匿名函数
	DataFormat func(item interface{}) map[string]interface{} // 单行数据格式化
}

// handleExportTask 导出cvs
func (s *TaskService) handleExportTask(taskId string, params map[string]interface{}) error {
	conf := params["config"].(ExportConfig)

	// 提取表头
	displayTitles := make([]string, 0, len(conf.Header))
	for _, item := range conf.Header {
		displayTitles = append(displayTitles, item.Title)
	}

	csvTool, err := utils.NewCsvExport(displayTitles, conf.FileName)
	if err != nil {
		return err
	}

	// 使用反射动态创建模型切片
	modelType := reflect.TypeOf(conf.Model)
	sliceType := reflect.SliceOf(modelType)
	results := reflect.New(sliceType).Interface()

	// 执行数据查询 (匿名函数)
	s.updateStatus(taskId, TaskStatus{Status: TaskStatusProcessing, Tips: "processing...", Process: 30})
	if err := conf.DataQuery(results); err != nil {
		logs.Error("导出任务数据查询错误: %v", err)
		return err
	}

	// 遍历处理并生成
	val := reflect.ValueOf(results).Elem()
	for i := 0; i < val.Len(); i++ {
		rowMap := conf.DataFormat(val.Index(i).Interface())
		record := make([]string, 0, len(conf.Header))
		for _, h := range conf.Header {
			record = append(record, fmt.Sprint(rowMap[h.Field]))
		}
		csvTool.AddRow(record)
	}
	downloadUrl := csvTool.GetFileLink("http://localhost:8082")

	// 完成任务
	s.updateStatus(taskId, TaskStatus{
		Status:       TaskStatusCompleted,
		TipType:      "showDownloadMsg",
		Tips:         "success",
		DownloadLink: downloadUrl,
	})
	return nil
}

// handleRefreshTask 刷新缓存
func (s *TaskService) handleRefreshTask(taskId string, params map[string]interface{}) error {
	var table string
	if params["table"] != nil {
		table = params["table"].(string)
	}

	s.updateStatus(taskId, TaskStatus{Status: TaskStatusProcessing, Tips: "processing...", Process: 50})

	packages, err := (&PackageService{}).AllPackage()
	if err != nil {
		logs.Error("获取所有包失败：%v", err)
		return fmt.Errorf("WeiZhiDeCuoWu")
	}

	if len(packages) == 0 {
		logs.Error("没有包存在")
		return fmt.Errorf("MeiYouBaoPeiZhiCunZai")
	}

	var errString string
	for _, packageInfo := range packages {
		err := s.refreshTable(packageInfo, table)
		if err != nil {
			errString = errString + packageInfo.Title + "\n"
		}
	}

	if errString != "" {
		logs.Error("刷新包缓存表失败：%v", errString)
		return fmt.Errorf("ShuaXinShiBai")
	}

	s.updateStatus(taskId, TaskStatus{
		Status:  TaskStatusCompleted,
		TipType: "showMsgBox",
		Tips:    "success",
	})
	return nil
}

// 刷新请求
func (s *TaskService) refreshTable(packageInfo models.Package, table string) error {
	secret := config.System.CurlSecretKey.Value
	params := map[string]interface{}{
		"package_id": packageInfo.Id,
		"table":      table,
	}
	request, _ := utils.SignMap(params, secret)
	packageService := PackageService{}
	resp := RequestPackageUrlRequest{
		PackageId: packageInfo.Id,
		Params:    request,
		Url:       "/data/refresh/refreshTables",
		Type:      int(config.CurlRequestTypePost),
	}
	response, err := packageService.RequestPackageUrl(resp)
	if err != nil {
		logs.Error("后台刷新缓存表失败：table%v", err)
		return fmt.Errorf("ShuXinHuanCunShiBai")
	}

	type Result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	ret := response.String()
	result, _ := utils.FromJson[Result](ret)
	if result.Code != 200 {
		logs.Error("后台刷新缓存表失败：table %v, err： %v", table, result.Msg)
		return fmt.Errorf("ShuXinHuanCunShiBai")
	}

	if !response.OK() {
		logs.Error("后台刷新缓存表请求失败：table %v", table)
		return fmt.Errorf("ShuXinHuanCunShiBai")
	}
	return nil
}

func (s *TaskService) updateStatus(taskId string, status TaskStatus) {
	data, _ := json.Marshal(status)
	utils.GetRedisClient().Set(context.Background(), config.RedisKeyName.TaskProcess+taskId, data, 3*24*time.Hour)
}

// GetMultiTaskProcess 获取多个任务进度
func (c *TaskService) GetMultiTaskProcess(taskIdsStr string) map[string]interface{} {
	ids := strings.Split(taskIdsStr, ",")
	results := make(map[string]interface{})

	redisCache := utils.GetRedisClient()
	for _, id := range ids {
		val, err := redisCache.Get(context.Background(), config.RedisKeyName.TaskProcess+id).Result()
		if err == nil {
			var status interface{}
			json.Unmarshal([]byte(val), &status)
			results[id] = status
		}
	}
	return results
}

// 导出示例(游戏表)
/*exportConfig := ExportConfig{
    FileName: "game",
	Model: models.Game{},
	Header: []HeaderItem{
		{Field: "id", Title: "ID"},
		{Field: "name", Title: "游戏名"},
		{Field: "code_rule", Title: "打码规则"},
	},
	DataQuery: func(res interface{}) error {
		_, err := gameModel.Where(gameModel.QueryTable(new(models.Game)).Filter("status", 1), condition).All(res) // (这里是res 不能是&res)
		return err
	},
	DataFormat: func(item interface{}) map[string]interface{} {
		data := item.(models.Game)

		codeRule := ""
		switch data.CodeRule {
		case 1:
			codeRule = "流水"
		case 2:
			codeRule = "净赢"
		case 3:
			codeRule = "赢金打码"
		}

		return map[string]interface{}{
			"id":        data.Id,
			"name":      data.Name,
			"code_rule": codeRule,
		}
	},
}
taskId, _ := (&TaskService{}).AddTask(TaskTypeExport, map[string]interface{}{"config": exportConfig})
fmt.Println("Task ID:", taskId)*/

// 刷新内存mysql表示例
// taskId, _ := (&TaskService{}).AddTask(TaskTypeRefresh, map[string]interface{}{})
// fmt.Println("Task ID:", taskId)
