package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/prettyirrelevant/kilishi/api/callbacks"
	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/api/playlists"
	"github.com/prettyirrelevant/kilishi/config"
	"github.com/prettyirrelevant/kilishi/pkg/aggregator"
)

func main() {
	config := setupConfiguration()
	db := setupDatabase(config)

	fiberConfig := fiber.Config{}
	if !config.Debug {
		fiberConfig.Prefork = true
		fiberConfig.DisableKeepalive = true
	}

	app := fiber.New(fiberConfig)
	setupMiddlewares(app)

	apiGroup := app.Group("/api")
	aggregatorService := aggregator.New(config)
	playlists.RouterV1(apiGroup, aggregatorService, db)
	callbacks.RouterV1(apiGroup, aggregatorService, db)

	log.Fatal(app.Listen(":8000"))
}

func setupMiddlewares(app *fiber.App) {
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())
}

func setupDatabase(config *config.Config) *database.Database {
	db, err := database.New(config.DatabaseURI)
	if err != nil {
		panic(err)
	}

	return db
}

func setupConfiguration() *config.Config {
	config, err := config.New()
	if err != nil {
		panic(err)
	}
	return config
}
