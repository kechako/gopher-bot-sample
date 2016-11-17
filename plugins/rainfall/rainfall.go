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
	cmd      *Command
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

	p := &plugin{
		appID:    appID,
		locStore: locStore,
	}

	p.cmd = NewCommand(p)

	return p, nil
}

func (p *plugin) Close() error {
	return p.locStore.Save()
}

func (p *plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	cmds := strings.SplitN(message, " ", 2)
	if len(cmds) == 0 || cmds[0] != RainfallPrefix {
		return false, message
	}

	if len(cmds) == 1 {
		return true, ""
	}

	return true, cmds[1]
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	result, err := p.cmd.Execute(message)
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

	rainfall list
	    登録された座標を一覧表示します。
    `
}

func (p *plugin) buildHelp() string {
	return fmt.Sprintf("```\n%s\n```", p.Help())
}
