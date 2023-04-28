package v1

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prettyirrelevant/waakye/api"
	"github.com/prettyirrelevant/waakye/api/aggregator"
	"github.com/prettyirrelevant/waakye/pkg/utils/cryptography"
)

// SpotifyOauthCallbackController handles Spotify OAuth callback requests.
func SpotifyOauthCallbackController(aggregator *aggregator.MusicStreamingPlatformsAggregator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Decrypt the state parameter received from Spotify to check if it's valid.
		stateParam := c.Query("state")
		stateParamDecrypted, err := cryptography.Decrypt(stateParam, aggregator.Config.SecretKey)
		if err != nil {
			return c.
				Status(http.StatusBadRequest).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		// Split the decrypted state parameter to obtain the time it was created.
		stateParamSlice := strings.Split(stateParamDecrypted, ":")
		stateParamTime, err := strconv.Atoi(stateParamSlice[0])
		if err != nil {
			return c.
				Status(http.StatusBadRequest).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		// Check if the state parameter is expired or invalid.
		if stateParamSlice[1] != string(api.Spotify) || time.Duration(time.Now().UnixMicro()-int64(stateParamTime)) >= 3*time.Minute {
			return c.
				Status(http.StatusBadRequest).
				JSON(fiber.Map{"status": false, "data": nil, "error": "provided `state` query parameter has expired."})
		}

		// Get the OAuth credentials from Spotify using the authorization code.
		oauthCredentials, err := aggregator.Spotify.GetAuthorizationCode(c.Query("code"))
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		// Store the OAuth credentials in the database.
		err = aggregator.Database.SetOauthCredentials(api.Spotify, oauthCredentials)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		// Return a success response if the OAuth credentials were successfully stored in the database.
		return c.
			Status(http.StatusOK).
			JSON(fiber.Map{"status": true, "data": "spotify token saved successfully", "error": nil})
	}
}

// DeezerOauthCallbackController handles Deezer OAuth callback requests.
func DeezerOauthCallbackController(aggregator *aggregator.MusicStreamingPlatformsAggregator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		oauthCredentials, err := aggregator.Deezer.GetAuthorizationCode(c.Query("code"))
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		err = aggregator.Database.SetOauthCredentials(api.Deezer, oauthCredentials)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		// Return a success response if the OAuth credentials were successfully stored in the database.
		return c.
			Status(http.StatusOK).
			JSON(fiber.Map{"status": true, "data": "deezer token saved successfully", "error": nil})
	}
}
