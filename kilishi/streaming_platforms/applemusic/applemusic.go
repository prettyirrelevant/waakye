package applemusic

import (
	"github.com/prettyirrelevant/kilishi/utils"
)

func New(opts InitialisationOpts) *AppleMusic {
	return &AppleMusic{}
}

func (a *AppleMusic) GetPlaylist(playlistURL string) (utils.Playlist, error) {
	return utils.Playlist{}, nil
}

func (a *AppleMusic) CreatePlaylist(playlist utils.Playlist, accessToken string) (string, error) {
	return "", nil
}

func (a *AppleMusic) RequiresAccessToken() bool {
	return true
}

func (a *AppleMusic) GetAuthorizationCode(code string) (utils.OauthCredentials, error) {
	return utils.OauthCredentials{}, nil
}
