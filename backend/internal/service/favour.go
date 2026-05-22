package service

import (
	"backend/internal/logging"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

type FavourService struct {
	repo    repository.Repository
	pushSvc *PushService
}

var (
	ErrFavourTooExpensive = errors.New("favour request exceeds available favour count")
)

func NewFavourService(repo repository.Repository, pushSvc *PushService) *FavourService {
	return &FavourService{repo: repo, pushSvc: pushSvc}
}

func (s *FavourService) ListFavourChoices(ctx context.Context) ([]models.FavourChoice, error) {
	return s.repo.ListFavourChoices(ctx)
}

func (s *FavourService) GetFavourCount(ctx context.Context) (models.FavourCount, error) {
	return s.repo.GetFavourCount(ctx)
}

func (s *FavourService) ListFavourRequests(ctx context.Context) ([]models.FavourRequest, error) {
	return s.repo.ListFavourRequests(ctx)
}

func (s *FavourService) RequestFavour(ctx context.Context, request *models.FavourRequest) error {
	return s.repo.WithTx(ctx, nil, func(txRepo repository.Repository) error {
		cost, err := txRepo.GetFavourCostByID(ctx, request.ChoiceID)
		if err != nil {
			return err
		}

		favourCount, err := txRepo.GetFavourCount(ctx)
		if err != nil {
			return err
		}

		if favourCount.Remaining < cost {
			return ErrFavourTooExpensive
		}

		return txRepo.CreateFavourRequest(ctx, request)
	})
}

func (s *FavourService) DeleteFavourRequest(ctx context.Context, requestID string) error {
	return s.repo.DeleteFavourRequest(ctx, requestID)
}

func (s *FavourService) CreateFavourChoice(ctx context.Context, choice *models.FavourChoice) error {
	return s.repo.CreateFavourChoice(ctx, choice)
}

func (s *FavourService) UpdateFavourChoice(ctx context.Context, choice models.FavourChoice) error {
	return s.repo.UpdateFavourChoice(ctx, choice)
}

func (s *FavourService) DeleteFavourChoice(ctx context.Context, choiceID string) error {
	return s.repo.DeleteFavourChoice(ctx, choiceID)
}

func (s *FavourService) UpdateFavourCount(ctx context.Context, count int) error {
	logger := logging.FromContext(ctx)

	prev, err := s.repo.GetFavourCount(ctx)
	if err != nil {
		return fmt.Errorf("failed to get existing favour count: %w", err)
	}

	deltaTotal := count - prev.Total
	newRemaining := prev.Remaining + deltaTotal
	if newRemaining < 0 {
		return fmt.Errorf("negative favour count: %d", newRemaining)
	}

	logger.DebugContext(ctx, "updating favour count", "prevTotal", prev.Total, "newTotal", count, "deltaTotal", deltaTotal, "prevRemaining", prev.Remaining, "newRemaining", newRemaining)
	if deltaTotal > 0 && s.pushSvc != nil {
		logger.DebugContext(ctx, "favour count increased, sending push notification to users")
		_, _ = s.pushSvc.SendMessage(ctx, models.PushRequest{
			Message: models.PushMessage{
				Title: "Favour Granted",
				Body:  fmt.Sprintf("Your favour count has increased! You now have %d favours.", newRemaining),
				Tag:   "favour-count-updated",
				Data: models.PushData{
					URL: "/favours",
				},
			},
			Topic:   "favour-count-updated",
			TTL:     5 * time.Hour,
			Urgency: models.PushUrgencyNormal,
		})
	}

	return s.repo.UpdateFavourCount(ctx, count)
}

func (s *FavourService) FulfilFavourRequest(ctx context.Context, requestID string) error {
	return s.repo.FulfilFavourRequest(ctx, requestID)
}

func (s *FavourService) UnfulfilFavourRequest(ctx context.Context, requestID string) error {
	return s.repo.UnfulfilFavourRequest(ctx, requestID)
}
