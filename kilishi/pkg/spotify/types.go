package spotify

import (
	"fmt"

	"github.com/imroc/req/v3"
)

type Spotify struct {
	RequestClient *req.Client
	Config        Config
}

type InitialisationOpts struct {
	RequestClient             *req.Client
	BaseAPIURI                string
	ClientID                  string
	UserID                    string
	ClientSecret              string
	AuthenticationURI         string
	AuthenticationRedirectURI string
}

type Config struct {
	BaseAPIURI                string
	UserID                    string
	ClientID                  string
	ClientSecret              string
	AuthenticationURI         string
	AuthenticationRedirectURI string
}

// /////////////////////////////////////////////////
// API Types
// /////////////////////////////////////////////////
type spotifyAPIGetPlaylistResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Images      []struct {
		URL string `json:"url"`
	} `json:"images"`
	Tracks spotifyAPITracksResponse `json:"tracks"`
}

type spotifyAPITracksResponse struct {
	Next  string `json:"next"`
	Total int    `json:"total"`
	Items []struct {
		Track struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"track"`
	} `json:"items"`
}

type spotifyAPISearchResponse struct {
	Tracks struct {
		Limit int    `json:"limit"`
		Next  string `json:"next"`
		Total int    `json:"total"`
		Items []struct {
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"items"`
	} `json:"tracks"`
}

type spotifyAPICreatePlaylistResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URI  string `json:"uri"`
}

type spotifyAPIClientCredentialsResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type spotifyAPIBearerCredentialsResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type spotifyAPIErrorResponse struct {
	APIError struct {
		Status  uint   `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *spotifyAPIErrorResponse) Error() string {
	return fmt.Sprintf("Spotify API Error: status: %v  reason: %s", e.APIError.Status, e.APIError.Message)
}
