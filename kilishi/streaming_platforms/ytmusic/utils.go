package ytmusic

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/imroc/req/v3"

	"github.com/prettyirrelevant/kilishi/utils"
)

// trackToSearchQuery takes a track and transforms it into a search query.
func trackToSearchQuery(track utils.Track) string {
	q := track.Title + " by"
	for _, artiste := range track.Artists {
		q += " " + artiste
		break
	}
	return q
}

func cleanTrackArtist(name string) string {
	re := regexp.MustCompile(`(?i)vevo`)
	cleanedName := strings.TrimSpace(re.ReplaceAllString(name, ""))
	return cleanedName
}

// parseGetPlaylistResponse transforms the playlist object returned from `ytmusicapi` into our internal object.
func parseGetPlaylistResponse(response ytmusicAPIGetPlaylistResponse) utils.Playlist {
	var tracks []utils.Track
	for _, entry := range response.Data.Tracks {
		var artists []string
		for _, artist := range entry.Artists {
			artists = append(artists, cleanTrackArtist(artist))
		}

		tracks = append(tracks, utils.Track{ID: entry.Identifier, Title: utils.CleanTrackTitle(entry.Title), Artists: artists})
	}

	return utils.Playlist{
		ID:          response.Data.Identifier,
		Title:       response.Data.Title,
		Description: response.Data.Description,
		Tracks:      tracks,
	}
}

func setupRequestClient(reqClient *req.Client, baseURL string) *req.Client {
	return reqClient.
		SetBaseURL(baseURL).
		SetCommonContentType(utils.ApplicationJSON).
		SetCommonErrorResult(&ytmusicAPIErrorResponse{}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			// There is an underlying error, e.g. network error or unmarshal error.
			if resp.Err != nil {
				return nil
			}
			// Server returns an error message, convert it to human-readable Go error.
			if apiErr, ok := resp.ErrorResult().(*ytmusicAPIErrorResponse); ok {
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
		SetCommonRetryCount(2)
}
