package rainfall

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

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
	w, ok := getObservationWeather(weathers)

	if !ok {
		messages = append(messages, "実測値が取得できません")
	} else if w.Rainfall > 0 {
		messages = append(messages, "雨降ってます")
	} else {
		messages = append(messages, "雨降ってないです")
	}

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

func getObservationWeather(weathers []Weather) (w Weather, ok bool) {
	for _, w = range weathers {
		if w.IsObservation() {
			ok = true
			return
		}
	}

	return
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
