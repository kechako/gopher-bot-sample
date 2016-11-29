package ppap

import (
	"bytes"
	"math/rand"
	"time"

	"github.com/kyokomi/slackbot/plugins"
)

const (
	pen        = "\u2712\ufe0f"
	pineapple  = "\U0001f34d"
	apple      = "\U0001F34E"
	ppapFinish = "ペンパイナッポーアッポーペン"
)

var (
	keywords = []string{"ppap", "PPAP"}
	random   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type plugin struct {
}

func NewPlugin() plugins.BotMessagePlugin {
	return &plugin{}
}

func (p *plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	return plugins.CheckMessageKeywords(message, keywords...)
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	ppapWords := [3]string{pen, pineapple, apple}
	good := [4]string{pen, pineapple, apple, pen}

	var current [4]string

	reply := bytes.NewBuffer(make([]byte, 0, 1024))

	for current != good {
		shift(&current)
		pa := ppapWords[random.Intn(3)]
		current[3] = pa
		reply.WriteString(pa)
	}

	reply.WriteString(ppapFinish)

	event.Reply(reply.String())

	return true
}

func (p *plugin) Help() string {
	return `PPAP:
	` + "\u2712\ufe0f\U0001f34d\U0001F34E" + `
    `
}

func shift(a *[4]string) {
	a[0], a[1], a[2], a[3] = a[1], a[2], a[3], ""
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
