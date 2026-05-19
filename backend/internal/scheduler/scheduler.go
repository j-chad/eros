package scheduler

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type internalTask struct {
	Task
	running atomic.Bool
	lastRun time.Time
}

type Scheduler struct {
	tasks []internalTask
	wg    sync.WaitGroup
}

func New() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, internalTask{
		Task: task,
	})
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// check immediately
	s.check(ctx, time.Now())

	for {
		select {
		case <-ticker.C:
			s.check(ctx, time.Now())
		case <-ctx.Done():
			s.waitForRunningTasks()
			return nil
		}
	}
}

func (s *Scheduler) check(ctx context.Context, now time.Time) {}

func (s *Scheduler) runTask(ctx context.Context, task *internalTask) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer task.running.Store(false)

		if err := task.Run(ctx); err != nil {
			// TODO: log the error.
		}
	}()
}

func (s *Scheduler) waitForRunningTasks() {
	s.wg.Wait()
}
