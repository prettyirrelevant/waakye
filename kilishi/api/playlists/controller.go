package playlists

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/api/presenter"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/aggregator"
)

func GetPlaylistController(aggregator *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody GetPlaylistRequest

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

		x := aggregator.GetStreamingPlatform(requestBody.Platform)
		playlist, err := x.GetPlaylist(requestBody.PlaylistURL)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("error retrieving playlist", err.Error()))
		}

		return c.Status(http.StatusOK).JSON(presenter.SuccessResponse("playlist retrieved successfully", playlist))
	}
}

func FindTrackController(aggregator *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody FindTrackRequest

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

		x := aggregator.GetStreamingPlatform(requestBody.Platform)
		track, err := x.LookupTrack(requestBody.Track)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("error searching for track", err.Error()))
		}

		return c.Status(http.StatusOK).JSON(presenter.SuccessResponse("track found successfully", track))
	}
}

func CreatePlaylistController(aggregator *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
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

		x := aggregator.GetStreamingPlatform(requestBody.Platform)
		playlistURL, err := x.CreatePlaylist(requestBody.Playlist, requestBody.AccessToken)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("error creating playlist", err.Error()))
		}

		return c.Status(http.StatusOK).JSON(presenter.SuccessResponse("playlist created successfully", playlistURL))
	}

}

// GetSupportedPlatformsController returns a handler function for getting the list of supported music streaming platforms.
func GetSupportedPlatformsController(aggregator *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(presenter.SuccessResponse("All supported music streaming platforms returned successfully!", aggregator.SupportedPlatforms()))
	}
}
