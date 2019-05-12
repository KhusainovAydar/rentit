package main

import (
	"context"
	"fmt"
	"log"

	"github.com/the-fusy/rentit/config"
	"github.com/the-fusy/rentit/flat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURL := fmt.Sprintf("mongodb://%s:%s@%s/%s", config.MongoUser, config.MongoPassword, config.MongoHost, config.MongoDatabase)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(config.MongoDatabase).Collection("rentit")

	cur, err := collection.Find(context.TODO(), bson.D{
		{"$or", bson.A{
			bson.D{{"processed", false}},
			bson.D{{"processed", bson.M{"$exists": false}}},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var flat flat.Flat
		err = cur.Decode(&flat)
		if err != nil {
			log.Print(flat)
			log.Fatal(err)
		}
		collection.UpdateOne(context.TODO(), flat, bson.D{
			{"$set", bson.D{{"processed", true}}},
		})
	}
	cur.Close(context.TODO())
}
