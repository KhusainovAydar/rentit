package parser

import "github.com/the-fusy/rentit/flat"

type FlatsRequest struct {
	City       uint8
	Rooms      []uint8
	MinPrice   uint64
	MaxPrice   uint64
	LastUpdate uint64
	FromOwner  bool
}

type Parser interface {
	getURL(req *FlatsRequest, page int) string
	parsePage(url *string, flatsChan chan []flat.Flat)
}

func (request *FlatsRequest) GetFlats(parser Parser, maxPage int) []flat.Flat {
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
