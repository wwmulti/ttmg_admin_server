package task

import (
	"os"
	"path/filepath"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type CleanupTask struct{}

func (ct *CleanupTask) CleanupExports() error {
	exportsDir := "public/exports"

	if _, err := os.Stat(exportsDir); os.IsNotExist(err) {
		logs.Info("导出目录不存在: %s", exportsDir)
		return nil
	}

	now := time.Now()
	cutoffTime := now.AddDate(0, 0, -2)

	err := filepath.Walk(exportsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			err := os.Remove(path)
			if err != nil {
				logs.Error("删除文件失败 %s: %v", path, err)
				return err
			}
			logs.Info("已删除过期文件: %s", path)
		}

		return nil
	})

	if err != nil {
		logs.Error("清理导出目录失败: %v", err)
		return err
	}

	logs.Info("导出目录清理完成")
	return nil
}
