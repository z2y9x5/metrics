package agent

import (
	"sync"

	models "github.com/z2y9x5/metrics/internal/model"
)

// Интерфейс хранилища.
type Storage interface {
	GetValues() []models.Metrics
	SetValue(metric models.Metrics)
}

// Хранилище данных в памяти.
type memoryStorage struct {
	sync.Mutex
	metrics map[string]models.Metrics
}

// Вернуть последние данные всех метрик.
func (m *memoryStorage) GetValues() []models.Metrics {
	res := make([]models.Metrics, 0, len(m.metrics))
	m.Lock()
	for k := range m.metrics {
		res = append(res, m.metrics[k])
	}
	m.Unlock()
	return res
}

// Добавить значение метрики.
func (m *memoryStorage) SetValue(metric models.Metrics) {
	m.Lock()
	m.metrics[metric.ID] = metric
	m.Unlock()
}

// Конструктор хранилища данных в памяти.
func NewMemoryStorage() Storage {
	return &memoryStorage{
		metrics: make(map[string]models.Metrics),
	}
}
