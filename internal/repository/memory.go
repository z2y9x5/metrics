package repository

import (
	models "github.com/z2y9x5/metrics/internal/model"
)

// Интерфейс хранилища данных.
type Repository interface {
	Append(metric models.Metrics)
	Replace(metric models.Metrics)
}

// Хранилище данных в памяти.
type memStorage struct {
	metrics map[string]models.Metrics
}

// Добавить данные метрики.
func (m *memStorage) Append(metric models.Metrics) {
	if v, ok := m.metrics[metric.ID]; ok {
		*v.Delta += *metric.Delta
		m.metrics[metric.ID] = v
		return
	}
	m.metrics[metric.ID] = metric
}

// Заменить данные метрики.
func (m *memStorage) Replace(metric models.Metrics) {
	m.metrics[metric.ID] = metric
}

// Конструктор хранилища данных в памяти.
func NewMemoryRepository() Repository {
	return &memStorage{
		metrics: make(map[string]models.Metrics),
	}
}
