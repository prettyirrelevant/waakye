package spotify

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/prettyirrelevant/kilishi/utils"
)

const (
	basePlaylistURL              = "https://open.spotify.com/playlist/"
	maximumNumOfTracksPerRequest = 100
)

// New initializes a `Spotify` object.
func New(opts *InitialisationOpts) *Spotify {
	return &Spotify{
		RequestClient: setupRequestClient(opts.RequestClient),
		Config: Config{
			UserID:                    opts.UserID,
			ClientID:                  opts.ClientID,
			BaseAPIURL:                opts.BaseAPIURL,
			ClientSecret:              opts.ClientSecret,
			AuthenticationURL:         opts.AuthenticationURL,
			AuthenticationRedirectURL: opts.AuthenticationRedirectURL,
		},
	}
}

func (s *Spotify) GetPlaylist(playlistURL string) (utils.Playlist, error) {
	playlistID, err := parsePlaylistURL(playlistURL)
	if err != nil {
		return utils.Playlist{}, err
	}

	clientAuthToken, err := s.getClientAuthenticationCredentials()
	if err != nil {
		return utils.Playlist{}, err
	}

	var response spotifyAPIGetPlaylistResponse
	err = s.RequestClient.
		Get(s.Config.BaseAPIURL + "/playlists/" + playlistID).
		SetBearerAuthToken(clientAuthToken).
		SetContentType(utils.ApplicationJSON).
		Do().
		Into(&response)

	if err != nil {
		return utils.Playlist{}, err
	}

	playlist := parseGetPlaylistResponse(&response)
	if response.Tracks.Next == "" {
		return playlist, nil
	}

	// spotify returns at most 100 tracks per request, this while loop gets the other items.
	// TODO: Since offset-limit is used, switch to goroutines to get all items in parallel.
	for {
		var playlistItemsResp spotifyAPITracksResponse
		err := s.RequestClient.
			Get(response.Tracks.Next).
			SetBearerAuthToken(clientAuthToken).
			SetContentType(utils.ApplicationJSON).
			Do().
			Into(&playlistItemsResp)

		if err != nil {
			break
		}

		playlist.Tracks = append(playlist.Tracks, parseTracksResponse(playlistItemsResp)...)
		if playlistItemsResp.Next == "" {
			break
		}

		response.Tracks = playlistItemsResp
	}

	return playlist, nil
}

// CreatePlaylist uses our internal playlist object to create a playlist on Spotify.
func (s *Spotify) CreatePlaylist(playlist utils.Playlist, accessToken string) (string, error) {
	var response spotifyAPICreatePlaylistResponse
	var trackURIs []string
	var wg sync.WaitGroup

	err := s.RequestClient.
		Post(s.Config.BaseAPIURL + "/users/" + s.Config.UserID + "/playlists").
		SetBearerAuthToken(accessToken).
		SetBodyJsonMarshal(map[string]any{
			"name":        playlist.Title,
			"description": playlist.Description,
			"public":      true,
		}).
		Do().
		Into(&response)

	if err != nil {
		return "", err
	}

	for _, entry := range playlist.Tracks {
		trackURIs = append(trackURIs, trackIDToURI(entry))
	}

	// https://github.com/golang/go/wiki/SliceTricks#batching-with-minimal-allocation
	requestsPayloads := make([][]string, 0, (len(trackURIs)+maximumNumOfTracksPerRequest-1)/maximumNumOfTracksPerRequest)
	for maximumNumOfTracksPerRequest < len(trackURIs) {
		trackURIs, requestsPayloads = trackURIs[maximumNumOfTracksPerRequest:], append(requestsPayloads, trackURIs[0:maximumNumOfTracksPerRequest:maximumNumOfTracksPerRequest])
	}
	requestsPayloads = append(requestsPayloads, trackURIs)
	for _, payload := range requestsPayloads {
		wg.Add(1)

		go func(entry []string) {
			defer wg.Done()
			err := s.RequestClient.
				Post(s.Config.BaseAPIURL + "/playlists/" + response.ID + "/tracks").
				SetBearerAuthToken(accessToken).
				SetBodyJsonMarshal(map[string]any{
					"uris": entry,
				}).
				Do()

			if err != nil {
				return
			}
		}(payload)
	}
	wg.Wait()

	return basePlaylistURL + response.ID, nil
}

