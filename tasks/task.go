package tasks

import (
	"log"
	"runtime"
	"sync"
	"time"
)

type Task struct {
	mu            sync.Mutex
	nextExecution time.Time
	interval      time.Duration
	action        func()
}

func NewTask(interval time.Duration, action func()) *Task {
	return &Task{
		nextExecution: time.Time{},
		interval:      interval,
		action:        action,
	}
}

func Run(tasks []*Task) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	for i := range tasks {
		tasks[i].nextExecution = startOfDay
	}

	for {
		now = time.Now()

		for i := range tasks {
			go func(t *Task) {
				// Skip if the previous execution is still running
				if !t.mu.TryLock() {
					return
				}
				defer t.mu.Unlock()

				defer func() {
					if r := recover(); r != nil {
						buf := make([]byte, 1024)
						stackSize := runtime.Stack(buf, false)
						stackTrace := string(buf[:stackSize])
						log.Printf("Task panic recovered: %v\nStack trace:\n%s", r, stackTrace)
					}
				}()

				if now.Equal(t.nextExecution) || now.After(t.nextExecution) {
					t.action()

					postExecuteTime := time.Now()

					for postExecuteTime.After(t.nextExecution) {
						t.nextExecution = t.nextExecution.Add(t.interval)
					}
				}
			}(tasks[i])
		}

		time.Sleep(time.Second)
	}
}
