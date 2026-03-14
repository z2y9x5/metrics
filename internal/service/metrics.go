package service

import (
	models "github.com/z2y9x5/metrics/internal/model"
	"github.com/z2y9x5/metrics/internal/repository"
)

// Интерфейс сервиса метрик.
type Metrics interface {
	Get(name string) (value models.Metrics, ok bool)
	GetAll() []models.Metrics
	Save(metric models.Metrics)
}

// Сервис метрик.
type metrics struct {
	repo repository.Repository
}

// Получить метрику.
func (s *metrics) Get(name string) (value models.Metrics, ok bool) {
	return s.repo.Get(name)
}

// Получить все метрики.
func (s *metrics) GetAll() []models.Metrics {
	return s.repo.GetAll()
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
