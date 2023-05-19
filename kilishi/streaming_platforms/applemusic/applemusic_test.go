package applemusic_test

import (
	"fmt"
	"testing"

	"github.com/prettyirrelevant/kilishi/streaming_platforms/applemusic"
	"github.com/stretchr/testify/assert"
)

func TestAppleMusic(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	apple := applemusic.New(applemusic.InitialisationOpts{})

	playlist, err := apple.GetPlaylist("https://music.apple.com/us/playlist/pl.2fc68f6d68004ae993dadfe99de83877")
	is.NoError(err)
	is.NotEmpty(playlist)

	fmt.Printf("Playlist Title: %s\n\n", playlist.Title)
	fmt.Printf("Playlist Description: %s\n\n", playlist.Description)
	for i, track := range playlist.Tracks {
		fmt.Printf("\nTrack(%d): %s\nArtistes: %+v\n\n\n", i+1, track.Title, track.Artists)
	}
}
