package main

import (
	"github.com/the-fusy/rentit/flat"
	"github.com/the-fusy/rentit/parser"
)

func main() {
	// address := "Москва Льва Толстого 16"
	// latitude, longitude, _ := maps.GetCoordinates(&address)

	cianParser := &parser.ParserCian{}
	flatsRequest := flat.FlatsRequest{
		City: flat.MOSCOW,
	}

	parser.GetFlats(cianParser, &flatsRequest, 50)
	// results := make(chan []interface{})

	// for i := range flats {
	// 	go flats[i].GetTravelTime(latitude, longitude, results)
	// }

	// for range flats {
	// 	result := <-results
	// 	switch travelTime := result[1].(type) {
	// 	case error:
	// 		fmt.Println(travelTime.Error())
	// 	case int16:
	// 		fmt.Println(result[0].(*flat.Flat).URL, travelTime/60, "minutes")
	// 	}
	// }
}
