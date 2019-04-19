package main

import (
	"fmt"

	"github.com/the-fusy/rentit/maps"

	"github.com/the-fusy/rentit/flat"
	"github.com/the-fusy/rentit/parser"
)

func main() {
	address := "Москва Льва Толстого 16"
	latitude, longitude, _ := maps.GetCoordinates(&address)
	lvaTolstogo := maps.Place{
		Latitude:  latitude,
		Longitude: longitude,
	}

	cianParser := &parser.ParserCian{}
	flatsRequest := flat.FlatsRequest{
		City: flat.MOSCOW,
	}
	flats := parser.GetFlats(cianParser, &flatsRequest, 1)

	for i := range flats {
		err := flats[i].FillCoordinates()
		if err != nil {
			fmt.Printf("ERROR %s\n", err.Error())
		} else {
			flat := maps.Place{
				Latitude:  flats[i].Latitude,
				Longitude: flats[i].Longitude,
			}
			travelTime, _ := maps.GetTravelTime(flat, lvaTolstogo)
			fmt.Println(flats[i].URL, travelTime/60, "minutes")
		}
	}
}
