package iyagoza

import (
	"fmt"
	"strings"

	"github.com/kyokomi/slackbot/plugins"
)

type plugin struct {
}

func NewPlugin() plugins.BotMessagePlugin {
	return &plugin{}
}

func (p *plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	return checkReplyMessage(event.BotID(), message)
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	event.Reply(fmt.Sprintf("<@%s>: いやでござる", event.SenderID()))
	return true
}

func (p *plugin) Help() string {
	return `iyagoza:
	reply 'いやでござる'
    `
}

func checkReplyMessage(botID string, message string) (bool, string) {
	keyword := fmt.Sprintf("<@%s>", botID)
	return strings.Index(message, keyword) >= 0, message
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
