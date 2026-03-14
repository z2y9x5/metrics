package handler

import (
	"net/http"
	"strconv"
	"text/template"

	models "github.com/z2y9x5/metrics/internal/model"
	"github.com/z2y9x5/metrics/internal/service"

	"github.com/go-chi/chi"
)

const rootPageTemplate = `
<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="UTF-8">
	<title>Метрики</title>
</head>
<body>
	<pre>
		{{range .}}
Name: {{.Name}}, Type: {{.Type}}, Value: {{.Value}}
		{{end}}
	</pre>
</body>
</html>
`

// Интерфейс эндпоинтов.
type Handlers interface {
	RootHandler(w http.ResponseWriter, r *http.Request)
	UpdateHandler(w http.ResponseWriter, r *http.Request)
	ValueHandler(w http.ResponseWriter, r *http.Request)
}

// Эндпоинты.
type handlers struct {
	metrics service.Metrics
}

// Возвращает HTML-страницу со списком имён и значений всех известных на текущий момент метрик.
func (h handlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(rootPageTemplate))
	metrics := h.metrics.GetAll()
	data := []struct {
		Name  string
		Type  string
		Value string
	}{}
	for _, v := range metrics {
		mValue := ""
		if v.MType == models.Gauge {
			mValue = strconv.FormatFloat(*v.Value, 'f', -1, 64)
		} else {
			mValue = strconv.FormatInt(*v.Delta, 10)
		}
		data = append(data, struct {
			Name  string
			Type  string
			Value string
		}{
			Name:  v.ID,
			Type:  v.MType,
			Value: mValue,
		})
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

// Принимает данные в формате /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>, Content-Type: text/plain.
// При успешном приёме возвращает http.StatusOK.
// При попытке передать запрос без имени метрики возвращает http.StatusNotFound.
// При попытке передать запрос с некорректным типом метрики или значением возвращает http.StatusBadRequest.
func (h handlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// if !(strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain")) {
	// 	http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
	// 	return
	// }

	metric := models.Metrics{}
	pathType := chi.URLParam(r, "type")
	pathName := chi.URLParam(r, "name")
	pathValue := chi.URLParam(r, "value")

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

// Принимает запрос в формате /value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>.
// Возвращает аккумулированное значение метрики в текстовом виде со статусом http.StatusOK.
// При попытке запроса неизвестной метрики возвращает http.StatusNotFound.
func (h handlers) ValueHandler(w http.ResponseWriter, r *http.Request) {
	pathType := chi.URLParam(r, "type")
	pathName := chi.URLParam(r, "name")

	val, ok := h.metrics.Get(pathName)
	if !ok {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}
	if pathType != val.MType {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	res := ""
	if val.MType == models.Gauge {
		res = strconv.FormatFloat(*val.Value, 'f', -1, 64)
	} else {
		res = strconv.FormatInt(*val.Delta, 10)
	}
	w.Write([]byte(res))
}

// Конструктор эндпоинтов.
func NewHandlers(service service.Metrics) Handlers {
	return &handlers{
		metrics: service,
	}
}
