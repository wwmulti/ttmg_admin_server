package task

import (
	"context"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/task"
)

type TaskManager struct {
	tasks map[string]*task.Task
}

func TaskJobManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*task.Task),
	}
}

// 注册任务
// spec 秒 分 时 日 月 周
func (tm *TaskManager) Register(name, spec string, fn func(ctx context.Context) error) {
	t := task.NewTask(name, spec, fn)
	tm.tasks[name] = t
	task.AddTask(name, t)
}

// 注册简单任务（不需要 context）
// spec 秒 分 时 日 月 周
func (tm *TaskManager) RegisterSimple(name, spec string, fn func() error) {
	wrappedFn := func(ctx context.Context) error {
		return fn()
	}
	tm.Register(name, spec, wrappedFn)
}

// 启动所有任务
func (tm *TaskManager) Start() {
	tm.setupTasks()
	task.StartTask()
	logs.Info("启动了 %d 个定时任务", len(tm.tasks))
	// 默认启动时候创建后期每天创建
}

// 停止所有任务
func (tm *TaskManager) Stop() {
	task.StopTask()
}

// 设置任务
func (tm *TaskManager) setupTasks() {
	// 零点检测
	tm.Register("ZeroTask", "0 0 0 */1 * *", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(1 * time.Second):
			logs.Info("每天零点执行任务")
			return nil
		}
	})
}
