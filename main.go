package main

import (
	"fmt"

	"github.com/the-fusy/rentit/flat"
	"github.com/the-fusy/rentit/parser"
)

func main() {
	cianParser := &parser.ParserCian{}
	flatsRequest := parser.FlatsRequest{
		City: flat.MOSCOW,
	}
	flats := flatsRequest.GetFlats(cianParser, 1)
	for i := range flats {
		err := flats[i].FillCoordinates()
		if err != nil {
			fmt.Printf("ERROR %s\n", err.Error())
		} else {
			fmt.Println(flats[i])
		}
	}
}
