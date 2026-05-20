package scheduler

import (
	"backend/internal/config"
	"backend/internal/logging"
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
	config  config.SchedulerConfig
	tasks   []internalTask
	running atomic.Bool
	wg      sync.WaitGroup
}

func New(conf config.SchedulerConfig, tasks ...Task) (*Scheduler, error) {
	sched := &Scheduler{config: conf}

	for _, task := range tasks {
		err := sched.AddTask(task)
		if err != nil {
			return nil, err
		}
	}

	return sched, nil
}

func (s *Scheduler) MustAddTask(task Task) {
	err := s.AddTask(task)
	if err != nil {
		panic(fmt.Sprintf("failed to add task %s: %v", task.Name, err))
	}
}

func (s *Scheduler) AddTask(task Task) error {
	if s.running.Load() {
		return fmt.Errorf("cannot add task %s: scheduler is already running", task.Name)
	}

	name := task.Name
	taskConfig, hasConfig := s.config.Tasks[name]

	if hasConfig {
		if taskConfig.Disabled {
			return nil
		}

		err := task.inheritConfig(taskConfig)
		if err != nil {
			return err
		}
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
	logger := logging.FromContext(ctx)

	if s.config.Disabled {
		logger.Warn("scheduler disabled")
		return
	}

	if !s.running.CompareAndSwap(false, true) {
		logger.WarnContext(ctx, "scheduler is already running")
	}

	logger.InfoContext(ctx, "starting scheduler", "task_count", len(s.tasks))

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// check immediately
	s.check(ctx, time.Now())

	for {
		select {
		case <-ticker.C:
			s.check(ctx, time.Now())
		case <-ctx.Done():
			logger.DebugContext(ctx, "stopping scheduler", "task_count", len(s.tasks))
			s.waitForRunningTasks()
			return
		}
	}
}

func (s *Scheduler) check(ctx context.Context, now time.Time) {
	logger := logging.FromContext(ctx)

	nowMinute := now.Truncate(time.Minute)
	for i, _ := range s.tasks {
		task := &s.tasks[i]

		if !task.Cron.Matches(now) {
			continue
		}

		if task.lastRunMinute.Equal(nowMinute) {
			logger.DebugContext(ctx, "skipping task because it has already run this minute", "task", task.Name)
			continue
		}

		if !task.running.CompareAndSwap(false, true) {
			logger.WarnContext(ctx, "skipping task because previous run is still in progress", "task", task.Name)
			continue
		}

		task.lastRunMinute = nowMinute
		s.runTask(ctx, task)
	}
}

func (s *Scheduler) runTask(ctx context.Context, task *internalTask) {
	logger := logging.FromContext(ctx).With("task", task.Name)
	logCtx := logging.NewContext(ctx, logger)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer task.running.Store(false)

		logger.DebugContext(ctx, "running scheduled task")
		if err := task.Run(logCtx); err != nil {
			logger.WarnContext(ctx, "error running scheduled task", "error", err)
			return
		}

		logger.DebugContext(ctx, "completed running scheduled task")
	}()
}

func (s *Scheduler) waitForRunningTasks() {
	s.wg.Wait()
}
