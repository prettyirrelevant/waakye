package ytmusic

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/prettyirrelevant/waakye/pkg/utils/queries"
	"github.com/prettyirrelevant/waakye/pkg/utils/types"
	"github.com/prettyirrelevant/ytmusicapi"
)

// New initialises a `YTMusic` object.
func New() *YTMusic {
	ytmusicapi.Setup()
	return &YTMusic{}
}

// GetPlaylist returns information about a playlist.
func (y *YTMusic) GetPlaylist(playlistURI string) (types.Playlist, error) {
	playlistID, err := y.parsePlaylistURI(playlistURI)
	if err != nil {
		return types.Playlist{}, err
	}

	playlist, err := ytmusicapi.GetPlaylist(playlistID, 0)
	if err != nil {
		return types.Playlist{}, fmt.Errorf("ytmusic: %s", err.Error())
	}

	return y.deserializePlaylist(playlist), nil
}

// CreatePlaylist creates a new playlist using the information provided.
func (y *YTMusic) CreatePlaylist(playlist types.Playlist) (string, error) {
	var foundTracks []types.Track
	var wg sync.WaitGroup

	for _, trackEntry := range playlist.Tracks {
		wg.Add(1)
		go func(payload types.Track) {
			defer wg.Done()
			y.lookupTrack(payload, &foundTracks)
		}(trackEntry)
	}
	wg.Wait()

	var trackIDs []string
	for _, entry := range foundTracks {
		trackIDs = append(trackIDs, entry.ID)
	}

	playlistID, err := ytmusicapi.CreatePlaylist(playlist.Title, playlist.Description, ytmusicapi.PUBLIC, "", trackIDs)
	if err != nil {
		return "", fmt.Errorf("ytmusic: %s", err.Error())
	}
	return playlistID, nil
}

// lookupTrack searches for track on YTMusic and appends the top result to a slice.
func (*YTMusic) lookupTrack(track types.Track, foundTracks *[]types.Track) {
	searchResults, err := ytmusicapi.SearchTracks(queries.TrackToSearchQuery(track), ytmusicapi.SongsFilter, ytmusicapi.NoScope, 5, false)
	if err == nil && len(searchResults) > 0 {
		*foundTracks = append(*foundTracks, types.Track{ID: searchResults[0].VideoID, Title: searchResults[0].Title, Artists: searchResults[0].Artistes})
	}
}

// deserializePlaylist transforms the playlist object returned from `ytmusicapi` into our internal object.
func (y *YTMusic) deserializePlaylist(playlist ytmusicapi.Playlist) types.Playlist {
	var tracks []types.Track

	for _, entry := range playlist.Tracks {
		tracks = append(tracks, types.Track{ID: entry.ID, Title: y.cleanTrackTitle(entry.Title), Artists: entry.Artistes})
	}

	return types.Playlist{
		ID:          playlist.ID,
		Title:       playlist.Title,
		Description: playlist.Description,
		Tracks:      tracks,
	}
}

// cleanTrackTitle removes noise from YTMusic track title such as [Official Audio], (Official Video), etc.
func (*YTMusic) cleanTrackTitle(title string) string {
	// An example usage of the regex can be found here regexr.com/7cf46
	re := regexp.MustCompile(`\s*\(_*\s*(Official Visualizer|Official Video|Official Audio|Official Music Video|Live)\s*_*\)\s*`)
	return re.ReplaceAllString(title, "")
}

// parsePlaylistURI validates a YTMusic playlist URI and returns the playlist ID.
func (*YTMusic) parsePlaylistURI(playlistURI string) (string, error) {
	re := regexp.MustCompile(`^https:\/\/music\.youtube\.com\/playlist\?list=([a-zA-Z0-9-_]+)$`)

	matches := re.FindStringSubmatch(playlistURI)
	if len(matches) < 1 {
		return "", fmt.Errorf("ytmusic: playlist url is invalid. check that it follows the format https://music.youtube.com/playlist?list=")
	}

	return matches[1], nil
}
