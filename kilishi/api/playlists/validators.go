package playlists

import (
	"fmt"

	"github.com/prettyirrelevant/kilishi/aggregator"
)

// ConvertPlaylistRequest is a struct that represents the request body for the ConvertPlaylistController function.
type ConvertPlaylistRequest struct {
	Source      aggregator.MusicStreamingPlatform `json:"source"`
	Destination aggregator.MusicStreamingPlatform `json:"destination"`
	PlaylistURL string                            `json:"playlistURL"`
}

func (c *ConvertPlaylistRequest) Validate() (bool, []string) {
	var foundErrors []string

	if ok := aggregator.AllMusicStreamingPlatforms[c.Source]; !ok {
		foundErrors = append(foundErrors, fmt.Sprintf("%s is not a supported streaming platform", c.Source))
	}
	if ok := aggregator.AllMusicStreamingPlatforms[c.Destination]; !ok {
		foundErrors = append(foundErrors, fmt.Sprintf("%s is not a supported streaming platform", c.Source))
	}

	if c.Source == c.Destination {
		foundErrors = append(foundErrors, "source cannot be the same as destination")
	}

	if len(foundErrors) > 0 {
		return false, foundErrors
	}
	return true, foundErrors
}
