package ytmusic

import (
	"fmt"

	"github.com/prettyirrelevant/kilishi/utils"
)

// New initializes a `YTMusic` object.
func New(opts *InitialisationOpts) *YTMusic {
	return &YTMusic{
		RequestClient: setupRequestClient(opts.RequestClient, opts.BaseAPIURL),
		Config: Config{
			BaseAPIURL:          opts.BaseAPIURL,
			AuthenticationToken: opts.AuthenticationToken,
		},
	}
}

// GetPlaylist returns information about a playlist.
func (y *YTMusic) GetPlaylist(playlistURL string) (utils.Playlist, error) {
	var response ytmusicAPIGetPlaylistResponse
	err := y.RequestClient.
		Post("/playlists").
		SetBody(map[string]string{"url": playlistURL}).
		Do().
		Into(&response)

	if err != nil {
		return utils.Playlist{}, err
	}
	return parseGetPlaylistResponse(response), nil
}

// CreatePlaylist creates a new playlist using the information provided.
func (y *YTMusic) CreatePlaylist(playlist utils.Playlist, _ string) (string, error) {
	var trackIDs []string
	for _, entry := range playlist.Tracks {
		if ok := utils.Contains(playlist.Tracks, entry); ok {
			trackIDs = append(trackIDs, entry.ID)
		}
	}

	var response ytmusicAPICreatePlaylistResponse
	err := y.RequestClient.
		Put("/playlists").
		SetBearerAuthToken(y.Config.AuthenticationToken).
		SetBody(map[string]interface{}{
			"title":       playlist.Title,
			"description": playlist.Description,
			"track_ids":   trackIDs,
		}).
		Do().
		Into(&response)

	if err != nil {
		return "", err
	}
	return response.Data, nil
}

// LookupTrack searches for track on YTMusic and appends the top result to a slice.
func (y *YTMusic) LookupTrack(track utils.Track) (utils.Track, error) {
	var foundTrack utils.Track
	var response ytmusicAPISearchResponse
	err := y.RequestClient.
		Post("/tracks/search").
		SetBody(map[string]interface{}{
			"q":               trackToSearchQuery(track),
			"ignore_spelling": true,
			"limit":           3,
		}).
		Do().
		Into(&response)

	if err != nil {
		return foundTrack, err
	}
	if len(response.Data) == 0 {
		return foundTrack, fmt.Errorf("ytmusic: no track found that matches %s", track.Title)
	}

	foundTrack.Title = response.Data[0].Title
	foundTrack.ID = response.Data[0].Identifier
	foundTrack.Artists = response.Data[0].Artists
	return foundTrack, nil
}

func (*YTMusic) GetAuthorizationCode(_ string) (utils.OauthCredentials, error) {
	return utils.OauthCredentials{}, nil // no-op
}

// RequiresAccessToken specifies if the streaming requires Oauth.
func (*YTMusic) RequiresAccessToken() bool {
	return false // no-op
}
