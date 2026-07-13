package statistics_service

import (
	"context"
	"time"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(ctx context.Context, userID *int, from, to *time.Time) ([]domain.Task, error)
}

func NewStatisticsService(repository StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: repository,
	}
}
