package flat

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/the-fusy/rentit/telegram"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/the-fusy/rentit/maps"
	"github.com/the-fusy/rentit/mongo"
)

type Flat struct {
	ID         primitive.ObjectID `bson:"_id"`
	URL        string
	Rooms      uint8
	Address    string
	Latitude   float64
	Longitude  float64
	Price      uint64
	Area       float64
	Fee        uint64
	Prepayment uint64
	Images     []string
	PlanImages []string `bson:"planImages"`
	Processed  bool
}

func (flat *Flat) Process(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Print("Start with flat ", flat.ID)

	flats, err := mongo.GetCollection("rentit")
	if err != nil {
		log.Print(err)
		return
	}

	if flat.Latitude == 0 || flat.Longitude == 0 {
		err = flat.FillCoordinates()
		if err != nil {
			log.Print(err)
			return
		}
	}

	address := "Москва Льва Толстого 16"
	latitude, longitude, _ := maps.GetCoordinates(&address)
	travelTime, err := flat.GetTravelTime(latitude, longitude)
	if err != nil {
		log.Print(err)
		return
	}

	travelTime /= 60

	if flat.Rooms >= 3 && flat.Rooms <= 4 && flat.Price <= 80000 {
		pereezhaem := telegram.Chat{Username: "pereezhaem"}
		text := fmt.Sprintf("nЕхать %v минут\n%v", travelTime, flat.URL)
		_, err = pereezhaem.SendMessage(text, false, false)
		if len(flat.Images) > 1 {
			_, err = pereezhaem.SendPhotos(&flat.Images)
		}
	}

	flat.Processed = true
	_, err = flats.ReplaceOne(context.TODO(), bson.D{{"_id", flat.ID}}, flat)
	if err != nil {
		log.Print(err)
		return
	}

}

func (flat *Flat) FillCoordinates() error {
	latitude, longitude, err := maps.GetCoordinates(&flat.Address)
	if err != nil {
		return err
	}
	flat.Latitude = latitude
	flat.Longitude = longitude
	return nil
}

func (flat *Flat) GetTravelTime(latitude, longitude float64) (uint16, error) {
	from := maps.Place{Latitude: flat.Latitude, Longitude: flat.Longitude}
	to := maps.Place{Latitude: latitude, Longitude: longitude}
	travelTime, err := maps.GetTravelTime(from, to)
	if err != nil {
		return 0, err
	}
	return travelTime, nil
}

type FlatsRequest struct {
	City       uint8
	Rooms      []uint8
	MinPrice   uint64
	MaxPrice   uint64
	LastUpdate uint64
	FromOwner  bool
}
