package config

import (
	"github.com/caarlos0/env/v7"
)

type Config struct {
	SecretKey               string `env:"SECRET_KEY,notEmpty"`
	MasaAuthUsername        string `env:"MASA_AUTH_USERNAME,notEmpty`
	MasaAuthPassword        string `env:"MASA_AUTH_PASSWORD,notEmpty`
	DatabaseURI             string `env:"DATABASE_URI,notEmpty"`
	Port                    int    `env:"PORT,notEmpty"`
	SpotifyClientID         string `env:"SPOTIFY_CLIENT_ID,notEmpty"`
	SpotifyClientSecret     string `env:"SPOTIFY_CLIENT_SECRET,notEmpty"`
	SpotifyClientAuthURL    string `env:"SPOTIFY_CLIENT_AUTH_URL,notEmpty"`
	SpotifyBaseApiURL       string `env:"SPOTIFY_BASE_API_URL,notEmpty"`
	SpotifyAuthEmail        string `env:"SPOTIFY_AUTH_EMAIL,notEmpty"`
	SpotifyAuthPassword     string `env:"SPOTIFY_AUTH_PASSWORD,notEmpty"`
	SpotifyAuthRedirectURI  string `env:"SPOTIFY_AUTH_REDIRECT_URI,notEmpty"`
	YTMusicBaseUrl          string `env:"YTMUSIC_BASE_API_URL,notEmpty"`
	YTMusicAuthCredentials  string `env:"YTMUSIC_AUTH_CREDENTIALS,notEmpty"`
	DeezerBaseApiURL        string `env:"DEEZER_BASE_API_URL,notEmpty"`
	DeezerAppID             string `env:"DEEZER_APP_ID,notEmpty"`
	DeezerClientSecret      string `env:"DEEZER_CLIENT_SECRET,notEmpty"`
	DeezerAuthenticationURI string `env:"DEEZER_AUTHENTICATION_URI,notEmpty"`
	DeezerAuthRedirectURI   string `env:"DEEZER_AUTH_REDIRECT_URI,notEmpty"`
	DeezerAuthEmail         string `env:"DEEZER_AUTH_EMAIL,notEmpty"`
	DeezerAuthPassword      string `env:"DEEZER_AUTH_PASSWORD,notEmpty"`
}

func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
