package flat

import (
	"context"
	"log"
	"strconv"
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

	log.Print("Processing with flat ", flat.ID)

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

	if flat.Rooms >= 3 && flat.Rooms <= 4 && flat.Price <= 80000 && travelTime <= 50 {
		pereezhaem := telegram.Chat{Username: "pereezhaem"}
		_, err = telegram.SendMessage(&pereezhaem, flat.GetDescription(travelTime), false, false)
		images := append(flat.PlanImages, flat.Images...)
		_, err = telegram.SendPhotos(&pereezhaem, &images)
	}

	flat.Processed = true
	_, err = flats.ReplaceOne(context.TODO(), bson.D{{"_id", flat.ID}}, flat)
	if err != nil {
		log.Print(err)
		return
	}

}

func (flat *Flat) GetDescription(travelTime uint16) *string {
	var text string

	if flat.Rooms == 0 {
		text += "Студия"
	} else {
		text += strconv.FormatUint(uint64(flat.Rooms), 10) + " "
		switch flat.Rooms % 10 {
		case 1:
			text += "комната"
		case 2, 3, 4:
			text += "комнаты"
		default:
			text += "комнат"
		}
	}
	text += "\n"

	text += "Стоит " + strconv.FormatUint(uint64(flat.Price), 10) + "\n"

	if flat.Prepayment == 0 {
		text += "Без залога"
	} else {
		text += "Залог " + strconv.FormatUint(uint64(flat.Prepayment), 10)
	}
	text += "\n"

	if flat.Fee == 0 {
		text += "Без комиссии"
	} else {
		text += "Комиссия " + strconv.FormatUint(uint64(flat.Fee), 10)
	}
	text += "\n"

	text += "Ехать до офиса " + strconv.FormatUint(uint64(travelTime), 10) + " "
	switch travelTime % 10 {
	case 1:
		text += "минуту"
	case 2, 3, 4:
		text += "минуты"
	default:
		text += "минут"
	}
	text += "\n"

	text += flat.URL + "\n"

	text += "Фоточки ⬇️"

	return &text
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
