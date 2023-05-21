package database

// OauthCredentialsInDB represents OAuth credentials stored in a database.
type OauthCredentialsInDB struct {
	Platform    string `redis:"platform"`
	Credentials []byte `redis:"credentials"`
	CreatedAt   int    `redis:"created_at"`
	UpdatedAt   int    `redis:"updated_at"`
}
