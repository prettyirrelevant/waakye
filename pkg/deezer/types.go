package deezer

import (
	"fmt"

	"github.com/imroc/req/v3"
)

type Deezer struct {
	RequestClient *req.Client
	Config        *Config
}

type InitialisationOpts struct {
	RequestClient     *req.Client
	AppID             string
	BaseAPIURI        string
	ClientSecret      string
	AuthenticationURI string
}

type Config struct {
	BaseAPIURI        string
	AppID             string
	ClientSecret      string
	AuthenticationURI string
}

// /////////////////////////////////////////////////
// API Types
// /////////////////////////////////////////////////
type deezerAPIGetPlaylistResponse struct {
	ID          int                         `json:"id"`
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	Tracks      deezerAPITracksDataResponse `json:"tracks`
}

type deezerAPISearchTrackResponse struct {
	Results deezerAPITracksDataResponse
}

type deezerAPITracksDataResponse struct {
	Data []struct {
		ID     int    `json:"id"`
		Title  string `json:"title"`
		Artist struct {
			Name string `json:"name"`
		} `json:"artist"`
	} `json:"data"`
}

type deezerAPICreatePlaylistResponse struct {
	ID string `json:"id"`
}

type deezerAPIBearerCredentialsResponse struct {
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires"`
}

type deezerAPIErrorResponse struct {
	APIError struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

func (e *deezerAPIErrorResponse) Error() string {
	return fmt.Sprintf("Deezer API Error: code %v type: %s  reason: %s", e.APIError.Code, e.APIError.Type, e.APIError.Message)
}