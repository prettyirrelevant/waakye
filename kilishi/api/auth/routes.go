package auth

import (
	"github.com/gofiber/fiber/v2"

	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/aggregator"
)

func RouterV1(router fiber.Router, aggregatorService *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) {
	router.Get("/v1/auth/deezer/callback", DeezerOauthCallbackController(aggregatorService, db))
	router.Get("/v1/auth/spotify/callback", SpotifyOauthCallbackController(aggregatorService, db))
	router.Post("/v1/auth/spotify/refresh", SpotifyRefreshAccessTokenController(aggregatorService, db))
}
