package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/xid"

	"github.com/prettyirrelevant/waakye/pkg/utils/types"
)

type Database struct {
	db *sqlx.DB
}

// New creates a `Database` struct.
func New(databaseURL string) (*Database, error) {
	db, err := sqlx.Connect("sqlite3", databaseURL)
	if err != nil {
		return nil, err
	}

	// TODO: log that database connection was successful!
	return &Database{db: db}, nil
}

// GenerateID returns an identifier for use as primary key in database.
func GenerateID() string {
	return xid.New().String()
}

// GetOauthCredentials retrieves oauth credentials from the database for a streaming platform.
func (d *Database) GetOauthCredentials(platform string) (OauthCredentialsInDB, error) {
	var credentials OauthCredentialsInDB
	err := d.db.Get(&credentials, "SELECT * FROM oauth_credentials WHERE platform=$1;", platform)
	if err != nil {
		return credentials, fmt.Errorf("api.database: could not fetch oauth credentials for %s due to %s", platform, err.Error())
	}

	return credentials, nil
}

// SetOauthCredentials saves the oauth credentials of a streaming platform in the database.
func (d *Database) SetOauthCredentials(platform string, entry types.OauthCredentials) error {
	authCredentialsString, err := entry.ToString()
	if err != nil {
		return err
	}

	d.db.Exec(
		`
			INSERT INTO oauth_credentials (id, platform, credentials, updated_at)
			VALUES ($1, $2, $3, $4) ON CONFLICT(platform) DO UPDATE SET credentials=excluded.credentials, updated_at=excluded.updated_at;
		`,
		GenerateID(), platform, authCredentialsString, time.Now().UnixMicro(),
	)
	return nil
}
