package callbacks

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prettyirrelevant/kilishi/aggregator"
	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/api/presenter"
	"github.com/prettyirrelevant/kilishi/utils"
)

// SpotifyOauthCallbackController handles Spotify OAuth callback requests.
func SpotifyOauthCallbackController(ag *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	// TODO: Log extensively
	return func(c *fiber.Ctx) error {
		state := c.Query("state")
		code := c.Query("code")
		if state == "" || code == "" {
			return c.
				Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("missing required query parameters: state and/or code"))
		}

		stateParamDecrypted, err := utils.Decrypt(state, ag.Config.SecretKey, ag.Config.InitializationVector)
		if err != nil {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(presenter.ErrorResponse("invalid state parameter provided", err.Error()))
		}

		// It takes the format of timeInMicroSecs:streamingPlatform
		stateParamSlice := strings.Split(stateParamDecrypted, ":")
		if len(stateParamSlice) < 2 {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(presenter.ErrorResponse("invalid state parameter provided", err.Error()))
		}

		stateParamTime, err := strconv.Atoi(stateParamSlice[0])
		if err != nil {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(presenter.ErrorResponse("invalid state parameter provided", err.Error()))
		}

		if stateParamSlice[1] != string(aggregator.Spotify) || time.Now().Unix()-int64(stateParamTime) >= 60 {
			return c.
				Status(http.StatusUnprocessableEntity).
				JSON(presenter.ErrorResponse("invalid/expired state parameter provided", err.Error()))
		}

		oauthCredentials, err := ag.Spotify.GetAuthorizationCode(c.Query("code"))
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("unable to retrieve authorization code", err.Error()))
		}

		err = db.SetOauthCredentials(aggregator.Spotify, oauthCredentials)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("unable to store authorization code in database", err.Error()))
		}

		return c.
			Status(http.StatusOK).
			JSON(presenter.SuccessResponse("spotify token saved successfully", nil))
	}
}

// DeezerOauthCallbackController handles Deezer OAuth callback requests.
func DeezerOauthCallbackController(ag *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		oauthCredentials, err := ag.Deezer.GetAuthorizationCode(c.Query("code"))
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("unable to retrieve authorization code", err.Error()))
		}

		err = db.SetOauthCredentials(aggregator.Deezer, oauthCredentials)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("unable to store authorization code in database", err.Error()))
		}

		return c.
			Status(http.StatusOK).
			JSON(presenter.SuccessResponse("deezer token saved successfully", nil))
	}
}
