package deezer

import (
	"strings"
	"sync"

	"github.com/prettyirrelevant/kilishi/pkg/utils"
)

// New initialises a `Spotify` object.
func New(opts InitialisationOpts) *Deezer {
	return &Deezer{
		RequestClient: setupRequestClient(opts.RequestClient),
		Config: Config{
			AppID:             opts.AppID,
			BaseAPIURI:        opts.BaseAPIURI,
			ClientSecret:      opts.ClientSecret,
			AuthenticationURI: opts.AuthenticationURI,
		},
	}
}

func (d *Deezer) GetPlaylist(playlistURI string) (utils.Playlist, error) {
	playlistID, err := parsePlaylistURI(playlistURI)
	if err != nil {
		return utils.Playlist{}, err
	}

	var response deezerAPIGetPlaylistResponse
	err = d.RequestClient.
		Get(d.Config.BaseAPIURI + "/playlist/" + playlistID).
		Do().
		Into(&response)

	if err != nil {
		return utils.Playlist{}, err
	}

	return parseGetPlaylistResponse(response), nil
}

func (d *Deezer) CreatePlaylist(playlist utils.Playlist, accessToken string) (string, error) {
	// first, look for the tracks on Spotify
	var tracksFound []utils.Track
	var wg sync.WaitGroup
	for _, entry := range playlist.Tracks {
		wg.Add(1)

		go func(track utils.Track) {
			defer wg.Done()
			d.lookupTrack(track, &tracksFound)
		}(entry)
	}
	wg.Wait()

	// then, create the playlist.
	var response deezerAPICreatePlaylistResponse
	err := d.RequestClient.
		Post(d.Config.BaseAPIURI + "/users/me/playlist").
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

	// finally, add the tracks in batches
	err = d.populatePlaylistWithTracks(playlist.Tracks, playlist.ID, accessToken)
	if err != nil {
		return "", err
	}

	return response.ID, nil
}

func (d *Deezer) GetAuthorizationCode(code string) (utils.OauthCredentials, error) {
	var response deezerAPIBearerCredentialsResponse
	err := d.RequestClient.
		Get(d.Config.AuthenticationURI).
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

func (d *Deezer) lookupTrack(track utils.Track, tracksFound *[]utils.Track) {
	var response deezerAPISearchTrackResponse
	err := d.RequestClient.
		Get(d.Config.BaseAPIURI + "/search/track").
		SetContentType(utils.ApplicationJSON).
		SetQueryParams(map[string]string{
			"q": trackToSearchQuery(track),
		}).
		Do().
		Into(&response)

	if err != nil {
		return
	}

	if len(response.Results.Data) == 0 {
		return
	}

	*tracksFound = append(*tracksFound, parseTracksResponse(response.Results)[0])
}

// populatePlaylistWithTracks adds tracks found on Deezer to a newly created playlist.
func (d *Deezer) populatePlaylistWithTracks(tracks []utils.Track, playlistID, accessToken string) error {
	var tracksURI []string
	for _, track := range tracks {
		tracksURI = append(tracksURI, track.ID)
	}

	var response any
	err := d.RequestClient.
		Post(d.Config.BaseAPIURI + "/playlists/" + playlistID + "/tracks").
		SetBearerAuthToken(accessToken).
		SetContentType(utils.ApplicationJSON).
		SetFormData(map[string]string{
			"songs": strings.Join(tracksURI, ","),
		}).
		Do().
		Into(&response)

	return err
}
