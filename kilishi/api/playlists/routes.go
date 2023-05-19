package playlists

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prettyirrelevant/kilishi/aggregator"
	"github.com/prettyirrelevant/kilishi/api/database"
)

func RouterV1(router fiber.Router, aggregatorService *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) {
	router.Post("/v1/playlists/convert", ConvertPlaylistController(aggregatorService, db))
	router.Get("/v1/playlists/supported", GetSupportedPlatformsController(aggregatorService, db))
}
