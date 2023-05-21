package database

import (
	"context"
	"fmt"
	"time"

	"github.com/prettyirrelevant/kilishi/streaming_platforms/aggregator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/prettyirrelevant/kilishi/utils"
)

// Database represents a connection to a SQLite database.
type Database struct {
	client *mongo.Client
	db     *mongo.Database
}

// New creates a new Database struct and connects to a MongoDB database using the provided URL.
func New(databaseURL string) (*Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(databaseURL))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Database{
		client: client,
		db:     client.Database("kilishi"),
	}, nil
}

// GetDBOauthCredentials retrieves the OAuth credentials for a given music streaming platform from the database.
func (d *Database) GetDBOauthCredentials(platform aggregator.MusicStreamingPlatform) (OauthCredentialsInDB, error) {
	var credentials OauthCredentialsInDB

	filter := bson.M{"platform": platform}
	err := d.db.Collection("oauth_credentials").FindOne(context.TODO(), filter).Decode(&credentials)
	if err != nil {
		return credentials, fmt.Errorf("database: could not fetch oauth credentials for %s due to %s", platform, err.Error())
	}

	return credentials, nil
}

// SetOauthCredentials saves the OAuth credentials for a given music streaming platform in the database.
func (d *Database) SetOauthCredentials(platform aggregator.MusicStreamingPlatform, credentials utils.OauthCredentials) error {
	authCredentialsString, err := credentials.ToString()
	if err != nil {
		return fmt.Errorf("database: could not set oauth credentials for %s due to %s", platform, err.Error())
	}

	now := time.Now().Unix()
	filterQuery := bson.M{"platform": platform}
	updateQuery := bson.M{
		"$set":         bson.D{{Key: "platform", Value: string(platform)}, {Key: "credentials", Value: authCredentialsString}, {Key: "updated_at", Value: now}},
		"$setOnInsert": bson.M{"created_at": now},
	}
	options := options.FindOneAndUpdate().SetUpsert(true)
	err = d.db.Collection("oauth_credentials").FindOneAndUpdate(context.TODO(), filterQuery, updateQuery, options).Err()
	if err != nil {
		return fmt.Errorf("database: could not set oauth credentials for %s due to %s", platform, err.Error())
	}

	return nil
}
