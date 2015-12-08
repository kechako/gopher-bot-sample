package rainfall

import (
	"fmt"
	"os"
	"regexp"
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
	w, err := getWeather(message)
	if err != nil {
		event.Reply("取得失敗")
	} else if w.Rainfall > 0 {
		event.Reply(fmt.Sprintf("雨降ってます (%f mm)", w.Rainfall))
	} else {
		event.Reply("雨降ってないです")
	}

	return true
}

func (p *plugin) Help() string {
	return `rainfall: 雨チェック
	指定された座標で雨が降っているかどうか表示します。

	rainfall <latitude> <longitude>
    `
}

func checkReplyMessage(botID string, message string) (bool, string) {
	keyword := fmt.Sprintf("<@%s>", botID)
	return strings.Index(message, keyword) >= 0, message
}

func getWeather(message string) (weather Weather, err error) {
	group := messageFormat.FindStringSubmatch(message)[1:]

	var latitude float32
	_, err = fmt.Sscan(group[0], &latitude)
	if err != nil {
		return
	}
	var longitude float32
	_, err = fmt.Sscan(group[1], &longitude)
	if err != nil {
		return
	}

	yw := NewYahooWeather(appID)

	yfd, err := yw.Place(latitude, longitude)
	if err != nil {
		return
	}

	err = fmt.Errorf("Error")
	for _, w := range yfd.Feature[0].Property.WeatherList.Weather {
		if w.Type == "observation" {
			weather = w
			err = nil
			break
		}
	}

	return
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
