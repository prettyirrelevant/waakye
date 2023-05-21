package database

import (
	"context"
	"fmt"
	"time"

	"github.com/prettyirrelevant/kilishi/streaming_platforms/aggregator"
	"github.com/redis/go-redis/v9"

	"github.com/prettyirrelevant/kilishi/utils"
)

var Ctx = context.TODO()

// Database represents a connection to a Redis instance.
type Database struct {
	client *redis.Client
}

// New creates a new Database struct and connects to a MongoDB database using the provided URL.
func New(databaseURL string) (*Database, error) {
	var db *Database

	opts, err := redis.ParseURL(databaseURL)
	if err != nil {
		return db, err
	}

	client := redis.NewClient(opts)
	if status := client.Ping(Ctx); status.Err() != nil {
		return db, status.Err()
	}

	db.client = client
	return db, nil
}

// GetDBOauthCredentials retrieves the OAuth credentials for a given music streaming platform from the database.
func (d *Database) GetDBOauthCredentials(platform aggregator.MusicStreamingPlatform) (OauthCredentialsInDB, error) {
	var dbCredentials OauthCredentialsInDB
	var hashKey = fmt.Sprintf("oauth_cred:%s", platform)

	err := d.client.HGetAll(Ctx, hashKey).Scan(&dbCredentials)
	if err != nil {
		return dbCredentials, fmt.Errorf("database: could not set oauth credentials for %s due to %s", platform, err.Error())
	}
	return dbCredentials, nil
}

// SetOauthCredentials saves the OAuth credentials for a given music streaming platform in the database.
func (d *Database) SetOauthCredentials(platform aggregator.MusicStreamingPlatform, credentials utils.OauthCredentials) error {
	var hashKey = fmt.Sprintf("oauth_cred:%s", platform)

	bytesCredentials, err := credentials.ToBytes()
	if err != nil {
		return fmt.Errorf("database: could not set oauth credentials for %s due to %s", platform, err.Error())
	}

	count, err := d.client.Exists(Ctx, hashKey).Result()
	if err != nil {
		return fmt.Errorf("database: could not set oauth credentials for %s due to %s", platform, err.Error())
	}
	if count == 1 {
		_, err = d.client.Pipelined(Ctx, func(p redis.Pipeliner) error {
			p.HSet(Ctx, hashKey, "credentials", bytesCredentials)
			p.HSet(Ctx, hashKey, "updated_at", time.Now().Unix())
			return nil
		})
		if err != nil {
			return fmt.Errorf("database: could not set oauth credentials for %s due to %s", platform, err.Error())
		}

		return nil
	}

	now := time.Now().Unix()
	_, err = d.client.Pipelined(Ctx, func(p redis.Pipeliner) error {
		p.HSet(Ctx, hashKey, "platform", string(platform))
		p.HSet(Ctx, hashKey, "credentials", bytesCredentials)
		p.HSet(Ctx, hashKey, "created_at", now)
		p.HSet(Ctx, hashKey, "updated_at", now)
		return nil
	})
	if err != nil {
		return fmt.Errorf("database: could not set oauth credentials for %s due to %s", platform, err.Error())
	}

	return nil
}
