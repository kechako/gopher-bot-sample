package rainfall

import (
	"errors"
	"strings"
)

type CommandType string

const (
	CommandTypeAddLocation    CommandType = "add"
	CommandTypeRemoveLocation CommandType = "rm"
	CommandTypeChangeLocation CommandType = "change"
	CommandTypeListLocation   CommandType = "list"
	CommandTypeNone           CommandType = "none"
)

var (
	CommandSyntaxError = errors.New("CommandSyntaxError")
)

type Commander interface {
	Execute(params []string) (string, error)
}

type Command struct {
	p            *plugin
	commanderMap map[CommandType]Commander
}

func NewCommand(p *plugin) *Command {
	commanderMap := map[CommandType]Commander{
		CommandTypeAddLocation:    NewAddCommand(p),
		CommandTypeRemoveLocation: NewRemoveCommand(p),
		CommandTypeChangeLocation: NewChangeCommand(p),
		CommandTypeListLocation:   NewListCommand(p),
		CommandTypeNone:           NewAskCommand(p),
	}

	return &Command{
		p:            p,
		commanderMap: commanderMap,
	}
}

func (c *Command) Execute(message string) (string, error) {
	params := make([]string, 0, 4)
	for _, param := range strings.Split(message, " ") {
		if param == "" {
			continue
		}
		params = append(params, param)
	}

	cmdType := CommandTypeNone
	if len(params) > 0 {
		cmdType = CommandType(params[0])
	}

	commander, ok := c.commanderMap[cmdType]
	if !ok {
		cmdType = CommandTypeNone
		commander = c.commanderMap[cmdType]
	}

	return commander.Execute(params)
}
