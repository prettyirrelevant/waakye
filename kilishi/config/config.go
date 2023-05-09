package config

import "github.com/caarlos0/env/v7"

type Config struct {
	Debug                   bool   `env:"DEBUG,notEmpty"`
	SecretKey               string `env:"SECRET_KEY,notEmpty"`
	RedisURI                string `env:"REDIS_URI,notEmpty"`
	InitializationVector    string `env:"INITIALIZATION_VECTOR,notEmpty"`
	DatabaseURI             string `env:"DATABASE_URI,notEmpty"`
	SpotifyClientID         string `env:"SPOTIFY_CLIENT_ID,notEmpty"`
	SpotifyClientSecret     string `env:"SPOTIFY_CLIENT_SECRET,notEmpty"`
	SpotifyClientAuthURL    string `env:"SPOTIFY_CLIENT_AUTH_URL,notEmpty"`
	SpotifyBaseApiURL       string `env:"SPOTIFY_BASE_API_URL,notEmpty"`
	SpotifyUserID           string `env:"SPOTIFY_USER_ID,notEmpty"`
	SpotifyAuthRedirectURI  string `env:"SPOTIFY_AUTH_REDIRECT_URI,notEmpty"`
	DeezerBaseApiURL        string `env:"DEEZER_BASE_API_URL,notEmpty"`
	DeezerAppID             string `env:"DEEZER_APP_ID,notEmpty"`
	DeezerClientSecret      string `env:"DEEZER_CLIENT_SECRET,notEmpty"`
	DeezerAuthenticationURI string `env:"DEEZER_AUTHENTICATION_URI,notEmpty"`
}

func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
