package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/kechako/gopher-bot/plugins/iyagoza"
	"github.com/kechako/gopher-bot/plugins/rainfall"
	"github.com/kechako/gopher-bot/plugins/zundoko"
	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins/akari"
	"github.com/kyokomi/slackbot/plugins/lgtm"
	"github.com/kyokomi/slackbot/plugins/naruhodo"
	"github.com/kyokomi/slackbot/plugins/suddendeath"
)

func main() {
	var token, appId, rainfallPath string
	flag.StringVar(&token, "token", os.Getenv("SLACK_BOT_TOKEN"), "Bot token.")
	flag.StringVar(&appId, "appid", os.Getenv("YAHOO_APP_ID"), "Yahoo App Id.")
	flag.StringVar(&rainfallPath, "rainfall-path", os.Getenv("RAINFALL_PATH"), "Rainfall plugin data store path.")
	flag.Parse()

	bot, err := slackbot.NewBotContext(token)
	if err != nil {
		panic(err)
	}
	// いやでござる
	bot.AddPlugin("iyagoza", iyagoza.NewPlugin())
	// 雨
	r, err := rainfall.NewPlugin(appId, rainfallPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer r.Close()
	bot.AddPlugin("rainfall", r)
	// ズンドコキヨシ
	bot.AddPlugin("zundoko", zundoko.NewPlugin())
	// あかり大好き
	bot.AddPlugin("akari", akari.NewPlugin())
	// LGTM
	bot.AddPlugin("lgtm", lgtm.NewPlugin())
	// なるほど
	bot.AddPlugin("naruhodo", naruhodo.NewPlugin())
	// 突然の死
	bot.AddPlugin("suddendeath", suddendeath.NewPlugin())

	bot.WebSocketRTM()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(":8000", nil)
}
