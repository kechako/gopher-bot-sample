package rainfall

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/kechako/yolp"
	"github.com/pkg/errors"
)

type askCommand struct {
	p *plugin
}

func NewAskCommand(p *plugin) Commander {
	return &askCommand{
		p: p,
	}
}

func (c *askCommand) Execute(params []string) (string, error) {
	loc, err := c.getLocation(params)
	if err != nil {
		return "", err
	}

	weathers, err := c.getWeathers(loc)
	if err != nil {
		return "", errors.Wrap(err, "取得失敗")
	}

	// 直近の天気情報
	w := getMostRecentWeather(weathers)

	var result string
	if w.IsRaining() {
		if w.IsObservation() {
			result = "雨降ってます"
		} else {
			result = "雨降ってるかも"
		}
	} else {
		if w.IsObservation() {
			result = "雨降ってないです"
		} else {
			result = "雨降ってないかも"
		}
	}

	resMessage := bytes.Buffer{}
	resMessage.WriteString(result)
	resMessage.WriteString("  ")
	resMessage.WriteString(getWeatherString(w))

	for _, w := range weathers {
		resMessage.WriteString("\n")
		resMessage.WriteString(getWeatherString(w))
	}

	return resMessage.String(), nil
}

func (c *askCommand) getLocation(params []string) (Location, error) {
	var loc Location

	switch len(params) {
	case 1:
		// by name
		loc, ok := c.p.locStore.Get(params[0])
		if !ok {
			return loc, errors.New("指定された名前のロケーションは未登録です。")
		}
	case 2:
		// latitude and longitude

		lat, err := strconv.ParseFloat(params[0], 32)
		if err != nil {
			return loc, CommandSyntaxError
		}
		loc.Latitude = float32(lat)

		long, err := strconv.ParseFloat(params[1], 32)
		if err != nil {
			return loc, CommandSyntaxError
		}
		loc.Longitude = float32(long)
	default:
		return loc, CommandSyntaxError
	}

	return loc, nil
}

func (c *askCommand) getWeathers(loc Location) (weathers []yolp.Weather, err error) {
	y := yolp.NewYOLP(c.p.appID)

	ydf, err := y.Place(loc.Latitude, loc.Longitude)
	if err != nil {
		return
	}

	if len(ydf.Feature) == 0 {
		err = fmt.Errorf("Could not get the weather data from the API response.")
		return
	}

	weathers = ydf.Feature[0].Property.WeatherList.Weather
	if len(weathers) == 0 {
		err = fmt.Errorf("Could not get the weather data from the API response.")
		return
	}

	return
}

func getMostRecentWeather(weathers []yolp.Weather) (weather yolp.Weather) {
	now := time.Now()

	var minDuration int64
	for i, w := range weathers {
		d := abs64(int64(now.Sub(w.Time())))
		if i == 0 || d < minDuration {
			minDuration = d
			weather = w
		}
	}

	return
}

func abs64(n int64) int64 {
	if n < 0 {
		return -n
	}

	return n
}

func getWeatherString(w yolp.Weather) string {
	str := fmt.Sprintf("[%s]  %.2f mm", w.Time().Format("15:04"), w.Rainfall)
	if w.IsObservation() {
		return str + "  (実測値)"
	} else if w.IsForecast() {
		return str + "  (予測値)"
	}
	return str
}
