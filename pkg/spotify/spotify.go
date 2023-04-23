package spotify

import (
	"strings"
	"sync"
	"time"

	"github.com/prettyirrelevant/waakye/pkg/utils/cache"
	"github.com/prettyirrelevant/waakye/pkg/utils/types"
)

// New initialises a `Spotify` object.
func New(opts *InitialisationOpts) *Spotify {
	return &Spotify{
		RequestClient: setupRequestClient(opts.RequestClient),
		Config: &Config{
			UserID:                    opts.UserID,
			ClientID:                  opts.ClientID,
			BaseAPIURI:                opts.BaseAPIURI,
			ClientSecret:              opts.ClientSecret,
			AuthenticationURI:         opts.AuthenticationURI,
			AuthenticationRedirectURI: opts.AuthenticationRedirectURI,
		},
	}
}

func (s *Spotify) GetPlaylist(playlistURI string) (types.Playlist, error) {
	playlistID, err := parsePlaylistURI(playlistURI)
	if err != nil {
		return types.Playlist{}, err
	}

	clientAuthToken, err := s.getClientAuthenticationCredentials()
	if err != nil {
		return types.Playlist{}, err
	}

	var response spotifyAPIGetPlaylistResponse
	err = s.RequestClient.
		Get(s.Config.BaseAPIURI + "/playlists/" + playlistID).
		SetBearerAuthToken(clientAuthToken).
		SetContentType("application/json").
		Do().
		Into(&response)

	if err != nil {
		return types.Playlist{}, err
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
			SetContentType("application/json").
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
func (s *Spotify) CreatePlaylist(playlist types.Playlist, accessToken string) (string, error) {
	// first, look for the tracks on Spotify
	var tracksFound []types.Track
	var wg sync.WaitGroup
	for _, entry := range playlist.Tracks {
		wg.Add(1)

		go func(track types.Track) {
			defer wg.Done()
			s.lookupTrack(track, &tracksFound)
		}(entry)
	}
	wg.Wait()

	// then, create the playlist.
	var response spotifyAPICreatePlaylistResponse
	err := s.RequestClient.
		Post(s.Config.BaseAPIURI + "/users/" + s.Config.UserID + "/playlists").
		SetContentType("application/json").
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

	// finally, add the tracks in batches
	s.populatePlaylistWithTracks(playlist.Tracks, playlist.ID, accessToken)

	return response.ID, nil
}

func (s *Spotify) GetAuthorizationCode(code string) (spotifyAPIBearerCredentialsResponse, error) {
	var response spotifyAPIBearerCredentialsResponse
	err := s.RequestClient.
		Post(s.Config.AuthenticationURI).
		SetFormData(map[string]string{"grant_type": "authorization_code", "code": code, "redirect_uri": s.Config.AuthenticationRedirectURI}).
		SetContentType("application/json").
		Do().
		Into(&response)

	if err != nil {
		return response, err
	}

	return response, nil
}

// populatePlaylistWithTracks adds tracks found on Spotify to a newly created playlist.
//
// More info can be found here https://developer.spotify.com/documentation/web-api/reference/#/operations/add-tracks-to-playlist
func (s *Spotify) populatePlaylistWithTracks(tracks []types.Track, playlistID, accessToken string) {
	var tracksURI []string
	var maximumNumOfTracksPerRequest = 100

	for _, entry := range tracks {
		tracksURI = append(tracksURI, trackIDToURI(entry))
	}

	// https://github.com/golang/go/wiki/SliceTricks#batching-with-minimal-allocation
	chunks := make([][]string, 0, (len(tracksURI)+maximumNumOfTracksPerRequest-1)/maximumNumOfTracksPerRequest)
	for maximumNumOfTracksPerRequest < len(tracksURI) {
		tracksURI, chunks = tracksURI[maximumNumOfTracksPerRequest:], append(chunks, tracksURI[0:maximumNumOfTracksPerRequest:maximumNumOfTracksPerRequest])
	}

	var wg sync.WaitGroup
	for _, entry := range chunks {
		wg.Add(1)

		go func(chunk []string) {
			defer wg.Done()
			err := s.RequestClient.
				Post(s.Config.BaseAPIURI + "/playlists/" + playlistID + "/tracks").
				SetContentType("application/json").
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

func (s *Spotify) lookupTrack(track types.Track, tracksFound *[]types.Track) {
	clientAuthToken, err := s.getClientAuthenticationCredentials()
	if err != nil {
		return
	}

	var response spotifyAPISearchResponse
	err = s.RequestClient.
		Get(s.Config.BaseAPIURI + "/search").
		SetBearerAuthToken(clientAuthToken).
		SetContentType("application/json").
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
	if token, ok := cache.GlobalCache.Get("spotifyClientAuthToken"); ok {
		return token.(string), nil
	}

	var response spotifyAPIClientCredentialsResponse
	err := s.RequestClient.
		Post(s.Config.AuthenticationURI).
		SetBasicAuth(s.Config.ClientID, s.Config.ClientSecret).
		SetFormData(map[string]string{"grant_type": "client_credentials"}).
		Do().
		Into(&response)

	if err != nil {
		cache.GlobalCache.Set("spotifyClientAuthToken", response.AccessToken, time.Second*time.Duration(response.ExpiresIn))
	}

	return response.AccessToken, err
}
