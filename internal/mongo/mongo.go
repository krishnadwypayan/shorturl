package mongo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/krishnadwypayan/shorturl/internal/logger"
	"github.com/krishnadwypayan/shorturl/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "shortify"
	collectionName = "url_mappings"
)

var (
	user     = os.Getenv("MONGO_USERNAME")
	password = os.Getenv("MONGO_PASSWORD")
	mongoUri = fmt.Sprintf("mongodb+srv://%s:%s@cluster0.obqvkgx.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", user, password)
)

func InsertUrlMapping(req model.ShortURLRequest, id string) error {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		mongoUri,
	))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// err = client.Ping(ctx, nil)

	// if err != nil {
	// 	fmt.Println("There was a problem connecting to your Atlas cluster. Check that the URI includes a valid username and password, and that your IP address has been added to the access list. Error: ")
	// 	panic(err)
	// }

	// logger.Debug().Msg("Connected to MongoDB!\n")

	collection := client.Database(dbName).Collection(collectionName)
	doc := UrlDocument{
		LongURL:   req.LongURL,
		ID:        id,
		IsAlias:   req.Alias != "",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(req.TTL) * time.Second),
	}

	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		logger.Error().Msg(fmt.Sprintf("Failed to insert URL mapping: %v", err))
		return fmt.Errorf("failed to insert URL mapping: %w", err)
	}
	logger.Info().Msg("URL mapping inserted successfully")
	return nil
}

func CheckAliasExists(alias string) (bool, error) {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		mongoUri,
	))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection := client.Database(dbName).Collection(collectionName)
	filter := map[string]interface{}{
		"id":       alias,
		"is_alias": true,
	}

	var result UrlDocument
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Alias does not exist
			return false, nil
		}
		logger.Error().Msg(fmt.Sprintf("Failed to check alias existence: %v", err))
		return false, fmt.Errorf("failed to check alias existence: %w", err)
	}
	logger.Info().Msg("Alias exists in the database")
	return true, nil
}

type UrlDocument struct {
	LongURL   string    `bson:"long_url"`
	ID        string    `bson:"id"`
	IsAlias   bool      `bson:"is_alias"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}
