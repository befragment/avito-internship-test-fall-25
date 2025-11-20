package main

import (
	"log"

	"avito-intern-test/internal/core"
	"avito-intern-test/internal/routing"
)

func main() {
	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbPool := core.MustInitPool()

	core.StartServer(
		dbPool,
		cfg.Port,
		routing.AppRoutes(
			routing.CommonRoutes(),
		),
	)
}