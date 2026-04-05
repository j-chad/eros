package service

import (
	"backend/internal/models"
	"backend/internal/testutil"
	"backend/internal/testutil/testdb"
	"context"
	"testing"
	"time"
)

func setupFavourTest(t *testing.T) (*FavourService, context.Context) {
	t.Helper()
	repo := testdb.New(t)
	ctx := context.Background()

	// Seed: 10 total favours, one choice costing 3
	repo.UpdateFavourCount(ctx, 10)
	repo.CreateFavourChoice(ctx, &models.FavourChoice{
		Label:     "Dinner",
		Cost:      3,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return NewFavourService(repo), ctx
}

func TestRequestFavour_Affordable(t *testing.T) {
	svc, ctx := setupFavourTest(t)

	choices, _ := svc.ListFavourChoices(ctx)
	testutil.Equal(t, len(choices), 1)

	req := &models.FavourRequest{
		ChoiceID:    choices[0].ID,
		RequestedAt: time.Now(),
	}
	testutil.NilErr(t, svc.RequestFavour(ctx, req))
	testutil.True(t, req.ID != "", "request should get an ID")

	// Verify remaining count decreased
	count, err := svc.GetFavourCount(ctx)
	testutil.NilErr(t, err)
	testutil.Equal(t, count.Remaining, 7)
}

func TestRequestFavour_TooExpensive(t *testing.T) {
	repo := testdb.New(t)
	ctx := context.Background()
	svc := NewFavourService(repo)

	// 2 total favours, choice costs 5
	repo.UpdateFavourCount(ctx, 2)
	repo.CreateFavourChoice(ctx, &models.FavourChoice{
		Label:     "Expensive",
		Cost:      5,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	choices, _ := svc.ListFavourChoices(ctx)
	req := &models.FavourRequest{
		ChoiceID:    choices[0].ID,
		RequestedAt: time.Now(),
	}
	testutil.ErrorIs(t, svc.RequestFavour(ctx, req), ErrFavourTooExpensive)
}

func TestRequestFavour_ExactlyEnough(t *testing.T) {
	repo := testdb.New(t)
	ctx := context.Background()
	svc := NewFavourService(repo)

	repo.UpdateFavourCount(ctx, 3)
	repo.CreateFavourChoice(ctx, &models.FavourChoice{
		Label:     "Exact",
		Cost:      3,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	choices, _ := svc.ListFavourChoices(ctx)
	req := &models.FavourRequest{
		ChoiceID:    choices[0].ID,
		RequestedAt: time.Now(),
	}
	testutil.NilErr(t, svc.RequestFavour(ctx, req))
}

func TestRequestFavour_MultipleRequestsDrainCount(t *testing.T) {
	svc, ctx := setupFavourTest(t)
	choices, _ := svc.ListFavourChoices(ctx)
	choiceID := choices[0].ID // costs 3, total is 10

	// Should be able to make 3 requests (3*3=9 <= 10)
	for i := 0; i < 3; i++ {
		req := &models.FavourRequest{ChoiceID: choiceID, RequestedAt: time.Now()}
		testutil.NilErr(t, svc.RequestFavour(ctx, req))
	}

	// 4th should fail (3*4=12 > 10)
	req := &models.FavourRequest{ChoiceID: choiceID, RequestedAt: time.Now()}
	testutil.ErrorIs(t, svc.RequestFavour(ctx, req), ErrFavourTooExpensive)
}

func TestRequestFavour_InvalidChoice(t *testing.T) {
	svc, ctx := setupFavourTest(t)
	req := &models.FavourRequest{
		ChoiceID:    "nonexistent-id",
		RequestedAt: time.Now(),
	}
	testutil.NotNilErr(t, svc.RequestFavour(ctx, req))
}
