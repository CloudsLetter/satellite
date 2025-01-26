package misc

import (
	"sync"
	"sync/atomic"
	"time"
)

var goroutinePool *GoroutinePool

type Task func()
type GoroutinePool struct {
	TaskQueue     chan Task
	MaxWorker     int32
	MinWorker     int32
	IdleTimeout   time.Duration
	CurrentWorker int32
	WG            sync.WaitGroup
}

func NewGoroutinePool(minWorker, maxWorker int32, idleTimeout time.Duration) {
	goroutinePool = &GoroutinePool{
		TaskQueue:     make(chan Task),
		MaxWorker:     maxWorker,
		MinWorker:     minWorker,
		IdleTimeout:   idleTimeout,
		CurrentWorker: 0,
	}
	for i := int32(0); i < minWorker; i++ {
		goroutinePool.addWorker()
	}
}

func (gp *GoroutinePool) addWorker() {
	atomic.AddInt32(&gp.CurrentWorker, 1)
	go func() {
		for {
			select {
			case task, ok := <-gp.TaskQueue:
				if !ok {
					atomic.AddInt32(&gp.CurrentWorker, -1)
					return
				}
				if task != nil {
					task()
					gp.WG.Done()
				}
			case <-time.After(gp.IdleTimeout):
				if atomic.LoadInt32(&gp.CurrentWorker) > gp.MinWorker {
					atomic.AddInt32(&gp.CurrentWorker, -1)
					return
				}
			}
		}
	}()
}

func (gp *GoroutinePool) Submit(task Task) {
	gp.WG.Add(1)
	select {
	case gp.TaskQueue <- task:
	default:
		if atomic.LoadInt32(&gp.CurrentWorker) < gp.MaxWorker {
			gp.addWorker()
			gp.TaskQueue <- task
		}
	}
}

func (gp *GoroutinePool) Shutdown() {
	close(gp.TaskQueue)
	gp.WG.Wait()
}

func GetGoroutinePool() *GoroutinePool {
	return goroutinePool
}
