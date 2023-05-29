package deezer_test

import (
	"testing"

	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/assert"

	"github.com/prettyirrelevant/kilishi/streaming_platforms/deezer"
)

func TestDeezer(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	deezer := deezer.New(&deezer.InitialisationOpts{
		RequestClient:     req.C(),
		BaseAPIURL:        "https://api.deezer.com",
		AuthenticationURL: "https://connect.deezer.com/oauth/access_token.php",
	})

	// Top Nigeria 2020
	playlist, err := deezer.GetPlaylist("https://www.deezer.com/en/playlist/8424712382")
	is.NoError(err)
	is.NotEmpty(playlist)
	is.Len(playlist.Tracks, 49)
}
