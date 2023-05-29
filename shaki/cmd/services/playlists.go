package services

import (
	"fmt"

	"github.com/imroc/req/v3"
)

var reqClient = req.C().
	SetBaseURL("https://kilishi.dokku.prettyirrelevant.wtf/api/v1").
	OnAfterResponse(func(client *req.Client, resp *req.Response) error {
		if resp.Err != nil {
			return nil
		}
		if apiErr, ok := resp.ErrorResult().(*apiError); ok {
			resp.Err = apiErr
			return nil
		}
		if !resp.IsSuccessState() {
			return fmt.Errorf("bad response, raw dump:\n%s", resp.Dump())
		}
		return nil
	}).
	SetCommonRetryCount(2)

func GetPlaylist(url, platform string) (APIGetPlaylistResponse, error) {
	var response APIGetPlaylistResponse
	err := reqClient.
		Get("/playlists").
		SetQueryParams(map[string]string{"platform": platform, "playlist_url": url}).
		Do().
		Into(&response)

	if err != nil {
		return response, err
	}

	return response, nil
}

func FindTrack(title, platform string, artists []string) (APIFindTrackResponse, error) {
	var response APIFindTrackResponse

	artistsMap := make(map[string]string)
	for _, artist := range artists {
		artistsMap["artists"] = artist
	}

	err := reqClient.
		Get("/playlists/tracks").
		SetQueryParams(map[string]string{"platform": platform, "title": title}).
		SetQueryParams(artistsMap).
		Do().
		Into(&response)

	if err != nil {
		return response, err
	}

	return response, nil
}

func CreatePlaylist(title, description, platform string, tracks []TrackResponse) (APICreatePlaylistResponse, error) {
	var response APICreatePlaylistResponse

	var tracksMap []map[string]any
	for _, track := range tracks {
		tracksMap = append(tracksMap, map[string]any{"id": track.ID, "title": track.Title, "artists": track.Artists})
	}

	payload := make(map[string]any)
	payload["platform"] = platform
	payload["playlist"] = map[string]any{
		"title":       title,
		"description": description,
		"tracks":      tracksMap,
	}

	err := reqClient.
		Post("/playlists").
		SetBodyJsonMarshal(payload).
		Do().
		Into(&response)

	if err != nil {
		return response, err
	}

	return response, nil
}

type APIGetPlaylistResponse struct {
	Data struct {
		ID          string          `json:"id"`
		Title       string          `json:"title"`
		Description string          `json:"description"`
		Tracks      []TrackResponse `json:"tracks"`
	} `json:"data"`
	Message string `json:"message"`
}

type APIFindTrackResponse struct {
	Data    TrackResponse `json:"data"`
	Message string        `json:"message"`
}

type TrackResponse struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Artists []string `json:"artists"`
}

type APICreatePlaylistResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
}

type apiError struct {
	Message string
	Errors  []string
}

func (a *apiError) Error() string {
	return fmt.Sprintf("API Error: message: %v  errors: %+v", a.Message, a.Errors)
}
