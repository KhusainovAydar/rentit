package main

import (
	"context"
	"log"
	"sync"

	"github.com/the-fusy/rentit/flat"
	"github.com/the-fusy/rentit/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	err := mongo.InitDataBase()
	if err != nil {
		log.Fatal(err)
	}

	flats, err := mongo.GetCollection("rentit")
	if err != nil {
		log.Fatal(err)
	}

	cur, err := flats.Find(context.TODO(), bson.D{
		{"$or", bson.A{
			bson.D{{"processed", false}},
			bson.D{{"processed", bson.M{"$exists": false}}},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for cur.Next(context.TODO()) {
		var flat flat.Flat
		err = cur.Decode(&flat)
		if err != nil {
			log.Print(flat)
			log.Print(err)
			continue
		}
		wg.Add(1)
		go flat.Process(&wg)
	}
	cur.Close(context.TODO())
	wg.Wait()
}
