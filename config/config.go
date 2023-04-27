package config

import (
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
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
	BaseDir                 string
}

func New() (*Config, error) {
	var cfg Config

	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return &cfg, err
	}

	cfg.BaseDir = currentWorkingDir
	err = godotenv.Load(filepath.Join(cfg.BaseDir, ".env"))
	if err != nil {
		return &cfg, err
	}

	if err := env.Parse(&cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
