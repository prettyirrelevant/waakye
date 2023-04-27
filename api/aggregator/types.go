package aggregator

import (
	"github.com/prettyirrelevant/waakye/api/database"
	"github.com/prettyirrelevant/waakye/pkg/deezer"
	"github.com/prettyirrelevant/waakye/pkg/spotify"
	"github.com/prettyirrelevant/waakye/pkg/utils/types"
	"github.com/prettyirrelevant/waakye/pkg/ytmusic"
)

type MusicStreamingPlatformsAggregator struct {
	Database *database.Database
	Spotify  *spotify.Spotify
	Deezer   *deezer.Deezer
	YTMusic  *ytmusic.YTMusic
}

type MusicStreamingPlatformInterface interface {
	CreatePlaylist(playlist types.Playlist, accessToken string) (string, error)
	GetPlaylist(playlistURL string) (types.Playlist, error)
	GetAuthorizationCode(code string) (types.OauthCredentials, error)
	RequiresAccessToken() bool
}
