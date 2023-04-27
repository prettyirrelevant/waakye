package aggregator

import (
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/prettyirrelevant/waakye/api"
	"github.com/prettyirrelevant/waakye/api/database"
	"github.com/prettyirrelevant/waakye/config"
	"github.com/prettyirrelevant/waakye/pkg/deezer"
	"github.com/prettyirrelevant/waakye/pkg/spotify"
	"github.com/prettyirrelevant/waakye/pkg/ytmusic"
)

func New(db *database.Database, config *config.Config) *MusicStreamingPlatformsAggregator {
	return &MusicStreamingPlatformsAggregator{
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
			BaseAPIURI:                config.SpotifyBaseApiURL,
			ClientID:                  config.SpotifyClientID,
			ClientSecret:              config.SpotifyClientSecret,
			AuthenticationURI:         config.SpotifyClientAuthURL,
			AuthenticationRedirectURI: config.SpotifyAuthRedirectURI,
		}),
	}
}

func (m *MusicStreamingPlatformsAggregator) ConvertPlaylist(source, destination api.MusicStreamingPlatform, playlistURL string) (string, error) {
	// this is checked at the API validation level but it does not hurt to check here also.
	if source == destination {
		return "", fmt.Errorf("api.aggregator: `source` must not be the same as `destination`")
	}

	sourcePlatform, destinationPlatform := m.getStreamingPlatform(source), m.getStreamingPlatform(destination)

	playlist, err := sourcePlatform.GetPlaylist(playlistURL)
	if err != nil {
		return "", err
	}

	var accessToken string

	// if the `toPlatform` requires access token, fetch it from the database.
	if destinationPlatform.RequiresAccessToken() {
		// m.Database
	}

	playlistURL, err = destinationPlatform.CreatePlaylist(playlist, accessToken)
	if err != nil {
		return "", err
	}

	return playlistURL, nil
}

// SupportedPlatforms returns a slice of supported music streaming platforms
func (m *MusicStreamingPlatformsAggregator) SupportedPlatforms() []api.MusicStreamingPlatform {
	return []api.MusicStreamingPlatform{api.Deezer, api.Spotify, api.YTMusic}
}

// getStreamingPlatform maps a `MusicStreamingPlatform` to one of `spotify.Spotify`, `deezer.Deezer` or `ytmusic.YTMusic`
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
