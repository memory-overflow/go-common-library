package misc

import (
	"fmt"
	"log"
	"runtime/debug"
)

// RecoverFunc 恢复函数
type RecoverFunc func() error

// GoroutineHelp ...
type GoroutineHelp struct {
}

// Recover panic 并打印日志, recoverFunc  用来恢复的执行函数, usgae: defer misc.Recover(nil)
func (GoroutineHelp) Recover(recoverFunc RecoverFunc) error {
	if p := recover(); p != nil {
		log.Printf("panic=%v stacktrace=%s\n", p, debug.Stack())
		if recoverFunc != nil {
			return recoverFunc()
		} else {
			return fmt.Errorf("panic=%v stacktrace=%s", p, debug.Stack())
		}
	}
	return nil
}

// GoroutinFunc ...
type GoroutinFunc func()

// SafeGo 安全启用携程，捕获 crash，并且重试
func (help GoroutineHelp) SafeGo(goFunc GoroutinFunc, recover bool) {
	go func() {
		var recoverFunc RecoverFunc = nil
		if recover {
			recoverFunc = func() error {
				help.SafeGo(goFunc, recover)
				return nil
			}
		}
		defer help.Recover(recoverFunc)
		goFunc()
	}()
}
