package misc

import (
	"time"
)

// RecoverFunc 重试函数
type RetryFunc func() error

// Retry 带重试的调用
func Retry(f RetryFunc) (err error) {
	sleepTime := 100 * time.Millisecond
	for i := 0; i < 3; i++ {
		if err = f(); err == nil {
			return nil
		}
		time.Sleep(sleepTime)
		sleepTime *= 2
	}
	return err
}
