package spotify

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/imroc/req/v3"

	"github.com/prettyirrelevant/kilishi/utils"
)

// parsePlaylistURL validates a Spotify playlist URL and returns the playlist ID.
func parsePlaylistURL(playlistURL string) (string, error) {
	re := regexp.MustCompile(`^(https:\/\/open\.spotify\.com\/playlist\/)([a-zA-Z0-9]+)(\?si=[a-zA-Z0-9]+)?`)

	matches := re.FindStringSubmatch(playlistURL)
	if len(matches) < 3 {
		return "", fmt.Errorf("spotify: playlist url is invalid. check that it follows the format https://open.spotify.com/playlist/<id>")
	}

	return matches[2], nil
}

// parseGetPlaylistResponse transforms the playlist object returned from Spotify API into our internal object.
func parseGetPlaylistResponse(response *spotifyAPIGetPlaylistResponse) utils.Playlist {
	tracks := parseTracksResponse(response.Tracks)
	return utils.Playlist{
		ID:          response.ID,
		Title:       response.Name,
		Description: response.Description,
		Tracks:      tracks,
	}
}

// parseTracksResponse transforms the playlist items returned from Spotify API into our internal object.
func parseTracksResponse(response spotifyAPITracksResponse) []utils.Track {
	var tracks []utils.Track

	for _, entry := range response.Items {
		artistes := []string{}
		for _, artiste := range entry.Track.Artists {
			artistes = append(artistes, artiste.Name)
		}
		tracks = append(
			tracks,
			utils.Track{
				ID:      entry.Track.ID,
				Title:   utils.CleanTrackTitle(entry.Track.Name),
				Artists: artistes,
			},
		)
	}

	return tracks
}

// parseTracksResponse transforms the playlist items returned from Spotify API into our internal object.
func parseSearchResponse(response spotifyAPISearchResponse) []utils.Track {
	var tracks []utils.Track

	for _, entry := range response.Tracks.Items {
		artistes := []string{}
		for _, artiste := range entry.Artists {
			artistes = append(artistes, artiste.Name)
		}
		tracks = append(tracks, utils.Track{ID: entry.ID, Title: entry.Name, Artists: artistes})
	}

	return tracks
}

// trackIDToURL transforms a Spotify ID into URL.
func trackIDToURL(track utils.Track) string {
	return "spotify:track:" + track.ID
}

// trackToSearchQuery transforms our internal track object into a Spotify search query.
func trackToSearchQuery(track utils.Track) string {
	q := fmt.Sprintf("track:%s", track.Title)
	for _, artistName := range track.Artists {
		q += fmt.Sprintf(" artist:%s", artistName)
		break // search with > 1 artiste fails.
	}

	return q
}

func setupRequestClient(reqClient *req.Client) *req.Client {
	return reqClient.
		SetCommonErrorResult(&spotifyAPIError{}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			// There is an underlying error, e.g. network error or unmarshal error.
			if resp.Err != nil {
				return nil
			}
			// Server returns an error message, convert it to human-readable Go error.
			if apiErr, ok := resp.ErrorResult().(*spotifyAPIError); ok {
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
			return resp.StatusCode == 429
		}).
		SetCommonRetryInterval(func(resp *req.Response, attempt int) time.Duration {
			// https://developer.spotify.com/documentation/web-api/guides/rate-limits/
			if resp.Response == nil {
				return 2 * time.Second
			}

			retryAfterHeader := resp.Header.Get("Retry-After")
			if retryAfter, err := strconv.Atoi(retryAfterHeader); retryAfterHeader != "" && err != nil {
				return time.Duration(retryAfter) * time.Second
			}

			return 2 * time.Second
		})
}
