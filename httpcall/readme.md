# http 最常用的 json 协议调用封装
4个方法, JsonPost、JsonGet、JsonPut、JsonDelete。

example: [httpcall_test.go](https://github.com/memory-overflow/go-common-library/blob/main/httpcall/httpcall_test.go)
```go
import (
	"context"
	"testing"
	"time"

	"github.com/memory-overflow/go-common-library/httpcall"
	service "github.com/memory-overflow/go-common-library/task_scheduler/example/add_service"
)

func TestJsonPost(t *testing.T) {
	// 控制请求超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	uri := "http://127.0.0.1:8000/add"
	req := &service.AddRequest{
		A: 1,
		B: 2,
	}
	rsp := &service.AddResponse{}
	err := httpcall.JsonPost(ctx, uri, nil, req, rsp)
	if err != nil {
		t.Fatalf("JsonPost err: %v", err)
	}
	t.Logf("add taskId: %d", rsp.TaskId)
}
```