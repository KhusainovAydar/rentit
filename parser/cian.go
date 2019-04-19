package parser

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/the-fusy/rentit/flat"
)

type ParserCian struct{}

func (parser *ParserCian) getURL(req *flat.FlatsRequest, page int) string {
	url := "https://www.cian.ru/cat.php?deal_type=rent&engine_version=2&offer_type=flat&type=4"
	url += "&region=" + strconv.FormatUint(uint64(req.City), 10)
	url += "&p=" + strconv.FormatInt(int64(page), 10)
	if req.FromOwner {
		url += "&is_by_homeowner=1"
	}
	if req.MaxPrice > 0 {
		url += "&maxprice=" + strconv.FormatUint(req.MaxPrice, 10)
	}
	if req.MinPrice > 0 {
		url += "&minprice=" + strconv.FormatUint(req.MinPrice, 10)
	}
	if req.LastUpdate > 0 {
		url += "&totime=" + strconv.FormatUint(req.LastUpdate, 10)
	}
	for _, n := range req.Rooms {
		url += "room" + strconv.FormatUint(uint64(n), 10) + "=1"
	}
	return url
}

func (parser *ParserCian) parsePage(url *string, flatsChan chan []flat.Flat) {
	resp, err := http.Get(*url)
	flats := make([]flat.Flat, 0)
	defer func() { flatsChan <- flats }()

	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return
	}

	// Needed to ignore random numbers in classes,
	// because probably they are used to prevent parsing
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		class, ok := s.Attr("class")
		if ok {
			class = strings.Replace(class, "--", " ", -1)
			s.SetAttr("class", class)
		}
	})

	doc.Find(".main").Each(func(i int, s *goquery.Selection) {
		flat := flat.Flat{}
		url, _ := s.Find("a").Attr("href")
		flat.URL = url

		findAndFilter := func(selector string, pattern string) string {
			re, _ := regexp.Compile(pattern)
			text := s.Find(selector).Text()
			return re.ReplaceAllString(text, "")
		}

		title := findAndFilter(".single_title", "[^0-9 ]")
		values := strings.Fields(title)

		if len(values) < 2 {
			title = findAndFilter(".subtitle", "[^0-9 ]")
			values = strings.Fields(title)
			if len(values) < 2 {
				return
			}
		}

		flat.Rooms = uint8(ParseUintOrDefault(values[0]))
		flat.Area = ParseUintOrDefault(values[1])

		price := findAndFilter(".header", "[^0-9]")
		flat.Price = ParseUintOrDefault(price)

		flat.Address = s.Find(".address-links").Text()

		if strings.ToLower(s.Find(".badge-container").Text()) == "собственник" {
			flat.FromOwner = true
		}
		flats = append(flats, flat)
	})
}
