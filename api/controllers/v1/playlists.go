package v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/prettyirrelevant/waakye/api"
	"github.com/prettyirrelevant/waakye/api/aggregator"
)

var validate = validator.New()

type ConvertPlaylistRequest struct {
	From        api.MusicStreamingPlatform `validate:"required,oneof=spotify ytmusic deezer"`
	To          api.MusicStreamingPlatform `validate:"required,oneof=spotify ytmusic deezer,nefield=From"`
	PlaylistURL string                     `validate:"required,url"`
}

func ConvertPlaylistController(aggregator *aggregator.MusicStreamingPlatformsAggregator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody ConvertPlaylistRequest
		err := c.BodyParser(&requestBody)
		if err != nil {
			return c.
				Status(http.StatusBadRequest).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		// validate the request struct.
		if err = validate.Struct(requestBody); err != nil {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		// convert the playlist
		playlistURL, err := aggregator.ConvertPlaylist(requestBody.From, requestBody.To, requestBody.PlaylistURL)
		if err != nil {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": true, "data": playlistURL, "error": nil})
	}
}

func GetSupportedPlatformsController(aggregator *aggregator.MusicStreamingPlatformsAggregator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(aggregator.SupportedPlatforms())
	}
}
