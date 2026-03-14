package repository

import (
	"sync"

	models "github.com/z2y9x5/metrics/internal/model"
)

// Интерфейс хранилища данных.
type Repository interface {
	Get(name string) (value models.Metrics, ok bool)
	GetAll() []models.Metrics
	Append(metric models.Metrics)
	Replace(metric models.Metrics)
}

// Хранилище данных в памяти.
type memStorage struct {
	sync.Mutex
	metrics map[string]models.Metrics
}

// Получить метрику.
func (m *memStorage) Get(name string) (value models.Metrics, ok bool) {
	m.Lock()
	defer m.Unlock()
	value, ok = m.metrics[name]
	return value, ok
}

// Получить все метрики.
func (m *memStorage) GetAll() []models.Metrics {
	m.Lock()
	defer m.Unlock()
	res := []models.Metrics{}
	for _, v := range m.metrics {
		res = append(res, v)
	}
	return res
}

// Добавить данные метрики.
func (m *memStorage) Append(metric models.Metrics) {
	m.Lock()
	defer m.Unlock()
	if v, ok := m.metrics[metric.ID]; ok {
		*v.Delta += *metric.Delta
		m.metrics[metric.ID] = v
		return
	}
	m.metrics[metric.ID] = metric
}

// Заменить данные метрики.
func (m *memStorage) Replace(metric models.Metrics) {
	m.Lock()
	defer m.Unlock()
	m.metrics[metric.ID] = metric
}

// Конструктор хранилища данных в памяти.
func NewMemoryRepository() Repository {
	return &memStorage{
		metrics: make(map[string]models.Metrics),
	}
}
