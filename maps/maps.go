package maps

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/the-fusy/rentit/config"
)

func parseFromHere(URL *string, params *url.Values, result interface{}) error {
	hereURL, _ := url.Parse(*URL)
	params.Set("app_id", config.HereAppID)
	params.Set("app_code", config.HereAppCode)
	params.Set("jsonAttributes", "0")
	hereURL.RawQuery = params.Encode()

	res, err := http.Get(hereURL.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, result)
	if err != nil {
		return err
	}

	return err
}

func GetCoordinates(address *string) (latitude float64, longitude float64, err error) {
	hereURL := "https://geocoder.api.here.com/6.2/geocode.json"
	params := url.Values{
		"searchtext": []string{*address},
	}
	result := geocode{}

	err = parseFromHere(&hereURL, &params, &result)
	if err != nil {
		return 0, 0, err
	}

	err = errors.New("Error to get coordinates")

	view := &result.Response.View
	if len(*view) == 0 {
		return 0, 0, err
	}

	viewResult := &(*view)[0].Result
	if len(*viewResult) == 0 {
		return 0, 0, err
	}

	navigationPosition := &(*viewResult)[0].Location.NavigationPosition
	if len(*navigationPosition) == 0 {
		return 0, 0, err
	}

	latitude = (*navigationPosition)[0].Latitude
	longitude = (*navigationPosition)[0].Longitude
	if latitude != 0.0 && longitude != 0.0 {
		err = nil
	}

	return latitude, longitude, err
}

type geocode struct {
	Response struct {
		View []struct {
			Result []struct {
				Location struct {
					NavigationPosition []struct {
						Latitude  float64
						Longitude float64
					}
				}
			}
		}
	}
}

func GetTravelTime(from, to Place) (int16, error) {
	hereURL := "https://route.api.here.com/routing/7.2/calculateroute.json"
	params := url.Values{
		"mode":           []string{"fastest;publicTransport"},
		"representation": []string{"overview"},
		"waypoint0":      []string{fmt.Sprintf("geo!%v,%v", from.Latitude, from.Longitude)},
		"waypoint1":      []string{fmt.Sprintf("geo!%v,%v", to.Latitude, to.Longitude)},
	}
	result := calculateroute{}

	err := parseFromHere(&hereURL, &params, &result)
	if err != nil {
		return 0, err
	}

	err = errors.New("Error to get travel time")

	route := result.Response.Route
	if len(route) == 0 {
		return 0, err
	}

	travelTime := route[0].Summary.TravelTime
	if travelTime != 0 {
		err = nil
	}

	return travelTime, err
}

type Place struct {
	Latitude  float64
	Longitude float64
}

type calculateroute struct {
	Response struct {
		Route []struct {
			Summary struct {
				TravelTime int16
			}
		}
	}
}
