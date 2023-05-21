package config

import "github.com/caarlos0/env/v7"

type Config struct {
	Debug                   bool   `env:"DEBUG,notEmpty"`
	Port                    int    `env:"PORT,notEmpty"`
	Address                 string `env:"ADDRESS,notEmpty"`
	SecretKey               string `env:"SECRET_KEY,notEmpty"`
	RedisURL                string `env:"REDIS_URL,notEmpty"`
	InitializationVector    string `env:"INITIALIZATION_VECTOR,notEmpty"`
	DatabaseURL             string `env:"DATABASE_URL,notEmpty"`
	SpotifyClientID         string `env:"SPOTIFY_CLIENT_ID,notEmpty"`
	SpotifyClientSecret     string `env:"SPOTIFY_CLIENT_SECRET,notEmpty"`
	SpotifyClientAuthURL    string `env:"SPOTIFY_CLIENT_AUTH_URL,notEmpty"`
	SpotifyBaseApiURL       string `env:"SPOTIFY_BASE_API_URL,notEmpty"`
	SpotifyUserID           string `env:"SPOTIFY_USER_ID,notEmpty"`
	SpotifyAuthRedirectURL  string `env:"SPOTIFY_AUTH_REDIRECT_URL,notEmpty"`
	DeezerBaseApiURL        string `env:"DEEZER_BASE_API_URL,notEmpty"`
	DeezerAppID             string `env:"DEEZER_APP_ID,notEmpty"`
	DeezerClientSecret      string `env:"DEEZER_CLIENT_SECRET,notEmpty"`
	DeezerAuthenticationURL string `env:"DEEZER_AUTHENTICATION_URL,notEmpty"`
	YTMusicApiBaseUrl       string `env:"YTMUSICAPI_BASE_URL,notEmpty"`
}

func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
