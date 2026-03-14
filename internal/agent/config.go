package agent

import (
	"flag"
	"time"
)

// Значения по умолчанию.
const (
	defaultServerAddr     = "localhost:8080" // Адрес сервера.
	defaultPollInterval   = 2                // Частота получения метрик в секундах.
	defaultReportInterval = 10               // Частота отправки данных в секундах.
)

// Интерфейс конфигурации.
type Config interface {
	GetAppConfig() configApp
	ApplyCLIArgs()
}

// Конфигурация приложения.
type configApp struct {
	ServerAddr     string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

// Конфигурация.
type config struct {
	App configApp
}

// Вернуть конфигурацию приложения.
func (c *config) GetAppConfig() configApp {
	return c.App
}

// Применить значения из аргументов командной строки.
func (c *config) ApplyCLIArgs() {
	a := flag.String("a", defaultServerAddr, "Адрес сервера в формате хост:порт. Пример: "+defaultServerAddr)
	p := flag.Int("p", defaultPollInterval, "Частота получения метрик в секундах.")
	r := flag.Int("r", defaultReportInterval, "Частота отправки метрик в секундах.")
	flag.Parse()
	c.App.ServerAddr = *a
	c.App.PollInterval = time.Duration(*p) * time.Second
	c.App.ReportInterval = time.Duration(*r) * time.Second
}

// Конструктор конфигурации.
func NewConfig() Config {
	return &config{
		App: configApp{
			ServerAddr:     defaultServerAddr,
			PollInterval:   time.Duration(defaultPollInterval) * time.Second,
			ReportInterval: time.Duration(defaultReportInterval) * time.Second,
		},
	}
}
