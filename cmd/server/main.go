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
	cnf := config.NewConfig()
	cnf.ApplyCLIArgs()
	cnfApp := cnf.GetAppConfig()

	db := repository.NewMemoryRepository()
	metrics := service.NewMetrics(db)
	handlers := handler.NewHandlers(metrics)

	mux := chi.NewRouter()
	mux.Use(handler.FixDoubleSlashes)
	mux.Get("/", handlers.RootHandler)
	mux.Post("/update/{type}/{name}/{value}", handlers.UpdateHandler)
	mux.Get("/value/{type}/{name}", handlers.ValueHandler)

	log.Fatal(http.ListenAndServe(cnfApp.ServerAddr, mux))
}
