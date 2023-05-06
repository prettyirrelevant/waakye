package ytmusic

import (
	"fmt"
	"regexp"

	"github.com/prettyirrelevant/kilishi/pkg/utils"
	"github.com/prettyirrelevant/kilishi/pkg/utils/types"
	"github.com/prettyirrelevant/ytmusicapi"
)

// lookupTrack searches for track on YTMusic and appends the top result to a slice.
func lookupTrack(track types.Track, foundTracks *[]types.Track) {
	searchResults, err := ytmusicapi.SearchTracks(trackToSearchQuery(track), ytmusicapi.SongsFilter, ytmusicapi.NoScope, 5, false)
	if err == nil && len(searchResults) > 0 {
		*foundTracks = append(*foundTracks, types.Track{ID: searchResults[0].VideoID, Title: searchResults[0].Title, Artists: searchResults[0].Artistes})
	}
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

// trackToSearchQuery takes a track and transforms it into a search query.
func trackToSearchQuery(track types.Track) string {
	searchQuery := track.Title + " by"
	for index, artiste := range track.Artists {
		if len(track.Artists) == 1 {
			searchQuery += " " + artiste
		} else if len(track.Artists) > 1 && len(track.Artists)-1 == index {
			searchQuery += " and " + artiste
		} else {
			searchQuery += " " + artiste
			if index < len(track.Artists)-2 {
				searchQuery += ","
			}
		}
	}
	return searchQuery
}

// parseGetPlaylistResponse transforms the playlist object returned from `ytmusicapi` into our internal object.
func parseGetPlaylistResponse(playlist ytmusicapi.Playlist) types.Playlist {
	var tracks []types.Track

	for _, entry := range playlist.Tracks {
		tracks = append(tracks, types.Track{ID: entry.ID, Title: utils.CleanTrackTitle(entry.Title), Artists: entry.Artistes})
	}

	return types.Playlist{
		ID:          playlist.ID,
		Title:       playlist.Title,
		Description: playlist.Description,
		Tracks:      tracks,
	}
}
