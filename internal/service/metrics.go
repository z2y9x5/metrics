package service

import (
	models "github.com/z2y9x5/metrics/internal/model"
	"github.com/z2y9x5/metrics/internal/repository"
)

// Интерфейс сервиса метрик.
type Metrics interface {
	Save(metric models.Metrics)
}

// Сервис метрик.
type metrics struct {
	repo repository.Repository
}

// Сохранить метрику.
func (s *metrics) Save(metric models.Metrics) {
	if metric.MType == models.Counter {
		s.repo.Append(metric)
		return
	}
	s.repo.Replace(metric)
}

// Конструктор сервиса метрик.
func NewMetrics(r repository.Repository) Metrics {
	return &metrics{
		repo: r,
	}
}
