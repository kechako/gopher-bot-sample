package rainfall

import (
	"fmt"
	"io"
	"strings"

	"github.com/kyokomi/slackbot/plugins"
)

const (
	RainfallPrefix = "rainfall"
)

type plugin struct {
	appID    string
	locStore *LocationStore
}

type BotMessagePluginCloser interface {
	plugins.BotMessagePlugin
	io.Closer
}

func NewPlugin(appID string, path string) (BotMessagePluginCloser, error) {
	locStore := NewLocationStore(path)
	err := locStore.Load()
	if err != nil {
		return nil, err
	}

	return &plugin{
		appID:    appID,
		locStore: locStore,
	}, nil
}

func (p *plugin) Close() error {
	return p.locStore.Save()
}

func (p *plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	cmds := strings.SplitN(message, " ", 2)
	if len(cmds) == 0 || cmds[0] != RainfallPrefix {
		return false, message
	}

	return true, cmds[1]
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	result, err := p.ExecuteCommand(message)
	if err != nil {
		if err == CommandSyntaxError {
			event.Reply(p.buildHelp())
		} else {
			event.Reply(err.Error())
		}
		return true
	}

	event.Reply(result)

	return true
}

func (p *plugin) Help() string {
	return `rainfall: 雨チェック
	指定された座標で雨が降っているかどうか表示します。

	rainfall <latitude> <longitude>
	    指定された座標で雨が降っているかどうか表示します。

	rainfall <name>
	    指定された名前の座標で雨が降っているかどうか表示します。

	rainfall add <name> <latitude> <longitude>
	    指定された名前で座標を登録します。

	rainfall change <name> <latitude> <longitude>
	    指定された名前の座標を変更します。

	rainfall rm <name>
	    指定された名前の座標を削除します。
    `
}

func (p *plugin) buildHelp() string {
	return fmt.Sprintf("```\n%s\n```", p.Help())
}
