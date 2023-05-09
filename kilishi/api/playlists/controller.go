package playlists

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/api/presenter"
	"github.com/prettyirrelevant/kilishi/pkg/aggregator"
	"github.com/prettyirrelevant/kilishi/pkg/utils"
)

// ConvertPlaylistController returns a handler function for converting a playlist
func ConvertPlaylistController(aggregator *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody ConvertPlaylistRequest

		err := c.BodyParser(&requestBody)
		if err != nil {
			return c.
				Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("validation error", err))
		}

		if ok, errors := requestBody.Validate(); !ok {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(presenter.ErrorResponse("validation error", errors...))
		}

		dbCredentials, err := db.GetOauthCredentials(requestBody.Destination)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("error retrieving credentials from database", err))
		}

		credentials, err := utils.OauthCredentialsFromDB(dbCredentials.Credentials)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("", err))
		}

		playlistURL, err := aggregator.ConvertPlaylist(requestBody.Source, requestBody.Destination, requestBody.PlaylistURL, credentials.AccessToken)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("error converting playlist", err))
		}

		return c.Status(http.StatusOK).JSON(presenter.SuccessResponse("Playlist converted successfully", playlistURL))
	}
}

// GetSupportedPlatformsController returns a handler function for getting the list of supported music streaming platforms.
func GetSupportedPlatformsController(aggregator *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(presenter.SuccessResponse("All supported music streaming platforms returned successfully!", aggregator.SupportedPlatforms()))
	}
}
