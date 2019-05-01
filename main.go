package main

import (
	"fmt"

	"github.com/the-fusy/rentit/maps"
	"github.com/the-fusy/rentit/telegram"

	"github.com/the-fusy/rentit/flat"
	"github.com/the-fusy/rentit/parser"
)

func main() {
	address := "Москва Льва Толстого 16"
	latitude, longitude, _ := maps.GetCoordinates(&address)
	telegramUser := telegram.User{ID: 12345}

	cianParser := &parser.ParserCian{}
	flatsRequest := flat.FlatsRequest{
		City:     flat.MOSCOW,
		MaxPrice: 120000,
		Rooms:    []uint8{5},
	}

	flats := parser.GetFlats(cianParser, &flatsRequest, 10)
	results := make(chan []interface{})

	for i := range flats {
		go flats[i].GetTravelTime(latitude, longitude, results)
	}

	for range flats {
		result := <-results
		switch travelTime := result[1].(type) {
		case error:
			fmt.Println(travelTime.Error())
		case int16:
			minutes := travelTime / 60
			if minutes <= 30 {
				telegramUser.SendMessage(fmt.Sprintf("%v %v minutes", result[0].(*flat.Flat).URL, minutes))
			}
		}
	}
}
