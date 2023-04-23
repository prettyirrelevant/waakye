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
