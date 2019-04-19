package parser

import "github.com/the-fusy/rentit/flat"

type Parser interface {
	getURL(req *flat.FlatsRequest, page int) string
	parsePage(url *string, flatsChan chan []flat.Flat)
}

func GetFlats(parser Parser, request *flat.FlatsRequest, maxPage int) []flat.Flat {
	flatsChan := make(chan []flat.Flat)
	flats := make([]flat.Flat, 0)
	defer close(flatsChan)
	for i := 1; i <= maxPage; i++ {
		url := parser.getURL(request, i)
		go parser.parsePage(&url, flatsChan)
	}
	for i := 1; i <= maxPage; i++ {
		flats = append(flats, <-flatsChan...)
	}
	return flats
}
