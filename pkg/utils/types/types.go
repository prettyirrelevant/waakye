package types

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

// ToSearchQuery takes a track and transforms it into a search query.
func (t *Track) ToSearchQuery() string {
	searchQuery := t.Title + " by"
	for index, artiste := range t.Artists {
		if len(t.Artists) == 1 {
			searchQuery += " " + artiste
		} else if len(t.Artists) > 1 && len(t.Artists)-1 == index {
			searchQuery += " and " + artiste
		} else {
			searchQuery += " " + artiste
			if index < len(t.Artists)-2 {
				searchQuery += ","
			}
		}
	}
	return searchQuery
}
