package ytmusic

import (
	"fmt"
	"sync"

	"github.com/prettyirrelevant/waakye/pkg/utils/types"
	"github.com/prettyirrelevant/ytmusicapi"
)

// New initialises a `YTMusic` object.
func New() *YTMusic {
	ytmusicapi.Setup()
	return &YTMusic{}
}

// GetPlaylist returns information about a playlist.
func (y *YTMusic) GetPlaylist(playlistURI string) (types.Playlist, error) {
	playlistID, err := parsePlaylistURI(playlistURI)
	if err != nil {
		return types.Playlist{}, err
	}

	playlist, err := ytmusicapi.GetPlaylist(playlistID, 0)
	if err != nil {
		return types.Playlist{}, fmt.Errorf("ytmusic: %s", err.Error())
	}

	return parseGetPlaylistResponse(playlist), nil
}

// CreatePlaylist creates a new playlist using the information provided.
func (*YTMusic) CreatePlaylist(playlist types.Playlist) (string, error) {
	var foundTracks []types.Track
	var wg sync.WaitGroup

	for _, trackEntry := range playlist.Tracks {
		wg.Add(1)
		go func(payload types.Track) {
			defer wg.Done()
			lookupTrack(payload, &foundTracks)
		}(trackEntry)
	}
	wg.Wait()

	var trackIDs []string
	for _, entry := range foundTracks {
		trackIDs = append(trackIDs, entry.ID)
	}

	playlistID, err := ytmusicapi.CreatePlaylist(playlist.Title, playlist.Description, ytmusicapi.PUBLIC, "", trackIDs)
	if err != nil {
		return "", fmt.Errorf("ytmusic: %s", err.Error())
	}
	return playlistID, nil
}
