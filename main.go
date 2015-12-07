package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/kechako/gopher-bot/plugins/iyagoza"
	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins/akari"
	"github.com/kyokomi/slackbot/plugins/lgtm"
	"github.com/kyokomi/slackbot/plugins/naruhodo"
	"github.com/kyokomi/slackbot/plugins/suddendeath"
)

func main() {
	var token string
	flag.StringVar(&token, "token", os.Getenv("SLACK_BOT_TOKEN"), "Bot token.")
	flag.Parse()

	bot, err := slackbot.NewBotContext(token)
	if err != nil {
		panic(err)
	}
	// いやでござる
	bot.AddPlugin("iyagoza", iyagoza.NewPlugin())
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
