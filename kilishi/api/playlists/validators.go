package playlists

import (
	"github.com/go-playground/validator/v10"
	"github.com/prettyirrelevant/kilishi/pkg/aggregator"
)

// create a new validator instance
var validate = validator.New()

// ConvertPlaylistRequest is a struct that represents the request body for the ConvertPlaylistController function.
type ConvertPlaylistRequest struct {
	Source      aggregator.MusicStreamingPlatform `json:"source"`
	Destination aggregator.MusicStreamingPlatform `json:"destination"`
	PlaylistURL string                            `json:"playlistURL"`
}
