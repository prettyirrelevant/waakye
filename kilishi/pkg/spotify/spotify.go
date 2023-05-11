package spotify

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/prettyirrelevant/kilishi/pkg/utils"
)

// New initialises a `Spotify` object.
func New(opts InitialisationOpts) *Spotify {
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

	playlist := parseGetPlaylistResponse(response)
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
	var tracksFound []utils.Track
	var wg sync.WaitGroup
	for _, entry := range playlist.Tracks {
		wg.Add(1)

		go func(track utils.Track) {
			defer wg.Done()
			s.lookupTrack(track, &tracksFound)
		}(entry)
	}
	wg.Wait()

	var response spotifyAPICreatePlaylistResponse
	err := s.RequestClient.
		Post(s.Config.BaseAPIURL + "/users/" + s.Config.UserID + "/playlists").
		SetBearerAuthToken(accessToken).
		SetContentType(utils.ApplicationJSON).
		SetFormData(map[string]string{
			"name":        playlist.Title,
			"description": playlist.Description,
			"public":      "true",
		}).
		Do().
		Into(&response)

	if err != nil {
		return "", err
	}

	s.populatePlaylistWithTracks(playlist.Tracks, playlist.ID, accessToken)

	return response.ID, nil
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

	fmt.Printf("Got response of %+v", response)
	return utils.OauthCredentials{AccessToken: response.AccessToken, RefreshToken: response.RefreshToken, ExpiresAt: int(response.ExpiresAt)}, nil
}

// RequiresAccessToken specifies if the streaming requires Oauth.
func (*Spotify) RequiresAccessToken() bool {
	return true
}

// populatePlaylistWithTracks adds tracks found on Spotify to a newly created playlist.
//
// More info can be found here https://developer.spotify.com/documentation/web-api/reference/#/operations/add-tracks-to-playlist
func (s *Spotify) populatePlaylistWithTracks(tracks []utils.Track, playlistID, accessToken string) {
	var tracksURL []string
	var maximumNumOfTracksPerRequest = 100

	for _, entry := range tracks {
		tracksURL = append(tracksURL, trackIDToURL(entry))
	}

	// https://github.com/golang/go/wiki/SliceTricks#batching-with-minimal-allocation
	chunks := make([][]string, 0, (len(tracksURL)+maximumNumOfTracksPerRequest-1)/maximumNumOfTracksPerRequest)
	for maximumNumOfTracksPerRequest < len(tracksURL) {
		tracksURL, chunks = tracksURL[maximumNumOfTracksPerRequest:], append(chunks, tracksURL[0:maximumNumOfTracksPerRequest:maximumNumOfTracksPerRequest])
	}

	var wg sync.WaitGroup
	for _, entry := range chunks {
		wg.Add(1)

		go func(chunk []string) {
			defer wg.Done()
			err := s.RequestClient.
				Post(s.Config.BaseAPIURL + "/playlists/" + playlistID + "/tracks").
				SetBearerAuthToken(accessToken).
				SetContentType(utils.ApplicationJSON).
				SetFormData(map[string]string{
					"uris": strings.Join(chunk, ","),
				}).
				Do()

			if err != nil {
				return
			}
		}(entry)
	}
	wg.Wait()
}

func (s *Spotify) lookupTrack(track utils.Track, tracksFound *[]utils.Track) {
	clientAuthToken, err := s.getClientAuthenticationCredentials()
	if err != nil {
		return
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
		return
	}

	if len(response.Tracks.Items) == 0 {
		return
	}
	*tracksFound = append(*tracksFound, parseSearchResponse(response)[0])
}

// getClientAuthenticationCredentials fetches the client credentials needed for Spotify authentication.
func (s *Spotify) getClientAuthenticationCredentials() (string, error) {
	if token, ok := utils.GlobalCache.Get("spotifyClientAuthToken"); ok {
		return token.(string), nil
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
