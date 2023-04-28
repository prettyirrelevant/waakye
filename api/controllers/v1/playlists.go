package v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/prettyirrelevant/waakye/api"
	"github.com/prettyirrelevant/waakye/api/aggregator"
)

// create a new validator instance
var validate = validator.New()

// ConvertPlaylistRequest is a struct that represents the request body for the ConvertPlaylistController function.
type ConvertPlaylistRequest struct {
	From        api.MusicStreamingPlatform `validate:"required,oneof=spotify ytmusic deezer"`
	To          api.MusicStreamingPlatform `validate:"required,oneof=spotify ytmusic deezer,nefield=From"`
	PlaylistURL string                     `validate:"required,url"`
}

// ConvertPlaylistController returns a handler function for converting a playlist
func ConvertPlaylistController(aggregator *aggregator.MusicStreamingPlatformsAggregator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// parse the request body into a ConvertPlaylistRequest struct
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

// GetSupportedPlatformsController returns a handler function for getting the list of supported music streaming platforms.
func GetSupportedPlatformsController(aggregator *aggregator.MusicStreamingPlatformsAggregator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		supportedPlatforms := aggregator.SupportedPlatforms()
		return c.JSON(supportedPlatforms)
	}
}
