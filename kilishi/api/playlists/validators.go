package playlists

import (
	"errors"
	"fmt"
	"strings"

	"github.com/prettyirrelevant/kilishi/streaming_platforms/aggregator"
	"github.com/prettyirrelevant/kilishi/utils"
)

type GetPlaylistRequest struct {
	Platform    aggregator.MusicStreamingPlatform `query:"platform"`
	PlaylistURL string                            `query:"playlist_url"`
}

func (g *GetPlaylistRequest) Validate() (bool, []string) {
	var foundErrors []string

	err := validateStreamingPlatform(g.Platform, fmt.Sprintf("%s is not a supported streaming platform", g.Platform))
	if err != nil {
		foundErrors = append(foundErrors, err.Error())
	}
	err = validateString(g.PlaylistURL, "`playlist_url` is required.")
	if err != nil {
		foundErrors = append(foundErrors, err.Error())
	}
	if len(foundErrors) > 0 {
		return false, foundErrors
	}

	return true, foundErrors
}

type FindTrackRequest struct {
	Platform aggregator.MusicStreamingPlatform `query:"platform"`
	Title    string                            `query:"title"`
	Artists  []string
}

func (f *FindTrackRequest) Validate() (bool, []string) {
	var foundErrors []string

	err := validateStreamingPlatform(f.Platform, fmt.Sprintf("%s is not a supported streaming platform", f.Platform))
	if err != nil {
		foundErrors = append(foundErrors, err.Error())
	}
	err = validateString(f.Title, "title is required")
	if err != nil {
		foundErrors = append(foundErrors, err.Error())
	}
	err = validateStringSlice(f.Artists, "artists is required")
	if err != nil {
		foundErrors = append(foundErrors, err.Error())
	}
	if len(foundErrors) > 0 {
		return false, foundErrors
	}

	return true, foundErrors
}

// ConvertPlaylistRequest is a struct that represents the request body for the ConvertPlaylistController function.
type ConvertPlaylistRequest struct {
	AccessToken string                            `json:"access_token"`
	Platform    aggregator.MusicStreamingPlatform `json:"platform"`
	Playlist    utils.Playlist                    `json:"playlist"`
}

func (c *ConvertPlaylistRequest) Validate() (bool, []string) {
	var foundErrors []string

	err := validateStreamingPlatform(c.Platform, fmt.Sprintf("%s is not a supported streaming platform", c.Platform))
	if err != nil {
		foundErrors = append(foundErrors, err.Error())
	}
	err = validateString(c.Playlist.Title, "`playlist` requires a title")
	if err != nil {
		foundErrors = append(foundErrors, err.Error())
	}
	if len(c.Playlist.Tracks) == 0 {
		foundErrors = append(foundErrors, "`playlist` requires at least one track")
	}

	for _, track := range c.Playlist.Tracks {
		err := validateString(track.ID, "one of the tracks is missing an identifier")
		if err != nil {
			foundErrors = append(foundErrors, err.Error())
			break
		}
		err = validateString(track.Title, "one of the tracks is missing a title")
		if err != nil {
			foundErrors = append(foundErrors, err.Error())
			break
		}
		err = validateStringSlice(track.Artists, "one of the tracks does not contain an artist")
		if err != nil {
			foundErrors = append(foundErrors, err.Error())
			break
		}
	}
	if len(foundErrors) > 0 {
		return false, foundErrors
	}
	return true, foundErrors
}

func validateString(m string, errMsg string) error {
	if strings.TrimSpace(m) == "" {
		return errors.New(errMsg)
	}
	return nil
}

func validateStringSlice(m []string, errMsg string) error {
	if len(m) == 0 {
		return errors.New(errMsg)
	}
	for _, i := range m {
		err := validateString(i, errMsg)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateStreamingPlatform(m aggregator.MusicStreamingPlatform, errMsg string) error {
	if err := validateString(string(m), errMsg); err != nil {
		return err
	}
	if ok := aggregator.AllMusicStreamingPlatforms[m]; !ok {
		return errors.New(errMsg)
	}

	return nil
}
