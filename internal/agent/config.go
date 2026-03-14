package agent

import "time"

const (
	defaultServerAddr     = "localhost:8080" // Адрес сервера.
	defaultPollInterval   = 2 * time.Second  // Частота получения метрик.
	defaultReportInterval = 10 * time.Second // Частота отправки данных.
)

// Конфигурация агента.
type config struct {
	ServerAddr     string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

// Конструктор конфигурации.
func NewConfig() config {
	return config{
		ServerAddr:     defaultServerAddr,
		PollInterval:   defaultPollInterval,
		ReportInterval: defaultReportInterval,
	}
}
