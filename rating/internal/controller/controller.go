package controller

import (
	"context"
	"errors"
	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

type ratingIngester interface {
	Ingest(ctx context.Context) (chan model.RatingEvent, error)
}

type Controller struct {
	repo     ratingRepository
	ingester ratingIngester
}

func New(repo ratingRepository, ingester ratingIngester) *Controller {
	return &Controller{repo: repo, ingester: ingester}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (ctrl *Controller) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := ctrl.repo.Get(ctx, recordID, recordType)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}

	return sum / float64(len(ratings)), nil
}

func (ctrl *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return ctrl.repo.Put(ctx, recordID, recordType, rating)
}

// StartIngestion starts the ingestion of rating events.
func (ctrl *Controller) StartIngestion(ctx context.Context) error {
	ch, err := ctrl.ingester.Ingest(ctx)
	if err != nil {
		return err
	}

	for e := range ch {
		if err := ctrl.PutRating(ctx, e.RecordID, model.RecordType(e.RecordType), &model.Rating{
			UserID: e.UserID,
			Value:  e.Value,
		}); err != nil {
			return err
		}
	}

	return nil
}
