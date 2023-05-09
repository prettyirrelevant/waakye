package utils

import "encoding/json"

const ApplicationJSON = "application/json"

// Playlist represents a playlist entry from any of the supported streaming platform internally.
type Playlist struct {
	ID          string
	Title       string
	Description string
	Tracks      []Track
}

// Track represents a song entry in a playlist from any of the supported streaming platform internally.
type Track struct {
	ID      string
	Title   string
	Artists []string
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

func OauthCredentialsFromDB(payload string) (OauthCredentials, error) {
	var credentials OauthCredentials
	err := json.Unmarshal([]byte(payload), &credentials)
	if err != nil {
		return credentials, err
	}

	return credentials, nil
}
