package app

import (
	"log"

	"github.com/bondarenkozahar80-hub/3.6/internal/cfg"
	"github.com/bondarenkozahar80-hub/3.6/internal/handlers"
	"github.com/bondarenkozahar80-hub/3.6/internal/service"
	"github.com/bondarenkozahar80-hub/3.6/internal/storage"
	"github.com/bondarenkozahar80-hub/3.6/internal/storage/postgres"
	"github.com/wb-go/wbf/ginext"
)

func Run() {
	config := cfg.Load()

	postgresStore, err := postgres.New(config.DatabaseURI)
	if err != nil {
		log.Fatalf("[app]failed to connect to PG DB: %v", err)
	}
	defer postgresStore.Close()
	log.Println("[app] Connected to Postgres successfully")

	store, err := storage.New(postgresStore)
	if err != nil {
		log.Fatalf("[app] failed to init unified storage: %v", err)
	}
	log.Println("[app]storage initialized successfully")

	service, err := service.New(store, store)
	if err != nil {
		log.Fatalf("[app] failed to init service: %v", err)
	}
	log.Println("[app] Service initialized successfully")

	engine := ginext.New("release")
	router := handlers.New(engine, service, service)
	router.Routes()

	log.Printf("[app] starting server on %s", config.ServerAddress)
	err = engine.Run(config.ServerAddress)
	if err != nil {
		log.Fatalf("[app] server failed: %v", err)
	}
}
