package scheduler

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type internalTask struct {
	Task
	running       atomic.Bool
	lastRunMinute time.Time
}

type Scheduler struct {
	tasks []internalTask
	wg    sync.WaitGroup
}

func New(tasks ...Task) *Scheduler {
	sched := &Scheduler{}

	for _, task := range tasks {
		sched.AddTask(task)
	}

	return sched
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

func (s *Scheduler) check(ctx context.Context, now time.Time) {
	nowMinute := now.Truncate(time.Minute)
	for i, _ := range s.tasks {
		task := &s.tasks[i]

		if !task.Cron.Matches(now) {
			continue
		}

		if task.lastRunMinute.Equal(nowMinute) {
			continue
		}

		if !task.running.CompareAndSwap(false, true) {
			// todo: log that the task is still running and will be skipped this time.
			continue
		}

		task.lastRunMinute = nowMinute
		s.runTask(ctx, task)
	}
}

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
