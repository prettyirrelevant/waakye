package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"

	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"

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
	setupMiddlewares(app, config)

	apiGroup := app.Group("/api")
	aggregatorService := aggregator.New(config)
	playlists.RouterV1(apiGroup, aggregatorService, db)
	callbacks.RouterV1(apiGroup, aggregatorService, db)
	apiGroup.Get("/ping", HealthCheckController)

	log.Fatal(app.Listen(":8000"))
}

func HealthCheckController(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(nil)
}

func setupMiddlewares(app *fiber.App, config *config.Config) {
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cache.New(cache.Config{
		Expiration: 168 * time.Hour,
		Methods:    []string{fiber.MethodPost},
		Next: func(c *fiber.Ctx) bool {
			ignoreCache := c.Query("ignoreCache", "false")
			if val, err := strconv.ParseBool(ignoreCache); err != nil || !val {
				return true
			}
			return false
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path()) + string(utils.CopyBytes(c.Body()))
		},
		Storage: redis.New(redis.Config{
			URL:   config.RedisURI,
			Reset: false,
		}),
	}))
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
