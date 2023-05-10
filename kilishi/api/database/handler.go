package database

import (
	"embed"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/rs/xid"

	"github.com/prettyirrelevant/kilishi/pkg/aggregator"
	"github.com/prettyirrelevant/kilishi/pkg/utils"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Database represents a connection to a SQLite database.
type Database struct {
	db *sqlx.DB
}

// New creates a new Database struct and connects to a Postgres database with the provided URL.
func New(databaseURL string) (*Database, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// TODO: log that migrations ran successfully!
	goose.SetBaseFS(embedMigrations)
	if err = goose.SetDialect("postgres"); err != nil {
		return nil, err
	}
	if err = goose.Up(db.DB, "migrations"); err != nil {
		return nil, err
	}

	// TODO: log that database connection was successful!
	return &Database{db: db}, nil
}

// GenerateID returns a new unique identifier to use as a primary key in the database.
func GenerateID() string {
	return xid.New().String()
}

// GetOauthCredentials retrieves the OAuth credentials for a given music streaming platform from the database.
func (d *Database) GetOauthCredentials(platform aggregator.MusicStreamingPlatform) (OauthCredentialsInDB, error) {
	var credentials OauthCredentialsInDB

	err := d.db.Get(&credentials, "SELECT * FROM oauth_credentials WHERE platform=$1;", platform)
	if err != nil {
		return credentials, fmt.Errorf("database: could not fetch OAuth credentials for %s due to %s", platform, err.Error())
	}

	return credentials, nil
}

// SetOauthCredentials saves the OAuth credentials for a given music streaming platform in the database.
func (d *Database) SetOauthCredentials(platform aggregator.MusicStreamingPlatform, credentials utils.OauthCredentials) error {
	authCredentialsString, err := credentials.ToString()
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	result, err := d.db.Exec(
		`
			INSERT INTO oauth_credentials (id, platform, credentials, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) ON CONFLICT(platform) DO UPDATE SET credentials=excluded.credentials, updated_at=excluded.updated_at;
		`,
		GenerateID(), platform, authCredentialsString, now, now,
	)
	if err != nil {
		return err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowCount != 1 {
		return fmt.Errorf("database: failed to add %s credentials", platform)
	}

	return nil
}
