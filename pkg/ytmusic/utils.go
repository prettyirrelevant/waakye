package ytmusic

import (
	"fmt"
	"regexp"

	"github.com/prettyirrelevant/waakye/pkg/utils/queries"
	"github.com/prettyirrelevant/waakye/pkg/utils/types"
	"github.com/prettyirrelevant/ytmusicapi"
)

// lookupTrack searches for track on YTMusic and appends the top result to a slice.
func lookupTrack(track types.Track, foundTracks *[]types.Track) {
	searchResults, err := ytmusicapi.SearchTracks(queries.TrackToSearchQuery(track), ytmusicapi.SongsFilter, ytmusicapi.NoScope, 5, false)
	if err == nil && len(searchResults) > 0 {
		*foundTracks = append(*foundTracks, types.Track{ID: searchResults[0].VideoID, Title: searchResults[0].Title, Artists: searchResults[0].Artistes})
	}
}

// cleanTrackTitle removes noise from YTMusic track title such as [Official Audio], (Official Video), etc.
func cleanTrackTitle(title string) string {
	// An example usage of the regex can be found here regexr.com/7cf46
	re := regexp.MustCompile(`\s*\(_*\s*(Official Visualizer|Official Video|Official Audio|Official Music Video|Live)\s*_*\)\s*`)
	return re.ReplaceAllString(title, "")
}

// parsePlaylistURI validates a YTMusic playlist URI and returns the playlist ID.
func parsePlaylistURI(playlistURI string) (string, error) {
	re := regexp.MustCompile(`^https:\/\/music\.youtube\.com\/playlist\?list=([a-zA-Z0-9-_]+)$`)

	matches := re.FindStringSubmatch(playlistURI)
	if len(matches) < 2 {
		return "", fmt.Errorf("ytmusic: playlist url is invalid. check that it follows the format https://music.youtube.com/playlist?list=")
	}

	return matches[1], nil
}
