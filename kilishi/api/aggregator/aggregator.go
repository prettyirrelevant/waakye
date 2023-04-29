package aggregator

import (
	"fmt"

	"github.com/imroc/req/v3"

	"github.com/prettyirrelevant/kilishi/api"
	"github.com/prettyirrelevant/kilishi/api/database"
	"github.com/prettyirrelevant/kilishi/config"
	"github.com/prettyirrelevant/kilishi/pkg/deezer"
	"github.com/prettyirrelevant/kilishi/pkg/spotify"
	"github.com/prettyirrelevant/kilishi/pkg/utils/types"
	"github.com/prettyirrelevant/kilishi/pkg/ytmusic"
)

// New creates a new MusicStreamingPlatformsAggregator instance.
func New(db *database.Database, config *config.Config) *MusicStreamingPlatformsAggregator {
	return &MusicStreamingPlatformsAggregator{
		Config:   config,
		Database: db,
		YTMusic:  ytmusic.New(),
		Deezer: deezer.New(&deezer.InitialisationOpts{
			RequestClient:     req.C(),
			AppID:             config.DeezerAppID,
			BaseAPIURI:        config.DeezerBaseApiURL,
			ClientSecret:      config.DeezerClientSecret,
			AuthenticationURI: config.DeezerAuthenticationURI,
		}),
		Spotify: spotify.New(&spotify.InitialisationOpts{
			RequestClient:             req.C(),
			UserID:                    config.SpotifyUserID,
			BaseAPIURI:                config.SpotifyBaseApiURL,
			ClientID:                  config.SpotifyClientID,
			ClientSecret:              config.SpotifyClientSecret,
			AuthenticationURI:         config.SpotifyClientAuthURL,
			AuthenticationRedirectURI: config.SpotifyAuthRedirectURI,
		}),
	}
}

// ConvertPlaylist converts a playlist from one music streaming platform to another.
func (m *MusicStreamingPlatformsAggregator) ConvertPlaylist(source, destination api.MusicStreamingPlatform, playlistURL string) (string, error) {
	if source == destination {
		return "", fmt.Errorf("api.aggregator: `source` must not be the same as `destination`")
	}

	sourcePlatform, destinationPlatform := m.getStreamingPlatform(source), m.getStreamingPlatform(destination)

	playlist, err := sourcePlatform.GetPlaylist(playlistURL)
	if err != nil {
		return "", err
	}

	var accessToken string
	if destinationPlatform.RequiresAccessToken() {
		dbCredentials, err := m.Database.GetOauthCredentials(destination)
		if err != nil {
			return "", err
		}

		credentials, err := types.OauthCredentialsFromDB(dbCredentials.Credentials)
		if err != nil {
			return "", err
		}

		accessToken = credentials.AccessToken
	}

	playlistURL, err = destinationPlatform.CreatePlaylist(playlist, accessToken)
	if err != nil {
		return "", err
	}

	return playlistURL, nil
}

// SupportedPlatforms returns a list of supported music streaming platforms.
func (m *MusicStreamingPlatformsAggregator) SupportedPlatforms() []api.MusicStreamingPlatform {
	return []api.MusicStreamingPlatform{api.Deezer, api.Spotify, api.YTMusic}
}

// getStreamingPlatform retrieves the music streaming platform from the MusicStreamingPlatformsAggregator.
func (m *MusicStreamingPlatformsAggregator) getStreamingPlatform(platform api.MusicStreamingPlatform) MusicStreamingPlatformInterface {
	switch platform {
	case api.Spotify:
		return m.Spotify
	case api.Deezer:
		return m.Deezer
	case api.YTMusic:
		return m.YTMusic
	default:
		return nil
	}
}
