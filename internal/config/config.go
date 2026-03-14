package config

import (
	"flag"
	"os"
)

// Значения по умолчанию.
const (
	defaultServerAddr = "localhost:8080"
)

// Интерфейс конфигурации.
type Config interface {
	GetAppConfig() app
	ApplyCLIArgs()
	ApplyEnvArgs()
}

// Конфигурация приложения.
type app struct {
	ServerAddr string
}

// Конфигурация.
type config struct {
	App app
}

// Вернуть конфигурацию приложения.
func (c *config) GetAppConfig() app {
	return c.App
}

// Применить значения из аргументов командной строки.
func (c *config) ApplyCLIArgs() {
	flag.StringVar(&c.App.ServerAddr, "a", defaultServerAddr, "Адрес сервера в формате хост:порт. Пример: "+defaultServerAddr)
	flag.Parse()
}

// Применить значения из переменных окружения.
func (c *config) ApplyEnvArgs() {
	if val := os.Getenv("ADDRESS"); val != "" {
		c.App.ServerAddr = val
	}
}

// Конструктор конфигурации.
func NewConfig() Config {
	return &config{
		App: app{
			ServerAddr: defaultServerAddr,
		},
	}
}
