package misc_test

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

func TestFlodLog(t *testing.T) {
	log := `{"RetrieveInputData":"其他","RetrieveInputType":0,"FilterSet":[],"PageNumber":1,"PageSize":6,"TIBusinessID":1,"TIProjectID":1}`
	fmt.Println(misc.FoldLog([]byte(log)))
}
