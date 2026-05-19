package scheduler

import (
	"backend/internal/config"
	"context"
	"fmt"
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
	config *config.SchedulerConfig
	tasks  []internalTask
	wg     sync.WaitGroup
}

func New(conf *config.SchedulerConfig, tasks ...Task) (*Scheduler, error) {
	sched := &Scheduler{config: conf}

	for _, task := range tasks {
		err := sched.AddTask(task)
		if err != nil {
			return nil, err
		}
	}

	return sched, nil
}

func (s *Scheduler) AddTask(task Task) error {
	name := task.Name
	taskConfig, _ := s.config.Tasks[name]

	if taskConfig.Disabled {
		return nil
	}

	err := task.inheritConfig(taskConfig)
	if err != nil {
		return err
	}

	if task.Cron.IsZero() {
		return fmt.Errorf("task %s has no cron expression", name)
	}

	s.tasks = append(s.tasks, internalTask{
		Task: task,
	})

	return nil
}

func (s *Scheduler) Run(ctx context.Context) {
	if s.config.Disabled {
		return
	}

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
			return
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
