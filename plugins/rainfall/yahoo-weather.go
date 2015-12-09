package rainfall

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type YDF struct {
	ResultInfo ResultInfo
	Feature    []Feature
}

type ResultInfo struct {
	Count       int
	Total       int
	Start       int
	Status      int
	Latency     float32
	Description string
	Copyright   string
}

type Feature struct {
	Id       string
	Name     string
	Geometry Geometry
	Property Property
}

type Geometry struct {
	Type        string
	Coordinates string
}

type Property struct {
	WeatherAreaCode int
	WeatherList     WeatherList
}

type WeatherList struct {
	Weather []Weather
}

type Weather struct {
	Type     string
	Date     string
	Rainfall float32
}

type yahooWeather struct {
	appID string
}

func NewYahooWeather(appID string) *yahooWeather {
	return &yahooWeather{
		appID: appID,
	}
}

const (
	yahooAPIPlaceUrl         = "http://weather.olp.yahooapis.jp/v1/place"
	yahooAPISearchZipCodeUrl = "http://search.olp.yahooapis.jp/OpenLocalPlatform/V1/zipCodeSearch"
)

func (y *yahooWeather) Place(latitude float32, longitude float32) (*YDF, error) {
	query := map[string]string{
		"coordinates": fmt.Sprintf("%f,%f", latitude, longitude),
	}

	return y.apiGet(y.makeUrl(yahooAPIPlaceUrl, query))
}

func (y *yahooWeather) SearchZipCode(zipCode string) (*YDF, error) {
	query := map[string]string{
		"query": zipCode,
	}

	return y.apiGet(y.makeUrl(yahooAPISearchZipCodeUrl, query))
}

func (y *yahooWeather) apiGet(url string) (*YDF, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	d := json.NewDecoder(res.Body)
	ydf := &YDF{}
	err = d.Decode(ydf)
	if err != nil {
		return nil, err
	}

	return ydf, nil
}

func (y *yahooWeather) makeUrl(baseUrl string, query map[string]string) string {
	u, err := url.Parse(baseUrl)
	if err != nil {
		log.Panic(err)
	}

	q := u.Query()

	for key, value := range query {
		q.Add(key, value)
	}
	q.Add("output", "json")
	q.Add("appid", y.appID)
	u.RawQuery = q.Encode()

	return u.String()
}
