package main

import (
	"fmt"

	"github.com/the-fusy/rentit/flat"
	"github.com/the-fusy/rentit/parser"
)

func main() {
	request := parser.FlatsRequest{
		City: flat.MOSCOW,
	}
	flats := request.GetFlats(&parser.ParserCian{}, 1)
	fmt.Println(flats)
}
