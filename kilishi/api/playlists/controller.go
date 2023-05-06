package playlists

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/pkg/aggregator"
	"github.com/prettyirrelevant/kilishi/pkg/utils/types"
)

// ConvertPlaylistController returns a handler function for converting a playlist
func ConvertPlaylistController(aggregator *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody ConvertPlaylistRequest

		err := c.BodyParser(&requestBody)
		if err != nil {
			return c.
				Status(http.StatusBadRequest).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		if err = validate.Struct(requestBody); err != nil {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		dbCredentials, err := db.GetOauthCredentials(requestBody.Destination)
		if err != nil {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		credentials, err := types.OauthCredentialsFromDB(dbCredentials.Credentials)
		if err != nil {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		playlistURL, err := aggregator.ConvertPlaylist(requestBody.Source, requestBody.Destination, requestBody.PlaylistURL, credentials.AccessToken)
		if err != nil {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": true, "data": playlistURL, "error": nil})
	}
}

// GetSupportedPlatformsController returns a handler function for getting the list of supported music streaming platforms.
func GetSupportedPlatformsController(aggregator *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		supportedPlatforms := aggregator.SupportedPlatforms()
		return c.JSON(supportedPlatforms)
	}
}
