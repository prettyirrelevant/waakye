package applemusic

import "github.com/imroc/req/v3"

type AppleMusic struct {
	RequestClient *req.Client
	Config        *Config
}

type InitialisationOpts struct {
	RequestClient             *req.Client
	BaseAPIURL                string
	ClientID                  string
	UserID                    string
	ClientSecret              string
	AuthenticationURL         string
	AuthenticationRedirectURL string
}

type Config struct {
	BaseAPIURL                string
	UserID                    string
	ClientID                  string
	ClientSecret              string
	AuthenticationURL         string
	AuthenticationRedirectURL string
}
