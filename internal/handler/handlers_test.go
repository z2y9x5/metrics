package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/z2y9x5/metrics/internal/config"
	models "github.com/z2y9x5/metrics/internal/model"
	"github.com/z2y9x5/metrics/internal/service"

	"github.com/stretchr/testify/assert"
)

// Заглушка хранилища данных.
type mockRepository struct{}

func (m mockRepository) Get(name string) (value models.Metrics, ok bool) {
	return models.Metrics{}, true
}
func (m mockRepository) GetAll() []models.Metrics {
	return []models.Metrics{}
}
func (m mockRepository) Append(metric models.Metrics)  {}
func (m mockRepository) Replace(metric models.Metrics) {}

func TestUpdateHandler(t *testing.T) {
	type have struct {
		pathType  string
		pathName  string
		pathValue string
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		have have
		want want
	}{
		{
			name: "UpdateHandler test 1",
			have: have{
				pathType:  "gauge",
				pathName:  "test",
				pathValue: "1",
			},
			want: want{
				statusCode: 200,
			},
		},
		{
			name: "UpdateHandler test 2",
			have: have{
				pathType:  "counter",
				pathName:  "test",
				pathValue: "1",
			},
			want: want{
				statusCode: 200,
			},
		},
		{
			name: "UpdateHandler test 3",
			have: have{
				pathType:  "wtf",
				pathName:  "test",
				pathValue: "1",
			},
			want: want{
				statusCode: 400,
			},
		},
		{
			name: "UpdateHandler test 4",
			have: have{
				pathType:  "",
				pathName:  "test",
				pathValue: "1",
			},
			want: want{
				statusCode: 404,
			},
		},
		{
			name: "UpdateHandler test 5",
			have: have{
				pathType:  "gauge",
				pathName:  "test",
				pathValue: "z",
			},
			want: want{
				statusCode: 400,
			},
		},
		{
			name: "UpdateHandler test 6",
			have: have{
				pathType:  "gauge",
				pathName:  "test",
				pathValue: "",
			},
			want: want{
				statusCode: 404,
			},
		},
		{
			name: "UpdateHandler test 7",
			have: have{
				pathType:  "counter",
				pathName:  "test",
				pathValue: "z",
			},
			want: want{
				statusCode: 400,
			},
		},
		{
			name: "UpdateHandler test 8",
			have: have{
				pathType:  "counter",
				pathName:  "test",
				pathValue: "",
			},
			want: want{
				statusCode: 404,
			},
		},
		{
			name: "UpdateHandler test 9",
			have: have{
				pathType:  "gauge",
				pathName:  "",
				pathValue: "1",
			},
			want: want{
				statusCode: 404,
			},
		},
	}

	cnf := config.NewConfig().GetAppConfig()
	db := mockRepository{}
	metrics := service.NewMetrics(db)
	handlers := NewHandlers(metrics)

	mux := chi.NewRouter()
	mux.Use(FixDoubleSlashes)
	mux.Post("/update/{type}/{name}/{value}", handlers.UpdateHandler)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/update/" + test.have.pathType + "/" + test.have.pathName + "/" + test.have.pathValue
			r := httptest.NewRequest(http.MethodPost, url, nil)
			r.Host = cnf.ServerAddr
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
		})
	}
}
