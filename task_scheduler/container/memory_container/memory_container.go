package memeorycontainer

import (
	"context"

	framework "github.com/memory-overflow/go-common-library/task_scheduler"
)

// MemeoryContainer 内存型任务容器，优先：快读快写，缺点：不可持久化，
type MemeoryContainer interface {
	framework.TaskContainer

	// AddRunningTask 向容器添加正在运行中的任务
	AddRunningTask(ctx context.Context, task framework.Task) (err error)
}
