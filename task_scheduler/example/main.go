package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	framework "github.com/memory-overflow/go-common-library/task_scheduler"
	memeorycontainer "github.com/memory-overflow/go-common-library/task_scheduler/container/memory_container"
	"github.com/memory-overflow/go-common-library/task_scheduler/example/actuator"
	service "github.com/memory-overflow/go-common-library/task_scheduler/example/add_service"
	"github.com/memory-overflow/go-common-library/task_scheduler/example/task"
)

func main() {
	go service.StartServer()
	container := memeorycontainer.MakeQueueContainer(10000, 100*time.Millisecond)
	actuator := actuator.MakeExampleActuator()
	sch := framework.MakeNewScheduler(
		context.Background(),
		container, actuator,
		framework.Config{
			TaskLimit:    5,
			ScanInterval: 0})

	var c chan os.Signal
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 1000; i++ {
		select {
		case <-c:
			return
		default:
			sch.AddTask(context.Background(),
				framework.Task{
					TaskId: strconv.Itoa(i),
					TaskItem: task.ExampleTask{
						TaskId: uint32(i),
						A:      r.Int31() % 1000,
						B:      r.Int31() % 1000,
					},
				})
		}
	}

	for range c {
		log.Println("stop Scheduling")
		sch.Close()
		return
	}
}
