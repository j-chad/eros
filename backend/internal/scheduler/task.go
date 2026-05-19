package scheduler

import (
	"backend/internal/config"
	"context"
	"time"
)

const defaultTimeout = 10 * time.Second

type TaskFunc func(ctx context.Context) error

type Task struct {
	Name    string
	Cron    CronExpression
	Timeout time.Duration
	Fn      TaskFunc
}

func (t *Task) Run(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, t.Timeout)
	defer cancel()

	return t.Fn(ctx)
}

func (t *Task) inheritConfig(cfg config.SchedulerTaskConfig) error {
	if cfg.Cron != "" {
		configCron, err := ParseCronExpression(cfg.Cron)
		if err != nil {
			return err
		}

		t.Cron = configCron
	}

	if cfg.Timeout != 0 {
		t.Timeout = cfg.Timeout
	}

	if cfg.Timeout == 0 {
		t.Timeout = defaultTimeout
	}

	return nil
}
