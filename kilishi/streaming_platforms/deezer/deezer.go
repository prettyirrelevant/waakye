package deezer

import (
	"fmt"
	"strings"

	"github.com/prettyirrelevant/kilishi/utils"
)

var basePlaylistURL = "https://www.deezer.com/en/playlist/"

// New initialises a `Spotify` object.
func New(opts InitialisationOpts) *Deezer {
	return &Deezer{
		RequestClient: setupRequestClient(opts.RequestClient),
		Config: Config{
			AppID:             opts.AppID,
			BaseAPIURL:        opts.BaseAPIURL,
			ClientSecret:      opts.ClientSecret,
			AuthenticationURL: opts.AuthenticationURL,
		},
	}
}

func (d *Deezer) GetPlaylist(playlistURL string) (utils.Playlist, error) {
	playlistID, err := parsePlaylistURL(playlistURL)
	if err != nil {
		return utils.Playlist{}, err
	}

	var response deezerAPIGetPlaylistResponse
	err = d.RequestClient.
		Get(d.Config.BaseAPIURL + "/playlist/" + playlistID).
		Do().
		Into(&response)

	if err != nil {
		return utils.Playlist{}, err
	}

	return parseGetPlaylistResponse(response), nil
}

func (d *Deezer) CreatePlaylist(playlist utils.Playlist, accessToken string) (string, error) {
	var response deezerAPICreatePlaylistResponse
	err := d.RequestClient.
		Post(d.Config.BaseAPIURL + "/users/me/playlist").
		SetContentType(utils.ApplicationJSON).
		SetBearerAuthToken(accessToken).
		SetQueryParams(map[string]string{
			"title": playlist.Title,
		}).
		Do().
		Into(&response)

	if err != nil {
		return "", err
	}

	err = d.populatePlaylistWithTracks(playlist.Tracks, playlist.ID, accessToken)
	if err != nil {
		return "", err
	}

	return basePlaylistURL + response.ID, nil
}

func (d *Deezer) GetAuthorizationCode(code string) (utils.OauthCredentials, error) {
	var response deezerAPIBearerCredentialsResponse
	err := d.RequestClient.
		Get(d.Config.AuthenticationURL).
		SetQueryParams(map[string]string{
			"app_id": d.Config.AppID,
			"secret": d.Config.ClientSecret,
			"code":   code,
			"output": "json",
		}).
		Do().
		Into(&response)

	if err != nil {
		return utils.OauthCredentials{}, err
	}

	return utils.OauthCredentials{AccessToken: response.AccessToken, ExpiresAt: int(response.Expires)}, nil
}

// RequiresAccessToken specifies if the streaming requires Oauth.
func (*Deezer) RequiresAccessToken() bool {
	return true
}

func (d *Deezer) LookupTrack(track utils.Track) (utils.Track, error) {
	var foundTrack utils.Track

	var response deezerAPISearchTrackResponse
	err := d.RequestClient.
		Get(d.Config.BaseAPIURL + "/search/track").
		SetContentType(utils.ApplicationJSON).
		SetQueryParams(map[string]string{
			"q": trackToSearchQuery(track),
		}).
		Do().
		Into(&response)

	if err != nil {
		return foundTrack, fmt.Errorf("deezer: %s", err.Error())
	}

	if len(response.Results.Data) == 0 {
		return foundTrack, fmt.Errorf("deezer: no track found that matches %s", track.Title)
	}

	foundTrack = parseTracksResponse(response.Results)[0]
	return foundTrack, nil
}

// populatePlaylistWithTracks adds tracks found on Deezer to a newly created playlist.
func (d *Deezer) populatePlaylistWithTracks(tracks []utils.Track, playlistID, accessToken string) error {
	var tracksURL []string
	for _, track := range tracks {
		tracksURL = append(tracksURL, track.ID)
	}

	var response any
	err := d.RequestClient.
		Post(d.Config.BaseAPIURL + "/playlists/" + playlistID + "/tracks").
		SetBearerAuthToken(accessToken).
		SetContentType(utils.ApplicationJSON).
		SetFormData(map[string]string{
			"songs": strings.Join(tracksURL, ","),
		}).
		Do().
		Into(&response)

	return err
}
