package config

// Значения по умолчанию.
const (
	defaultServerAddr = "localhost:8080"
)

// Интерфейс конфигурации.
type Config interface {
	GetAppConfig() app
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

// Конструктор конфигурации.
func NewConfig() Config {
	return &config{
		App: app{
			ServerAddr: defaultServerAddr,
		},
	}
}
