package flat

import (
	"errors"
	"sync"

	"github.com/the-fusy/rentit/maps"
)

const (
	MOSCOW uint8 = 1
	SPB
)

type Flat struct {
	URL       string
	Rooms     uint8
	Address   string
	Latitude  float64
	Longitude float64
	Price     uint64
	Area      uint64
	FromOwner bool
}

func (flat *Flat) FillCoordinates(wg *sync.WaitGroup) error {
	defer wg.Done()
	latitude, longitude, err := maps.GetCoordinates(&flat.Address)
	if err != nil {
		return err
	}
	flat.Latitude = latitude
	flat.Longitude = longitude
	return nil
}

func (flat *Flat) GetTravelTime(latitude, longitude float64, result chan []interface{}) error {
	data := []interface{}{flat, errors.New("Error to get travel time")}
	defer func() { result <- data }()

	from := maps.Place{Latitude: flat.Latitude, Longitude: flat.Longitude}
	to := maps.Place{Latitude: latitude, Longitude: longitude}
	travelTime, err := maps.GetTravelTime(from, to)
	if err != nil {
		return err
	}

	data[1] = travelTime
	return nil
}

type FlatsRequest struct {
	City       uint8
	Rooms      []uint8
	MinPrice   uint64
	MaxPrice   uint64
	LastUpdate uint64
	FromOwner  bool
}
