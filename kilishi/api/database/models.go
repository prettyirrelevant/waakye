package database

// PlaylistConversionHistoryInDB represents playlist conversion history stored in a database.
type PlaylistConversionHistoryInDB struct {
	ID              string `db:"id"`
	PlaylistURL     string `db:"playlist_url"`
	Source          string `db:"source"`
	Destination     string `db:"destination"`
	ConversionCount int    `db:"conversion_count"`
	CreatedAt       int    `db:"created_at"`
	UpdatedAt       int    `db:"updated_at"`
}

// OauthCredentialsInDB represents OAuth credentials stored in a database.
type OauthCredentialsInDB struct {
	ID          string `db:"id"`
	Platform    string `db:"platform"`
	Credentials string `db:"credentials"`
	CreatedAt   int    `db:"created_at"`
	UpdatedAt   int    `db:"updated_at"`
}
