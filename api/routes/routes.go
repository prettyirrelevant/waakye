package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/prettyirrelevant/waakye/api/aggregator"
	v1 "github.com/prettyirrelevant/waakye/api/controllers/v1"
)

func RouterV1(app fiber.Router, aggregatorService *aggregator.MusicStreamingPlatformsAggregator) {
	app.Post("/playlists/convert", v1.ConvertPlaylistController(aggregatorService))
	app.Get("/playlists/supported", v1.GetSupportedPlatformsController(aggregatorService))

	app.Get("/spotify/callback", v1.SpotifyOauthCallbackController(aggregatorService))
	app.Get("/deezer/callback", v1.DeezerOauthCallbackController(aggregatorService))
}
