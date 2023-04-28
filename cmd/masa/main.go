package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/prettyirrelevant/waakye/config"
)

func main() {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	// app.Use(basicauth.New(basicauth.Config{
	// 	Users: map[string]string{"eniola": "eniola"},
	// }))
	app.Get("/metrics", monitor.New())

	apiGroup := app.Group("/api")
	apiGroup.Post("/auth/spotify", spotifyAuthController(config))
	apiGroup.Post("/auth/deezer", deezerAuthController(config))

	log.Fatal(app.Listen(":5001"))
}

func spotifyAuthController(config *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, allocatorCancel, cancel := setupChromeDp()
		defer allocatorCancel()
		defer cancel()

		chromedp.Run(ctx, spotifyAuthenticationTask(config))
		return c.Status(http.StatusOK).Send(nil)
	}
}

func deezerAuthController(config *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, allocatorCancel, cancel := setupChromeDp()
		defer allocatorCancel()
		defer cancel()

		chromedp.Run(ctx, deezerAuthenticationTask(config))
		return c.Status(http.StatusOK).Send(nil)
	}
}

// spotifyAuthenticationTask represents an automation task in chromedp for Spotify authentication
func spotifyAuthenticationTask(config *config.Config) chromedp.Tasks {
	queryParams := url.Values{}
	queryParams.Add("response_type", "code")
	queryParams.Add("client_id", config.SpotifyClientID)
	queryParams.Add("redirect_uri", config.SpotifyAuthRedirectURI)
	queryParams.Add("scope", "playlist-modify-public")

	return chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("https://accounts.spotify.com/authorize?%s", queryParams.Encode())),
		chromedp.SendKeys(`#login-username`, config.SpotifyAuthEmail),
		chromedp.SendKeys(`#login-password`, config.SpotifyAuthPassword),
		chromedp.Click(`#login-button`, chromedp.ByID),
		chromedp.WaitVisible(`//*[contains(text(), "spotify token saved successfully")]`),
		// chromedp.FullScreenshot(screenshotData, 100),
	}
}

// deezerAuthenticationTask represents an automation task in chromedp for Deezer authentication
func deezerAuthenticationTask(config *config.Config) chromedp.Tasks {
	queryParams := url.Values{}
	queryParams.Add("app_id", config.DeezerAppID)
	queryParams.Add("redirect_uri", config.DeezerAuthRedirectURI)
	queryParams.Add("perms", "manage_library,offline_access")

	return chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("https://connect.deezer.com/oauth/auth.php?%s", queryParams.Encode())),
		chromedp.SendKeys(`#login_mail`, config.DeezerAuthEmail),
		chromedp.SendKeys(`#login_password`, config.DeezerAuthPassword),
		chromedp.Click(`#login_form_submit`, chromedp.ByID),
		chromedp.WaitVisible(`//*[contains(text(), "deezer token saved successfully")]`),
		// chromedp.FullScreenshot(screenshotData, 100),
	}
}

func setupChromeDp() (context.Context, context.CancelFunc, context.CancelFunc) {
	actx, acancel := chromedp.NewExecAllocator(context.Background(), chromedp.DisableGPU, chromedp.Headless)
	ctx, cancel := chromedp.NewContext(actx,
		chromedp.WithLogf(log.Printf),
		chromedp.WithDebugf(log.Printf),
		chromedp.WithErrorf(log.Printf),
	)

	return ctx, acancel, cancel
}
