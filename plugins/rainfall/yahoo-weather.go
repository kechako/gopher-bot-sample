package rainfall

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type WeatherType string

const (
	Observation WeatherType = "observation"
	Forecast    WeatherType = "forecast"
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

func (w *Weather) IsObservation() bool {
	return w.Type == "observation"
}

func (w *Weather) IsForecast() bool {
	return w.Type == "forecast"
}

func (w *Weather) IsRaining() bool {
	return w.Rainfall > 0
}

func (w *Weather) Time() time.Time {
	var year, month, day, hour, min int

	fmt.Println(w.Date)
	_, err := fmt.Sscanf(w.Date, "%4d%2d%2d%2d%2d", &year, &month, &day, &hour, &min)
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}

	return time.Date(year, time.Month(month), day, hour, min, 0, 0, time.Local)
}

func (w *Weather) String() string {
	str := fmt.Sprintf("[%s]  %.2f mm", w.Time().Format("2006-01-02 15:04"), w.Rainfall)
	if w.IsObservation() {
		return str + "  (実測値)"
	} else if w.IsForecast() {
		return str + "  (予測値)"
	}
	return str
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
		"interval":    "5",
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
