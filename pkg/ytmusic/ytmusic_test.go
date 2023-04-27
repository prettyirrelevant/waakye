package ytmusic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prettyirrelevant/waakye/pkg/utils/types"
	"github.com/prettyirrelevant/waakye/pkg/ytmusic"
	"github.com/prettyirrelevant/ytmusicapi"
)

func TestYTMusic(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Test get playlist
	ytmusic := ytmusic.New()
	playlist, err := ytmusic.GetPlaylist("https://music.youtube.com/playlist?list=PL4fGSI1pDJn5dHScZlGIe6TEoGzFv_qZE")
	is.NoError(err)
	is.Len(playlist.Tracks, 100)

	// Test create playlist
	id, err := ytmusic.CreatePlaylist(types.Playlist{Title: "Hello!", Description: "Hahahahahahaha!", Tracks: playlist.Tracks[:20]}, "")
	is.NoError(err)
	is.NotEmpty(id)
	// delete the created playlist
	err = ytmusicapi.DeletePlaylist(id)
	is.NoError(err)
}
