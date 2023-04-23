package deezer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/prettyirrelevant/waakye/pkg/utils/types"
)

// parsePlaylistURI validates a Deezer playlist URI and returns the playlist ID.
func parsePlaylistURI(playlistURI string) (string, error) {
	re := regexp.MustCompile(`^https:\/\/www\.deezer\.com\/..\/playlist\/(\d+)$`)
	matches := re.FindStringSubmatch(playlistURI)
	if len(matches) < 2 {
		return "", fmt.Errorf("deezer: playlist url is invalid. check that it follows the format https://www.deezer.com/<country_code>/playlist/<id>")
	}

	return matches[1], nil
}

// trackToSearchQuery transforms our internal track object into a Deezer search query.
func trackToSearchQuery(track types.Track) string {
	query := "track:" + track.Title
	for _, artist := range track.Artists {
		query += " artist:" + artist
	}

	return query
}

// parseGetPlaylistResponse transforms the playlist object returned from Deezer API into our internal object.
func parseGetPlaylistResponse(response deezerAPIGetPlaylistResponse) types.Playlist {
	tracks := parseTracksResponse(response.Tracks)
	return types.Playlist{
		ID:          strconv.Itoa(response.ID),
		Title:       response.Title,
		Description: response.Description,
		Tracks:      tracks,
	}
}

func parseTracksResponse(response deezerAPITracksDataResponse) []types.Track {
	var tracks []types.Track
	for _, track := range response.Data {
		tracks = append(tracks, types.Track{
			ID:      strconv.Itoa(track.ID),
			Title:   track.Title,
			Artists: []string{track.Artist.Name},
		})
	}

	return tracks
}

func setupRequestClient(reqClient *req.Client) *req.Client {
	return reqClient.
		SetCommonErrorResult(&deezerAPIErrorResponse{}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			// There is an underlying error, e.g. network error or unmarshal error.
			if resp.Err != nil {
				return nil
			}
			// Server returns an error message, convert it to human-readable Go error.
			if apiErr, ok := resp.ErrorResult().(*deezerAPIErrorResponse); ok {
				resp.Err = apiErr
				return nil
			}
			// Edge case: neither an error state response nor a success state response,
			// dump content to help troubleshoot.
			if !resp.IsSuccessState() {
				return fmt.Errorf("bad response, raw dump:\n%s", resp.Dump())
			}
			return nil
		}).
		SetCommonRetryCount(3).
		AddCommonRetryCondition(func(resp *req.Response, err error) bool {
			return strings.Contains(resp.String(), "Quota limit exceeded") && resp.StatusCode == 200
		}).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second)
}
