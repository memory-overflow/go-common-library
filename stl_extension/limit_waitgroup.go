package stlextension

import "sync"

// LimitWaitGroup ...
type LimitWaitGroup struct {
	wg     sync.WaitGroup
	bucket chan struct{}
}

// NewLimitWaitGroup ...
func NewLimitWaitGroup(limit uint) *LimitWaitGroup {
	return &LimitWaitGroup{
		wg:     sync.WaitGroup{},
		bucket: make(chan struct{}, limit),
	}
}

// Add ...
func (w *LimitWaitGroup) Add(delta int) {
	w.wg.Add(delta)
	for i := 0; i < delta; i++ {
		w.bucket <- struct{}{}
	}
}

// Done ...
func (w *LimitWaitGroup) Done() {
	<-w.bucket
	w.wg.Done()
}

// Wait ...
func (w *LimitWaitGroup) Wait() {
	w.wg.Wait()
}
