package maps

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/the-fusy/rentit/config"
)

func GetCoordinates(address *string) (latitude float64, longitude float64, err error) {
	hereURL, _ := url.Parse("https://geocoder.api.here.com/6.2/geocode.json")
	hereParams := url.Values{}
	hereParams.Set("app_id", config.HereAppID)
	hereParams.Set("app_code", config.HereAppCode)
	hereParams.Set("searchtext", *address)
	hereURL.RawQuery = hereParams.Encode()

	res, err := http.Get(hereURL.String())
	if err != nil {
		return
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	result := geocode{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return
	}

	err = errors.New("Wrong address")
	view := &result.Response.View
	if len(*view) == 0 {
		return
	}
	viewResult := &(*view)[0].Result
	if len(*viewResult) == 0 {
		return
	}
	navigationPosition := &(*viewResult)[0].Location.NavigationPosition
	if len(*navigationPosition) == 0 {
		return
	}
	err = nil

	latitude = (*navigationPosition)[0].Latitude
	longitude = (*navigationPosition)[0].Longitude

	return
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
