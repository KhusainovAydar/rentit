package parser

import (
	"fmt"

	"github.com/the-fusy/rentit/flat"
)

type Parser interface {
	getURL(req *flat.FlatsRequest, page int) string
	parsePage(url *string, flatsChan chan []interface{})
}

func GetFlats(parser Parser, request *flat.FlatsRequest, maxPage int) []flat.Flat {
	flatsChan := make(chan []interface{})
	flats := make([]flat.Flat, 0)
	defer close(flatsChan)
	for i := 1; i <= maxPage; i++ {
		url := parser.getURL(request, i)
		go parser.parsePage(&url, flatsChan)
	}
	for i := 1; i <= maxPage; i++ {
		fmt.Println(i)
		newFlats := <-flatsChan
		fmt.Println(newFlats)
		// wg := sync.WaitGroup{}
		// wg.Add(len(newFlats))
		// for i := range newFlats {
		// 	go newFlats[i].FillCoordinates(&wg)
		// }
		// wg.Wait()

		// flats = append(flats, newFlats...)
	}
	return flats
}
