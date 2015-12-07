package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/kechako/gopher-bot/plugins/iyagoza"
	"github.com/kyokomi/slackbot"
)

func main() {
	var token string
	flag.StringVar(&token, "token", os.Getenv("SLACK_BOT_TOKEN"), "Bot token.")
	flag.Parse()

	bot, err := slackbot.NewBotContext(token)
	if err != nil {
		panic(err)
	}
	bot.AddPlugin("iyagoza", iyagoza.NewPlugin())

	bot.WebSocketRTM()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(":8000", nil)
}
