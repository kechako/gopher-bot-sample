package rainfall

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kyokomi/slackbot/plugins"
)

var (
	messageFormat = regexp.MustCompile(`^rainfall\s([+-]?\d+\.\d+)\s+([+-]?\d+\.\d+)`)
	appID         string
)

func init() {
	appID = os.Getenv("YAHOO_APP_ID")
}

type plugin struct {
}

func NewPlugin() plugins.BotMessagePlugin {
	return &plugin{}
}

func (p *plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	return messageFormat.MatchString(message), message
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	weathers, err := getWeathers(message)
	if err != nil {
		event.Reply(fmt.Sprintf("取得失敗 : %v", err))
	}

	messages := make([]string, 0, 10)

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
	messages = append(messages, result+"  "+w.String())

	for _, w := range weathers {
		messages = append(messages, w.String())
	}

	event.Reply(strings.Join(messages, "\n"))

	return true
}

func (p *plugin) Help() string {
	return `rainfall: 雨チェック
	指定された座標で雨が降っているかどうか表示します。

	rainfall <latitude> <longitude>
    `
}

func getWeathers(message string) (weathers []Weather, err error) {
	group := messageFormat.FindStringSubmatch(message)[1:]

	latitude, err := strconv.ParseFloat(group[0], 32)
	if err != nil {
		return
	}
	longitude, err := strconv.ParseFloat(group[1], 32)
	if err != nil {
		return
	}

	yw := NewYahooWeather(appID)

	ydf, err := yw.Place(float32(latitude), float32(longitude))
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

func getMostRecentWeather(weathers []Weather) (weather Weather) {
	now := time.Now()

	var minDuration int64
	for i, w := range weathers {
		d := Abs64(int64(now.Sub(w.Time())))
		if i == 0 || d < minDuration {
			minDuration = d
			weather = w
		}
	}

	return
}

func Abs64(n int64) int64 {
	if n < 0 {
		return -n
	}

	return n
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
