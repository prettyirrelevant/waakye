package database

// PlaylistConversionHistoryInDB represents playlist conversion history stored in a database.
type PlaylistConversionHistoryInDB struct {
	PlaylistURL     string `bson:"playlist_url"`
	Source          string `bson:"source"`
	Destination     string `bson:"destination"`
	ConversionCount int    `bson:"conversion_count"`
	CreatedAt       int    `bson:"created_at"`
	UpdatedAt       int    `bson:"updated_at"`
}

// OauthCredentialsInDB represents OAuth credentials stored in a database.
type OauthCredentialsInDB struct {
	Platform    string `bson:"platform"`
	Credentials string `bson:"credentials"`
	CreatedAt   int    `bson:"created_at"`
	UpdatedAt   int    `bson:"updated_at"`
}
