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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int    `json:"expires_at"`
}

func (o *OauthCredentials) ToBytes() ([]byte, error) {
	var result []byte

	result, err := json.Marshal(o)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (o *OauthCredentials) FromDB(payload []byte) error {
	err := json.Unmarshal(payload, &o)
	if err != nil {
		return err
	}

	return nil
}
