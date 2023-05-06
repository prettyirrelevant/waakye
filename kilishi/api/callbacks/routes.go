package callbacks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/pkg/aggregator"
)

func RouterV1(router fiber.Router, aggregatorService *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) {
	router.Get("/auth/spotify/callback", SpotifyOauthCallbackController(aggregatorService, db))
	router.Get("/auth/deezer/callback", DeezerOauthCallbackController(aggregatorService, db))
}
