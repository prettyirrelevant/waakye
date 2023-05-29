package aggregator

import (
	"github.com/imroc/req/v3"

	"github.com/prettyirrelevant/kilishi/config"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/deezer"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/spotify"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/ytmusic"
)

// New creates a new MusicStreamingPlatformsAggregator instance.
func New(configuration *config.Config) *MusicStreamingPlatformsAggregator {
	return &MusicStreamingPlatformsAggregator{
		Config: configuration,
		YTMusic: ytmusic.New(&ytmusic.InitialisationOpts{
			RequestClient:       createRequestClient(configuration),
			BaseAPIURL:          configuration.YTMusicAPIBaseURL,
			AuthenticationToken: configuration.SecretKey,
		}),
		Deezer: deezer.New(&deezer.InitialisationOpts{
			RequestClient:     createRequestClient(configuration),
			AppID:             configuration.DeezerAppID,
			BaseAPIURL:        configuration.DeezerBaseAPIURL,
			ClientSecret:      configuration.DeezerClientSecret,
			AuthenticationURL: configuration.DeezerAuthenticationURL,
		}),
		Spotify: spotify.New(&spotify.InitialisationOpts{
			RequestClient:             createRequestClient(configuration),
			UserID:                    configuration.SpotifyUserID,
			BaseAPIURL:                configuration.SpotifyBaseAPIURL,
			ClientID:                  configuration.SpotifyClientID,
			ClientSecret:              configuration.SpotifyClientSecret,
			AuthenticationURL:         configuration.SpotifyClientAuthURL,
			AuthenticationRedirectURL: configuration.SpotifyAuthRedirectURL,
		}),
	}
}

func createRequestClient(configuration *config.Config) *req.Client {
	client := req.C()
	if configuration.Debug {
		return client.DevMode()
	}
	return client
}

// SupportedPlatforms returns a list of supported music streaming platforms.
func (m *MusicStreamingPlatformsAggregator) SupportedPlatforms() []MusicStreamingPlatform {
	var platforms []MusicStreamingPlatform
	for k := range AllMusicStreamingPlatforms {
		platforms = append(platforms, k)
	}
	return platforms
}

// GetStreamingPlatform retrieves the music streaming platform from the MusicStreamingPlatformsAggregator.
func (m *MusicStreamingPlatformsAggregator) GetStreamingPlatform(platform MusicStreamingPlatform) MusicStreamingPlatformInterface {
	switch platform {
	case Spotify:
		return m.Spotify
	case Deezer:
		return m.Deezer
	case YTMusic:
		return m.YTMusic
	default:
		return nil
	}
}
