package applemusic

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/prettyirrelevant/kilishi/utils"
)

func New(opts InitialisationOpts) *AppleMusic {
	return &AppleMusic{}
}

func (a *AppleMusic) GetPlaylist(playlistURL string) (utils.Playlist, error) {
	var foundError error
	var playlist utils.Playlist

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.AllowedDomains("music.apple.com"),
	)

	c.OnHTML("div.songs-list-row", func(e *colly.HTMLElement) {
		var artistes []string
		e.ForEach("div.songs-list__col.songs-list__col--secondary > div.songs-list__song-link-wrapper > span > a.click-action", func(i int, h *colly.HTMLElement) {
			artistes = append(artistes, h.Text)
		})

		playlist.Tracks = append(playlist.Tracks, utils.Track{
			Title:   utils.CleanTrackTitle(e.ChildText("div.songs-list-row__song-name")),
			Artists: artistes,
		})
	})

	c.OnHTML("div.container-detail-header > div.description p.content", func(e *colly.HTMLElement) {
		playlist.Description = e.Text
	})

	c.OnHTML("div.container-detail-header > div.headings > h1.headings__title", func(e *colly.HTMLElement) {
		playlist.Title = e.Text
	})

	c.OnError(func(r *colly.Response, err error) {
		foundError = fmt.Errorf("applemusic: playlist retrieval failed with %s", err.Error())
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s ...", r.URL)
	})

	c.Visit(playlistURL)

	if foundError != nil {
		return utils.Playlist{}, foundError
	}

	return playlist, nil
}

func (a *AppleMusic) CreatePlaylist(playlist utils.Playlist, accessToken string) (string, error) {
	return "", nil
}

func (a *AppleMusic) RequiresAccessToken() bool {
	return true
}

func (a *AppleMusic) GetAuthorizationCode(code string) (utils.OauthCredentials, error) {
	return utils.OauthCredentials{}, nil
}
