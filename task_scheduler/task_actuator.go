package framework

import "context"

// TaskActuator 任务执行器接口
type TaskActuator interface {

	// Init 任务在被调度前的初始化工作
	Init(ctx context.Context, task *Task) (newTask *Task, err error)

	// Start 开始执行任务
	// ignoreErr 是否忽略任务调度的错误，等待恢复，如果 ignoreErr = false, Start 返回 error 任务会失败
	Start(ctx context.Context, task *Task) (newTask *Task, ignoreErr bool, err error)

	// ExportOutput 导出任务输出，自行处理任务结果
	ExportOutput(ctx context.Context, task *Task) error

	// GetOutput 获取任务数据
	GetOutput(ctx context.Context, task *Task) (data interface{}, err error)

	// Stop 停止任务
	Stop(ctx context.Context, task *Task) error

	// GetTaskStatus 获取异步执行中的任务的状态
	GetAsyncTaskStatus(ctx context.Context, tasks []Task) (status []AsyncTaskStatus, err error)

	// Delete 删除任务
	Delete(ctx context.Context, task *Task) error
}
