package aggregator

import (
	"github.com/prettyirrelevant/kilishi/config"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/applemusic"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/deezer"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/spotify"
	"github.com/prettyirrelevant/kilishi/streaming_platforms/ytmusic"
	"github.com/prettyirrelevant/kilishi/utils"
)

const (
	Spotify    MusicStreamingPlatform = "spotify"
	Deezer     MusicStreamingPlatform = "deezer"
	YTMusic    MusicStreamingPlatform = "ytmusic"
	AppleMusic MusicStreamingPlatform = "apple music"
)

var AllMusicStreamingPlatforms = map[MusicStreamingPlatform]bool{Spotify: true, Deezer: true, YTMusic: true, AppleMusic: true}

// MusicStreamingPlatformsAggregator is a struct that represents an aggregator of different music streaming platforms.
type MusicStreamingPlatformsAggregator struct {
	Config     *config.Config
	Spotify    *spotify.Spotify
	Deezer     *deezer.Deezer
	YTMusic    *ytmusic.YTMusic
	AppleMusic *applemusic.AppleMusic
}

type MusicStreamingPlatformInterface interface {
	// CreatePlaylist creates a new playlist on the platform.
	// It takes a utils.Playlist object and an access token string
	// and returns the URL of the newly created playlist and an error, if any.
	CreatePlaylist(playlist utils.Playlist, accessToken string) (string, error)

	// GetPlaylist returns a utils.Playlist object for a given playlist URL.
	// It takes a playlist URL string and returns the corresponding playlist object and an error, if any.
	GetPlaylist(playlistURL string) (utils.Playlist, error)

	// GetAuthorizationCode returns an oauth credentials object for the given authorization code.
	// It takes an authorization code string and returns an oauth credentials object and an error, if any.
	GetAuthorizationCode(code string) (utils.OauthCredentials, error)

	// RequiresAccessToken returns a boolean indicating whether the platform requires an access token for API calls.
	RequiresAccessToken() bool
}

type MusicStreamingPlatform string
