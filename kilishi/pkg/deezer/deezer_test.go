package deezer_test

import (
	"fmt"
	"testing"

	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/assert"

	"github.com/prettyirrelevant/kilishi/pkg/deezer"
)

func TestDeezer(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	deezer := deezer.New(deezer.InitialisationOpts{
		RequestClient:     req.C(),
		BaseAPIURI:        "https://api.deezer.com",
		AuthenticationURI: "https://connect.deezer.com/oauth/access_token.php",
	})

	// Top Nigeria 2020
	playlist, err := deezer.GetPlaylist("https://www.deezer.com/en/playlist/8424712382")
	is.NoError(err)
	is.NotEmpty(playlist)
	is.Len(playlist.Tracks, 49)
	fmt.Printf("Playlist Title: %s\n\n", playlist.Title)
	fmt.Printf("Playlist Description: %s\n\n", playlist.Description)
	for i, track := range playlist.Tracks {
		fmt.Printf("\nTrack(%d): %s\nArtistes: %+v\n\n\n", i+1, track.Title, track.Artists)
	}
}
