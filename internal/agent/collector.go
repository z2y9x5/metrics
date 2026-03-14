package agent

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"time"

	models "github.com/z2y9x5/metrics/internal/model"
)

// Интерфейс сборщика метрик.
type Collector interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

// Сборщик метрик.
type collector struct {
	metrics      Storage
	memStats     runtime.MemStats
	pollCount    int64
	pollInterval time.Duration
}

// Запустить сбор метрик.
func (c *collector) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.collect()
		case <-ctx.Done():
			return
		}
	}
}

// Сохранить значение типа counter.
func (c *collector) saveCounter(id string, value int64) {
	c.metrics.SetValue(models.Metrics{
		MType: models.Counter,
		ID:    id,
		Delta: &value,
	})
}

// Сохранить значение типа gauge.
func (c *collector) saveGauge(id string, value float64) {
	c.metrics.SetValue(models.Metrics{
		MType: models.Gauge,
		ID:    id,
		Value: &value,
	})
}

// Собрать метрики.
func (c *collector) collect() {
	// Метрики runtime.
	runtime.ReadMemStats(&c.memStats)
	c.saveGauge("Alloc", float64(c.memStats.Alloc))
	c.saveGauge("BuckHashSys", float64(c.memStats.BuckHashSys))
	c.saveGauge("Frees", float64(c.memStats.Frees))
	c.saveGauge("GCCPUFraction", float64(c.memStats.GCCPUFraction))
	c.saveGauge("GCSys", float64(c.memStats.GCSys))
	c.saveGauge("HeapAlloc", float64(c.memStats.HeapAlloc))
	c.saveGauge("HeapIdle", float64(c.memStats.HeapIdle))
	c.saveGauge("HeapInuse", float64(c.memStats.HeapInuse))
	c.saveGauge("HeapObjects", float64(c.memStats.HeapObjects))
	c.saveGauge("HeapReleased", float64(c.memStats.HeapReleased))
	c.saveGauge("HeapSys", float64(c.memStats.HeapSys))
	c.saveGauge("LastGC", float64(c.memStats.LastGC))
	c.saveGauge("Lookups", float64(c.memStats.Lookups))
	c.saveGauge("MCacheInuse", float64(c.memStats.MCacheInuse))
	c.saveGauge("MCacheSys", float64(c.memStats.MCacheSys))
	c.saveGauge("MSpanInuse", float64(c.memStats.MSpanInuse))
	c.saveGauge("MSpanSys", float64(c.memStats.MSpanSys))
	c.saveGauge("Mallocs", float64(c.memStats.Mallocs))
	c.saveGauge("NextGC", float64(c.memStats.NextGC))
	c.saveGauge("NumForcedGC", float64(c.memStats.NumForcedGC))
	c.saveGauge("NumGC", float64(c.memStats.NumGC))
	c.saveGauge("OtherSys", float64(c.memStats.OtherSys))
	c.saveGauge("PauseTotalNs", float64(c.memStats.PauseTotalNs))
	c.saveGauge("StackInuse", float64(c.memStats.StackInuse))
	c.saveGauge("StackSys", float64(c.memStats.StackSys))
	c.saveGauge("Sys", float64(c.memStats.Sys))
	c.saveGauge("TotalAlloc", float64(c.memStats.TotalAlloc))
	// Дополнительные метрики.
	c.saveCounter("PollCount", c.pollCount+1)
	c.saveGauge("RandomValue", rand.Float64())
}

// Конструктор сборщика метрик.
func NewCollector(cnf config, db Storage) Collector {
	return &collector{
		metrics:      db,
		pollInterval: cnf.PollInterval,
	}
}
