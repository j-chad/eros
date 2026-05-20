package main

import (
	"backend/internal/logging"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/scheduler"
	"backend/internal/service"
	"context"
	"fmt"
	"time"
)

// MustRegisterTasks is used to register scheduled tasks that cut across multiple services
func MustRegisterTasks(sched *scheduler.Scheduler, repo repository.Repository, push *service.PushService) {
	sched.MustAddTask(notifyUnlockedGraphs(repo, push))
}

func notifyUnlockedGraphs(repo repository.Repository, push *service.PushService) scheduler.Task {
	return scheduler.Task{
		Name:    "notify_unlocked_graphs",
		Cron:    scheduler.MustParseCronExpression("*/5 * * * *"),
		Timeout: 30 * time.Second,
		Fn: func(ctx context.Context) error {
			logger := logging.FromContext(ctx)

			if !push.IsEnabled() {
				logger.DebugContext(ctx, "push notifications not enabled, skipping graph unlock notifications")
				return nil
			}

			graphs, err := repo.GetGraphsPendingNotification(ctx)
			if err != nil {
				return fmt.Errorf("get graphs pending notification: %w", err)
			}
			if len(graphs) == 0 {
				logger.DebugContext(ctx, "no graphs pending notification")
				return nil
			}

			logger.DebugContext(ctx, fmt.Sprintf("%d graphs pending notification", len(graphs)))
			for _, graph := range graphs {
				result, err := push.SendMessage(ctx, models.PushRequest{
					Topic:   "graph_updates",
					TTL:     24 * time.Hour,
					Urgency: models.PushUrgencyNormal,
					Message: models.PushMessage{
						Title: "Graph Unlocked",
						Body:  "The graph \"" + graph.Title + "\" has been unlocked and is now available to play!",
						Tag:   "graph_unlocked",
					},
				})
				if err != nil {
					logger.WarnContext(ctx, "failed to send push notification for graph", "graph_id", graph.ID, "error", err)
					continue
				}

				err = repo.MarkGraphNotified(ctx, graph.ID)
				if err != nil {
					logger.ErrorContext(ctx, "failed to mark graph as notified", "graph_id", graph.ID, "error", err)
					continue
				}

				logger.DebugContext(ctx, "successfully sent push notification for graph", "graph_id", graph.ID, "sent", result.Sent, "cleaned", result.Cleaned, "failed", result.Failed)
			}

			return nil
		},
	}
}
