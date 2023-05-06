package applemusic

import "github.com/imroc/req/v3"

type AppleMusic struct {
	RequestClient *req.Client
	Config        *Config
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
