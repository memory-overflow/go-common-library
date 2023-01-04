package timerelated

// Timed 统计执行时间，ms
func Timed(f func() error) (ms int64, err error) {
	tick := BuildTimeTick()
	err = f()
	ms = tick.Tick()
	return ms, err
}
