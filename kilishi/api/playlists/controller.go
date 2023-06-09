package playlists

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/api/presenter"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/aggregator"
	"github.com/prettyirrelevant/kilishi/utils"
)

func GetPlaylistController(ag *aggregator.MusicStreamingPlatformsAggregator, _ *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var queryParams GetPlaylistRequest

		err := c.QueryParser(&queryParams)
		if err != nil {
			return c.
				Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("validation error", err.Error()))
		}

		if ok, errors := queryParams.Validate(); !ok {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(presenter.ErrorResponse("validation error", errors...))
		}

		x := ag.GetStreamingPlatform(queryParams.Platform)
		playlist, err := x.GetPlaylist(queryParams.PlaylistURL)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("error retrieving playlist", err.Error()))
		}

		return c.Status(http.StatusOK).JSON(presenter.SuccessResponse("playlist retrieved successfully", playlist))
	}
}

func FindTrackController(ag *aggregator.MusicStreamingPlatformsAggregator, _ *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var queryParams FindTrackRequest

		err := c.QueryParser(&queryParams)
		if err != nil {
			return c.
				Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("validation error", err.Error()))
		}

		if ok, errors := queryParams.Validate(); !ok {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(presenter.ErrorResponse("validation error", errors...))
		}

		x := ag.GetStreamingPlatform(queryParams.Platform)
		track, err := x.LookupTrack(utils.Track{Title: queryParams.Title, Artists: queryParams.Artists})
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("error searching for track", err.Error()))
		}

		return c.Status(http.StatusOK).JSON(presenter.SuccessResponse("track found successfully", track))
	}
}

func CreatePlaylistController(ag *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody ConvertPlaylistRequest

		err := c.BodyParser(&requestBody)
		if err != nil {
			return c.
				Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("validation error", err.Error()))
		}

		if ok, errors := requestBody.Validate(); !ok {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(presenter.ErrorResponse("validation error", errors...))
		}

		x := ag.GetStreamingPlatform(requestBody.Platform)
		accessToken := strings.TrimSpace(requestBody.AccessToken)
		if accessToken == "" && x.RequiresAccessToken() {
			credentialsInDB, _err := db.GetDBOauthCredentials(requestBody.Platform)
			if _err != nil {
				return c.
					Status(http.StatusInternalServerError).
					JSON(presenter.ErrorResponse("error fetching access token from db", _err.Error()))
			}

			var credentials utils.OauthCredentials
			err = credentials.FromDB(credentialsInDB.Credentials)
			if err != nil {
				return c.
					Status(http.StatusInternalServerError).
					JSON(presenter.ErrorResponse("error parsing access token from db", err.Error()))
			}

			accessToken = credentials.AccessToken
		}

		playlistURL, err := x.CreatePlaylist(requestBody.Playlist, accessToken)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("error creating playlist", err.Error()))
		}

		return c.Status(http.StatusOK).JSON(presenter.SuccessResponse("playlist created successfully", playlistURL))
	}
}

// GetSupportedPlatformsController returns a handler function for getting the list of supported music streaming platforms.
func GetSupportedPlatformsController(ag *aggregator.MusicStreamingPlatformsAggregator, _ *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(presenter.SuccessResponse("All supported music streaming platforms returned successfully!", ag.SupportedPlatforms()))
	}
}
