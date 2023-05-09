package ytmusic

import (
	"fmt"
	"sync"

	"github.com/prettyirrelevant/kilishi/pkg/utils"
	"github.com/prettyirrelevant/ytmusicapi"
)

// New initialises a `YTMusic` object.
func New() *YTMusic {
	ytmusicapi.Setup()
	return &YTMusic{}
}

// GetPlaylist returns information about a playlist.
func (y *YTMusic) GetPlaylist(playlistURI string) (utils.Playlist, error) {
	playlistID, err := parsePlaylistURI(playlistURI)
	if err != nil {
		return utils.Playlist{}, err
	}

	playlist, err := ytmusicapi.GetPlaylist(playlistID, 0)
	if err != nil {
		return utils.Playlist{}, fmt.Errorf("ytmusic: %s", err.Error())
	}

	return parseGetPlaylistResponse(playlist), nil
}

// CreatePlaylist creates a new playlist using the information provided.
func (*YTMusic) CreatePlaylist(playlist utils.Playlist, accessToken string) (string, error) {
	var foundTracks []utils.Track
	var wg sync.WaitGroup

	for _, trackEntry := range playlist.Tracks {
		wg.Add(1)
		go func(payload utils.Track) {
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

func (*YTMusic) GetAuthorizationCode(code string) (utils.OauthCredentials, error) {
	return utils.OauthCredentials{}, nil // no-op
}

// RequiresAccessToken specifies if the streaming requires Oauth.
func (*YTMusic) RequiresAccessToken() bool {
	return false // no-op
}
