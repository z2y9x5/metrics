package main

import (
	"log"
	"net/http"

	"github.com/z2y9x5/metrics/internal/config"
	"github.com/z2y9x5/metrics/internal/handler"
	"github.com/z2y9x5/metrics/internal/repository"
	"github.com/z2y9x5/metrics/internal/service"

	"github.com/go-chi/chi"
)

func main() {
	cnfApp := config.NewConfig().GetAppConfig()
	db := repository.NewMemoryRepository()
	metrics := service.NewMetrics(db)
	handlers := handler.NewHandlers(metrics)

	mux := chi.NewRouter()
	mux.Use(handler.FixDoubleSlashes)
	mux.Post("/update/{type}/{name}/{value}", handlers.UpdateHandler)

	log.Fatal(http.ListenAndServe(cnfApp.ServerAddr, mux))
}