func (s *Spotify) RefreshAccessToken(payload utils.OauthCredentials) (utils.OauthCredentials, error) {
	var response spotifyAPIRefreshTokenResponse

	err := s.RequestClient.
		Post(s.Config.AuthenticationURL).
		SetFormData(map[string]string{"grant_type": "refresh_token", "refresh_token": payload.RefreshToken}).
		SetBasicAuth(s.Config.ClientID, s.Config.ClientSecret).
		SetContentType("application/x-www-form-urlencoded").
		Do().
		Into(&response)

	if err != nil {
		return utils.OauthCredentials{}, err
	}

	return utils.OauthCredentials{AccessToken: response.AccessToken, ExpiresAt: response.ExpiresIn}, nil
}

func (s *Spotify) GetAuthorizationCode(code string) (utils.OauthCredentials, error) {
	var response spotifyAPIBearerCredentialsResponse
	err := s.RequestClient.
		Post(s.Config.AuthenticationURL).
		SetFormData(map[string]string{"grant_type": "authorization_code", "code": code, "redirect_uri": s.Config.AuthenticationRedirectURL}).
		SetBasicAuth(s.Config.ClientID, s.Config.ClientSecret).
		SetContentType("application/x-www-form-urlencoded").
		Do().
		Into(&response)

	if err != nil {
		return utils.OauthCredentials{}, err
	}

	return utils.OauthCredentials{AccessToken: response.AccessToken, RefreshToken: response.RefreshToken, ExpiresAt: response.ExpiresAt}, nil
}

// RequiresAccessToken specifies if the streaming requires Oauth.
func (*Spotify) RequiresAccessToken() bool {
	return true
}

func (s *Spotify) LookupTrack(track utils.Track) (utils.Track, error) {
	var foundTrack utils.Track

	clientAuthToken, err := s.getClientAuthenticationCredentials()
	if err != nil {
		return foundTrack, fmt.Errorf("spotify: %s", err.Error())
	}

	var response spotifyAPISearchResponse
	err = s.RequestClient.
		Get(s.Config.BaseAPIURL + "/search").
		SetBearerAuthToken(clientAuthToken).
		SetContentType(utils.ApplicationJSON).
		SetQueryParams(map[string]string{
			"q":     trackToSearchQuery(track),
			"type":  "track",
			"limit": "5",
		}).
		Do().
		Into(&response)

	if err != nil {
		return foundTrack, fmt.Errorf("spotify: %s", err.Error())
	}
	if len(response.Tracks.Items) == 0 {
		return foundTrack, fmt.Errorf("spotify: no track found that matches %s", track.Title)
	}

	foundTrack = parseSearchResponse(response)[0]
	return foundTrack, nil
}

// getClientAuthenticationCredentials fetches the client credentials needed for Spotify authentication.
func (s *Spotify) getClientAuthenticationCredentials() (string, error) {
	if token, ok := utils.GlobalCache.Get("spotifyClientAuthToken"); ok {
		if val, ok := token.(string); ok {
			return val, nil
		}

		return "", errors.New("client authentication credentials corrupted in cache")
	}

	var response spotifyAPIClientCredentialsResponse
	err := s.RequestClient.
		Post(s.Config.AuthenticationURL).
		SetBasicAuth(s.Config.ClientID, s.Config.ClientSecret).
		SetFormData(map[string]string{"grant_type": "client_credentials"}).
		Do().
		Into(&response)

	if err != nil {
		utils.GlobalCache.Set("spotifyClientAuthToken", response.AccessToken, time.Second*time.Duration(response.ExpiresIn))
	}

	return response.AccessToken, err
}
