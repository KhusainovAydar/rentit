package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/the-fusy/rentit/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func InitDataBase() error {
	mongoURL := fmt.Sprintf("mongodb://%s:%s@%s/%s", config.MongoUser, config.MongoPassword, config.MongoHost, config.MongoDatabase)
	var err error

	DB, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		return errors.New("failed to connect when init database")
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = DB.Ping(ctx, nil)
	if err != nil {
		return errors.New("failed to ping when init database")
	}

	return nil
}

func GetCollection(name string) (*mongo.Collection, error) {
	if DB == nil {
		return nil, errors.New("database is not initialized")
	}
	return DB.Database(config.MongoDatabase).Collection(name), nil
}
