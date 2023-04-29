package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/prettyirrelevant/kilishi/api/aggregator"
	v1 "github.com/prettyirrelevant/kilishi/api/controllers/v1"
)

// RouterV1 registers v1 routes for the given fiber.Router and aggregator service
func RouterV1(router fiber.Router, aggregatorService *aggregator.MusicStreamingPlatformsAggregator) {
	router.Post("/playlists/convert", v1.ConvertPlaylistController(aggregatorService))

	router.Get("/playlists/supported", v1.GetSupportedPlatformsController(aggregatorService))

	router.Get("/auth/spotify/callback", v1.SpotifyOauthCallbackController(aggregatorService))

	router.Get("/auth/deezer/callback", v1.DeezerOauthCallbackController(aggregatorService))
}
