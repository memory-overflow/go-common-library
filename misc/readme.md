[TOC]

# misc 模块说明

## time_related 模块
时间处理相关的函数
- TimeTick —— 计时器，通过 TimeTick.Tick() 打点计时间，返回距离上一次 Tick 过去的时间，单位 ms。
- Timed —— 函数调用计时。

example:
```go
import (
	"testing"
	"time"

	timerelated "github.com/memory-overflow/go-common-library/misc/time_related"
)

func f() error {
	time.Sleep(3 * time.Second)
	return nil
}

func TestTimeRelated(t *testing.T) {
	// test TimeTick
	// build tick
	tick := timerelated.BuildTimeTick()
	time.Sleep(time.Second)
	// 距离上一次打点的时间间隔
	lastTickPast := tick.Tick()
	t.Logf("after last tick time pass: %dms", lastTickPast)
	// 距离上一次打点的时间间隔
	time.Sleep(2 * time.Second)
	lastTickPast = tick.Tick()
	t.Logf("after last tick time pass: %dms", lastTickPast)

	// test Timed
	cost, _ := timerelated.Timed(f)
	t.Logf("call f time cost %dms", cost)
}

```

## download 模块
- DownloadFile 下载网络文件
- DownloadFileWithLimit 下载网络文件（带文件大小限制）

example: 
```go
import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/memory-overflow/go-common-library/misc"
)

func TestDownload(t *testing.T) {
	ctx := context.Background()
	uri := "http://www.baidu.com"
	dirs := "/root"
	filename := "baidu.html"
	defer os.Remove(path.Join(dirs, filename))
	// 目录不存在会自动创建
	if err := misc.DownloadFile(ctx, uri, dirs, filename); err != nil {
		t.Fatalf("download file: %v", err)
	}

	// 目录不存在会自动创建
	if err := misc.DownloadFileWithLimit(ctx, uri, dirs, filename, 1000); err != nil {
		t.Fatalf("download file: %v", err)
	}
}
```


## goroutine_help 模块
协程相关的处理模块。
- GoroutineHelp
    - GoroutineHelp.Recoverd：协程 recover 处理函数。
    - GoroutineHelp.SafeGo：安全的启用携程。
  
example:

```go
import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/memory-overflow/go-common-library/misc"
)

func goFunc() {
	fmt.Println("start goFunc")
	panic("test panic")
}

func TestGoroutineHelp(t *testing.T) {
	help := misc.GoroutineHelp{}
	// SafeGo, 如果异常,重启协程
	help.SafeGo(goFunc, true)

	// 后面的处理，如果出现异常，捕获异常防止线程退出，recoverFunc 用来恢复携程的处理。
	defer help.Recover(func() error {
		// recover code
		return nil
	})

	// other code
	for {

	}
}
```
## retry
Retry——带重试调用。


## id 生成器
IDGenerator 对象封装了很多生成 id 的函数。
- IDGenerator.GenerateRandomString：生成随机字符串。
- IDGenerator.GenerateUUID：生成UUID。
- IDGenerator.BighumpToUnderscore：大驼峰参数转换成下划线。
- IDGenerator.UnderscoreToBighump：下划线参数转换成大驼峰。

## log
- FoldLog: 折叠日志中数据过长的字段，输入仅支持 json 字符串。

example:
```go
import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/memory-overflow/go-common-library/misc"
)

func TestFlodLog(t *testing.T) {
	log := `{"RetrieveInputData":"其他","RetrieveInputType":0,"FilterSet":[],"PageNumber":1,"PageSize":6,"TIBusinessID":1,"TIProjectID":1}`
	fmt.Println(misc.FoldLog([]byte(log)))
}
```