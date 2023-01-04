package timerelated

import "time"

// TimeTick 计时器，单位 ms
type TimeTick struct {
	TimeList  []int64 // 打点时间点
	BuildTime int64   // 开始时间
}

// BuildTimeTick 构建一个计时器
func BuildTimeTick() TimeTick {
	return TimeTick{BuildTime: time.Now().UnixNano() / 1000000}
}

// Tick 计时器打点
func (tick *TimeTick) Tick() int64 {
	now := time.Now().UnixNano() / 1000000
	defer func() { tick.TimeList = append(tick.TimeList, now) }()
	if len(tick.TimeList) == 0 {
		return now - tick.BuildTime
	} else {
		return now - tick.TimeList[len(tick.TimeList)-1]
	}
}
