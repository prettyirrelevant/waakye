package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis/v2"

	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"

	"github.com/prettyirrelevant/kilishi/api/auth"
	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/api/playlists"
	"github.com/prettyirrelevant/kilishi/config"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/aggregator"
)

func main() {
	cfg := setupConfiguration()
	db := setupDatabase(cfg)

	fiberConfig := fiber.Config{}
	if !cfg.Debug {
		fiberConfig.Prefork = true
		fiberConfig.DisableKeepalive = true
	}

	app := fiber.New(fiberConfig)
	setupMiddlewares(app, cfg)

	apiGroup := app.Group("/api")
	aggregatorService := aggregator.New(cfg)

	playlists.RouterV1(apiGroup, aggregatorService, db)
	auth.RouterV1(apiGroup, aggregatorService, db)

	apiGroup.Get("/v1/ping", HealthCheckController)

	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)))
}

func HealthCheckController(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func setupMiddlewares(app *fiber.App, cfg *config.Config) {
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cache.New(cache.Config{
		Expiration: 24 * time.Hour,
		Methods:    []string{fiber.MethodPost, fiber.MethodGet},
		Next: func(c *fiber.Ctx) bool {
			noCacheEndpoints := map[string]bool{"/api/v1/playlists": true, "/api/v1/auth/spotify/refresh": true}
			if _, ok := noCacheEndpoints[c.Path()]; ok && c.Method() == fiber.MethodPost {
				return true
			}

			return false
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.OriginalURL())
		},
		Storage: redis.New(redis.Config{
			URL:   cfg.DatabaseURL,
			Reset: false,
		}),
	}))
}

func setupDatabase(cfg *config.Config) *database.Database {
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	return db
}

func setupConfiguration() *config.Config {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	return cfg
}
