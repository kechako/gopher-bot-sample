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

const yahooAPIUrl = "http://weather.olp.yahooapis.jp/v1/place"

func (y *yahooWeather) Place(latitude float32, longitude float32) (*YDF, error) {

	res, err := http.Get(makeUrl(y.appID, latitude, longitude))
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

func makeUrl(appID string, latitude float32, longitude float32) string {
	u, err := url.Parse(yahooAPIUrl)
	if err != nil {
		log.Panic(err)
	}

	q := u.Query()
	q.Add("coordinates", fmt.Sprintf("%f,%f", latitude, longitude))
	q.Add("output", "json")
	q.Add("appid", appID)
	u.RawQuery = q.Encode()

	return u.String()
}
