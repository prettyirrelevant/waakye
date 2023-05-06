package playlists

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/pkg/aggregator"
)

func RouterV1(router fiber.Router, aggregatorService *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) {
	router.Post("/playlists/convert", ConvertPlaylistController(aggregatorService, db))
	router.Get("/playlists/supported", GetSupportedPlatformsController(aggregatorService, db))
}
