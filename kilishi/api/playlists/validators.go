package playlists

import (
	"fmt"

	"github.com/prettyirrelevant/kilishi/pkg/aggregator"
)

// ConvertPlaylistRequest is a struct that represents the request body for the ConvertPlaylistController function.
type ConvertPlaylistRequest struct {
	Source      aggregator.MusicStreamingPlatform `json:"source"`
	Destination aggregator.MusicStreamingPlatform `json:"destination"`
	PlaylistURL string                            `json:"playlistURL"`
}

func (c *ConvertPlaylistRequest) Validate() (bool, []error) {
	var foundErrors []error

	if ok := aggregator.AllMusicStreamingPlatforms[c.Source]; !ok {
		foundErrors = append(foundErrors, fmt.Errorf("`source` is not one of the supported streaming platforms"))
	}
	if ok := aggregator.AllMusicStreamingPlatforms[c.Destination]; !ok {
		foundErrors = append(foundErrors, fmt.Errorf("`destination` is not one of the supported streaming platforms"))
	}

	if len(foundErrors) > 0 {
		return false, foundErrors
	}
	return true, foundErrors
}
