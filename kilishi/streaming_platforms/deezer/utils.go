package deezer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req/v3"

	"github.com/prettyirrelevant/kilishi/utils"
)

// parsePlaylistURL validates a Deezer playlist URL and returns the playlist ID.
func parsePlaylistURL(playlistURL string) (string, error) {
	re := regexp.MustCompile(`^https://www\.deezer\.com/../playlist/(\d+)$`)
	matches := re.FindStringSubmatch(playlistURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("deezer: playlist url is invalid. check that it follows the format https://www.deezer.com/<country_code>/playlist/<id>")
	}

	return matches[1], nil
}

// trackToSearchQuery transforms our internal track object into a Deezer search query.
func trackToSearchQuery(track utils.Track) string {
	q := fmt.Sprintf("track:%q", track.Title)
	for _, artistName := range track.Artists {
		q += fmt.Sprintf(" artist:%q", artistName)
		break
	}

	return q
}

// parseGetPlaylistResponse transforms the playlist object returned from Deezer API into our internal object.
func parseGetPlaylistResponse(response *deezerAPIGetPlaylistResponse) utils.Playlist {
	tracks := parsePlaylistTracksResponse(response)
	return utils.Playlist{
		ID:          strconv.Itoa(response.ID),
		Title:       response.Title,
		Description: response.Description,
		Tracks:      tracks,
	}
}
func parsePlaylistTracksResponse(data *deezerAPIGetPlaylistResponse) []utils.Track {
	var tracks []utils.Track
	for _, track := range data.Tracks.Data {
		tracks = append(tracks, utils.Track{
			ID:      strconv.Itoa(track.ID),
			Title:   utils.CleanTrackTitle(track.Title),
			Artists: []string{track.Artist.Name},
		})
	}

	return tracks
}

func parseSearchTracksResponse(data deezerAPISearchTrackResponse) []utils.Track {
	var tracks []utils.Track
	for _, track := range data.Data {
		tracks = append(tracks, utils.Track{
			ID:      strconv.Itoa(track.ID),
			Title:   utils.CleanTrackTitle(track.Title),
			Artists: []string{track.Artist.Name},
		})
	}

	return tracks
}

func setupRequestClient(reqClient *req.Client) *req.Client {
	return reqClient.
		EnableDumpEachRequest().
		SetCommonErrorResult(&deezerAPIError{}).
		SetResultStateCheckFunc(func(resp *req.Response) req.ResultState {
			if resp.StatusCode >= 400 {
				return req.ErrorState
			}
			if resp.StatusCode == 200 && strings.HasPrefix(resp.String(), "{\"error\":{\"type\":") {
				return req.ErrorState
			}
			return req.SuccessState
		}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			// There is an underlying error, e.g. network error or unmarshal error.
			if resp.Err != nil {
				return nil
			}
			// Server returns an error message, convert it to human-readable Go error.
			if apiErr, ok := resp.ErrorResult().(*deezerAPIError); ok {
				resp.Err = apiErr
				return nil
			}
			// Edge case: neither an error state response nor a success state response,
			// dump content to help troubleshoot.
			if !resp.IsSuccessState() {
				resp.Err = fmt.Errorf("bad status: %s\nraw content:\n%s", resp.Status, resp.Dump())
				return nil
			}
			return nil
		}).
		SetCommonRetryCount(2).
		AddCommonRetryCondition(func(resp *req.Response, err error) bool {
			return strings.Contains(resp.String(), "Quota limit exceeded") && resp.StatusCode == 200
		}).
		SetCommonRetryBackoffInterval(2*time.Second, 5*time.Second)
}
