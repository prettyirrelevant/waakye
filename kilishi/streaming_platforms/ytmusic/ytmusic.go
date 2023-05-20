package ytmusic

import (
	"fmt"

	"github.com/prettyirrelevant/kilishi/utils"
	"github.com/prettyirrelevant/ytmusicapi"
)

// New initialises a `YTMusic` object.
func New() *YTMusic {
	ytmusicapi.Setup()
	return &YTMusic{}
}

// GetPlaylist returns information about a playlist.
func (y *YTMusic) GetPlaylist(playlistURL string) (utils.Playlist, error) {
	playlistID, err := parsePlaylistURL(playlistURL)
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
func (y *YTMusic) CreatePlaylist(playlist utils.Playlist, accessToken string) (string, error) {
	var trackIDs []string
	for _, entry := range playlist.Tracks {
		trackIDs = append(trackIDs, entry.ID)
	}

	playlistID, err := ytmusicapi.CreatePlaylist(playlist.Title, playlist.Description, ytmusicapi.PUBLIC, "", trackIDs)
	if err != nil {
		return "", fmt.Errorf("ytmusic: %s", err.Error())
	}
	return playlistID, nil
}

// LookupTrack searches for track on YTMusic and appends the top result to a slice.
func (*YTMusic) LookupTrack(track utils.Track) (utils.Track, error) {
	var foundTrack utils.Track

	searchResults, err := ytmusicapi.SearchTracks(trackToSearchQuery(track), ytmusicapi.SongsFilter, ytmusicapi.NoScope, 5, false)
	if err != nil {
		return foundTrack, fmt.Errorf("ytmusic: %s", err.Error())
	}

	return utils.Track{ID: searchResults[0].VideoID, Title: searchResults[0].Title, Artists: searchResults[0].Artistes}, nil
}

func (*YTMusic) GetAuthorizationCode(code string) (utils.OauthCredentials, error) {
	return utils.OauthCredentials{}, nil // no-op
}

// RequiresAccessToken specifies if the streaming requires Oauth.
func (*YTMusic) RequiresAccessToken() bool {
	return false // no-op
}
