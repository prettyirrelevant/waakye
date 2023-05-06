package callbacks

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/pkg/aggregator"
	ag "github.com/prettyirrelevant/kilishi/pkg/aggregator"
	"github.com/prettyirrelevant/kilishi/pkg/utils/cryptography"
)

// SpotifyOauthCallbackController handles Spotify OAuth callback requests.
func SpotifyOauthCallbackController(aggregator *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		state := c.Query("state")
		code := c.Query("code")
		if state == "" || code == "" {
			return c.
				Status(http.StatusBadRequest).
				JSON(fiber.Map{"status": false, "data": nil, "error": "missing required query parameters, `state` & `code`"})
		}

		stateParamDecrypted, err := cryptography.Decrypt(state, aggregator.Config.SecretKey, aggregator.Config.InitializationVector)
		if err != nil {
			return c.
				Status(http.StatusBadRequest).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		// It takes the format of timeInMicroSecs:streamingPlatform
		stateParamSlice := strings.Split(stateParamDecrypted, ":")
		if len(stateParamSlice) < 2 {
			return c.
				Status(http.StatusBadRequest).
				JSON(fiber.Map{"status": false, "data": nil, "error": "invalid `state` parameter"})
		}

		stateParamTime, err := strconv.Atoi(stateParamSlice[0])
		if err != nil {
			return c.
				Status(http.StatusBadRequest).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		if stateParamSlice[1] != string(ag.Spotify) || time.Now().UnixMilli()-int64(stateParamTime) >= 60000 {
			return c.
				Status(http.StatusBadRequest).
				JSON(fiber.Map{"status": false, "data": nil, "error": "provided `state` query parameter has expired."})
		}

		oauthCredentials, err := aggregator.Spotify.GetAuthorizationCode(c.Query("code"))
		if err != nil {
			fmt.Printf("Encountered error: %s", err.Error())
			return c.
				Status(http.StatusInternalServerError).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		err = db.SetOauthCredentials(ag.Spotify, oauthCredentials)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		return c.
			Status(http.StatusOK).
			JSON(fiber.Map{"status": true, "data": "spotify token saved successfully", "error": nil})
	}
}

// DeezerOauthCallbackController handles Deezer OAuth callback requests.
func DeezerOauthCallbackController(aggregator *aggregator.MusicStreamingPlatformsAggregator, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		oauthCredentials, err := aggregator.Deezer.GetAuthorizationCode(c.Query("code"))
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		err = db.SetOauthCredentials(ag.Deezer, oauthCredentials)
		if err != nil {
			return c.
				Status(http.StatusInternalServerError).
				JSON(fiber.Map{"status": false, "data": nil, "error": err.Error()})
		}

		return c.
			Status(http.StatusOK).
			JSON(fiber.Map{"status": true, "data": "deezer token saved successfully", "error": nil})
	}
}
