package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/prettyirrelevant/waakye/api/aggregator"
	"github.com/prettyirrelevant/waakye/api/database"
	"github.com/prettyirrelevant/waakye/api/routes"
	"github.com/prettyirrelevant/waakye/config"
)

func main() {
	config := setupConfiguration()
	db := setupDatabase(config)

	app := fiber.New()
	setupMiddlewares(app)

	aggregatorService := aggregator.New(db, config)

	apiGroup := app.Group("/api")
	routes.RouterV1(apiGroup, aggregatorService)

	log.Fatal(app.Listen(":8000"))
}

func setupMiddlewares(app *fiber.App) {
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Get("/metrics", monitor.New())
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
