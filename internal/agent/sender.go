package agent

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	models "github.com/z2y9x5/metrics/internal/model"
)

// Интерфейс отправителя данных.
type Sender interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

// Отправитель данных.
type sender struct {
	metrics        Storage
	serverAddr     string
	reportInterval time.Duration
}

// Запустить отправитель данных.
func (s *sender) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(s.reportInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.send()
		case <-ctx.Done():
			return
		}
	}
}

// Отправить данные.
func (s *sender) send() {
	data := s.metrics.GetValues()
	for _, v := range data {
		pathType := v.MType
		pathName := v.ID
		pathValue := ""
		if v.MType == models.Gauge && v.Value != nil {
			pathValue = strconv.FormatFloat(*v.Value, 'f', -1, 64)
		}
		if v.MType == models.Counter && v.Delta != nil {
			pathValue = strconv.FormatInt(*v.Delta, 10)
		}
		url := "http://" + s.serverAddr + "/update/" + pathType + "/" + pathName + "/" + pathValue
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			log.Println("request error", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Println("response code error", resp.StatusCode, "на запрос", url)
			continue
		}
		log.Println("successful sending", pathName, pathValue)
	}
}

// Конструктор отправителя данных.
func NewSender(cnf config, db Storage) Sender {
	return &sender{
		metrics:        db,
		serverAddr:     cnf.ServerAddr,
		reportInterval: cnf.ReportInterval,
	}
}
