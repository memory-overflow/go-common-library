package timerelated_test

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
