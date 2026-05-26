package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type CsvExport struct {
	File     *os.File
	Writer   *csv.Writer
	FileName string // 导出的文件名
	SubPath  string // 相对路径 e.g. exports/260327/
}

// NewCsvExport
func NewCsvExport(head []string, fileNamePrefix string) (*CsvExport, error) {
	// 构建路径
	dateDir := time.Now().Format("060102")
	subPath := filepath.Join("public", "exports", dateDir)

	// 确保目录存在
	if _, err := os.Stat(subPath); os.IsNotExist(err) {
		err = os.MkdirAll(subPath, 0777)
		if err != nil {
			logs.Error("导出csv创建目录失败: %v", err)
			return nil, err
		}
	}

	// 生成文件名
	fileName := fmt.Sprintf("%s_%s.csv", fileNamePrefix, time.Now().Format("060102150405"))
	fullPath := filepath.Join(subPath, fileName)

	// 打开文件
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logs.Error("导出csv打开文件失败: %v", err)
		return nil, err
	}

	// 写入 UTF-8 BOM (防止 Excel 打开乱码)
	f.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(f)

	// 写入表头
	if err := writer.Write(head); err != nil {
		f.Close()
		logs.Error("导出csv写入表头失败: %v", err)
		return nil, err
	}
	writer.Flush()

	return &CsvExport{
		File:     f,
		Writer:   writer,
		FileName: fileName,
		SubPath:  filepath.Join("public/exports", dateDir),
	}, nil
}

// AddRow
func (c *CsvExport) AddRow(data []string) error {
	if c.Writer != nil {
		if err := c.Writer.Write(data); err != nil {
			return err
		}
		c.Writer.Flush() // 保证数据实时写入磁盘
	}
	return nil
}

// AddList
func (c *CsvExport) AddList(data [][]string) error {
	for _, row := range data {
		if err := c.AddRow(row); err != nil {
			return err
		}
	}
	return nil
}

// GetFileLink （host 传域名，如 http://api.xxx.com）
func (c *CsvExport) GetFileLink(host string) string {
	if c.File != nil {
		c.File.Close() // 返回链接前关闭文件流
	}

	// 处理路径斜杠
	urlPath := filepath.Join(c.SubPath, c.FileName)
	urlPath = strings.ReplaceAll(urlPath, "\\", "/")

	// 构建最终下载链接
	link := fmt.Sprintf("%s/%s", strings.TrimRight(host, "/"), urlPath)

	// 加上端口参数
	link = fmt.Sprintf("%s?p=%s", link, "8888")

	return link
}
