package callbacks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prettyirrelevant/kilishi/aggregator"
	"github.com/prettyirrelevant/kilishi/api/database"
)

func RouterV1(router fiber.Router, aggregatorService *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) {
	router.Get("/v1/auth/spotify/callback", SpotifyOauthCallbackController(aggregatorService, db))
	router.Get("/v1/auth/deezer/callback", DeezerOauthCallbackController(aggregatorService, db))
}
