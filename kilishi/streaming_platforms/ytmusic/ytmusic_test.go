package ytmusic_test

import (
	"os"
	"testing"

	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/assert"

	"github.com/prettyirrelevant/kilishi/streaming_platforms/ytmusic"
	"github.com/prettyirrelevant/kilishi/utils"
)

func TestYTMusic(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Test get playlist
	ytmusic := ytmusic.New(&ytmusic.InitialisationOpts{
		RequestClient:       req.C(),
		BaseAPIURL:          os.Getenv("YTMUSICAPI_BASE_URL"),
		AuthenticationToken: os.Getenv("SECRET_KEY"),
	})
	playlist, err := ytmusic.GetPlaylist("https://music.youtube.com/playlist?list=PL4fGSI1pDJn5dHScZlGIe6TEoGzFv_qZE")
	is.NoError(err)
	is.Len(playlist.Tracks, 100)

	// Test create playlist
	id, err := ytmusic.CreatePlaylist(utils.Playlist{Title: "Hello!", Description: "Hahahahahahaha!", Tracks: playlist.Tracks[:20]}, "")
	is.NoError(err)
	is.NotEmpty(id)
}
