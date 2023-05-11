package aggregator

import (
	"fmt"

	"github.com/imroc/req/v3"

	"github.com/prettyirrelevant/kilishi/config"
	"github.com/prettyirrelevant/kilishi/pkg/deezer"
	"github.com/prettyirrelevant/kilishi/pkg/spotify"
	"github.com/prettyirrelevant/kilishi/pkg/ytmusic"
)

// New creates a new MusicStreamingPlatformsAggregator instance.
func New(config *config.Config) *MusicStreamingPlatformsAggregator {
	return &MusicStreamingPlatformsAggregator{
		Config:  config,
		YTMusic: ytmusic.New(),
		Deezer: deezer.New(deezer.InitialisationOpts{
			RequestClient:     createRequestClient(config),
			AppID:             config.DeezerAppID,
			BaseAPIURL:        config.DeezerBaseApiURL,
			ClientSecret:      config.DeezerClientSecret,
			AuthenticationURL: config.DeezerAuthenticationURL,
		}),
		Spotify: spotify.New(spotify.InitialisationOpts{
			RequestClient:             createRequestClient(config),
			UserID:                    config.SpotifyUserID,
			BaseAPIURL:                config.SpotifyBaseApiURL,
			ClientID:                  config.SpotifyClientID,
			ClientSecret:              config.SpotifyClientSecret,
			AuthenticationURL:         config.SpotifyClientAuthURL,
			AuthenticationRedirectURL: config.SpotifyAuthRedirectURL,
		}),
	}
}

func createRequestClient(config *config.Config) *req.Client {
	client := req.C()
	if config.Debug {
		return client.DevMode().Clone()
	}
	return client
}

// ConvertPlaylist converts a playlist from one music streaming platform to another.
func (m *MusicStreamingPlatformsAggregator) ConvertPlaylist(source, destination MusicStreamingPlatform, playlistURL, accessToken string) (string, error) {
	if source == destination {
		return "", fmt.Errorf("aggregator: `source` must not be the same as `destination`")
	}

	sourcePlatform, destinationPlatform := m.getStreamingPlatform(source), m.getStreamingPlatform(destination)
	playlist, err := sourcePlatform.GetPlaylist(playlistURL)
	if err != nil {
		return "", err
	}

	playlistURL, err = destinationPlatform.CreatePlaylist(playlist, accessToken)
	if err != nil {
		return "", err
	}

	return playlistURL, nil
}

// SupportedPlatforms returns a list of supported music streaming platforms.
func (m *MusicStreamingPlatformsAggregator) SupportedPlatforms() []MusicStreamingPlatform {
	var platforms []MusicStreamingPlatform
	for k := range AllMusicStreamingPlatforms {
		platforms = append(platforms, k)
	}
	return platforms
}

// getStreamingPlatform retrieves the music streaming platform from the MusicStreamingPlatformsAggregator.
func (m *MusicStreamingPlatformsAggregator) getStreamingPlatform(platform MusicStreamingPlatform) MusicStreamingPlatformInterface {
	switch platform {
	case Spotify:
		return m.Spotify
	case Deezer:
		return m.Deezer
	case YTMusic:
		return m.YTMusic
	case AppleMusic:
		return m.AppleMusic
	default:
		return nil
	}
}
