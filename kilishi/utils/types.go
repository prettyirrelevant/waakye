package utils

import "encoding/json"

const ApplicationJSON = "application/json"

// Playlist represents a playlist entry from any of the supported streaming platform internally.
type Playlist struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Tracks      []Track `json:"tracks"`
}

// Track represents a song entry in a playlist from any of the supported streaming platform internally.
type Track struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Artists []string `json:"artists"`
}

type OauthCredentials struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    int
}

func (o *OauthCredentials) ToString() (string, error) {
	stringFormat, err := json.Marshal(o)
	if err != nil {
		return "", err
	}

	return string(stringFormat), nil
}

func (o OauthCredentials) FromDB(payload string) error {
	err := json.Unmarshal([]byte(payload), &o)
	if err != nil {
		return err
	}

	return nil
}
