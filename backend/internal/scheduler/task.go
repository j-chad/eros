package scheduler

import (
	"context"
	"time"
)

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
