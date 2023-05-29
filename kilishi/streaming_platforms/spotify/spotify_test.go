package spotify_test

import (
	"os"
	"testing"

	"github.com/imroc/req/v3"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/spotify"
	"github.com/stretchr/testify/assert"
)

func TestSpotify(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	spotify := spotify.New(&spotify.InitialisationOpts{
		RequestClient:             req.C(),
		BaseAPIURL:                "https://api.spotify.com/v1",
		AuthenticationURL:         "https://accounts.spotify.com/api/token",
		UserID:                    os.Getenv("SPOTIFY_USER_ID"),
		ClientID:                  os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret:              os.Getenv("SPOTIFY_CLIENT_SECRET"),
		AuthenticationRedirectURL: os.Getenv("SPOTIFY_AUTHENTICATION_REDIRECT_URL"),
	})

	playlist, err := spotify.GetPlaylist("https://open.spotify.com/playlist/0lLDqS7JwiuvpEgLuc92Dw?si=0180cdf3acb8427a")
	is.NoError(err)
	is.Len(playlist.Tracks, 200)
}
