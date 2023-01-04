package actuator

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	framework "github.com/memory-overflow/go-common-library/task_scheduler"
	service "github.com/memory-overflow/go-common-library/task_scheduler/example/add_service"
	exampletask "github.com/memory-overflow/go-common-library/task_scheduler/example/task"
)

// ExampleActuator 示例执行器
type ExampleActuator struct {
}

// MakeExampleActuator 构造执行器
func MakeExampleActuator() *ExampleActuator {
	return &ExampleActuator{}
}

// Init 任务在被调度前的初始化工作
func (e *ExampleActuator) Init(ctx context.Context, task *framework.Task) (
	newTask *framework.Task, err error) {
	return task, nil
}

// Run 执行任务
func (e *ExampleActuator) Start(ctx context.Context, task *framework.Task) (
	newTask *framework.Task, ignoreErr bool, err error) {
	body := fmt.Sprintf("{\"A\": %d, \"B\": %d}",
		task.TaskItem.(exampletask.ExampleTask).A, task.TaskItem.(exampletask.ExampleTask).B)
	rsp, err := http.Post("http://127.0.0.1:8000/add", "application/json", strings.NewReader(string(body)))
	if err != nil {
		return nil, false, err
	}
	bs, _ := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	var r service.AddResponse
	json.Unmarshal(bs, &r)
	rtask := task.TaskItem.(exampletask.ExampleTask)
	rtask.TaskId = uint32(r.TaskId)
	task.TaskItem = rtask
	log.Printf("start success, taskid: %d\n", r.TaskId)
	return task, false, nil
}

// ExportOutput 导出任务输出，自行处理任务结果
func (e *ExampleActuator) ExportOutput(ctx context.Context, task *framework.Task) error {
	rtask := task.TaskItem.(exampletask.ExampleTask)
	body := fmt.Sprintf("{\"TaskId\": %d}", rtask.TaskId)
	rsp, err := http.Post("http://127.0.0.1:8000/result", "application/json", strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	bs, _ := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	var r service.ResultResponse
	json.Unmarshal(bs, &r)
	log.Printf("finished task %d: %d + %d = %d \n", rtask.TaskId, rtask.A, rtask.B, r.Res)
	return nil
}

// Stop 停止任务
func (e *ExampleActuator) Stop(ctx context.Context, task *framework.Task) error {
	return nil
}

// GetAsyncTaskStatus 获取任务状态
func (e *ExampleActuator) GetAsyncTaskStatus(ctx context.Context, tasks []framework.Task) (
	status []framework.AsyncTaskStatus, err error) {
	status = make([]framework.AsyncTaskStatus, len(tasks))
	wg := sync.WaitGroup{}
	for i := range tasks {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			body := fmt.Sprintf("{\"TaskId\": %d}", tasks[index].TaskItem.(exampletask.ExampleTask).TaskId)
			rsp, err := http.Post("http://127.0.0.1:8000/status", "application/json", strings.NewReader(string(body)))
			if err != nil {
				status[index].TaskStatus = framework.TASK_STATUS_RUNNING
				return
			}
			bs, _ := ioutil.ReadAll(rsp.Body)
			defer rsp.Body.Close()
			var r service.StatusResponse
			json.Unmarshal(bs, &r)
			if r.Status == "not found" {
				status[index].TaskStatus = framework.TASK_STATUS_FAILED
				status[index].FailedReason = fmt.Errorf("没找到任务")
			} else if r.Status == "running" {
				status[index].TaskStatus = framework.TASK_STATUS_RUNNING
			} else if r.Status == "success" {
				status[index].TaskStatus = framework.TASK_STATUS_SUCCESS
			} else {
				status[index].TaskStatus = framework.TASK_STATUS_INVALID
			}
			status[index].Progress = 0
		}(i)
	}
	wg.Wait()
	return status, nil
}

// GetOutput ...
func (e *ExampleActuator) GetOutput(ctx context.Context, task *framework.Task) (
	data interface{}, err error) {
	return nil, nil
}

// GetOutput ...
func (e *ExampleActuator) Delete(ctx context.Context, task *framework.Task) (err error) {
	return nil
}
