package service

import (
	"context"
	"errors"
	"time"

	"github.com/bondarenkozahar80-hub/3.6/internal/model"
)

type transactionCRUD interface {
	Create(ctx context.Context, t *model.Transaction) (int, error)
	GetByID(ctx context.Context, id int) (*model.Transaction, error)
	GetAll(ctx context.Context) ([]model.Transaction, error)
	Update(ctx context.Context, t *model.Transaction) error
	Delete(ctx context.Context, id int) error
	GetByPeriod(ctx context.Context, from, to time.Time) ([]model.Transaction, error)
	GetAllSorted(ctx context.Context, sortField, order string) ([]model.Transaction, error)
}

type transactionAnalytics interface {
	GetSum(ctx context.Context, from, to time.Time) (int64, error)
	GetAvg(ctx context.Context, from, to time.Time) (float64, error)
	GetCount(ctx context.Context, from, to time.Time) (int64, error)
	GetMedian(ctx context.Context, from, to time.Time) (float64, error)
	GetPercentile90(ctx context.Context, from, to time.Time) (float64, error)
	GroupByDay(ctx context.Context, from, to time.Time) (map[string]int64, error)
	GroupByWeek(ctx context.Context, from, to time.Time) (map[string]int64, error)
	GroupByMonth(ctx context.Context, from, to time.Time) (map[string]int64, error)
	GroupByCategory(ctx context.Context, from, to time.Time) (map[string]int64, error)
}

type Service struct {
	storageCRUD      transactionCRUD
	storageAnalytics transactionAnalytics
}

func New(crud transactionCRUD, analytics transactionAnalytics) (*Service, error) {
	if crud == nil || analytics == nil {
		return nil, errors.New("[service] storage is nil")
	}
	return &Service{
		storageCRUD:      crud,
		storageAnalytics: analytics,
	}, nil
}
