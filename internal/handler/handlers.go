package handler

import (
	"net/http"
	"strconv"
	"strings"

	models "github.com/z2y9x5/metrics/internal/model"
	"github.com/z2y9x5/metrics/internal/service"
)

// Интерфейс эндпоинтов.
type Handlers interface {
	UpdateHandler(w http.ResponseWriter, r *http.Request)
}

// Эндпоинты.
type handlers struct {
	metrics service.Metrics
}

// Принимает данные в формате /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>, Content-Type: text/plain.
// При успешном приёме возвращает http.StatusOK.
// При попытке передать запрос без имени метрики возвращает http.StatusNotFound.
// При попытке передать запрос с некорректным типом метрики или значением возвращает http.StatusBadRequest.
func (h handlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if !(strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain")) {
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	metric := models.Metrics{}
	pathType := r.PathValue("type")
	pathName := r.PathValue("name")
	pathValue := r.PathValue("value")

	switch pathType {
	case models.Counter:
		metric.MType = models.Counter
		val, err := strconv.ParseInt(pathValue, 10, 64)
		if err != nil {
			http.Error(w, "Incorrect metric value for the counter", http.StatusBadRequest)
			return
		}
		metric.Delta = &val
	case models.Gauge:
		metric.MType = models.Gauge
		val, err := strconv.ParseFloat(pathValue, 64)
		if err != nil {
			http.Error(w, "Incorrect metric value for the gauge", http.StatusBadRequest)
			return
		}
		metric.Value = &val
	default:
		http.Error(w, "Incorrect metric type", http.StatusBadRequest)
		return
	}
	if pathName == "" {
		http.Error(w, "Incorrect metric name", http.StatusNotFound)
		return
	}
	metric.ID = pathName

	h.metrics.Save(metric)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// Конструктор эндпоинтов.
func NewHandlers(service service.Metrics) Handlers {
	return &handlers{
		metrics: service,
	}
}
