package deezer_test

import (
	"testing"

	"github.com/imroc/req/v3"
	"github.com/prettyirrelevant/waakye/pkg/deezer"
	"github.com/stretchr/testify/assert"
)

func TestDeezer(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	deezer := deezer.New(&deezer.InitialisationOpts{
		RequestClient:     req.C(),
		BaseAPIURI:        "https://api.deezer.com",
		AuthenticationURI: "https://connect.deezer.com/oauth/access_token.php",
	})

	// Top Nigeria 2020
	playlist, err := deezer.GetPlaylist("https://www.deezer.com/en/playlist/8424712382")
	is.NoError(err)
	is.NotEmpty(playlist)
	is.Len(playlist.Tracks, 49)
}
