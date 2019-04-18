package flat

import (
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

func (flat *Flat) FillCoordinates() error {
	latitude, longitude, err := maps.GetCoordinates(&flat.Address)
	if err != nil {
		return err
	}
	flat.Latitude = latitude
	flat.Longitude = longitude
	return nil
}
