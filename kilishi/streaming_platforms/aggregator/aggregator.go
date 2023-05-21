package aggregator

import (
	"github.com/imroc/req/v3"

	"github.com/prettyirrelevant/kilishi/config"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/deezer"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/spotify"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/ytmusic"
)

// New creates a new MusicStreamingPlatformsAggregator instance.
func New(config *config.Config) *MusicStreamingPlatformsAggregator {
	return &MusicStreamingPlatformsAggregator{
		Config: config,
		YTMusic: ytmusic.New(ytmusic.InitialisationOpts{
			RequestClient:       createRequestClient(config),
			BaseAPIURL:          config.YTMusicApiBaseUrl,
			AuthenticationToken: config.SecretKey,
		}),
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
