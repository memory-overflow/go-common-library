package persistcontainer

import framework "github.com/memory-overflow/go-common-library/task_scheduler"

// PersistContainer 可持久化任务容器，优点：可持久化存储，缺点：依赖db、需要扫描表，对 db 压力比较大。
type PersistContainer framework.TaskContainer
