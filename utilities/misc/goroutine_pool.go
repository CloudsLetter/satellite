package misc

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Task func()
type GoroutinePool struct {
	TaskQueue     chan Task
	MaxWorker     int32
	MinWorker     int32
	IdleTimeout   time.Duration
	CurrentWorker int32
	WG            sync.WaitGroup
}

var goroutinePool *GoroutinePool

// NewGoroutinePool initializes a new GoroutinePool instance.
func NewGoroutinePool(minWorker, maxWorker int32, idleTimeout time.Duration) *GoroutinePool {
	goroutinePool = &GoroutinePool{
		TaskQueue:     make(chan Task, maxWorker),
		MaxWorker:     maxWorker,
		MinWorker:     minWorker,
		IdleTimeout:   idleTimeout,
		CurrentWorker: 0,
	}
	for i := int32(0); i < minWorker; i++ {
		goroutinePool.addWorker()
	}
	return goroutinePool
}

// addWorker adds a new worker to the pool.
func (gp *GoroutinePool) addWorker() {
	atomic.AddInt32(&gp.CurrentWorker, 1)
	go func() {
		timer := time.NewTimer(gp.IdleTimeout)
		defer atomic.AddInt32(&gp.CurrentWorker, -1)
		defer timer.Stop()
		for {
			select {
			case task, ok := <-gp.TaskQueue:
				if !ok {
					return
				}
				if task != nil {
					task()
					gp.WG.Done()
				}
				timer.Reset(gp.IdleTimeout)
			case <-timer.C:
				if atomic.LoadInt32(&gp.CurrentWorker) > gp.MinWorker {
					return
				}
			}
		}
	}()
}

// Submit submits a task to the pool.
func (gp *GoroutinePool) Submit(task Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}
	gp.WG.Add(1)
	select {
	case gp.TaskQueue <- task:
	default:
		if atomic.LoadInt32(&gp.CurrentWorker) < gp.MaxWorker {
			gp.addWorker()
		}
		gp.TaskQueue <- task
	}
	return nil
}

// Shutdown waits for all tasks to complete and then closes the task queue.
func (gp *GoroutinePool) Shutdown() {
	close(gp.TaskQueue)
	gp.WG.Wait()
}

// GetGoroutinePool returns the singleton instance of GoroutinePool.
func GetGoroutinePool() *GoroutinePool {
	return goroutinePool
}
